package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "exabgp"

var (
	upHelp            = `is exabgp up`
	upName            = `up`
	parseHelp         = `number of errors while parsing output`
	parseName         = `exporter_parse_failures`
	totalScrapesName  = `exporter_total_scrapes`
	totalScrapesHelp  = `current total exabgp scrapes`
	summaryHelp       = `shows the state of a bgp peer`
	summaryLabelNames = []string{"peer_ip", "peer_asn"}
	ribHelp           = `shows the state of a given nlri`
	ribLabelNames     = []string{"peer_ip", "peer_asn", "local_ip", "local_asn", "nlri", "family"}
	exabgpUp          = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of exabgp successful.", nil, nil)
)

func newSummaryMetric(metricName string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "state", metricName), summaryHelp, summaryLabelNames, nil)
}

func newRibMetric(metricName string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(namespace, "state", metricName), ribHelp, ribLabelNames, nil)
}

// BaseExporter is common data between the two types of exporters
type BaseExporter struct {
	up            prometheus.Gauge
	totalScrapes  prometheus.Counter
	parseFailures prometheus.Counter
}

// NewBaseExporter returns a BaseExporter for embedding
func NewBaseExporter() BaseExporter {
	return BaseExporter{
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      upName,
			Help:      upHelp,
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      totalScrapesName,
			Help:      totalScrapesHelp,
		}),
		parseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      parseName,
			Help:      parseHelp,
		}),
	}
}

// Describe describes all the metrics ever exported by the exabgp exporter.
// It implements prometheus.Collector.
func (e *BaseExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- exabgpUp
	ch <- e.totalScrapes.Desc()
	ch <- e.parseFailures.Desc()
}

func (e *BaseExporter) setExabgpStatus(ch chan<- prometheus.Metric, i int) {
	ch <- prometheus.MustNewConstMetric(exabgpUp, prometheus.GaugeValue, float64(i))
}
