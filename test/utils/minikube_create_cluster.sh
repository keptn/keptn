#!/bin/bash

# download and install minikube
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v1.2.0/minikube-linux-amd64 && chmod +x minikube && sudo mv minikube /usr/local/bin/


# create cluster (with size depending on params)
if [[ "$USE_CASE" == "quality-gates" ]]; then
  # create a cluster with vm-driver=none
  echo "Creating Minikube Cluster with 2 vCPUS and 2 GB memory"
  sudo minikube start --vm-driver=none --cpus 2 --memory 2048
else
  echo "Creating Minikube Cluster with 6 vCPUs and ~6 GB memory"
  sudo minikube start --vm-driver=none --cpus 6 --memory 6144
fi

# make sure kubeconfig has the right permissions
sudo chown -R $USER $HOME/.kube $HOME/.minikube