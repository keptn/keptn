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

# Grant cluster admin rights to gcloud user
export GCLOUD_USER=$(gcloud config get-value account)
kubectl create clusterrolebinding dynatrace-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER

# Create K8s namespaces
kubectl create -f ../manifests/k8s-namespaces.yml

# Create container registry
kubectl create -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl create -f ../manifests/container-registry/k8s-docker-registry-service.yml

echo "wait 30s for docker service to get ip..."
sleep 30

# Create a route for the docker registry service
# Store the docker registry route in a variable
export REGISTRY_URL=$(kubectl describe svc docker-registry -n cicd | grep IP: | sed 's~IP:[ \t]*~~')

echo "REGISTRY_URL: " $REGISTRY_URL

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
kubectl create namespace dynatrace
kubectl create -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/kubernetes.yaml --validate=false
sleep 60
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"
rm -f ../manifests/gen/oneagent-cr.yml
curl -o ../manifests/dynatrace/oneagent-cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/oneagent-cr.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/dynatrace/cr_tmp.yml
mv ../manifests/dynatrace/cr_tmp.yml ../manifests/gen/oneagent-cr.yml
kubectl create -f ../manifests/gen/oneagent-cr.yml

# Create a Bearer token for authenticating against the Kubernetes API
kubectl apply -f kubernetes-monitoring-service-account.yaml

# Apply auto tagging rules in Dynatrace
echo "--------------------------"
echo "Apply auto tagging rules in Dynatrace "
echo "--------------------------"

./applyAutoTaggingRules.sh

echo "--------------------------"
echo "End applying auto tagging rules in Dynatrace "
echo "--------------------------"

echo "--------------------------"
echo "Apply Dynatrace Request Attributes"
echo "--------------------------"

./applyRequestAttributeRules.sh

echo "--------------------------"
echo "End applying Dynatrace Request Attributes"
echo "--------------------------"

# Deploy sockshop application
echo "--------------------------"
echo "Deploy SockShop "
echo "--------------------------"

./deploySockshop.sh

echo "--------------------------"
echo "End Deploy Sockshop "
echo "--------------------------"

echo "wait for changes to apply..."
sleep 150

# Set up credentials in Jenkins
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
echo "Set up Istio "
echo "--------------------------"

./setupIstio.sh $DT_TENANT_ID $DT_PAAS_TOKEN

echo "--------------------------"
echo "End set up Istio "
echo "--------------------------"

# Create Ansible Tower

echo "--------------------------"
echo "Setup Ansible Tower "
echo "--------------------------"

kubectl create -f ../manifests/ansible-tower/namespace.yml
kubectl create -f ../manifests/ansible-tower/deployment.yml
kubectl create -f ../manifests/ansible-tower/service.yml

echo "wait 2 minutes for changes to apply..."
sleep 120

./configureAnsible.sh

echo "--------------------------"
echo "End set up Ansible Tower "
echo "--------------------------"

echo "----------------------------------------------------"
echo "Finished setting up infrastructure "
echo "----------------------------------------------------"
