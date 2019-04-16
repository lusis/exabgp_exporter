# exabgp_exporter

This is a prometheus exporter for exabgp. It currently ONLY works with exabgp 4.0.10 as that had some critical fixes for json encoding.

## Usage

The exporter runs as an api process inside exabgp.

It currently only captures ipv4 route announce/withdraw as well as peer error counters. Support for flows is in place but not being exported yet.

### exabgp.conf

You need to add some stuff to your exabgp.conf:

```text
process prometheus_exporter {
        run /path/to/exabgp_exporter;
        encoder json;
}
```

Either in your template or neighbor definition you'll need to add a new api section

```text
        api logging {
                processes [ prometheus_exporter ];
                neighbor-changes;
                signal;
                receive {
                        parsed;
                        notification;
                        update;
                        refresh;
                }
                send {
                        parsed;
                        notification;
                        update;
                        refresh;
                }
        }
```

## metrics

### `peer_resets`

```text
# HELP peer_resets tracks number of resets communicating with a peer
# TYPE peer_resets gauge
peer_resets{peer_asn="64496",peer_ip="192.168.1.158",reason="peer reset, message (closing connection) error(Broken TCP connection)",self_asn="64496",self_ip="192.168.1.184"} 9
```

### `peer_state`

```text
# HELP peer_state shows the current peer state
# TYPE peer_state gauge
peer_state{peer_asn="64496",peer_ip="192.168.1.158",self_asn="64496",self_ip="192.168.1.184"} 1
```

`1` for up. `0` for down

### `peer_route_state`

```text
# HELP peer_route_state shows the current peer state for a given route
# TYPE peer_route_state gauge
peer_route_state{family="ipv4 unicast",peer_asn="64496",peer_ip="192.168.1.158",nlri="192.168.88.0/24",self_asn="64496",self_ip="192.168.1.184"} 1
```

`1` for up. `0` for down