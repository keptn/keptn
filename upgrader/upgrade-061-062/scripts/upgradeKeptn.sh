#!/bin/bash
source ./utils.sh

KEPTN_VERSION="0.6.2"
print_debug "Upgrading from Keptn 0.6.1 to $KEPTN_VERSION"

manifests=(
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/logging/mongodb-datastore/k8s/mongodb-datastore.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/logging/mongodb-datastore/mongodb-datastore-distributor.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/core.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/keptn-api-virtualservice.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/keptn-ingress.yaml"
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


print_debug "Check if Keptn 0.6.1 is currently installed"
API_IMAGE=$(kubectl get deployment -n keptn api -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

  if [[ $API_IMAGE != 'keptn/api:0.6.1' ]]; then
    print_error "Installed Keptn version does not match 0.6.1. aborting."
    exit 1
  fi


print_debug "Updating Keptn core."
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/core.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml

kubectl -n keptn delete deployment api
kubectl -n keptn delete service api

# check if full installation is available
kubectl -n keptn get svc gatekeeper-service

  if [[ $? == '0' ]]; then
      print_debug "Full installation detected. Upgrading CD and CO services"
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/keptn-api-virtualservice.yaml
  else
      print_debug "Quality gates installation detected. Upgrading Nginx ingress"
      kubectl apply -f https://raw.githubusercontent.com/keptn/keptn/release-$KEPTN_VERSION/installer/manifests/keptn/keptn-ingress.yaml
      kubectl -n keptn delete ingress api-ingress

      kubectl get namespace openshift

      if [[ $? == '0' ]]; then
        print_debug "OpenShift platform detected. Updating routes"
        oc delete route api -n keptn

        print_info "Creating route to api-gateway-nginx"
        oc create route edge api --service=api-gateway-nginx --port=http --insecure-policy='None' -n keptn

        BASE_URL=$(oc get route -n keptn api -oyaml | yq r - spec.host | sed 's~api-keptn.~~')
        DOMAIN=$BASE_URL

        print_info "Used domain for api OC route ${DOMAIN}"
        oc delete route api -n keptn

        oc create route edge api --service=api-gateway-nginx --port=http --insecure-policy='None' -n keptn --hostname="api.keptn.$BASE_URL"
        oc create route edge api2 --service=api-gateway-nginx --port=http --insecure-policy='None' -n keptn --hostname="api.keptn"
      fi
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
      print_debug "Prometheus-service detected. Upgrading to 0.3.2"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.3.2/deploy/service.yaml
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


