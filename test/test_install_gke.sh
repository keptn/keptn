#!/bin/bash

source test/utils.sh

echo "Building the CLI..."
# Build and install keptn CLI
cd cli/
dep ensure
go build -o keptn
sudo mv keptn /usr/local/bin/keptn
cd ..
# ToDo: We should really just download the nightly build of the CLI here

# Prepare creds.json file
cd ./installer/scripts

export GITU=$GITHUB_USER_NAME_NIGHTLY	
export GITAT=$GITHUB_TOKEN_NIGHTLY	
export CLN=$CLUSTER_NAME_NIGHTLY	
export CLZ=$CLOUDSDK_COMPUTE_ZONE	
export PROJ=$PROJECT_NAME	
export GITO=$GITHUB_ORG_NIGHTLY	

source ./gke/defineCredentialsHelper.sh
replaceCreds

echo "Installing keptn on cluster"

# Install keptn (using the develop version, which should point the :latest docker images)
keptn install --keptn-version=develop --creds=creds.json --verbose

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
verify_deployment_in_namespace "gatekeeper-service" "keptn"

# verify the pods within the keptn-datastore namespace
verify_deployment_in_namespace "mongodb" "keptn-datastore"
verify_deployment_in_namespace "mongodb-datastore" "keptn-datastore"

# verify the pods within istio-system
verify_deployment_in_namespace "istio-ingressgateway" "istio-system"
verify_deployment_in_namespace "istio-pilot" "istio-system"
verify_deployment_in_namespace "istio-citadel" "istio-system"
verify_deployment_in_namespace "istio-sidecar-injector" "istio-system"


cd ../..

echo "Installation done!"

exit 0
