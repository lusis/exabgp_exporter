#!/bin/bash
# this is an `standalone` mode exporter that should NOT be able to successfully call exabgpcli listening on a different port
exec /exabgp/exabgp_exporter --web.listen-address=":9571" --log.format="logger:stderr?json=true" standalone