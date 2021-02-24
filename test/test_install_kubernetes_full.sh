#!/bin/bash

KEPTN_INSTALLER_REPO=${KEPTN_INSTALLER_REPO:-https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz}

# shellcheck disable=SC1091
source test/utils.sh

# install istio
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.6.5 sh -
cd istio-1.6.5 || exit
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo

# verify the pods within istio-system
verify_deployment_in_namespace "istio-ingressgateway" "istio-system"
verify_deployment_in_namespace "istio-pilot" "istio-system"
verify_deployment_in_namespace "istio-citadel" "istio-system"
verify_deployment_in_namespace "istio-sidecar-injector" "istio-system"

echo "Installing Keptn on cluster"
echo "{}" > creds.json # empty credentials file

# install Keptn using the develop version, which refers to the :latest docker images
keptn install --chart-repo="${KEPTN_INSTALLER_REPO}" --platform=kubernetes --creds=creds.json --endpoint-service-type=NodePort --verbose --use-case=continuous-delivery --hide-sensitive-data
verify_test_step $? "keptn install --chart-repo=${KEPTN_INSTALLER_REPO} - failed"

echo "Verifying that services and namespaces have been created"

# verify the deployments within the keptn namespace
verify_deployment_in_namespace "api-gateway-nginx" "keptn"
verify_deployment_in_namespace "api-service" "keptn"
verify_deployment_in_namespace "bridge" "keptn"
verify_deployment_in_namespace "configuration-service" "keptn"
verify_deployment_in_namespace "approval-service" "keptn"
verify_deployment_in_namespace "jmeter-service" "keptn"
verify_deployment_in_namespace "lighthouse-service" "keptn"
verify_deployment_in_namespace "mongodb" "keptn"
verify_deployment_in_namespace "mongodb-datastore" "keptn"

# authenticate at Keptn API
API_PORT=$(kubectl get svc api-gateway-nginx -n keptn -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
INTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
KEPTN_ENDPOINT=http://${INTERNAL_NODE_IP}:${API_PORT}/api
# shellcheck disable=SC1083
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)
auth_at_keptn "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN"
#keptn auth --endpoint=$KEPTN_ENDPOINT --api-token=$KEPTN_API_TOKEN

verify_test_step $? "Could not authenticate at Keptn API"

echo "Keptn installed in version:"
keptn version

cd ../..

echo "Installing Keptn on cluster done ✓"

exit 0
