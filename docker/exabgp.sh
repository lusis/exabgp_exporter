#!/bin/bash
env exabgp.log.destination=syslog \
    exabgp --root=/exabgp/ /exabgp/etc/exabgp/exabgp.conf