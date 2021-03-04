#!/bin/bash

# install micro k8s via snap
MICROK8S_VERSION=${MICROK8S_VERSION:-"1.18"}
sudo snap install microk8s --classic --channel="$MICROK8S_VERSION"

sudo iptables -P FORWARD ACCEPT
sudo microk8s.enable dns
sudo microk8s.enable dns ingress
sudo microk8s.enable storage

# store kubeconfig
# shellcheck disable=SC2024
sudo /snap/bin/microk8s.config > ~/kubeconfig
export KUBECONFIG=~/kubeconfig
kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.8.3/manifests/metallb.yaml

