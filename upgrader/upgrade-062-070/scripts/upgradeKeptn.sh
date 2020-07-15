#!/bin/bash
source ./utils.sh

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

PREVIOUS_KEPTN_VERSION="0.6.2"
KEPTN_VERSION=${KEPTN_VERSION:-"release-0.7.0"}
HELM_CHART_URL=${HELM_CHART_URL:-"https://storage.googleapis.com/keptn-installer/0.7.0"}
MONGODB_SOURCE_URL=${MONGODB_SOURCE_URL:-"mongodb://user:password@mongodb.keptn-datastore:27017/keptn"}
MONGODB_TARGET_URL=${MONGODB_TARGET_URL:-"mongodb://user:password@mongodb.keptn:27017/keptn"}

print_debug "Upgrading from Keptn 0.6.2 to $KEPTN_VERSION"

KEPTN_API_URL=https://api.keptn.$(kubectl get cm keptn-domain -n keptn -ojsonpath={.data.app_domain})
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 -d)
KEPTN_DOMAIN=$(kubectl get cm keptn-domain -n keptn -ojsonpath={.data.app_domain})

print_debug "Check if Keptn 0.6.2 is currently installed"
API_IMAGE=$(kubectl get deployment -n keptn api-service -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

  if [[ $API_IMAGE != 'keptn/api:0.6.2' ]]; then
    print_error "Installed Keptn version does not match 0.6.2. aborting."
    exit 1
  fi

USE_CASE=""
# check if full installation is available
kubectl -n keptn get svc gatekeeper-service

  if [[ $? == '0' ]]; then
      print_debug "Full installation detected. Upgrading CD and CO services"
      USE_CASE="continuous-delivery"
  fi

old_manifests=(
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/core.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/keptn-api-virtualservice.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/keptn-ingress.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/uniform-services-openshift.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/uniform-distributors-openshift.yaml"
  )

for manifest in "${old_manifests[@]}"
do
   :
   if curl --head --silent -k --fail $manifest 2> /dev/null;
     then
      kubectl delete -f $manifest
      continue
     else
      print_error "Required manifest $manifest not available. Aborting upgrade."
      exit 1
    fi
done

# delete resource that have been generated procedurally by the installer to avoid conflicts with helm install
kubectl deletee configmap -n keptn api-nginx-config
kubectl delete secret -n keptn keptn-api-token

BRIDGE_USERNAME=""
kubectl get secret -n keptn bridge-credentials
  if [[ $? == '0' ]]; then
    print_debug "Bridge credentials detected. Fetching credentials"
    BRIDGE_USERNAME=$(kubectl get secret bridge-credentials -n keptn -ojsonpath={.data.BASIC_AUTH_USERNAME} | base64 -d)
    BRIDGE_PASSWORD=$(kubectl get secret bridge-credentials -n keptn -ojsonpath={.data.BASIC_AUTH_PASSWORD} | base64 -d)
    kubectl -n keptn delete secret bridge-credentials
  fi


# install helm chart

helm3 repo add keptn $HELM_CHART_URL

if [[ $USECASE == "continuous-delivery" ]]; then
  kubectl get namespace openshift
  if [[ $? == '0' ]]; then
    print_debug "OpenShift platform detected. Updating OpenShift core services"
    helm3 install keptn keptn/keptn -n keptn --set continuous-delivery.enabled=true --set continuous-delivery.openshift.enabled=true
  else
    helm3 install keptn keptn/keptn -n keptn --set continuous-delivery.enabled=true
  fi
else
  helm3 install keptn keptn/keptn -n keptn --set continuous-delivery.enabled=false
fi

kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN" -oyaml --dry-run | kubectl replace -f -

if [[ $BRIDGE_USERNAME == "" ]]; then
  echo "No previous bridge credentials found. No need to update"
else
  echo "Setting bridge credentials to previous values"
  kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=$BRIDGE_USERNAME" --from-literal="BASIC_AUTH_PASSWORD=$BRIDGE_PASSWORD" -oyaml --dry-run | kubectl replace -f -
fi

kubectl -n keptn set env deployment/configuration-service MONGO_DB_CONNECTION_STRING='mongodb://user:password@mongodb.keptn-datastore:27017/keptn'
kubectl delete pod -n keptn -lrun=configuration-service
sleep 30

MONGO_TARGET_USER=$(kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.user} | base64 -d)
MONGO_TARGET_PASSWORD=$(kubectl get secret mongodb-credentials -n keptn -ojsonpath={.data.password} | base64 -d)

./upgradecollections $MONGODB_SOURCE_URL "mongodb://user:password@${MONGODB_TARGET_URL}" $CONFIGURATION_SERVICE_URL

kubectl -n keptn set env deployment/configuration-service MONGO_DB_CONNECTION_STRING='mongodb://user:password@mongodb:27017/keptn'
kubectl delete pod -n keptn -lrun=configuration-service

#print_debug "Deleting outdated keptn-datastore namespace"
#kubectl delete namespace keptn-datastore

if [[ $USECASE == "continuous-delivery" ]]; then
  # set values for the ingress-config to reflect the previous installation
  kubectl create configmap -n keptn ingress-config --from-literal=ingress_hostname_suffix=${KEPTN_DOMAIN} --from-literal=ingress_port="" --from-literal=ingress_protocol="" --from-literal=istio_gateway="public-gateway.istio-system" -oyaml --dry-run | kubectl replace -f -
fi

# check for keptn-contrib services
kubectl -n keptn get svc dynatrace-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-service detected. Upgrading to 0.7.0"
      kubectl -n keptn create secret generic keptn-credentials --from-literal="KEPTN_API_URL=${KEPTN_API_URL}" --from-literal="KEPTN_API_TOKEN=${KEPTN_API_TOKEN}"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-service/release-0.7.0/deploy/manifests/dynatrace-service/dynatrace-service.yaml
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

kubectl delete ClusterRoleBinding keptn-rbac
kubectl delete ClusterRoleBinding rbac-service-account
