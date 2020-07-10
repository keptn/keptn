#!/bin/bash

source test/utils.sh

echo "Installing keptn on cluster"
echo "{}" > creds.json # empty credentials file
# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --keptn-installer-image="${KEPTN_INSTALLER_IMAGE}" --platform=kubernetes --creds=creds.json --gateway=NodePort --verbose

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

# verify the datastore deployments
verify_deployment_in_namespace "mongodb" "keptn"
verify_deployment_in_namespace "mongodb-datastore" "keptn"


cd ../..

echo "Installation done!"

exit 0
