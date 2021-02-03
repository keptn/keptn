#!/bin/bash

OC_VERSION=3.11.0
OPENSHIFT_FILENAME="openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz"
# check if file exists

if [[ ! -f ~/downloads/${OPENSHIFT_FILENAME} ]]; then
  echo "Downloading ${OPENSHIFT_FILENAME}"
  wget "https://github.com/openshift/origin/releases/download/v3.11.0/${OPENSHIFT_FILENAME}" -O ~/downloads/${OPENSHIFT_FILENAME}
fi

tar xzvf ~/downloads/${OPENSHIFT_FILENAME} && \
  sudo cp openshift-origin-client-tools-*/oc /bin/oc && \
  sudo mv openshift-origin-client-tools-*/oc /usr/local/bin && \
  rm -rf openshift-origin-client-tools-* && \
  rm -rf openshift*tar.gz
