#!/bin/bash

# download and install minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v1.2.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/

# create a cluster with vm-driver=none
sudo minikube start --vm-driver=none
