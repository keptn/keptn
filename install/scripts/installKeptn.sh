#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

source ./utils.sh

echo "[keptn|INFO] [Fri Sep 09 10:42:29.902022 2011] Starting installation of keptn"

# Variables for gcloud
if [[ -z "${CLUSTER_NAME}" ]]; then
  print_debug "CLUSTER_NAME not set, take it from creds.json"
  CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
  verify_variable $CLUSTER_NAME "CLUSTER_NAME is empty, stop installation." 
fi

if [[ -z "${CLUSTER_ZONE}" ]]; then
  print_debug "CLUSTER_ZONE not set, take it from creds.json"
  CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
  verify_variable $CLUSTER_ZONE "CLUSTER_ZONE is empty, stop installation." 
fi

# Variables for installing Istio and Knative
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  print_debug "CLUSTER_IPV4_CIDR not set, retrieve it using gcloud."
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
  if [[ $? != 0 ]]; then
    print_error "gloud failed, stop installation." && exit 1
  fi
  verify_variable $CLUSTER_IPV4_CIDR "CLUSTER_IPV4_CIDR is empty, stop installation." 
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  print_debug "SERVICES_IPV4_CIDR not set, retrieve it using gcloud"
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
  if [[ $? != 0 ]]; then
    print_error "gloud failed, stop installation." && exit 1
  fi
  verify_variable $SERVICES_IPV4_CIDR "SERVICES_IPV4_CIDR is empty, stop installation." 
fi

# Variables for creating cluster role binding
if [[ -z "${GCLOUD_USER}" ]]; then
  print_debug "GCLOUD_USER not set, retrieve it using gcloud"
  GCLOUD_USER=$(gcloud config get-value account)
  if [[ $? != 0 ]]; then
    print_error "gloud failed, stop installation." && exit 1
  fi
  verify_variable $GCLOUD_USER "GCLOUD_USER is empty, stop installation." 
fi

# Test connection to cluster
print_info "Test connection to cluster"
./testConnection.sh $CLUSTER_NAME $CLUSTER_ZONE
verify_install_step $? "Could not connect to cluster, stop installation. Please check the values for your Cluster Name, GKE Project, and Cluster Zone during the credentials setup."
print_info "Connection to cluster successful"

# Grant cluster admin rights to gcloud user
kubectl create clusterrolebinding keptn-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER
verify_kubectl $? "Cluster role binding could not be created, stop installation."

# Create K8s namespaces
kubectl apply -f ../manifests/keptn/keptn-namespace.yml
verify_kubectl $? "Namespace could not be created, stop installation."

# Create container registry
print_info "Creating container registry"
./setupContainerRegistry.sh
verify_install_step $? "Creating container registry failed, stop installation."
print_info "Creating container registry done"

# Install Istio service mesh
print_info "Installing Istio"
./setupIstio.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Istio failed, stop installation."
print_info "Installing Istio done"

# Install knative core components
print_info "Installing Knative"
./setupKnative.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Knative failed, stop installation."
print_info "Installing Knative done"

# # Install keptn core services - Install keptn channels
# print_info "Install keptn"
# ./setupKeptn.sh
# verify_install_step $? "Installing keptn failed, stop installation."
# print_info "Install keptn done"

# # Install keptn services
# print_info "Wear uniform"
# ./wearUniform.sh
# verify_install_step $? "Installing keptn's uniform failed, stop installation."
# print_info "Keptn wears uniform"

# # Install done
# print_info "Installation of keptn complete."

# print_info "To retrieve the Keptn API Token, please execute the following command:"
# print_info "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
