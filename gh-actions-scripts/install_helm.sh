#!/bin/bash

# download
curl https://get.helm.sh/helm-v3.2.4-linux-amd64.tar.gz > helm.tar.gz
tar -zxvf helm.tar.gz
if [ ! -d "$GITHUB_WORKSPACE/bin" ]; then
  mkdir "$GITHUB_WORKSPACE/bin"
fi

mv ./linux-amd64/helm .
