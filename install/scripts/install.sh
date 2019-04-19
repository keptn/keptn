#!/bin/bash

kubectl apply -f ./rbac.yaml

# Update installer.yaml
CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')

CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
GCLOUD_USER=$(gcloud config get-value account)

# For uniform: (soon obsolete)
JENKINS_USER=$(cat creds.json | jq -r '.jenkinsUser')
JENKINS_PASSWORD=$(cat creds.json | jq -r '.jenkinsPassword')
GITHUB_PERSONAL_ACCESS_TOKEN=$(cat creds.json | jq -r '.githubPersonalAccessToken')
GITHUB_USER_EMAIL=$(cat creds.json | jq -r '.githubUserEmail')
GITHUB_ORGANIZATION=$(cat creds.json | jq -r '.githubOrg')
DT_TENANT_ID=$(cat creds_dt.json | jq -r '.dynatraceTenant')
DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')

cat ./installer.yaml | \
  sed 's~value: CLUSTER_IPV4_CIDR~'"value: $CLUSTER_IPV4_CIDR"'~' | \
  sed 's~value: SERVICES_IPV4_CIDR~'"value: $SERVICES_IPV4_CIDR"'~' | \
  sed 's~value: GCLOUD_USER~'"value: $GCLOUD_USER"'~' | \
  sed 's~value: JENKINS_USER~'"value: $JENKINS_USER"'~' | \
  sed 's~value: JENKINS_PASSWORD~'"value: $JENKINS_PASSWORD"'~' | \
  sed 's~value: GITHUB_PERSONAL_ACCESS_TOKEN~'"value: $GITHUB_PERSONAL_ACCESS_TOKEN"'~' | \
  sed 's~value: GITHUB_USER_EMAIL~'"value: $GITHUB_USER_EMAIL"'~' | \
  sed 's~value: GITHUB_ORGANIZATION~'"value: $GITHUB_ORGANIZATION"'~' >> ./installer-gen.yaml

# Roll-out installer on cluster
kubectl apply -f ./installer-gen.yaml
