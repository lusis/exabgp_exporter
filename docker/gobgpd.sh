#!/bin/bash
cd /gobgp || exit 1
./gobgpd -f gobgp.yaml --disable-stdlog  --syslog yes --log-plain yes
