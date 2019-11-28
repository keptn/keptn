#!/bin/bash

source ./common/utils.sh

print_info "Starting installation of Keptn"

# Install Tiller for Helm
if [[ "$USE_CASE" == "all" ]]; then
  print_info "Installing Tiller"
  kubectl apply -f ../manifests/tiller/tiller.yaml
  helm init --service-account tiller
  print_info "Installing Tiller done"
else
  print_debug "Installing Tiller is skipped since use case ${USE_CASE} does not need it." 
fi

# Install Keptn core services
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