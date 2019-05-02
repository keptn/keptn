#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/upgradeKeptn.log)
exec 2>&1

source ./utils.sh

echo "Starting upgrade to keptn 0.2.1"

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
  echo "Please provide the GitHub personal access token as 2nd parameter"
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

SERVICENAME=control
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/keptn-$SERVICENAME:0.2.1
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml


SERVICENAME=authenticator
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/keptn-$SERVICENAME:0.2.1
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml


SERVICENAME=event-broker
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/keptn-$SERVICENAME:0.2.1
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml


SERVICENAME=event-broker-ext
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/keptn-$SERVICENAME:0.2.1
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml


SERVICENAME=github-service
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/$SERVICENAME:0.1.1
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml


SERVICENAME=pitometer-service
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/$SERVICENAME:0.1.1
kubectl apply -f $SERVICENAME-deployment.yaml
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml

SERVICENAME=jenkins-service
print_debug "Update $SERVICENAME service"
SERVICE_REVISION=$(kubectl get revisions --namespace=keptn | grep $SERVICENAME | cut -d' ' -f1)
kubectl get ksvc $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.runLatest.configuration.revisionTemplate.spec.container.image keptn/$SERVICENAME:0.2.0
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
print_debug "Removing old revision of $SERVICENAME service"
kubectl delete revision $SERVICE_REVISION -n keptn
rm $SERVICENAME-deployment.yaml

SERVICENAME=jenkins-deployment
print_debug "Update $SERVICENAME"
kubectl get deployment $SERVICENAME -n keptn -o=yaml > $SERVICENAME-deployment.yaml
verify_kubectl $? "$SERVICENAME could not be retrieved."
yq w -i $SERVICENAME-deployment.yaml spec.template.spec.containers[0].image keptn/jenkins:0.5.0
kubectl apply -f $SERVICENAME-deployment.yaml
verify_kubectl $? "Updating of $SERVICENAME failed."
rm $SERVICENAME-deployment.yaml


GATEWAY=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

JENKINS_URL="jenkins.keptn.$GATEWAY.xip.io"

JENKINS_USER=$(kubectl get secret jenkins-secret -n keptn -o=yaml | yq - r data.user | base64 --decode)
JENKINS_PASSWORD=$(kubectl get secret jenkins-secret -n keptn -o=yaml | yq - r data.password | base64 --decode)

# Configure Jenkins with GitHub credentials
RETRY=0; RETRY_MAX=12; 
sleep 5
while [[ $RETRY -lt $RETRY_MAX ]]; do
  response=$(curl --write-out %{http_code} --silent --output /dev/null $JENKINS_URL)
  #echo $response
  if [[ $response -eq 200 ]]
  then
    break
  fi 
  echo "Jenkins not yet available. Retrying..."
  sleep 15
  RETRY=$[$RETRY+1]
done

if [[ $RETRY == $RETRY_MAX ]]; then
  print_error "Jenkins not available, thus Git credentials could not be created in Jenkins."
  exit 1
fi

RETRY=0

while [[ $RETRY -lt $RETRY_MAX ]]; do
  curl -X POST http://$JENKINS_URL/credentials/store/system/domain/_/createCredentials \
    --user $JENKINS_USER:$JENKINS_PASSWORD \
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

  if [[ $? == '0' ]]
  then
    print_debug "Git credentials in Jenkins created."
    break
  fi
  RETRY=$[$RETRY+1]
  print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for creating git credentials in Jenkins ..."
  sleep 10
done

if [[ $RETRY == $RETRY_MAX ]]; then
  print_error "Git credentials could not be created in Jenkins."
  exit 1
fi

echo "Upgrade to keptn 0.2.1 done."
echo "You can find your Jenkins here: http://jenkins.keptn.$GATEWAY.xip.io/configure"
echo "Please verify and save your Jenkins configuration."
