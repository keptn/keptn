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
kubectl create -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-service.yml

echo "Wait 100s for docker service to get public ip..."
sleep 100

# Create a route for the docker registry service
# Store the docker registry route in a variable
export REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

# Create Jenkins
rm -f ../manifests/gen/k8s-jenkins-deployment.yml

cat ../manifests/jenkins/k8s-jenkins-deployment.yml | \
  sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITHUB_USER_EMAIL"'~' | \
  sed 's~GITHUB_ORGANIZATION_PLACEHOLDER~'"$GITHUB_ORGANIZATION"'~' | \
  sed 's~DOCKER_REGISTRY_IP_PLACEHOLDER~'"$REGISTRY_URL"'~' | \
  sed 's~DT_TENANT_URL_PLACEHOLDER~'"$DT_TENANT_URL"'~' | \
  sed 's~DT_API_TOKEN_PLACEHOLDER~'"$DT_API_TOKEN"'~' >> ../manifests/gen/k8s-jenkins-deployment.yml

kubectl create -f ../manifests/jenkins/k8s-jenkins-pvcs.yml 
kubectl create -f ../manifests/gen/k8s-jenkins-deployment.yml
kubectl create -f ../manifests/jenkins/k8s-jenkins-rbac.yml

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

# Setup credentials in Jenkins
echo "--------------------------"
echo "Setup Credentials in Jenkins "
echo "--------------------------"

# Export Jenkins route in a variable
export JENKINS_URL=$(kubectl describe svc jenkins -n cicd | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

curl -X POST http://$JENKINS_URL:24711/credentials/store/system/domain/_/createCredentials --user $JENKINS_USER:$JENKINS_PASSWORD \
--data-urlencode 'json={
  "": "0",
  "credentials": {
    "scope": "GLOBAL",
    "id": "registry-creds",
    "username": "user",
    "password": "'$TOKEN_VALUE'",
    "description": "Token used by Jenkins to push to the OpenShift container registry",
    "$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
  }
}'

curl -X POST http://$JENKINS_URL:24711/credentials/store/system/domain/_/createCredentials --user $JENKINS_USER:$JENKINS_PASSWORD \
--data-urlencode 'json={
  "": "0",
  "credentials": {
    "scope": "GLOBAL",
    "id": "git-credentials-acm",
    "username": "'$GITHUB_USER_NAME'",
    "password": "'$GITHUB_PERSONAL_ACCESS_TOKEN'",
    "description": "Token used by Jenkins to access the GitHub repositories",
    "$class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl"
  }
}'

curl -X POST http://$JENKINS_URL:24711/credentials/store/system/domain/_/createCredentials --user $JENKINS_USER:$JENKINS_PASSWORD \
--data-urlencode 'json={
  "": "0",
  "credentials": {
    "scope": "GLOBAL",
    "id": "perfsig-api-token",
    "apiToken": "'$DT_API_TOKEN'",
    "description": "Dynatrace API Token used by the Performance Signature plugin",
    "$class": "de.tsystems.mms.apm.performancesignature.dynatracesaas.model.DynatraceApiTokenImpl"
  }
}'

echo "--------------------------"
echo "End setup credentials in Jenkins "
echo "--------------------------"

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

./setupKnative.sh $JENKINS_USER $JENKINS_PASSWORD $REGISTRY_URL $CLUSTER_NAME $CLUSTER_ZONE

echo "--------------------------"
echo "End setup Knative components "
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
