#!/bin/bash

source test/utils.sh

echo "{
  \"openshiftUrl\": \"https://127.0.0.1:8443\",
  \"openshiftUser\": \"admin\",
  \"openshiftPassword\": \"admin\"
}" > creds.json

echo "Installing keptn on minishift cluster"

# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --platform=openshift --keptn-installer-image="${KEPTN_INSTALLER_IMAGE}" --use-case=quality-gates --creds=creds.json --verbose

verify_test_step $? "keptn install failed"

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

# verify the pods within the keptn-datastore namespace
verify_deployment_in_namespace "mongodb" "keptn-datastore"
verify_deployment_in_namespace "mongodb-datastore" "keptn-datastore"

echo "Installation done!"

exit 0
