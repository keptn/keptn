#!/bin/bash

source test/utils.sh

function cleanup() {
  # scale the helm-service back up again
  kubectl -n keptn scale deployment.v1.apps/helm-service --replicas=1

  keptn delete project delivery-assistant-project
  kubectl delete ns delivery-assistant-project-combi1
  kubectl delete ns delivery-assistant-project-combi2
  kubectl delete ns delivery-assistant-project-combi3
  kubectl delete ns delivery-assistant-project-combi4

}
trap cleanup EXIT

# get keptn api details
KEPTN_ENDPOINT=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/api
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
check_no_open_approvals $PROJECT combi1

check_number_open_approvals $PROJECT combi2 1

combi2ApprovalId=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi2 | awk '{if(NR>1)print}' | jq -r '.[0].id')
keptn_context_id=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi2 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')
keptn send event approval.finished --id=${combi2ApprovalId} --project=delivery-assistant-project --stage=combi2
sleep 5
check_no_open_approvals $PROJECT combi2

response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.configuration.change $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

# print the response
echo "Resulting configuration.change event by approval:"
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "gatekeeper-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi2"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.canary.action" "set"
verify_using_jq "$response" ".data.canary.value" "100"
verify_using_jq "$response" ".data.valuesCanary.image" "docker.io/keptnexamples/carts:0.11.1"

check_number_open_approvals $PROJECT combi3 1

combi3ApprovalId=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}' | jq -r '.[0].id')
keptn_context_id=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')
keptn send event approval.finished --id=${combi3ApprovalId} --project=delivery-assistant-project --stage=combi3
sleep 5
check_no_open_approvals $PROJECT combi3
combi3EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}')

response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.configuration.change $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

# print the response
echo "Resulting configuration.change event by approval:"
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "gatekeeper-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi3"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.canary.action" "set"
verify_using_jq "$response" ".data.canary.value" "100"
verify_using_jq "$response" ".data.valuesCanary.image" "docker.io/keptnexamples/carts:0.11.1"

check_number_open_approvals $PROJECT combi4 2

combi4ApprovalId1=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].id')
keptn_context_id_1=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')
combi4ApprovalId2=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[1].id')
keptn_context_id_2=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')

keptn send event approval.finished --id=${combi4ApprovalId1} --project=delivery-assistant-project --stage=combi4
sleep 5
check_number_open_approvals $PROJECT combi4 1

response=$(get_keptn_event $PROJECT $keptn_context_id_1 sh.keptn.event.configuration.change $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

# print the response
echo "Resulting configuration.change event by approval:"
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "gatekeeper-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi4"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.canary.action" "set"
verify_using_jq "$response" ".data.canary.value" "100"
verify_using_jq "$response" ".data.valuesCanary.image" "docker.io/keptnexamples/carts:0.11.1"


keptn send event approval.finished --id=${combi4ApprovalId2} --project=delivery-assistant-project --stage=combi4
sleep 5
check_no_open_approvals $PROJECT combi4


response=$(get_keptn_event $PROJECT $keptn_context_id_2 sh.keptn.event.configuration.change $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

# print the response
echo "Resulting configuration.change event by approval:"
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "gatekeeper-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi4"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.canary.action" "set"
verify_using_jq "$response" ".data.canary.value" "100"
verify_using_jq "$response" ".data.valuesCanary.image" "docker.io/keptnexamples/carts:0.11.1"

