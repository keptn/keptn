#!/bin/bash

source ./common/utils.sh

if [[ "$USE_CASE" == "continuous-delivery" ]]; then
  # Install istio
  source ./common/setupIstio.sh
  setupKeptnDomain "istio" "istio-ingressgateway" "istio-system"

  # Note: We need to use DOMAIN_PLACEHOLDER (does not contain a port) here
  cat ../manifests/keptn/keptn-api-virtualservice.yaml | \
    sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' | kubectl apply -f -
  verify_kubectl $? "Deploying keptn api virtualservice failed."
else
  # Install NGINX
  source ./common/setupNginx.sh
  setupKeptnDomain "nginx" "ingress-nginx" "ingress-nginx"

  # Add config map in keptn namespace that contains the domain - this will be used by other services as well
  # Update ingress with updated hosts
  # Note: we need to use INGRESS_HOST (does not contain a port) here
  cat ../manifests/keptn/keptn-ingress.yaml | \
    sed 's~domain.placeholder~'"$INGRESS_HOST"'~' | sed 's~ingress.placeholder~nginx~' | kubectl apply -f -
  verify_kubectl $? "Deploying ingress failed."
fi

print_info "Starting installation of Keptn"

./common/setupKeptn.sh
verify_install_step $? "Installing Keptn failed."
print_info "Installing Keptn done"

# Install Keptn services
if [[ "$USE_CASE" == "continuous-delivery" ]]; then
  print_info "Wear uniform"
  ./common/wearUniform.sh
  verify_install_step $? "Installing Keptn's uniform failed."
  print_info "Keptn wears uniform"
else
  print_debug "Wear uniform is skipped since continuous-delivery use case has not been activated."
fi

# Install done
print_info "Installation of Keptn complete."

# wait a few seconds to make sure the last log output is captured by the CLI before the pod is deleted
sleep 10
