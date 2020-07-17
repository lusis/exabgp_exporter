package exporter

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lusis/exabgp_exporter/pkg/exabgp"
	"github.com/prometheus/common/log"
)

type EmbeddedExporter struct {
	mutex   sync.RWMutex
	summary *prometheus.GaugeVec
	rib     *prometheus.GaugeVec
	BaseExporter
}

func NewEmbeddedExporter() (*EmbeddedExporter, error) {
	be := NewBaseExporter()
	be.up.Set(float64(1))
	sm := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "peer",
		Namespace: namespace,
		Subsystem: "state",
		Help:      summaryHelp,
	}, summaryLabelNames)
	rm := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "route",
		Namespace: namespace,
		Subsystem: "state",
		Help:      ribHelp,
	}, ribLabelNames)
	prometheus.MustRegister(sm)
	prometheus.MustRegister(rm)
	return &EmbeddedExporter{
		summary:      sm,
		rib:          rm,
		BaseExporter: be,
	}, nil
}

// Run starts the background reader for populating metrics
func (e *EmbeddedExporter) Run(reader *bufio.Reader) {
	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err != nil && err != io.EOF {
				log.Errorf("unknown error: %s", err.Error())
				e.BaseExporter.parseFailures.Inc()
				continue
			}
			evt, err := exabgp.ParseEvent(line)
			if err != nil {
				log.Errorf("unable to parse line: %s", err.Error())
				log.Errorf("failed line: %s", line)
				e.BaseExporter.parseFailures.Inc()
				continue
			}
			var labels = map[string]string{
				"peer_ip":  evt.Peer.IP,
				"peer_asn": fmt.Sprintf("%d", evt.Peer.ASN),
			}
			switch evt.Peer.State {
			case "down":
				e.summary.With(labels).Set(float64(0))
			default:
				e.summary.With(labels).Set(float64(1))
			}
			if evt.Direction == "send" {
				announcements := evt.GetAnnouncements()
				if announcements != nil {
					labels["local_ip"] = evt.Self.IP
					labels["local_asn"] = fmt.Sprintf("%d", evt.Self.ASN)
					for _, v := range announcements.IPV4Unicast {
						labels["family"] = "ipv4 unicast"
						for _, r := range v.NLRI {
							labels["nlri"] = r
							e.rib.With(labels).Set(float64(1))
						}
					}
					for _, v := range announcements.IPV6Unicast {
						labels["family"] = "ipv6 unicast"
						for _, r := range v.NLRI {
							labels["nlri"] = r
							e.rib.With(labels).Set(float64(1))
						}
					}
				}
				withdraws := evt.GetWithdrawals()
				if withdraws != nil {
					labels["local_ip"] = evt.Self.IP
					labels["local_asn"] = fmt.Sprintf("%d", evt.Self.ASN)
					for _, w := range withdraws.IPv4Unicast {
						for _, r := range w.NLRI {
							labels["family"] = "ipv4 unicast"
							labels["nlri"] = r
							e.rib.With(labels).Set(float64(0))
						}
					}
					for _, w := range withdraws.IPv6Unicast {
						for _, r := range w.NLRI {
							labels["family"] = "ipv6 unicast"
							labels["nlri"] = r
							e.rib.With(labels).Set(float64(0))
						}
					}
				}
			}
		}
	}()
}

// Collect delivers all seen stats as Prometheus metrics
// It implements prometheus.Collector.
func (e *EmbeddedExporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.BaseExporter.totalScrapes.Inc()
	ch <- e.BaseExporter.totalScrapes
	ch <- e.BaseExporter.parseFailures
	ch <- e.BaseExporter.up
}

// Describe describes all the metrics ever exported by the exabgp exporter
// It implements prometheus.Collector
func (e *EmbeddedExporter) Describe(ch chan<- *prometheus.Desc) {
	e.BaseExporter.Describe(ch)
}
