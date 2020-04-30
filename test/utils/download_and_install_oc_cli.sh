#!/bin/bash

OC_VERSION=3.11.0
wget https://github.com/openshift/origin/releases/download/v3.11.0/openshift-origin-client-tools-v$OC_VERSION-0cbc58b-linux-64bit.tar.gz && \
  tar xzvf openshift*tar.gz && \
  sudo cp openshift-origin-client-tools-*/oc /bin/oc && \
  sudo mv openshift-origin-client-tools-*/oc /usr/local/bin && \
  rm -rf openshift-origin-client-tools-* && \
  rm -rf openshift*tar.gz
