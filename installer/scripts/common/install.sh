#!/bin/bash

source ./common/utils.sh

print_info "Starting installation of Keptn"

# Test kubectl get namespaces
print_info "Testing connection to Kubernetes API"
kubectl get namespaces
verify_kubectl $? "Could not connect to Kubernetes API."
print_info "Connection to Kubernetes API successful"

# Create Keptn namespaces
kubectl apply -f ../manifests/keptn/namespace.yaml
verify_kubectl $? "Creating Keptn namespace failed."

# Install Istio service mesh
print_info "Installing Istio (this might take a while)"
./common/setupIstio.sh
verify_install_step $? "Installing Istio failed."
print_info "Installing Istio done"

# Install Tiller for Helm
print_info "Installing Tiller"
kubectl apply -f ../manifests/tiller/tiller.yaml
helm init --service-account tiller
print_info "Installing Tiller done"

# Install Keptn core services - Install Keptn channels
print_info "Installing Keptn"
./common/setupKeptn.sh
verify_install_step $? "Installing Keptn failed."
print_info "Installing Keptn done"

# Install Keptn services
print_info "Wear uniform"
./common/wearUniform.sh
verify_install_step $? "Installing Keptn's uniform failed."
print_info "Keptn wears uniform"

# Install done
print_info "Installation of Keptn complete."