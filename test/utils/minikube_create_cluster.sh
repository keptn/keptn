#!/bin/bash

# download and install minikube
MINIKUBE_VERSION=${MINIKUBE_VERSION:-"v1.2.0"}
echo "Downloading and installing Minikube in Version ${MINIKUBE_VERSION}"
curl -Lo minikube "https://storage.googleapis.com/minikube/releases/${MINIKUBE_VERSION}/minikube-linux-amd64" && chmod +x minikube && sudo mv minikube /usr/local/bin/


# create cluster (with size depending on params)
if [[ "$USE_CASE" == "continuous-delivery" ]]; then
  echo "Creating Minikube Cluster with 6 vCPUs and ~6 GB memory"
  sudo minikube start --vm-driver=none --cpus 6 --memory 6144
else
  # create a cluster with vm-driver=none
  echo "Creating Minikube Cluster with 2 vCPUS and 2 GB memory"
  sudo minikube start --vm-driver=none --cpus 2 --memory 2048
fi

# make sure kubeconfig has the right permissions
sudo chown -R $USER $HOME/.kube $HOME/.minikube
