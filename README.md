# exabgp_exporter

[![CircleCI](https://circleci.com/gh/lusis/exabgp_exporter.svg?style=svg)](https://circleci.com/gh/lusis/exabgp_exporter)

This is a prometheus exporter for exabgp. It currently works with the following exabgp versions (tested as part of CI):

- 4.0.4
- 4.0.5
- 4.0.6
- 4.0.8
- 4.0.9
- 4.0.10

installed from pypi under python3.

Additionally, we test can test exabgp master git but it is not run in CI

## Usage

The exporter can run in two modes:

### standalone

This is the default mode. Each scrape invokes `exabgpcli` twice - once to gather the outbound rib and again to get neighbor status

### stream

The exporter reads from stdin. This mode is appropriate for embedding inside the exabgp process itself as a `processs` definition

#### stream mode configuration

If you want to embed the exporter inside exabgp you'll need to make the following changes to your config:

```text
process prometheus_exporter {
        run /path/to/exabgp_exporter --log.format="logger:stderr?json=true" stream;
	# alternately run /path/to/exabgp_exporter --log-format="logger:syslog?appname=exabgp_exporter" stream;
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

### Similarities between the two modes

Both modes listen on the documented port of `9576`. The scraped output is the same between each.
Logging can be configured through the flags but if you're using the `stream` mode you should never set that to `stdout` as the data will attempt to be interpreted by exabgp.

### Differences between the modes

In `stream` mode, we see events as they happen. This means for routes we've seen we can explicitly mark them down if they are withdrawn (set the value to `0`)

In standalone mode, however, we rely on what data we can get from `exabgpcli`.
To get the rib, we call `exabgpcli show adj-rib out extensive`. This ONLY shows announced routes. If a route is withdrawn it simply doesn't get output.
In standalone mode, we do *NOT* maintain any state.

This means that detecting if a given announcement has been withdrawn means checking if a metric is present or not.

Having said that, realistically, you would need to be more explicit in your checks ANYWAY since either mode has no rib state at startup.
For rib entries you consider critical, you should match explicit on the labels.

## metrics

### `exabgp_up`

```text
# HELP exabgp_up is exabgp up
# TYPE exabgp_up gauge
exabgp_up 1
```

In `stream` mode, this is always `1` as we are likely embedded in the `exabgp` process itself.
In `standalone` mode, this is based on if `exabgpcli` exit code.

### `exabgp_exporter_parse_failures`

```text
# HELP exabgp_exporter_parse_failures number of errors while parsing output
# TYPE exabgp_exporter_parse_failures counter
exabgp_exporter_parse_failures 0
```

Tracks parsing failures in both `stream` and `standalone` mode. In `standalone` mode counter is increased for each `exabgpcli` invocation if parsing fails.

### `exabgp_exporter_total_scrapes`

```text
# HELP exabgp_exporter_total_scrapes current total exabgp scrapes
# TYPE exabgp_exporter_total_scrapes counter
exabgp_exporter_total_scrapes 3
```

Tracks scrapes of the exporter

### `exabgp_state_peer`

```text
# HELP exabgp_state_peer shows the state of a bgp peer
# TYPE exabgp_state_peer gauge
exabgp_state_peer{peer_asn="64496",peer_ip="127.0.0.1"} 1
```

Tracks the connectivity to BGP peers from exabgp. `1` for up. `0` for down.
In `standalone` mode, this is a result of calling `exabgpcli show neighbor summary`

### `exabgp_state_route`

```text
# HELP exabgp_state_route shows the state of a given nlri
# TYPE exabgp_state_route gauge
exabgp_state_route{family="ipv4 unicast",local_asn="64496",local_ip="127.0.0.1",nlri="192.168.88.0/29",peer_asn="64496",peer_ip="127.0.0.1"} 0
```

Tracks the state of a given nlri per family (formatted from `afi` + `safi`) for a given peer+local combination.
The cardinality is high here because exabgp can have multiple peers each with different local and peer ASNs.

`0` (or missing/stale) for down, `1` for up

***WARNING**

As stated above, this metric can be missing (either due to exporter restart or based on the mode you're using).
If you care about knowing the state of a specific set of labels, you should check not only that the value is `1` but also the age or even presence of the metric.
