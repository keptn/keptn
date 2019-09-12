#!/bin/bash
source ./common/utils.sh

# Create subscriptions
kubectl apply -f ../manifests/keptn/sub-cm.yaml

# Deploy uniform
kubectl apply -f ../manifests/keptn/uniform-services.yaml
verify_kubectl $? "Deploying keptn's uniform-services failed."
wait_for_deployment_in_namespace "servicenow-service" "keptn"

kubectl apply -f ../manifests/keptn/uniform-distributors.yaml
verify_kubectl $? "Deploying keptn's uniform-destributors failed."
wait_for_deployment_in_namespace "servicenow-service-problem-distributor" "keptn"

##############################################
## Start validation of keptn's uniform      ##
##############################################
wait_for_all_pods_in_namespace "keptn"