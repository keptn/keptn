#!/bin/bash

source test/utils.sh

KEPTN_INSTALLER_IMAGE=${KEPTN_INSTALLER_IMAGE:-keptn/installer:latest}
KEPTN_INSTALLATION_TYPE=${KEPTN_INSTALLATION_TYPE:-FULL}
PROJECT_NAME=${PROJECT_NAME:-sockshop}

# Prepare creds.json file
cd ./installer/scripts

export CLN=$CLUSTER_NAME_NIGHTLY
export CLZ=$CLOUDSDK_COMPUTE_ZONE	
export PROJ=$PROJECT_NAME

source ./gke/defineCredentialsHelper.sh
replaceCreds

echo "Installing keptn on cluster"

INGRESS_CONFIG=""

# use a different installation method if re-using istio
if [[ "$KEPTN_INSTALLATION_TYPE" == "REUSE-ISTIO" ]]; then
  INGRESS_CONFIG="--ingress-install-option=Reuse"
  echo "... using existing istio"
fi

# Install keptn (using the develop version, which should point the :latest docker images)
keptn install ${INGRESS_CONFIG} --keptn-installer-image="${KEPTN_INSTALLER_IMAGE}" --creds=creds.json --verbose
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
verify_deployment_in_namespace "gatekeeper-service" "keptn"
verify_deployment_in_namespace "jmeter-service" "keptn"
verify_deployment_in_namespace "lighthouse-service" "keptn"

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
