#!/bin/bash
source ./common/utils.sh

# Create subscriptions
kubectl apply -f ../manifests/keptn/sub-cm-openshift.yaml

# Deploy uniform
kubectl apply -f ../manifests/keptn/uniform-services-openshift.yaml --wait
verify_kubectl $? "Deploying keptn's uniform-services failed."
wait_for_deployment_in_namespace "openshift-route-service" "keptn"

kubectl apply -f ../manifests/keptn/uniform-distributors-openshift.yaml
verify_kubectl $? "Deploying keptn's uniform-destributors-openshift failed."
wait_for_deployment_in_namespace "openshift-route-service-create-project-distributor" "keptn"

##############################################
## Start validation of keptn's uniform      ##
##############################################
wait_for_all_pods_in_namespace "keptn"