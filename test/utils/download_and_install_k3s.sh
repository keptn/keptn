#!/bin/bash

# install k3s version from github (see https://github.com/rancher/k3s/releases)
K3S_VERSION=${K3S_VERSION:-"v1.18.3+k3s1"}
K3S_FILENAME="k3s-${K3S_VERSION}"

if [[ ! -f ~/downloads/${K3S_FILENAME} ]]; then
  echo "Downloading and installing K3s in Version ${K3S_VERSION}"
  echo "Downloading ${K3S_FILENAME}"
  wget "https://github.com/rancher/k3s/releases/download/${K3S_VERSION}/k3s" -O ~/downloads/${K3S_FILENAME}
fi

# see https://rancher.com/docs/k3s/latest/en/installation/
cp ~/downloads/${K3S_FILENAME} k3s
chmod +x k3s && sudo mv k3s /usr/local/bin/
