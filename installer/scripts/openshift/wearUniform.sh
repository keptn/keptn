#!/bin/bash
source ./common/utils.sh

# Deploy uniform

kubectl apply -f ../manifests/openshift/uniform-distributors-openshift.yaml
verify_kubectl $? "Deploying keptn's uniform-destributors-openshift failed."
wait_for_deployment_in_namespace "openshift-route-service-create-project-distributor" "keptn"

##############################################
## Start validation of keptn's uniform      ##
##############################################
wait_for_all_pods_in_namespace "keptn"