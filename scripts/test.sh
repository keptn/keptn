export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_TENANT_URL=$(echo $DT_TENANT_ID | sed 's,/,///,g')
echo $DT_TENANT_URL
cat ../manifests/jenkins/k8s-jenkins-deployment-defaults.yml | \
  sed 's~GITHUB_USER_EMAIL_PLACEHOLDER~'"$GITHUB_USER_EMAIL"'~' | \
  sed 's~GITHUB_ORGANIZATION_PLACEHOLDER~'"$GITHUB_ORGANIZATION"'~' | \
  sed 's~DOCKER_REGISTRY_IP_PLACEHOLDER~'"$REGISTRY_URL"'~' | \
  sed 's~DT_TENANT_URL_PLACEHOLDER~'"$DT_TENANT_URL"'~' | \
  sed 's~DT_API_TOKEN_PLACEHOLDER~'"$DT_API_TOKEN"'~' >> ../manifests/jenkins/k8s-jenkins-deployment_tmp.yml
