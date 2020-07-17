package exporter

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"

	"github.com/lusis/exabgp_exporter/pkg/exabgp/messages/text"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	showAdjRibSubcommand  = []string{"show", "adj-rib", "out", "extensive"}
	showSummarySubcommand = []string{"show", "neighbor", "summary"}
)

// StandaloneExporter is a prometheus exporter that gathers metrics via calling exabgpcli
type StandaloneExporter struct {
	ExaBGPCLI  string
	ExaBGPRoot string
	mutex      sync.RWMutex
	BaseExporter
}

// NewStandaloneExporter returns an initialized TextExporter.
func NewStandaloneExporter(exabgpcli string, exabgproot string) (*StandaloneExporter, error) {
	be := NewBaseExporter()
	return &StandaloneExporter{
		ExaBGPCLI:    exabgpcli,
		ExaBGPRoot:   exabgproot,
		BaseExporter: be,
	}, nil
}

// Describe describes all the metrics ever exported by the exabgp exporter
// It implements prometheus.Collector
func (e *StandaloneExporter) Describe(ch chan<- *prometheus.Desc) {
	e.BaseExporter.Describe(ch)
}

// Collect fetches the stats from configured exabpcli command and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *StandaloneExporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	e.BaseExporter.totalScrapes.Inc()
	ribs, peers, err := e.scrape(ch)
	if err != nil {
		log.Error(err.Error())
	} else {
		for _, u := range peers {
			desc := newSummaryMetric("peer")
			isUp := 0
			if u.Status != "down" {
				isUp = 1
			}
			m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(isUp), u.IPAddress, u.AS)
			ch <- m
		}
		for _, r := range ribs {
			switch r.Family() {
			case "ipv4 unicast":
				v4u, _ := r.IPv4Unicast()
				desc := newRibMetric("route")
				m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), r.PeerIP, r.PeerAS, r.LocalIP, r.LocalAS, v4u.NLRI, r.Family())
				ch <- m
			case "ipv6 unicast":
				v6u, _ := r.IPv6Unicast()
				desc := newRibMetric("route")
				m := prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), r.PeerIP, r.PeerAS, r.LocalIP, r.LocalAS, v6u.NLRI, r.Family())
				ch <- m
			default:
				log.Errorf("unable to handle family %s", r.Family())
			}
		}
	}

	ch <- e.BaseExporter.totalScrapes
	ch <- e.BaseExporter.parseFailures
}

func (e *StandaloneExporter) scrape(ch chan<- prometheus.Metric) ([]*text.RIBMessage, []*text.NeighborSummary, error) {
	var ns []*text.NeighborSummary
	var rs []*text.RIBMessage

	res, err := e.getSummary()
	if err != nil {
		e.BaseExporter.setExabgpStatus(ch, 0)
		e.BaseExporter.parseFailures.Inc()
		return rs, ns, fmt.Errorf("stdout: %s, error: %s", string(res), err.Error())
	}
	e.BaseExporter.setExabgpStatus(ch, 1)
	status, err := text.SummariesFromBytes(res)
	if err != nil {
		e.BaseExporter.parseFailures.Inc()
		return rs, ns, err
	}
	ribres, riberr := e.getRIB()
	if riberr != nil {
		e.BaseExporter.setExabgpStatus(ch, 0)
		e.BaseExporter.parseFailures.Inc()
		return rs, ns, fmt.Errorf("stdout: %s, error: %s", string(ribres), riberr.Error())
	}
	ribs, ribserr := text.RibFromBytes(ribres)
	if ribserr != nil {
		return rs, ns, ribserr
	}
	return ribs, status, nil
}
func (e *StandaloneExporter) getSummary() ([]byte, error) {
	return e.runExaBGPCLI(showSummarySubcommand)
}

func (e *StandaloneExporter) getRIB() ([]byte, error) {
	return e.runExaBGPCLI(showAdjRibSubcommand)
}

func (e *StandaloneExporter) runExaBGPCLI(subcommand []string) ([]byte, error) {
	args := []string{"--root", e.ExaBGPRoot}
	args = append(args, subcommand...)
	cmd := exec.Command(e.ExaBGPCLI, args...)
	var se, so bytes.Buffer
	cmd.Stderr = &se
	cmd.Stdout = &so
	err := cmd.Run()
	return so.Bytes(), err
}
