# TODO items and notes

## Add flows to exported stats

We capture the data now in the global announcement/withdraw tracking.
It just needs to have an Gauge added for it

## Consider reworking how we track routes being available

So this is a slight limitation of exabgp.
Announce messages look like this

```json
{
    "exabgp": "4.0.1",
    "time": 1554843223.5592246,
    "host": "node1",
    "pid": 31372,
    "ppid": 1,
    "counter": 11,
    "type": "update",
    "neighbor": {
        "address": {
            "local": "192.168.1.184",
            "peer": "192.168.1.2"
        },
        "asn": {
            "local": 64496,
            "peer": 64496
        },
        "direction": "send",
        "message": {
            "update": {
                "attribute": {
                    "origin": "igp",
                    "med": 100,
                    "local-preference": 100
                },
                "announce": {
                    "ipv4 unicast": {
                        "192.168.1.184": ["192.168.88.2/32"]
                    }
                }
            }
        }
    }
}
```

The indexing of a route announce is based on the `ipv4 unicast` key under `announce`.
However, withdraw messages don't have that information:

```json
{
    "exabgp": "4.0.1",
    "time": 1554843266.7842073,
    "host": "node1",
    "pid": 31372,
    "ppid": 1,
    "counter": 14,
    "type": "update",
    "neighbor": {
        "address": {
            "local": "192.168.1.184",
            "peer": "192.168.1.2"
        },
        "asn": {
            "local": 64496,
            "peer": 64496
        },
        "direction": "send",
        "message": {
            "update": {
                "attribute": {
                    "origin": "igp",
                    "local-preference": 100
                },
                "withdraw": {
                    "ipv4 unicast": ["192.168.88.2/32"]
                }
            }
        }
    }
}
```

Under the `withdraw` key we only have the network being withdrawn.

However, if you announce the same network with two different `next-hop`s, those are distinct routing entries on the peer.
In essence, the peer still has a route you announced to that network (just under a different `next-hop`)

A better solution may be change this:

```go
                announcements := evt.GetAnnouncements()
                if announcements != nil {
                    for _, v := range announcements.IPV4Unicast {
                        labels["family"] = "ipv4 unicast"
                        for _, r := range v.Routes {
                            labels["route"] = r
                            exporter.PeerRouteStateMetric.With(labels).Set(float64(1))
                        }
                    }
                }

                withdraws := evt.GetWithdrawals()
                if withdraws != nil {
                    for _, w := range withdraws.IPv4Unicast {
                        for _, r := range w.Routes {
                            labels["family"] = "ipv4 unicast"
                            labels["route"] = r
                            exporter.PeerRouteStateMetric.With(labels).Set(float64(0))
                        }
                    }
                }
```

to this:

```go
                announcements := evt.GetAnnouncements()
                if announcements != nil {
                    for _, v := range announcements.IPV4Unicast {
                        labels["family"] = "ipv4 unicast"
                        for _, r := range v.Routes {
                            labels["route"] = r
                            exporter.PeerRouteStateMetric.With(labels).Add(float64(1))
                        }
                    }
                }

                withdraws := evt.GetWithdrawals()
                if withdraws != nil {
                    for _, w := range withdraws.IPv4Unicast {
                        for _, r := range w.Routes {
                            labels["family"] = "ipv4 unicast"
                            labels["route"] = r
                            exporter.PeerRouteStateMetric.With(labels).Dec(float64(0))
                        }
                    }
                }
```

The challenge is we can't really tell which next hop is being decremented because we simply don't have that information.
It probably doesn't matter if all you care about is if this exabgp node is announcing at least `n` paths to a given network.