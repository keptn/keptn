#!/bin/bash

source ./common/utils.sh

print_info "Starting installation of Keptn"

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