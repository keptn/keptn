#!/bin/bash

# install micro k8s via snap
sudo snap install microk8s --classic

sudo iptables -P FORWARD ACCEPT
sudo microk8s.enable dns
sudo microk8s.enable dns ingress
sudo microk8s.enable storage

# store kubeconfig
sudo /snap/bin/microk8s.config > ~/kubeconfig
export KUBECONFIG=~/kubeconfig
kubectl apply -f https://raw.githubusercontent.com/google/metallb/v0.8.3/manifests/metallb.yaml

