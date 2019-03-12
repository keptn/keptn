#!/bin/sh
REGISTRY_URI=$1

export JENKINS_USER=$(cat creds.json | jq -r '.jenkinsUser')
export JENKINS_PASSWORD=$(cat creds.json | jq -r '.jenkinsPassword')
export GITHUB_USER_EMAIL=$(cat creds.json | jq -r '.githubUserEmail')
export GITHUB_ORGANIZATION=$(cat creds.json | jq -r '.githubOrg')
export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds.json | jq -r '.dynatraceApiToken')
export DT_TENANT_URL="$DT_TENANT_ID.live.dynatrace.com"

# Deploy Jenkins - see keptn/install/setupInfrastructure.sh:
rm -f ../manifests/gen/k8s-jenkins-deployment.yml

export GATEWAY=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

cat config/jenkins/k8s-jenkins-deployment.yml | \
  sed 's~GATEWAY_PLACEHOLDER~'"$GATEWAY"'~' | \
  sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITHUB_USER_EMAIL"'~' | \
  sed 's~GITHUB_ORGANIZATION_PLACEHOLDER~'"$GITHUB_ORGANIZATION"'~' | \
  sed 's~DOCKER_REGISTRY_IP_PLACEHOLDER~'"$REGISTRY_URL"'~' | \
  sed 's~DT_TENANT_URL_PLACEHOLDER~'"$DT_TENANT_URL"'~' | \
  sed 's~DT_API_TOKEN_PLACEHOLDER~'"$DT_API_TOKEN"'~' >> ../manifests/gen/k8s-jenkins-deployment.yml

kubectl create -f ../manifests/jenkins/k8s-jenkins-pvcs.yml 
kubectl create -f ../manifests/gen/k8s-jenkins-deployment.yml
kubectl create -f ../manifests/jenkins/k8s-jenkins-rbac.yml
kubectl create -f ../manifests/jenkins/k8s-jenkins-service-entry.yml

echo "Wait 100s for Jenkins..."
sleep 100

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

# Create secret and deploy jenkins-service
kubectl create secret generic -n keptn jenkins-secret --from-literal=jenkinsurl="jenkins.keptn.svc.cluster.local" --from-literal=user="$JENKINS_USER" --from-literal=password="$JENKINS_PASSWORD"

rm -f ../manifests/gen/service.yaml

cat config/service/service.yaml | \
  sed 's~REGISTRY_URI_PLACEHOLDER~'"$REGISTRY_URI"'~' >> ../manifests/gen/service.yaml

kubectl apply -f ../manifests/gen/service.yaml

# Deploy Tiller for Helm
kubectl -n kube-system create serviceaccount tiller
kubectl create clusterrolebinding tiller --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller

#kubectl -n kube-system  rollout status deploy/tiller-deploy