#!/bin/bash

# see https://rancher.com/docs/k3s/latest/en/installation/
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
