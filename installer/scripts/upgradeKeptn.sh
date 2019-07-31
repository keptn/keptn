#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/upgradeKeptn.log)
exec 2>&1

source ./common/utils.sh

echo "Starting upgrade to keptn 0.3.0"

GITHUB_USER_NAME=$1
GITHUB_PERSONAL_ACCESS_TOKEN=$2

if [ -z $1 ]
then
  echo "Please provide the github username as first parameter"
  echo ""
  echo "Usage: ./upgradeKeptn.sh GitHub_username GitHub_personal_access_token"
  exit 1
fi

if [ -z $2 ]
then
  echo "Please provide the GitHub personal access token as second parameter"
  echo ""
  echo "Usage: ./upgradeKeptn.sh GitHub_username GitHub_personal_access_token"
  exit 1
fi

if [[ $GITHUB_USER_NAME = '' ]]
then
  echo "GitHub username not set."
  exit 1
fi

if [[ $GITHUB_PERSONAL_ACCESS_TOKEN = '' ]]
then
  echo "GitHub personal access token not set."
  exit 1
fi

kubectl delete -f ../manifests/keptn/uniform-subscriptions.yaml --ignore-not-found
verify_kubectl $? "Removing keptn uniform subscriptions failed."
kubectl delete -f ../manifests/keptn/uniform-services.yaml --ignore-not-found
verify_kubectl $? "Removing keptn uniform services failed."

kubectl delete -f ../manifests/keptn/core.yaml --ignore-not-found
verify_kubectl $? "Removing keptn core services failed."

# Remove subscriptions of Jenkins service
kubectl delete subscription jenkins-configuration-changed-subscription -n keptn --ignore-not-found
kubectl delete subscription jenkins-deployment-finished-subscription -n keptn --ignore-not-found
kubectl delete subscription jenkins-evaluation-done-subscription -n keptn --ignore-not-found

# Install tiller for helm
print_info "Installing Tiller"
kubectl apply -f ../manifests/tiller/tiller.yaml
helm init --service-account tiller
print_info "Installing Tiller done"

print_info "Upgrading keptn core"
./common/setupKeptn.sh
print_info "Upgrading keptn core done"

print_info "Upgrading keptn uniform"
./common/wearUniform.sh
print_info "Upgrading keptn uniform done"

echo "Upgrade to keptn 0.3.0 done."
