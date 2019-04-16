#!/bin/bash
## start rsyslog
nohup ./rsyslog.sh &
sleep 2
## start exabgp so we get peer errors due to gobgpd not started yet
nohup ./exabgp.sh &
sleep 5
## start gobgpd
nohup ./gobgpd.sh &
tail -f /var/log/syslog