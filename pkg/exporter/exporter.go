package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PeerRouteStateMetric represents a prometheus metric for the current state of a route announcement to a peer
var PeerRouteStateMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "peer_route_state",
	Help: "shows the current peer state for a given route",
}, []string{"peer_ip", "self_ip", "peer_asn", "self_asn", "nlri", "family"})

// PeerStateMetric represents a prometheus metric for the current state of an exabgp neighbor
var PeerStateMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "peer_state",
	Help: "shows the current peer state",
}, []string{"peer_ip", "self_ip", "peer_asn", "self_asn"})

// PeerResetCountMetric tracks peer errors
var PeerResetCountMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "peer_resets",
	Help: "tracks number of resets communicating with a peer",
}, []string{"peer_ip", "self_ip", "peer_asn", "self_asn", "reason"})

func init() {
	prometheus.MustRegister(PeerStateMetric)
	prometheus.MustRegister(PeerRouteStateMetric)
	prometheus.MustRegister(PeerResetCountMetric)
}

// StartHandler starts the prometheus exporter
func StartHandler(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(port, nil)
}
