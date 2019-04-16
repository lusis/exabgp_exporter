#!/bin/bash
echo "neighbor 127.0.0.1 announce route 192.168.88.0/24 next-hop 192.168.1.2 split /29" > /tmp/exabgp.cmd 
