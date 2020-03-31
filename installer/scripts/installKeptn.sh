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

# check if Keptn has already been installed
kubectl get ns keptn
KEPTN_NS_EXISTING=$?

if [[ "$KEPTN_NS_EXISTING" == 0 ]]; then
  print_error "Existing Keptn installation found in namespace keptn."
  KEPTN_DOMAIN=$(kubectl get cm -n keptn keptn-domain -ojsonpath={.data.app_domain})
  echo "Existing Keptn installation found in namespace keptn (${KEPTN_DOMAIN}). Aborting installation..."
  exit 1
fi

# Create Keptn namespace 
# The Keptn namespace is needed before installing the ingress as it is installed into the Keptn namespace
kubectl apply -f ../manifests/keptn/namespace.yaml
verify_kubectl $? "Creating Keptn namespace failed."
print_info "Keptn Namespace created"

case $PLATFORM in
  aks|eks|gke|pks|kubernetes)
    ./common/install.sh
    ;;
  openshift)
    echo "Install Keptn on OpenShift"
    ./openshift/installOnOpenshift.sh
    ;;
  *)
    echo "Platform not provided"
    echo "Installation aborted, please provide platform when executing keptn install --platform="
    exit 1
    ;;
esac
