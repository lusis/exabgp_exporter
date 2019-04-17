#!/bin/bash
if [[ "${EXABGP_VERSION}" == "master" ]]; then
	echo "Installing exabgp from git master"
	cd /tmp || exit 1
	git clone https://github.com/Exa-Networks/exabgp.git || exit 1
	cd exabgp || exit 1
	pip3 install -r requirements.txt || exit 1
	pip3 install . || exit 1
else
	pip3 install exabgp==${EXABGP_VERSION} || exit 1
fi
exec /init
