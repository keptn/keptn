#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

RED='' #'\033[0;31m'
GREEN='' #'\033[0;32m'
BLUE='' #'\033[0;34m'
NC='' #'\033[0m' # No Color

echo -e "${BLUE}--------------------------"
echo "Starting keptn installation"
echo -e "--------------------------${NC}"

# Environment variables for connecting to cluster
export GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')
export CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
export CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')

./testConnection.sh $GKE_PROJECT $CLUSTER_NAME $CLUSTER_ZONE

# Grant cluster admin rights to gcloud user
if [[ -z "${GCLOUD_USER}" ]]; then
  GCLOUD_USER=$(gcloud config get-value account)
fi
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER

# Create K8s namespaces
kubectl apply -f ../manifests/k8s-namespaces.yml

# Create container registry
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml

# Create a route for the docker registry service
# Store the docker registry route in a variable
export REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

# Install Istio service mesh
echo -e "${BLUE}--------------------------"
echo "Setup Istio "
echo -e "--------------------------${NC}"

./setupIstio.sh $CLUSTER_NAME $CLUSTER_ZONE

echo -e "${BLUE}--------------------------"
echo "End setup Istio "
echo -e "--------------------------${NC}"

# Install knative based core components
echo -e "${BLUE}--------------------------"
echo "Setup Knative components "
echo -e "--------------------------${NC}"

./setupKnative.sh $REGISTRY_URL $CLUSTER_NAME $CLUSTER_ZONE

echo -e "${BLUE}--------------------------"
echo "End setup Knative components "
echo -e "--------------------------${NC}"

# Create Jenkins
echo -e "${BLUE}--------------------------"
echo "Setup CD Services "
echo -e "--------------------------${NC}"

./deployServices.sh $REGISTRY_URL

echo -e "${BLUE}--------------------------"
echo "End Setup CD Services"
echo -e "--------------------------${NC}"

echo -e "${GREEN}----------------------------------------------------"
echo "Finished setting up infrastructure "
echo "----------------------------------------------------"

echo "To retrieve the Keptn API Token, please execute the following command"
echo -e "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode${NC}"
