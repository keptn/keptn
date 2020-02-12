#!/bin/bash

source test/utils.sh

echo "Download the CLI"
# Download latest KEPTN cli for linux
wget https://storage.googleapis.com/keptn-cli/latest/keptn-linux.zip
unzip keptn-linux.zip

sudo mv keptn /usr/local/bin/keptn


echo "Installing keptn on cluster"
echo "{}" > creds.json # empty credentials file
# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --keptn-installer-image=keptn/installer:latest --platform=kubernetes --use-case=quality-gates --creds=creds.json --gateway=NodePort --verbose

verify_test_step $? "keptn install failed"

# verify that the keptn CLI has successfully authenticated
echo "Checking that keptn is authenticated..."
ls -la ~/.keptn/.keptn
verify_test_step $? "Could not find keptn credentials in ~/.keptn folder"

echo "Verifying that services and namespaces have been created"

# verify the deployments within the keptn namespace
verify_deployment_in_namespace "api" "keptn"
verify_deployment_in_namespace "bridge" "keptn"
verify_deployment_in_namespace "configuration-service" "keptn"
verify_deployment_in_namespace "lighthouse-service" "keptn"

# verify the pods within the keptn-datastore namespace
verify_deployment_in_namespace "mongodb" "keptn-datastore"
verify_deployment_in_namespace "mongodb-datastore" "keptn-datastore"


cd ../..

echo "Installation done!"

exit 0
