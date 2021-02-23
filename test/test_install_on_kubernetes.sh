#!/bin/bash

KEPTN_INSTALLER_REPO=${KEPTN_INSTALLER_REPO:-https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz}
KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

source test/utils.sh

echo "Installing keptn on cluster"
echo "{}" > creds.json # empty credentials file

# install Keptn using the develop version, which refers to the :latest docker images
keptn install --namespace=${KEPTN_NAMESPACE} --chart-repo="${KEPTN_INSTALLER_REPO}" --platform=kubernetes --creds=creds.json --endpoint-service-type=NodePort --use-case=continuous-delivery --verbose --hide-sensitive-data
verify_test_step $? "keptn install --chart-repo=${KEPTN_INSTALLER_REPO} - failed"

echo "Verifying that services and namespaces have been created"

# verify the deployments within the keptn namespace
verify_deployment_in_namespace "api-gateway-nginx" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "api-service" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "bridge" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "configuration-service" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "lighthouse-service" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "shipyard-controller" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "approval-service" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "remediation-service" ${KEPTN_NAMESPACE}

# verify the datastore deployments
verify_deployment_in_namespace "mongodb" ${KEPTN_NAMESPACE}
verify_deployment_in_namespace "mongodb-datastore" ${KEPTN_NAMESPACE}

# authenticate at Keptn API
API_PORT=$(kubectl get svc api-gateway-nginx -n ${KEPTN_NAMESPACE} -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
INTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
KEPTN_ENDPOINT="http://${INTERNAL_NODE_IP}:${API_PORT}"/api
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)

echo "Trying to authenticate at ${KEPTN_ENDPOINT}/api"
auth_at_keptn $KEPTN_ENDPOINT $KEPTN_API_TOKEN
#keptn auth --endpoint=${KEPTN_ENDPOINT} --api-token=$KEPTN_API_TOKEN

verify_test_step $? "Could not authenticate at Keptn API"

echo "Keptn installed in version:"
keptn version

cd ../..

echo "Installing Keptn on cluster done ✓"

exit 0
