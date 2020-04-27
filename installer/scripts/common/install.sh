#!/bin/bash

source ./common/utils.sh

if [[ "$USE_CASE" == "all" ]]; then
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

# Install Tiller for Helm
if [[ "$USE_CASE" == "all" ]]; then
  print_info "Installing Tiller"
  kubectl apply -f ../manifests/tiller/tiller.yaml
  verify_kubectl $? "Applying Tiller manifest failed."
  kubectl get pods
  print_info "Initializing helm/tiller"
  helm init --service-account tiller
  verify_install_step $? "Installing Tiller failed"
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