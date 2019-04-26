#!/bin/bash
REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

# Environment variables jenkins-service
if [[ -z "${JENKINS_USER}" ]]; then
  echo "[keptn|DEBUG] JENKINS_USER not set, take it from creds.json"
  JENKINS_USER=$(cat creds.json | jq -r '.jenkinsUser')
fi

if [[ -z "${JENKINS_PASSWORD}" ]]; then
  echo "[keptn|DEBUG] JENKINS_PASSWORD not set, take it from creds.json"
  JENKINS_PASSWORD=$(cat creds.json | jq -r '.jenkinsPassword')
fi

if [[ -z "${GITHUB_USER_NAME}" ]]; then
  echo "[keptn|DEBUG] GITHUB_USER_NAME not set, take it from creds.json"
  GITHUB_USER_NAME=$(cat creds.json | jq -r '.githubUserName')
fi

if [[ -z "${GITHUB_PERSONAL_ACCESS_TOKEN}" ]]; then
  echo "[keptn|DEBUG] GITHUB_PERSONAL_ACCESS_TOKEN not set, take it from creds.json"
  GITHUB_PERSONAL_ACCESS_TOKEN=$(cat creds.json | jq -r '.githubPersonalAccessToken')
fi

if [[ -z "${GITHUB_USER_EMAIL}" ]]; then
  echo "[keptn|DEBUG] GITHUB_USER_EMAIL not set, take it from creds.json"
  GITHUB_USER_EMAIL=$(cat creds.json | jq -r '.githubUserEmail')
fi

if [[ -z "${GITHUB_ORGANIZATION}" ]]; then
  echo "[keptn|DEBUG] GITHUB_ORGANIZATION not set, take it from creds.json"
  GITHUB_ORGANIZATION=$(cat creds.json | jq -r '.githubOrg')
fi

if [[ -z "${DT_TENANT_ID}" ]]; then
  echo "[keptn|DEBUG] DT_TENANT_ID not set, use value from creds_dt.json"
  DT_TENANT_ID=$(cat creds_dt.json | jq -r '.dynatraceTenant')
fi

if [[ -z "${DT_API_TOKEN}" ]]; then
  echo "[keptn|DEBUG] DT_API_TOKEN not set, use value from creds_dt.json"
  DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')
fi

if [[ -z "${DT_TENANT_URL}" ]]; then
  echo "[keptn|DEBUG] DT_TENANT_URL not set, define it based on DT_TENANT_ID"
  DT_TENANT_URL="$DT_TENANT_ID.live.dynatrace.com"
fi

# Install services

rm -rf keptn-services
mkdir keptn-services
cd keptn-services

git clone --branch 0.1.0 https://github.com/keptn/jenkins-service.git --single-branch
cd jenkins-service
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $JENKINS_USER $JENKINS_PASSWORD $GITHUB_USER_NAME $GITHUB_USER_EMAIL $GITHUB_ORGANIZATION $GITHUB_PERSONAL_ACCESS_TOKEN $DT_API_TOKEN $DT_TENANT_URL
cd ..

git clone --branch 0.1.0 https://github.com/keptn/github-service.git --single-branch
cd github-service
chmod +x deploy.sh
./deploy.sh
cd ..

git clone --branch 0.1.0 https://github.com/keptn/servicenow-service.git --single-branch
cd servicenow-service
chmod +x deploy.sh
./deploy.sh
cd ..

git clone --branch 0.1.0 https://github.com/keptn/pitometer-service.git --single-branch
cd pitometer-service
chmod +x deploy.sh
./deploy.sh
cd ..
