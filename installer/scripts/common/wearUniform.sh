#!/bin/bash
source ./common/utils.sh

# Deploy uniform

#### The uniform is empty: uniform-services.yaml and uniform-distributors.yaml
#  kubectl apply -f ../manifests/keptn/uniform-services.yaml
#  verify_kubectl $? "Deploying keptn's uniform-services failed."

#  kubectl apply -f ../manifests/keptn/uniform-distributors.yaml
#  verify_kubectl $? "Deploying keptn's uniform-destributors failed."

##############################################
## Start validation of keptn's uniform      ##
##############################################
wait_for_all_pods_in_namespace "keptn"