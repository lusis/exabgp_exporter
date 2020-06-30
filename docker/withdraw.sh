#!/bin/bash

echo "neighbor 127.0.0.1 withdraw route 192.168.88.0/24 next-hop 192.168.1.2 split /29" > /tmp/exabgp.cmd
echo "neighbor 127.0.0.1 withdraw route 192.168.0.0/24 next-hop 192.168.1.2" > /tmp/exabgp.cmd
echo "neighbor 127.0.0.1 withdraw route 2001:db8:1000::/64 next-hop 2001:db8:ffff::1" > /tmp/exabgp.cmd
