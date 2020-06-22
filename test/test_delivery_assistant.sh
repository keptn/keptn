#!/bin/bash

source test/utils.sh

# get keptn api details
KEPTN_ENDPOINT=https://api.keptn.$(kubectl get cm keptn-domain -n keptn -ojsonpath={.data.app_domain})
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

# test configuration
PROJECT="delivery-assistant-project"
SERVICE="carts"

########################################################################################################################
# Pre-requesits
########################################################################################################################

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

echo "Testing delivery assistant for project $PROJECT ..."

echo "Creating a new project without git upstream"
keptn create project $PROJECT --shipyard=./test/assets/delivery_assistant_shipyard.yaml
verify_test_step $? "keptn create project command failed."
sleep 10

# verify that the project has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" != "${PROJECT}" ]]; then
  echo "Failed to check that the project exists via the API."
  echo "${response}"
  exit 2
else
  echo "Verified that Project exists via api"
fi

###########################################
# create service frontend                #
###########################################

rm -rf examples
git clone --branch master https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

keptn onboard service carts --project=$PROJECT --chart=./carts
verify_test_step $? "keptn onboard service ${SERVICE} failed."
sleep 10

cd ../..

# scale down the helm service to avoid deployments
kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=0

# Send 3 evaluation-done events (result: pass, warning, failed) for each stage (dev, combi1, combi2, combi3) using the CLI

send_evaluation_done_event $PROJECT dev $SERVICE pass
send_evaluation_done_event $PROJECT dev $SERVICE warning
send_evaluation_done_event $PROJECT dev $SERVICE failed

send_evaluation_done_event $PROJECT combi1 $SERVICE pass
send_evaluation_done_event $PROJECT combi1 $SERVICE warning
send_evaluation_done_event $PROJECT combi1 $SERVICE failed

send_evaluation_done_event $PROJECT combi2 $SERVICE pass
send_evaluation_done_event $PROJECT combi2 $SERVICE warning
send_evaluation_done_event $PROJECT combi2 $SERVICE failed

send_evaluation_done_event $PROJECT combi3 $SERVICE pass
send_evaluation_done_event $PROJECT combi3 $SERVICE warning
send_evaluation_done_event $PROJECT combi3 $SERVICE failed

send_evaluation_done_event $PROJECT combi4 $SERVICE pass
send_evaluation_done_event $PROJECT combi4 $SERVICE warning
send_evaluation_done_event $PROJECT combi4 $SERVICE failed


# verify the number of open approval events

combi1EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi1 | awk '{if(NR>1)print}' | jq -r '.' | jq -r --slurp 'length')
if [[ "$combi1EventLength" != "0" ]]; then
  echo "Received number of approval.triggered events"
  echo "${response}"
  # scale the helm-service back up again
  kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1
  exit 2
else
  echo "Verified number of approval.triggered events"
fi

combi2EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi2 | awk '{if(NR>1)print}' | jq -r '.' | jq -r --slurp 'length')
if [[ "$combi2EventLength" != "1" ]]; then
  echo "Received number of approval.triggered events"
  echo "${response}"
  # scale the helm-service back up again
  kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1
  exit 2
else
  echo "Verified number of approval.triggered events"
fi

combi3EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}' | jq -r '.' | jq -r --slurp 'length')
if [[ "$combi3EventLength" != "1" ]]; then
  echo "Received number of approval.triggered events"
  echo "${response}"
  # scale the helm-service back up again
  kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1
  exit 2
else
  echo "Verified number of approval.triggered events"
fi

combi4EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.' | jq -r --slurp 'length')
if [[ "$combi4EventLength" != "2" ]]; then
  echo "Received number of approval.triggered events"
  echo "${response}"
  # scale the helm-service back up again
  kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1
  exit 2
else
  echo "Verified number of approval.triggered events"
fi




# scale the helm-service back up again
kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1


