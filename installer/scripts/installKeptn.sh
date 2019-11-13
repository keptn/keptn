#!/bin/bash

source ./common/utils.sh

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

# Test kubectl get namespaces
print_info "Testing connection to Kubernetes API"
kubectl get namespaces
verify_kubectl $? "Could not connect to Kubernetes API."
print_info "Connection to Kubernetes API successful"

# Create Keptn namespaces (the Keptn namespace is needed before
# installing the ingress as it is installed into the Keptn namespace)
kubectl apply -f ../manifests/keptn/namespace.yaml
verify_kubectl $? "Creating Keptn namespace failed."


source ./installIngress.sh

case $PLATFORM in
  aks)    
    echo "Installing Keptn on AKS"
    ./common/install.sh
    ;;
  eks)
    echo "Install Keptn on EKS"
    ./common/install.sh
    ;;
  openshift)
    echo "Install Keptn on OpenShift"
    ./openshift/installOnOpenshift.sh
    ;;
  gke)    
    echo "Install Keptn on GKE"
    ./common/install.sh
    ;;
  pks)
    echo "Install Keptn on PKS"
    ./common/install.sh
    ;;
  kubernetes)
    echo "Install Keptn on Kubernetes"
    ./common/install.sh
    ;;
  *)
    echo "Platform not provided"
    echo "Installation aborted, please provide platform when executing keptn install --platform="
    exit 1
    ;;
esac
