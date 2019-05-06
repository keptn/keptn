#!/bin/bash
source ./utils.sh

GITHUB_USER=$(cat ../install/scripts/creds.json | jq -r '.githubUserName')
GITHUB_ORG=$(cat ../install/scripts/creds.json | jq -r '.githubOrg')
GITHUB_TOKEN=$(cat ../install/scripts/creds.json | jq -r '.githubPersonalAccessToken')

KEPTN_ENDPOINT=https://$(kubectl get ksvc -n keptn control -o=yaml | yq r - status.domain)
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode)

PROJECT=sockshop

# Delete old project
git ls-remote https://github.com/$GITHUB_ORG_NIGHTLY/$PROJECT > /dev/null 2>&1
if [ $? = 0 ]; then 
    echo "Delete project $PROJECT" 
    hub delete -y $GITHUB_ORG_NIGHTLY/$PROJECT
    echo "Finished deleting project $PROJECT"
else 
    echo "No project to delete"
fi

# Authenticate keptn CLI
keptn auth --endpoint=$KEPTN_ENDPOINT --api-token=$KEPTN_API_TOKEN
verify_test_step $? "Could not authenticate at keptn API."

keptn configure --org=$GITHUB_ORG --user=$GITHUB_USER --token=$GITHUB_TOKEN
verify_test_step $? "keptn config command failed."

# Test keptn config result
RETRY=0; RETRY_MAX=12;
while [[ $RETRY -lt $RETRY_MAX ]]; do
  sleep 10
  STORED_GITHUB_USER=$(kubectl get secret github-credentials -n keptn -o=yaml | yq - r data.user | base64 --decode)

  if [ "$STORED_GITHUB_USER" == "$GITHUB_USER" ]; then
      echo "Keptn config succeeded."
      break
  fi
  RETRY=$[$RETRY+1]
  echo "Expected value user=$GITHUB_USER not yet stored in cluster. Actual value is $STORED_GITHUB_USER. Trying again in 10 seconds."
  sleep 10
  if [ $RETRY -eq $RETRY_MAX ]; then
    echo "keptn config failed."
  fi
done

# Test keptn create-project and onboard
rm -rf examples
git clone https://github.com/keptn/examples
cd examples
cd onboarding-carts

keptn create project $PROJECT shipyard.yaml
verify_test_step $? "keptn create project command failed."

sleep 30
keptn onboard service --project=$PROJECT --values=values_carts.yaml

sleep 30
keptn onboard service --project=$PROJECT --values=values_carts_db.yaml --deployment=deployment_carts_db.yaml --service=service_carts_db.yaml

sleep 60
cd ../..
npm install newman
yq w keptn.postman_environment.json values[0].value $GITHUB_ORG | yq  - w values[1].value $PROJECT | yq - r -j > keptn.postman_environment_tmp.json
rm keptn.postman_environment.json
mv keptn.postman_environment_tmp.json keptn.postman_environment.json
node_modules/.bin/newman run keptn.postman_collection.json -e keptn.postman_environment.json
