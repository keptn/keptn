#!/bin/bash
source ./utils.sh

KEPTN_VERSION="0.6.1"
print_debug "Upgrading from Keptn 0.6.0 to $KEPTN_VERSION"

manifests=(
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/pvc.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/svc.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/core.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml"
  )

for manifest in "${manifests[@]}"
do
   :
   if curl --head --silent --fail $manifest 2> /dev/null;
     then
      continue
     else
      print_error "Required manifest $manifest not available. Aborting upgrade."
      exit 1
    fi
done


print_debug "Check if Keptn 0.6.0 is currently installed"
API_IMAGE=$(kubectl get deployment -n keptn api -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

  if [[ $API_IMAGE != 'keptn/api:0.6.0' ]]; then
    print_error "Installed Keptn version does not match 0.6.0. aborting."
    exit 1
  fi

# export data from 0.6.0 mongodb
print_debug "Exporting events from previous Keptn installation."
mongoexport --uri="mongodb://user:password@mongodb.keptn-datastore.svc.cluster.local:27017/keptn"  --collection=events  --out=events.json
verify_install_step $? "Mongodb export failed."

# delete old deployment of mongodb
print_debug "Updating MongoDB."
kubectl delete deployment -n keptn-datastore mongodb
# deploy 0.6.1 mongodb
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/pvc.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/deployment.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/logging/mongodb/svc.yaml
wait_for_deployment_in_namespace "mongodb" "keptn-datastore"

# import previous data to new mongodb
print_debug "Importing events from previous installation to updated MongoDB."
mongoimport --uri="mongodb://user:password@mongodb.keptn-datastore.svc.cluster.local:27017/keptn" --collection events --file events.json


print_debug "Updating mongodb-datastore."
kubectl -n keptn-datastore set image deployment/mongodb-datastore mongodb-datastore=keptn/mongodb-datastore:$KEPTN_VERSION --record
kubectl -n keptn-datastore set image deployment/mongodb-datastore-distributor distributor=keptn/distributor:$KEPTN_VERSION --record


print_debug "Updating Keptn core."
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/core.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml

# check if full installation is available
kubectl -n keptn get svc gatekeeper-service

  if [[ $? == '0' ]]; then
      print_debug "Full installation detected. Upgrading CD and CO services"
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/$KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml
  fi

# check for keptn-contrib services
kubectl -n keptn get svc dynatrace-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-service detected. Upgrading to 0.6.2"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-service/release-0.6.2/deploy/manifests/dynatrace-service/dynatrace-service.yaml
  fi

kubectl -n keptn get svc dynatrace-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-sli-service detected. Upgrading to 0.3.1"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/release-0.3.1/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-service detected. Upgrading to 0.3.1"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.1/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-sli-service detected. Upgrading to 0.2.1"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/release-0.2.1/deploy/service.yaml
  fi

kubectl -n keptn get svc servicenow-service

  if [[ $? == '0' ]]; then
      print_debug "Servicenow-service detected. Upgrading to 0.2.0"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/servicenow-service/release-0.2.0/deploy/service.yaml
  fi


