#!/bin/bash

KEPTN_INSTALLER_REPO=${KEPTN_INSTALLER_REPO:-https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz}

source test/utils.sh

echo "Installing keptn on cluster"
echo "{}" > creds.json # empty credentials file
# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --chart-repo="${KEPTN_INSTALLER_REPO}" --platform=kubernetes --creds=creds.json --keptn-api-service-type=NodePort --verbose

verify_test_step $? "keptn install failed"

echo "Verifying that services and namespaces have been created"

# verify the deployments within the keptn namespace
verify_deployment_in_namespace "api-gateway-nginx" "keptn"
verify_deployment_in_namespace "api-service" "keptn"
verify_deployment_in_namespace "bridge" "keptn"
verify_deployment_in_namespace "configuration-service" "keptn"
verify_deployment_in_namespace "lighthouse-service" "keptn"

# verify the datastore deployments
verify_deployment_in_namespace "mongodb" "keptn"
verify_deployment_in_namespace "mongodb-datastore" "keptn"

# authenticate at Keptn API
API_PORT=$(kubectl get svc api-gateway-nginx -n keptn -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
EXTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="ExternalIP")].address }')
KEPTN_ENDPOINT=http://${INTERNAL_NODE_IP}:${API_PORT}
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)
keptn auth --endpoint=http://$KEPTN_ENDPOINT/api --api-token=$KEPTN_API_TOKEN

verify_test_step $? "Could not authenticate at Keptn API"

# verify that the keptn CLI has successfully authenticated
echo "Checking that keptn is authenticated..."
ls -la ~/.keptn/.keptn
verify_test_step $? "Could not find keptn credentials in ~/.keptn folder"

cd ../..

echo "Installation done!"

exit 0
