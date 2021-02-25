#!/bin/bash
# shellcheck disable=SC2181

# shellcheck disable=SC1091
source ./utils.sh

if [ "$HELM_RELEASE_UPGRADE" == "true" ];
then
  # Upgrade from Helm v2 to Helm v3
  helm init --client-only
  verify_install_step $? "Helm init failed."
  RELEASES=$(helm list -aq)
  verify_install_step $? "Helm list failed."
  echo "$RELEASES"

  helm3 plugin install https://github.com/helm/helm-2to3
  verify_install_step $? "Helm-2to3 plugin installation failed."
  yes y | helm3 2to3 move config
  verify_install_step $? "Helm-2to3 move of config failed."

  for release in $RELEASES; do
    helm3 2to3 convert "$release" --dry-run
    verify_install_step $? "Helm2-to3 release convertion dry-run failed"
    helm3 2to3 convert "$release"
    verify_install_step $? "Helm2-to3 release convertion failed"
  done

  yes y | helm3 2to3 cleanup --tiller-cleanup
  verify_install_step $? "Helm2-to3 cleanup failed"

fi

PREVIOUS_KEPTN_VERSION="0.6.2"
KEPTN_VERSION=${KEPTN_VERSION:-"0.7.0"}
HELM_CHART_URL=${HELM_CHART_URL:-"https://storage.googleapis.com/keptn-installer/0.7.0"}
MONGODB_SOURCE_URL=${MONGODB_SOURCE_URL:-"mongodb://user:password@mongodb.keptn-datastore:27017/keptn"}
MONGODB_TARGET_URL=${MONGODB_TARGET_URL:-"mongodb.keptn:27017/keptn"}

print_debug "Upgrading from Keptn 0.6.2 to $KEPTN_VERSION"

KEPTN_API_URL=https://$(kubectl get cm keptn-domain -n keptn -o jsonpath='{.data.app_domain}')/api
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o jsonpath='{.data.keptn-api-token}' | base64 -d)
KEPTN_DOMAIN=$(kubectl get cm keptn-domain -n keptn -o jsonpath='{.data.app_domain}')

print_debug "Check if Keptn 0.6.2 is currently installed"
API_IMAGE=$(kubectl get deployment -n keptn api-service -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

  if [[ $API_IMAGE != 'keptn/api:0.6.2' ]]; then
    print_error "Installed Keptn version does not match 0.6.2. aborting."
    exit 1
  fi

USE_CASE=""
# check if full installation is available
kubectl -n keptn get svc approval-service

  if [[ $? == '0' ]]; then
      print_debug "Full installation detected. Upgrading CD and CO services"
      USE_CASE="continuous-delivery"
  fi

./upgradecollections "$MONGODB_SOURCE_URL" "mongodb://user:password@${MONGODB_TARGET_URL}" "$CONFIGURATION_SERVICE_URL" "store-projects-mv"

# copy content from previous configuration-service PVC
mkdir config-svc-backup
CONFIG_SERVICE_POD=$(kubectl get pods -n keptn -lrun=configuration-service -ojsonpath='{.items[0].metadata.name}')
kubectl cp "keptn/$CONFIG_SERVICE_POD:/data" ./config-svc-backup/ -c configuration-service

old_manifests=(
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/core.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/api-gateway-nginx.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/quality-gates.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/continuous-deployment.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/continuous-operations.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/keptn-api-virtualservice.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/keptn-ingress.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/uniform-services-openshift.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/keptn/uniform-distributors-openshift.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/nats/nats-cluster.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/nats/nats-operator-deploy.yaml"
  "https://raw.githubusercontent.com/keptn/keptn/release-$PREVIOUS_KEPTN_VERSION/installer/manifests/nats/nats-operator-prereqs.yaml"
  )

for manifest in "${old_manifests[@]}"
do
   :
   if curl --head --silent -k --fail "$manifest" > /dev/null;
     then
      kubectl delete -f "$manifest"
      continue
     else
      print_error "Required manifest $manifest not available. Aborting upgrade."
      exit 1
    fi
done

# delete resource that have been generated procedurally by the installer to avoid conflicts with helm install
kubectl delete secret -n keptn keptn-api-token
kubectl delete configmap -n keptn keptn-domain

BRIDGE_USERNAME=""
kubectl get secret -n keptn bridge-credentials
  if [[ $? == '0' ]]; then
    print_debug "Bridge credentials detected. Fetching credentials"
    BRIDGE_USERNAME=$(kubectl get secret bridge-credentials -n keptn -o jsonpath='{.data.BASIC_AUTH_USERNAME}' | base64 -d)
    BRIDGE_PASSWORD=$(kubectl get secret bridge-credentials -n keptn -o jsonpath='{.data.BASIC_AUTH_PASSWORD}' | base64 -d)
    kubectl -n keptn delete secret bridge-credentials
  fi


# install helm chart
helm3 repo add keptn "$HELM_CHART_URL"

if [[ $USE_CASE == "continuous-delivery" ]]; then
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

wait_for_all_pods_in_namespace "keptn"

kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN" -oyaml --dry-run | kubectl replace -f -

if [[ $BRIDGE_USERNAME == "" ]]; then
  echo "No previous bridge credentials found. No need to update"
else
  echo "Setting bridge credentials to previous values"
  kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=$BRIDGE_USERNAME" --from-literal="BASIC_AUTH_PASSWORD=$BRIDGE_PASSWORD" -oyaml --dry-run | kubectl replace -f -
fi

CONFIG_SERVICE_POD=$(kubectl get pods -n keptn -lapp.kubernetes.io/name=configuration-service -ojsonpath='{.items[0].metadata.name}')
kubectl cp ./config-svc-backup/* "keptn/$CONFIG_SERVICE_POD:/data" -c configuration-service

kubectl cp ./reset-git-repos.sh "keptn/$CONFIG_SERVICE_POD:/" -c configuration-service
kubectl exec -n keptn "$CONFIG_SERVICE_POD" -c configuration-service -- chmod +x -R ./reset-git-repos.sh
kubectl exec -n keptn "$CONFIG_SERVICE_POD" -c configuration-service -- ./reset-git-repos.sh

# shellcheck disable=SC2034
MONGO_TARGET_USER=$(kubectl get secret mongodb-credentials -n keptn -o jsonpath='{.data.user}' | base64 -d)
MONGO_TARGET_PASSWORD=$(kubectl get secret mongodb-credentials -n keptn -o jsonpath='{.data.password}' | base64 -d)

./upgradecollections "$MONGODB_SOURCE_URL" "mongodb://user:${MONGO_TARGET_PASSWORD}@${MONGODB_TARGET_URL}" "$CONFIGURATION_SERVICE_URL"

if [[ $USE_CASE == "continuous-delivery" ]]; then
  # set values for the ingress-config to reflect the previous installation
  kubectl create configmap -n keptn ingress-config --from-literal="ingress_hostname_suffix=${KEPTN_DOMAIN}" --from-literal=ingress_port="" --from-literal=ingress_protocol="" --from-literal=istio_gateway="public-gateway.istio-system" -oyaml --dry-run | kubectl replace -f -
fi

# check for keptn-contrib services
kubectl -n keptn get svc dynatrace-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-service detected. Upgrading to 0.8.0"
      DT_TENANT=$(kubectl get secret dynatrace -n keptn -o jsonpath='{.data.DT_TENANT}' | base64 -d)
      DT_API_TOKEN=$(kubectl get secret dynatrace -n keptn -o jsonpath='{.data.DT_API_TOKEN}' | base64 -d)

      kubectl -n keptn create secret generic dynatrace --from-literal="KEPTN_API_URL=${KEPTN_API_URL}" --from-literal="KEPTN_API_TOKEN=${KEPTN_API_TOKEN}" --from-literal="DT_API_TOKEN=${DT_API_TOKEN}" --from-literal="DT_TENANT=${DT_TENANT}" -oyaml --dry-run | kubectl replace -f -
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-service/0.8.0/deploy/service.yaml
  fi

kubectl -n keptn get svc dynatrace-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Dynatrace-sli-service detected. Upgrading to 0.5.0"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/0.5.0/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-service detected. Upgrading to 0.3.5"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/0.3.5/deploy/service.yaml
  fi

kubectl -n keptn get svc prometheus-sli-service

  if [[ $? == '0' ]]; then
      print_debug "Prometheus-sli-service detected. Upgrading to 0.2.2"
      kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-sli-service/0.2.2/deploy/service.yaml
  fi

# delete all pods in keptn namespace to make sure all secret references are updated
kubectl delete pods -n keptn -l 'app notin (upgrader)'
wait_for_all_pods_in_namespace "keptn"

kubectl delete ClusterRoleBinding keptn-rbac --ignore-not-found
kubectl delete ClusterRoleBinding rbac-service-account --ignore-not-found
