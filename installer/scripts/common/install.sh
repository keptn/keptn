#!/bin/bash

source ./common/utils.sh

print_info "Starting installation of keptn"

# Test kubectl get namespaces
print_info "Testing connection to Kubernetes API"
kubectl get namespaces
verify_kubectl $? "Could not connect to Kubernetes API."
print_info "Connection to Kubernetes API successful"

# Create keptn namespaces
kubectl apply -f ../manifests/keptn/namespace.yaml
verify_kubectl $? "Creating keptn namespace failed."

# Install Istio service mesh
print_info "Installing Istio"
./common/setupIstio.sh
verify_install_step $? "Installing Istio failed."
print_info "Installing Istio done"

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