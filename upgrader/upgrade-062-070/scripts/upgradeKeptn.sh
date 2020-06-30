#!/bin/bash
source ./utils.sh

./upgradecollections $MONGODB_URL $CONFIGURATION_SERVICE_URL

if [ $HELM_RELEASE_UPGRADE == "true" ]; 
then
  # Upgrade from Helm v2 to Helm v3
  helm init --client-only
  verify_install_step $? "Helm init failed."
  RELEASES=$(helm list -aq)
  verify_install_step $? "Helm list failed."
  echo $RELEASES

  helm3 plugin install https://github.com/helm/helm-2to3
  verify_install_step $? "Helm-2to3 plugin installation failed."
  yes y | helm3 2to3 move config
  verify_install_step $? "Helm-2to3 move of config failed."

  for release in $RELEASES; do
    helm3 2to3 convert $release --dry-run
    verify_install_step $? "Helm2-to3 release convertion dry-run failed"
    helm3 2to3 convert $release
    verify_install_step $? "Helm2-to3 release convertion failed"
  done

  yes y | helm3 2to3 cleanup --tiller-cleanup
  verify_install_step $? "Helm2-to3 cleanup failed"
  
fi


KEPTN_VERSION=${KEPTN_VERSION:-"release-0.7.0"}
print_debug "Upgrading from Keptn 0.6.2 to $KEPTN_VERSION"

manifests=(
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/rbac.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/rbac.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/secret.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb-datastore/mongodb-datastore.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb-datastore/mongodb-datastore-distributor.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/core.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/keptn-api-virtualservice.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/keptn-ingress.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/uniform-services-openshift.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/uniform-distributors-openshift.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/api-gateway-nginx.yaml"
  )

for manifest in "${manifests[@]}"
do
   :
   if curl --head --silent -k --fail $manifest 2> /dev/null;
     then
      continue
     else
      print_error "Required manifest $manifest not available. Aborting upgrade."
      exit 1
    fi
done


print_debug "Check if Keptn 0.6.2 is currently installed"
API_IMAGE=$(kubectl get deployment -n keptn api-service -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

  if [[ $API_IMAGE != 'keptn/api:0.6.2' ]]; then
    print_error "Installed Keptn version does not match 0.6.1. aborting."
    exit 1
  fi

print_debug "Updating MongoDB and mongodb-datastore."
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/secret.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb-datastore/k8s/mongodb-datastore.yaml

print_debug "Updating Keptn core."
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/api-gateway-nginx.yaml
kubectl -n keptn delete pod -lrun=api-gateway-nginx
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/core.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml
# remove the remediation-service-problem-distributor deployment since the remediation service now has a new distributor for multiple types of evetns
kubectl delete deployment -n keptn remediation-service-problem-distributor

kubectl get namespace openshift
  if [[ $? == '0' ]]; then
    print_debug "OpenShift platform detected. Updating OpenShift core services"
    kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/uniform-services-openshift.yaml
    kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/uniform-distributors-openshift.yaml
  fi

DOMAIN=$(kubectl get configmap -n keptn keptn-domain -ojsonpath="{.data.app_domain}")

# check if full installation is available
kubectl -n keptn get svc gatekeeper-service

  if [[ $? == '0' ]]; then
      print_debug "Full installation detected. Upgrading CD and CO services"
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml
  fi

# check for keptn-contrib services
kubectl -n keptn get svc dynatrace-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-service detected. Upgrading to 0.6.3"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-service/release-0.6.3/deploy/manifests/dynatrace-service/dynatrace-service.yaml
  fi

kubectl -n keptn get svc dynatrace-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-sli-service detected. Upgrading to 0.3.2"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/release-0.3.2/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-service detected. Upgrading to 0.3.3"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.3/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-sli-service detected. Upgrading to 0.2.2"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/release-0.2.2/deploy/service.yaml
  fi

kubectl -n keptn get svc servicenow-service

kubectl delete ClusterRoleBinding keptn-rbac
kubectl delete ClusterRoleBinding rbac-service-account
