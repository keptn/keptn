#!/bin/bash

source ./common/utils.sh

if [[ "$USE_CASE" == "all" ]]; then
  # Install istio
  source ./common/setupIstio.sh
  setupKeptnDomain "istio" "istio-ingressgateway" "istio-system"

  cat ../manifests/keptn/keptn-api-virtualservice.yaml | \
    sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' | kubectl apply -f -
  verify_kubectl $? "Deploying keptn api virtualservice failed."
else
  # Install NGINX
  source ./common/setupNginx.sh
  setupKeptnDomain "nginx" "ingress-nginx" "ingress-nginx"

  # Add config map in keptn namespace that contains the domain - this will be used by other services as well
  # Update ingress with updated hosts
  cat ../manifests/keptn/keptn-ingress.yaml | \
    sed 's~domain.placeholder~'"$DOMAIN"'~' | sed 's~ingress.placeholder~nginx~' | kubectl apply -f -
  verify_kubectl $? "Deploying ingress failed."
fi

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