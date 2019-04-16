#!/bin/bash
cd /gobgp || exit 1
# sleep so that gobgp isn't ready yet
sleep 5
exec ./gobgpd -f gobgp.yaml
