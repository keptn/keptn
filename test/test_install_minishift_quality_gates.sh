#!/bin/bash

KEPTN_INSTALLER_REPO=${KEPTN_INSTALLER_REPO:-https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz}

source test/utils.sh

echo "{
  \"openshiftUrl\": \"https://127.0.0.1:8443\",
  \"openshiftUser\": \"admin\",
  \"openshiftPassword\": \"admin\"
}" > creds.json

echo "Installing keptn on minishift cluster"

# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --platform=openshift --chart-repo="${KEPTN_INSTALLER_REPO}" --creds=creds.json --verbose

verify_test_step $? "keptn install failed"

oc expose svc/api-gateway-nginx -n keptn --hostname=api.keptn.127.0.0.1.nip.io

sleep 30

KEPTN_API_URL=api.keptn.127.0.0.1.nip.io
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)
keptn auth --endpoint=http://$KEPTN_API_URL/api --api-token=$KEPTN_API_TOKEN

# verify that the keptn CLI has successfully authenticated
echo "Checking that keptn is authenticated..."
ls -la ~/.keptn/.keptn
verify_test_step $? "Could not find keptn credentials in ~/.keptn folder"


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

echo "Installation done!"

exit 0
