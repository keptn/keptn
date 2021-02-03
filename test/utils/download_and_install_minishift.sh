#!/bin/bash

# download and install minishift
MINISHIFT_VERSION=1.34.3
MINISHIFT_FILENAME=minishift-${MINISHIFT_VERSION}-linux-amd64.tgz

if [[ ! -f ~/downloads/${MINISHIFT_FILENAME} ]]; then
  echo "Downloading ${MINISHIFT_FILENAME}"
  wget "https://github.com/minishift/minishift/releases/download/v${MINISHIFT_VERSION}/${MINISHIFT_FILENAME}" -O ~/downloads/${MINISHIFT_FILENAME}
fi

tar zxvf ~/downloads/${MINISHIFT_FILENAME} && \
  sudo mv minishift-*-linux-amd64/minishift /usr/local/bin/
