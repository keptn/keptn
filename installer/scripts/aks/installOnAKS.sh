#!/bin/bash


# kubectl apply --filename https://github.com/knative/serving/releases/download/v0.4.0/istio-crds.yaml
# verify_kubectl $? "Error applying Istio Credentials"
# kubectl apply --filename https://github.com/knative/serving/releases/download/v0.4.0/istio.yaml
# verify_kubectl $? "Error applying Istio"
# kubectl label namespace default istio-injection=enabled --overwrite=true
# verify_kubectl $? "Error setting istio-injection flag "
# wait_for_all_pods_in_namespace "istio-system"

source ./common/utils.sh

print_info "Starting installation of keptn"

if [[ -z "${KEPTN_INSTALL_ENV}" ]]; then
  # Variables for az
  if [[ -z "${CLUSTER_NAME}" ]]; then
    print_debug "CLUSTER_NAME is not set, take it from creds.json"
    CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
    verify_variable "$CLUSTER_NAME" "CLUSTER_NAME is not defined in environment variable nor in creds.json file." 
  fi

  if [[ -z "${AZURE_RESOURCEGROUP}" ]]; then
    print_debug "AZURE_RESOURCEGROUP is not set, reading it from creds.json"
    export AZURE_RESOURCEGROUP=$(cat creds.json | jq -r '.azureResourceGroup')
  fi 

  if [[ -z "${AZURE_SUBSCRIPTION}" ]]; then
    print_debug "AZURE_SUBSCRIPTION is not set, reading it from creds.json"
    export AZURE_SUBSCRIPTION=$(cat creds.json | jq -r '.azureSubscription')
  fi 

  # Test connection to cluster
  print_info "Test connection to cluster"
  ./aks/testConnection.sh $CLUSTER_NAME $AZURE_RESOURCEGROUP $AZURE_SUBSCRIPTION
  # verify_install_step $? "Could not connect to cluster. Please check the values for your Cluster Name, GKE Project, and Cluster Zone during the credentials setup."
  print_info "Connection to cluster successful"
fi

# Variables for installing Istio and Knative
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  print_debug "CLUSTER_IPV4_CIDR is not set, retrieve it using az."
  CLUSTER_IPV4_CIDR=$(az aks show --name ${CLUSTER_NAME} --resource-group ${AZURE_RESOURCEGROUP} --query=networkProfile.podCidr --output tsv)
  if [[ $? != 0 ]]; then
    print_error "az failed to describe the ${CLUSTER_NAME} cluster for retrieving the ${CLUSTER_IPV4_CIDR} property." && exit 1
  fi
  verify_variable "$CLUSTER_IPV4_CIDR" "CLUSTER_IPV4_CIDR is not defined in environment variable nor could it be retrieved using az." 
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  print_debug "SERVICES_IPV4_CIDR is not set, retrieve it using az."
  SERVICES_IPV4_CIDR=$(az aks show --name ${CLUSTER_NAME} --resource-group ${AZURE_RESOURCEGROUP} --query=networkProfile.serviceCidr --output tsv)
  if [[ $? != 0 ]]; then
    print_error "az failed to describe the ${CLUSTER_NAME} cluster for retrieving the ${SERVICES_IPV4_CIDR} property." && exit 1
  fi
  verify_variable "$SERVICES_IPV4_CIDR" "SERVICES_IPV4_CIDR is not defined in environment variable nor could it be retrieved using az." 
fi

# Variables for creating cluster role binding
if [[ -z "${USER}" ]]; then
  print_debug "USER is not set, retrieve it using az."
  USER=$(az account show --query=user.name --output=tsv)
  if [[ $? != 0 ]]; then
    print_error "az failed to get account values." && exit 1
  fi
  verify_variable "$USER" "USER is not defined in environment variable nor could it be retrieved using az." 
fi

# Test kubectl get namespaces
print_info "Testing connection to Kubernetes API"
kubectl get namespaces
verify_kubectl $? "Could not connect to Kubernetes API."
print_info "Connection to Kubernetes API successful"

# Grant cluster admin rights to gcloud user
kubectl create clusterrolebinding keptn-cluster-admin-binding --clusterrole=cluster-admin --user=$USER
verify_kubectl $? "Cluster role binding could not be created."

# Create keptn namespaces
kubectl apply -f ../manifests/keptn/namespace.yaml
verify_kubectl $? "Creating keptn namespace failed."

# Install Istio service mesh
print_info "Installing Istio"
./common/setupIstio.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Istio failed."
print_info "Installing Istio done"

# Install knative core components
print_info "Installing Knative"
./common/setupKnative.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR
verify_install_step $? "Installing Knative failed."
print_info "Installing Knative done"

# Enable fluentd 
kubectl label nodes --all beta.kubernetes.io/fluentd-ds-ready="true"

# Install tiller for helm
print_info "Installing Tiller"
kubectl apply -f ../manifests/tiller/tiller.yaml
helm init --service-account tiller
print_info "Installing Tiller done"

# Install keptn core services - Install keptn channels
print_info "Installing keptn"
./common/setupKeptn.sh
verify_install_step $? "Installing keptn failed."
print_info "Installing keptn done"

# Install keptn services
print_info "Wear uniform"
./common/wearUniform.sh
verify_install_step $? "Installing keptn's uniform failed."
print_info "Keptn wears uniform"

# Install done
print_info "Installation of keptn complete."

# Retrieve keptn endpoint and api-token
KEPTN_ENDPOINT=https://api.keptn.$(kubectl get cm -n keptn keptn-domain -oyaml | yq - r data.app_domain)
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode)

print_info "keptn endpoint: $KEPTN_ENDPOINT"
print_info "keptn api-token: $KEPTN_API_TOKEN"

#print_info "To retrieve the keptn API token, please execute the following command:"
#print_info "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
