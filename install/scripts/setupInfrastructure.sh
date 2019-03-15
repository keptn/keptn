#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/setupInfrastructure.log)
exec 2>&1

echo "--------------------------"
echo "Setup Infrastructure "
echo "--------------------------"

# Script if you don't want to apply all yaml files manually

export JENKINS_USER=$(cat creds.json | jq -r '.jenkinsUser')
export JENKINS_PASSWORD=$(cat creds.json | jq -r '.jenkinsPassword')
export GITHUB_PERSONAL_ACCESS_TOKEN=$(cat creds.json | jq -r '.githubPersonalAccessToken')
export GITHUB_USER_NAME=$(cat creds.json | jq -r '.githubUserName')
export GITHUB_USER_EMAIL=$(cat creds.json | jq -r '.githubUserEmail')
export GITHUB_ORGANIZATION=$(cat creds.json | jq -r '.githubOrg')
export CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
export CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
export CLUSTER_REGION=$(cat creds.json | jq -r '.clusterRegion')
export GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')

set -e
gcloud container clusters get-credentials $CLUSTER_NAME --zone $CLUSTER_ZONE --project $GKE_PROJECT
set +e

# Grant cluster admin rights to gcloud user
export GCLOUD_USER=$(gcloud config get-value account)
kubectl create clusterrolebinding dynatrace-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER

# Create K8s namespaces
kubectl create -f ../manifests/k8s-namespaces.yml 

# Create container registry
kubectl create -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-service.yml

echo "Wait 100s for docker service to get public ip..."
sleep 100

# Create a route for the docker registry service
# Store the docker registry route in a variable
export REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

# Install Istio service mesh
echo "--------------------------"
echo "Setup Istio "
echo "--------------------------"

./setupIstio.sh $CLUSTER_NAME $CLUSTER_ZONE

echo "--------------------------"
echo "End setup Istio "
echo "--------------------------"

# Install knative based core components
echo "--------------------------"
echo "Setup Knative components "
echo "--------------------------"

./setupKnative.sh $REGISTRY_URL $CLUSTER_NAME $CLUSTER_ZONE

echo "--------------------------"
echo "End setup Knative components "
echo "--------------------------"

# Create Jenkins
echo "--------------------------"
echo "Setup CD Services "
echo "--------------------------"

./deployServices.sh $REGISTRY_URL

echo "--------------------------"
echo "End Setup CD Services"
echo "--------------------------"

echo "----------------------------------------------------"
echo "Finished setting up infrastructure "
echo "----------------------------------------------------"

echo "To retrieve the Keptn API Token, please execute the following command"
echo "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
