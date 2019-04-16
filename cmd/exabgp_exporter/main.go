package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lusis/exabgp_exporter/pkg/exabgp"
	"github.com/lusis/exabgp_exporter/pkg/exporter"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			line, _, rerr := reader.ReadLine()
			if rerr != nil && rerr != io.EOF {
				log.Printf("error: %s", rerr.Error())
				continue
			}
			evt, perr := exabgp.ParseEvent(line)
			if perr != nil {
				log.Printf("unable to parse line: %s", perr.Error())
				log.Printf("failing line: %s", line)
				continue
			}
			// parse into metrics
			status := exabgp.GetStatus()
			var labels = map[string]string{
				"peer_ip":  evt.Peer.IP,
				"self_ip":  evt.Self.IP,
				"peer_asn": fmt.Sprintf("%d", evt.Peer.ASN),
				"self_asn": fmt.Sprintf("%d", evt.Self.ASN),
			}
			var metric int
			switch status {
			case "down":
				metric = 0
				if strings.Contains(exabgp.GetStatusReason(), "reset") {
					rlabels := map[string]string{}
					for k, v := range labels {
						rlabels[k] = v
					}
					rlabels["reason"] = exabgp.GetStatusReason()
					exporter.PeerResetCountMetric.With(rlabels).Add(float64(1))
				}
			case "unknown":
				metric = 0
			default:
				metric = 1
			}
			exporter.PeerStateMetric.With(labels).Set(float64(metric))

			// right now we only care about messages we send
			if evt.Direction == "send" {
				// uniqueness is peer_ip, local_ip, peer_as, remote_as, nlri (an nlri is essentially a given network/netmask)
				// we set it to either 0 or 1 based on withdrawn vs announced
				announcements := evt.GetAnnouncements()
				if announcements != nil {
					for _, v := range announcements.IPV4Unicast {
						labels["family"] = "ipv4 unicast"
						for _, r := range v.NLRI {
							labels["nlri"] = r
							exporter.PeerRouteStateMetric.With(labels).Set(float64(1))
						}
					}
				}

				withdraws := evt.GetWithdrawals()
				if withdraws != nil {
					for _, w := range withdraws.IPv4Unicast {
						for _, r := range w.NLRI {
							labels["family"] = "ipv4 unicast"
							labels["nlri"] = r
							exporter.PeerRouteStateMetric.With(labels).Set(float64(0))
						}
					}
				}
			}
		}
	}()
	// Start the exporter
	err := exporter.StartHandler(":9569")
	if err != nil {
		log.Fatalf("unable to start exporter: %s", err.Error())
	}
}
