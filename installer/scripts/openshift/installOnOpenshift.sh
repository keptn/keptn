#!/bin/bash

source ./openshift/installationFunctions.sh
source ./common/utils.sh

# configure the host path volume plugin
oc apply -f ../manifests/openshift/oc-scc-hostpath.yaml
verify_kubectl $? "Deploying hostpath SCC failed."
oc patch scc hostpath -p '{"allowHostDirVolumePlugin": true}'
# verify_install_step "Patching hostpath plugin failed."
oc adm policy add-scc-to-group hostpath system:authenticated
#verify_install_step "Creating hostpath SCC failed."

# Install monitoring
#oc adm policy add-scc-to-user privileged -z elasticsearch-logging -n knative-monitoring
#oc adm policy add-scc-to-user anyuid system:serviceaccount:knative-monitoring:fluentd-ds
#oc adm policy add-scc-to-user privileged system:serviceaccount:knative-monitoring:fluentd-ds
#kubectl label nodes --all beta.kubernetes.io/fluentd-ds-ready="true"
#verify_kubectl $? "Labelling nodes failed."
#kubectl apply -f ../manifests/knative/monitoring.yaml
#verify_kubectl $? "Applying knative monitoring components failed."
#wait_for_all_pods_in_namespace "knative-monitoring"


# Install Tiller for Helm
if [[ "$USE_CASE" == "all" ]]; then
  print_info "Installing Istio (this might take a while)"
  source ./common/setupIstio.sh
  verify_install_step $? "Installing Istio failed."
  print_info "Installing Istio done"

  print_info "Installing Tiller"
  kubectl apply -f ../manifests/tiller/tiller.yaml
  helm init --service-account tiller
  print_info "Installing Tiller done"
  oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:kube-system:tiller
else
  print_debug "Installing Tiller is skipped since use case ${USE_CASE} does not need it." 
fi

# Install keptn core services - Install keptn channels
print_info "Installing Keptn"
./openshift/setupKeptn.sh
verify_install_step $? "Installing Keptn failed."
print_info "Installing Keptn done"

# Install keptn services
if [[ "$USE_CASE" == "all" ]]; then
  print_info "Wear uniform"
  ./common/wearUniform.sh
  verify_install_step $? "Installing Keptn's uniform failed."
  print_info "Keptn wears uniform"
fi

# Install additional keptn services for openshift
print_info "Wear Openshift uniform"
./openshift/wearUniform.sh
verify_install_step $? "Installing Keptn's Openshift uniform failed."
print_info "Keptn wears Openshift uniform"

# Install done
print_info "Installation of Keptn complete."

