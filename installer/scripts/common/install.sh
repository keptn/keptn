#!/bin/bash

source ./common/utils.sh

print_info "Starting installation of Keptn"

./common/setupKeptn.sh
verify_install_step $? "Installing Keptn failed."
print_info "Installing Keptn done"

# Install Keptn services
if [[ "$USE_CASE" == "all" ]]; then
  print_info "Wear uniform"
  ./common/wearUniform.sh
  verify_install_step $? "Installing Keptn's uniform failed."
  print_info "Keptn wears uniform"
else
  print_debug "Wear uniform is skipped since use case ${USE_CASE} does not need it." 
fi

# Install done
print_info "Installation of Keptn complete."