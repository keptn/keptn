#!/bin/bash
source ./common/utils.sh

# Create subscriptions
kubectl apply -f ../manifests/keptn/sub-cm.yaml

# Deploy uniform
kubectl apply -f ../manifests/keptn/uniform-services.yaml
verify_kubectl $? "Deploying keptn's uniform-services failed."
wait_for_deployment_in_namespace "gatekeeper-service" "keptn"
wait_for_deployment_in_namespace "jmeter-service" "keptn"
wait_for_deployment_in_namespace "helm-service" "keptn"
wait_for_deployment_in_namespace "github-service" "keptn"
wait_for_deployment_in_namespace "pitometer-service" "keptn"
wait_for_deployment_in_namespace "servicenow-service" "keptn"

kubectl apply -f ../manifests/keptn/uniform-distributors.yaml
verify_kubectl $? "Deploying keptn's uniform-destributors failed."
wait_for_deployment_in_namespace "github-service-create-project-distributor" "keptn"
wait_for_deployment_in_namespace "github-service-onboard-service-distributor" "keptn"
wait_for_deployment_in_namespace "github-service-configure-distributor" "keptn"
wait_for_deployment_in_namespace "github-service-new-artifact-distributor" "keptn"
wait_for_deployment_in_namespace "helm-service-configuration-changed-distributor" "keptn"
wait_for_deployment_in_namespace "jmeter-service-deployment-distributor" "keptn"
wait_for_deployment_in_namespace "pitometer-service-tests-finished-distributor" "keptn"
wait_for_deployment_in_namespace "gatekeeper-service-evaluation-done-distributor" "keptn"
wait_for_deployment_in_namespace "servicenow-service-problem-distributor" "keptn"

##############################################
## Start validation of keptn's uniform      ##
##############################################
wait_for_all_pods_in_namespace "keptn"