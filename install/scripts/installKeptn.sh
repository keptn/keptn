#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

source ./utils.sh

print_info "Starting installation of keptn"

if [[ -z "${KEPTN_INSTALL_ENV}" ]]; then
  # Variables for gcloud
  if [[ -z "${CLUSTER_NAME}" ]]; then
    print_debug "CLUSTER_NAME is not set, take it from creds.json"
    CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
    verify_variable "$CLUSTER_NAME" "CLUSTER_NAME is not defined in environment variable nor in creds.json file." 
  fi

  if [[ -z "${CLUSTER_ZONE}" ]]; then
    print_debug "CLUSTER_ZONE is not set, take it from creds.json"
    CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
    verify_variable "$CLUSTER_ZONE" "CLUSTER_NAME is not defined in environment variable nor in creds.json file." 
  fi

  # Test connection to cluster
  print_info "Test connection to cluster"
  ./testConnection.sh $CLUSTER_NAME $CLUSTER_ZONE
  # verify_install_step $? "Could not connect to cluster. Please check the values for your Cluster Name, GKE Project, and Cluster Zone during the credentials setup."
  print_info "Connection to cluster successful"
fi

# Variables for installing Istio and Knative
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  print_debug "CLUSTER_IPV4_CIDR is not set, retrieve it using gcloud."
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
  if [[ $? != 0 ]]; then
    print_error "gcloud failed to describe the ${CLUSTER_NAME} cluster for retrieving the ${CLUSTER_IPV4_CIDR} property." && exit 1
  fi
  verify_variable "$CLUSTER_IPV4_CIDR" "CLUSTER_IPV4_CIDR is not defined in environment variable nor could it be retrieved using gcloud." 
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  print_debug "SERVICES_IPV4_CIDR is not set, retrieve it using gcloud."
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
  if [[ $? != 0 ]]; then
    print_error "gcloud failed to describe the ${CLUSTER_NAME} cluster for retrieving the ${SERVICES_IPV4_CIDR} property." && exit 1
  fi
  verify_variable "$SERVICES_IPV4_CIDR" "SERVICES_IPV4_CIDR is not defined in environment variable nor could it be retrieved using gcloud." 
fi

# Variables for creating cluster role binding
if [[ -z "${GCLOUD_USER}" ]]; then
  print_debug "GCLOUD_USER is not set, retrieve it using gcloud."
  GCLOUD_USER=$(gcloud config get-value account)
  if [[ $? != 0 ]]; then
    print_error "gloud failed to get account values." && exit 1
  fi
  verify_variable "$GCLOUD_USER" "GCLOUD_USER is not defined in environment variable nor could it be retrieved using gcloud." 
fi

# Test kubectl get namespaces
print_info "Testing connection to Kubernetes API"
kubectl get namespaces
verify_kubectl $? "Could not connect to Kubernetes API."
print_info "Connection to Kubernetes API successful"

# Grant cluster admin rights to gcloud user
# TODO create vs apply
kubectl create clusterrolebinding keptn-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER
verify_kubectl $? "Cluster role binding could not be created."

# Create K8s namespaces
kubectl apply -f ../manifests/keptn/keptn-namespace.yml
verify_kubectl $? "Creating keptn namespace failed."

# Create container registry
print_info "Creating container registry"
./setupContainerRegistry.sh
verify_install_step $? "Creating container registry failed."
print_info "Creating container registry done"

# Install Istio service mesh
print_info "Installing Istio"
./setupIstio.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Istio failed."
print_info "Installing Istio done"

# Install knative core components
print_info "Installing Knative"
./setupKnative.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Knative failed."
print_info "Installing Knative done"

# Install keptn core services - Install keptn channels
print_info "Installing keptn"
./setupKeptn.sh
verify_install_step $? "Installing keptn failed."
print_info "Installing keptn done"

# Install keptn services
print_info "Wear uniform"
./wearUniform.sh
verify_install_step $? "Installing keptn's uniform failed."
print_info "Keptn wears uniform"

# Install done
print_info "Installation of keptn complete."

# Retrieve keptn endpoint and api-token
KEPTN_ENDPOINT=https://$(kubectl get ksvc -n keptn control -o=yaml | yq r - status.domain)
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode)

print_info "keptn endpoint: $KEPTN_ENDPOINT"
print_info "keptn api-token: $KEPTN_API_TOKEN"

#print_info "To retrieve the keptn API token, please execute the following command:"
#print_info "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
