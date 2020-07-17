#!/bin/bash

# announce 192.168.88.0/24 as 32 /29s
echo "neighbor 127.0.0.1 announce route 192.168.88.0/24 next-hop 192.168.1.2 split /29" > /tmp/exabgp.cmd

# announce single ipv4 route
echo "neighbor 127.0.0.1 announce route 192.168.0.0/24 next-hop 192.168.1.2" > /tmp/exabgp.cmd

# announce single ipv6 route
echo "neighbor 127.0.0.1 announce route 2001:db8:1000::/64 next-hop 2001:db8:ffff::1" > /tmp/exabgp.cmd
