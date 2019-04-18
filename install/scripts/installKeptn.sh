#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

echo 'Starting installation of keptn' 

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
echo "Installing Istio"
./setupIstio.sh $CLUSTER_NAME $CLUSTER_ZONE
echo "Installing Istio done"

# Install knative based core components
echo "Installing Knative"
./setupKnative.sh $CLUSTER_NAME $CLUSTER_ZONE
echo "Installing Knative done"

# Install keptn core services - Install keptn channels
echo "Install keptn"
./setupKeptn.sh $REGISTRY_URL
echo "Install keptn done"

# Install services
echo "Wear uniform"
./deployServices.sh $REGISTRY_URL
echo "Wear uniform done"

echo "Installation of keptn complete."

echo "To retrieve the Keptn API Token, please execute the following command"
echo -e "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode${NC}"
