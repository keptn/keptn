#!/bin/bash

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

sleep 100

# Create a route for the docker registry service
# Store the docker registry route in a variable
export REGISTRY_URL=$(kubectl describe svc docker-registry -n cicd | grep IP: | sed 's/IP:[ \t]*//')

# Create Jenkins
cat ../manifests/jenkins/k8s-jenkins-deployment.yml | sed 's/GITHUB_USER_EMAIL_PLACEHOLDER/'"$GITHUB_USER_EMAIL"'/' | \
  sed 's/GITHUB_ORGANIZATION_PLACEHOLDER/'"$GITHUB_ORGANIZATION"'/' | \
  sed 's/DOCKER_REGISTRY_IP_PLACEHOLDER/'"$REGISTRY_URL"'/' | \
  sed 's/DT_TENANT_URL_PLACEHOLDER/'"$DT_TENANT_URL"'/' | \
  sed 's/DT_API_TOKEN_PLACEHOLDER/'"$DT_API_TOKEN"'/' >> ../manifests/jenkins/k8s-jenkins-deployment_tmp.yml

kubectl create -f ../manifests/jenkins/k8s-jenkins-pvcs.yml 
kubectl create -f ../manifests/jenkins/k8s-jenkins-deployment_tmp.yml
kubectl create -f ../manifests/jenkins/k8s-jenkins-rbac.yml

rm ../manifests/jenkins/k8s-jenkins-deployment_tmp.yml

# Deploy Dynatrace operator
kubectl create namespace dynatrace
kubectl create -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/master/deploy/kubernetes.yaml

sleep 60

kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"
cat ../manifests/dynatrace/oneagent.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/dynatrace/oneagent_tmp.yml

kubectl create -f ../manifests/dynatrace/oneagent_tmp.yml
rm ../manifests/dynatrace/oneagent_tmp.yml

# Deploy sockshop application
./deploySockshop.sh

# Set up credentials in Jenkins
sleep 150

# Export Jenkins route in a variable
export JENKINS_URL=$(kubectl describe svc jenkins -n cicd | grep "LoadBalancer Ingress:" | sed 's/LoadBalancer Ingress:[ \t]*//')

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

# Install Istio service mesh
./setupIstio.sh $DT_TENANT_ID $DT_PAAS_TOKEN

# Create Ansible Tower
kubectl create -f ../manifests/ansible-tower/namespace.yml
kubectl create -f ../manifests/ansible-tower/deployment.yml
kubectl create -f ../manifests/ansible-tower/service.yml

sleep 120

./configureAnsible.sh