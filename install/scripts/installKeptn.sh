#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

echo 'Starting installation of keptn' 

# Environment variables for test connection to cluster
if [[ -z "${GKE_PROJECT}" ]]; then
  echo "GKE_PROJECT not set, take it from creds.json"
  GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')
fi

if [[ -z "${CLUSTER_NAME}" ]]; then
  echo "CLUSTER_NAME not set, take it from creds.json"
  CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
fi

if [[ -z "${CLUSTER_ZONE}" ]]; then
  echo "CLUSTER_ZONE not set, take it from creds.json"
  CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
fi

./testConnection.sh $GKE_PROJECT $CLUSTER_NAME $CLUSTER_ZONE

# Grant cluster admin rights to gcloud user
if [[ -z "${GCLOUD_USER}" ]]; then
  echo "GCLOUD_USER not set, retrieve it using gcloud"
  GCLOUD_USER=$(gcloud config get-value account)
fi
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER

# Create K8s namespaces
kubectl apply -f ../manifests/k8s-namespaces.yml

# Create container registry
echo "Creating container registry"
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml
echo "Creating container registry done"

# Install Istio service mesh
echo "Installing Istio"
./setupIstio.sh $CLUSTER_NAME $CLUSTER_ZONE
echo "Installing Istio done"

# Install knative core components
echo "Installing Knative"
./setupKnative.sh $CLUSTER_NAME $CLUSTER_ZONE
echo "Installing Knative done"

# Install keptn core services - Install keptn channels
echo "Install keptn"
./setupKeptn.sh
echo "Install keptn done"

# Install services
echo "Wear uniform"
./wearUniform.sh
echo "Keptn wears uniform"

# Install done
echo "Installation of keptn complete."

echo "To retrieve the Keptn API Token, please execute the following command"
echo -e "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode${NC}"
