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
export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds.json | jq -r '.dynatraceApiToken')
export DT_PAAS_TOKEN=$(cat creds.json | jq -r '.dynatracePaaSToken')
export GITHUB_ORGANIZATION=$(cat creds.json | jq -r '.githubOrg')
export DT_TENANT_URL="$DT_TENANT_ID.live.dynatrace.com"
export CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
export CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
export CLUSTER_REGION=$(cat creds.json | jq -r '.clusterRegion')
export GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')

gcloud container clusters get-credentials $CLUSTER_NAME --zone $CLUSTER_ZONE --project $GKE_PROJECT

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


# Deploy Dynatrace operator
export LATEST_RELEASE=$(curl -s https://api.github.com/repos/dynatrace/dynatrace-oneagent-operator/releases/latest | grep tag_name | cut -d '"' -f 4)
echo "Installing Dynatrace Operator $LATEST_RELEASE"
kubectl create -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/kubernetes.yaml
sleep 60
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"

rm -f ../manifests/gen/cr.yml

curl -o ../manifests/dynatrace/cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/cr.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/gen/cr.yml

kubectl create -f ../manifests/gen/cr.yml

# Apply auto tagging rules in Dynatrace
echo "--------------------------"
echo "Apply auto tagging rules in Dynatrace "
echo "--------------------------"

./applyAutoTaggingRules.sh $DT_TENANT_ID $DT_API_TOKEN

echo "--------------------------"
echo "End applying auto tagging rules in Dynatrace "
echo "--------------------------"

echo "Wait 150s for changes to apply..."
sleep 150

# Install Istio service mesh
echo "--------------------------"
echo "Setup Istio "
echo "--------------------------"

./setupIstio.sh $DT_TENANT_ID $DT_PAAS_TOKEN $CLUSTER_NAME $CLUSTER_ZONE

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

echo "Wait 10s for changes to apply..."
sleep 10

# Create Ansible Tower

echo "--------------------------"
echo "Setup Ansible Tower "
echo "--------------------------"

kubectl create -f ../manifests/ansible-tower/namespace.yml
kubectl create -f ../manifests/ansible-tower/deployment.yml
kubectl create -f ../manifests/ansible-tower/service.yml

echo "--------------------------"
echo "End setup Ansible Tower "
echo "--------------------------"

echo "----------------------------------------------------"
echo "Finished setting up infrastructure "
echo "----------------------------------------------------"

echo "To retrieve the Keptn API Token, please execute the following command"
echo "kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
