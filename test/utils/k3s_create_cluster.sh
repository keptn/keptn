#!/bin/bash

# install k3s version from github (see https://github.com/rancher/k3s/releases)
K3S_VERSION=${K3S_VERSION:-"v1.18.3+k3s1"}

# see https://rancher.com/docs/k3s/latest/en/installation/
echo "Downloading and installing K3s in Version ${K3S_VERSION}"
curl -Lo k3s "https://github.com/rancher/k3s/releases/download/${K3S_VERSION}/k3s" && chmod +x k3s && sudo mv k3s /usr/local/bin/

sudo k3s server --no-deploy=traefik --write-kubeconfig-mode=644 --log /k3s.log &
# wait a bit for the server to start
sleep 30
echo "Waiting until K3s is ready ..."
sleep 30
echo "Still waiting ..."
sleep 30
cat /k3s.log

# Kubeconfig is written to /etc/rancher/k3s/k3s.yaml
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

kubectl get nodes
kubectl get services --all-namespaces
