#!/bin/bash

  # download
curl https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz > helm.tar.gz
tar -zxvf helm.tar.gz
mv linux-amd64/helm /usr/local/bin/helm