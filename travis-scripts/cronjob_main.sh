#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# prints the full command before output of the command.
set -x

install_hub
install_yq
install_helm

setup_gcloud
setup_gcloud_nightly

uninstall_keptn

delete_nightly_cluster
create_nightly_cluster

install_sed

cd ./install/scripts

source ./defineCredentialsUtils.sh

# Set enviornment variables used in replaceCreds function
export GITU=$GITHUB_USER_NAME_NIGHTLY
export GITAT=$GITHUB_TOKEN_NIGHTLY
export GITE=$GITHUB_EMAIL_NIGHTLY
export CLN=$CLUSTER_NAME_NIGHTLY
export CLZ=$CLOUDSDK_COMPUTE_ZONE
export PROJ=$PROJECT_NAME
export GITO=$GITHUB_ORG_NIGHTLY

replaceCreds

# Add execution right because there 
# is a rights problem with e.g. testConnection
find . -type f -exec chmod +x {} \;
source ./installKeptn.sh
cd ../..

# Test front-end keptn v.0.1
# export FRONT_END_DEV=$(kubectl describe svc front-end -n dev | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
# export FRONT_END_STAGING=$(kubectl describe svc front-end -n staging | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
# export FRONT_END_PRODUCTION=$(kubectl describe svc front-end -n production | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

export ISTIO_INGRESS=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
export_names

# Execute unit tests
execute_cli_tests

build_and_install_cli

# Execute end-to-end test
cd test
source ./testOnboarding.sh

#- cat ../test/keptn.postman_environment.json |sed 's~FRONT_END_DEV_PLACEHOLDER~'"$FRONT_END_DEV"'~' |sed 's~FRONT_END_STAGING_PLACEHOLDER~'"$FRONT_END_STAGING"'~' |sed 's~FRONT_END_PRODUCTION_PLACEHOLDER~'"$FRONT_END_PRODUCTION"'~' |sed 's~ISTIO_INGRESS_PLACEHOLDER~'"$ISTIO_INGRESS"'~' >> ../test/env.json
#- npm install newman
#- node_modules/.bin/newman run ../test/keptn.postman_collection.json -e ../test/env.json


