#!/bin/bash
pip3 install exabgp==${EXABGP_VERSION} || exit 1
exec /init
