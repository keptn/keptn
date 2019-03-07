#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# Causes the shell to exit immediately if a simple command exits with a nonzero exit value as well as
# prints the full command before output of the command.
set -e -x

install_hub
install_yq
setup_gcloud
setup_glcoud_master
install_sed

# Run scripts in ~/keptn/install/scripts
cd install/scripts
./forkGitHubRepositories.sh $GITHUB_ORG

cat ./creds.sav |sed 's~DYNATRACE_TENANT_PLACEHOLDER~'"$DT_TENANT"'~' |sed 's~DYNATRACE_API_TOKEN~'"$DT_API_TOKEN"'~' |sed 's~DYNATRACE_PAAS_TOKEN~'"$DT_PAAS_TOKEN"'~' |sed 's~GITHUB_USER_NAME_PLACEHOLDER~'"$GITHUB_USER_NAME"'~' |sed 's~PERSONAL_ACCESS_TOKEN_PLACEHOLDER~'"$GITHUB_TOKEN"'~' |sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITHUB_EMAIL"'~' |sed 's~GITHUB_ORG_PLACEHOLDER~'"$GITHUB_ORG"'~' >> creds.json
./setupInfrastructure.sh
cd ../..

# Test front-end keptn v.0.1
# export FRONT_END_DEV=$(kubectl describe svc front-end -n dev | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
# export FRONT_END_STAGING=$(kubectl describe svc front-end -n staging | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
# export FRONT_END_PRODUCTION=$(kubectl describe svc front-end -n production | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

export ISTIO_INGRESS=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
export_names

#- cat ../test/keptn.postman_environment.json |sed 's~FRONT_END_DEV_PLACEHOLDER~'"$FRONT_END_DEV"'~' |sed 's~FRONT_END_STAGING_PLACEHOLDER~'"$FRONT_END_STAGING"'~' |sed 's~FRONT_END_PRODUCTION_PLACEHOLDER~'"$FRONT_END_PRODUCTION"'~' |sed 's~ISTIO_INGRESS_PLACEHOLDER~'"$ISTIO_INGRESS"'~' >> ../test/env.json
#- npm install newman
#- node_modules/.bin/newman run ../test/keptn.postman_collection.json -e ../test/env.json
        
# Clean up cluster
#- ./cleanupCluster.sh

# Delete K8s cluster
gcloud container clusters delete $CLUSTER_NAME --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME --quiet
