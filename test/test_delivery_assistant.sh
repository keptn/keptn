#!/bin/bash

source test/utils.sh

function cleanup() {
  keptn delete project delivery-assistant-project
}
trap cleanup EXIT

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)

# test configuration
PROJECT="delivery-assistant-project"
SERVICE="carts"

########################################################################################################################
# Pre-requesites
########################################################################################################################

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

echo "Testing delivery assistant for project $PROJECT ..."

echo "Creating a new project without Git upstream"
keptn create project $PROJECT --shipyard=./test/assets/delivery_assistant_shipyard.yaml
verify_test_step $? "keptn create project ${PROJECT} - failed"
sleep 10

# verify that the project has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" != "${PROJECT}" ]]; then
  echo "Failed to check that the project exists via the API"
  echo "${response}"
  exit 2
else
  echo "Verified that project exists via API"
fi

###########################################
# create service frontend                #
###########################################

keptn create service $SERVICE --project=$PROJECT
verify_test_step $? "keptn create service ${SERVICE} - failed"

# Send 3 approval.triggered events (result: pass, warning, failed) for each stage (dev, combi1, combi2, combi3) using the CLI

send_approval_triggered_event $PROJECT combi1 $SERVICE pass
send_approval_triggered_event $PROJECT combi1 $SERVICE warning
send_approval_triggered_event $PROJECT combi1 $SERVICE fail

send_approval_triggered_event $PROJECT combi2 $SERVICE pass
send_approval_triggered_event $PROJECT combi2 $SERVICE warning
send_approval_triggered_event $PROJECT combi2 $SERVICE fail

send_approval_triggered_event $PROJECT combi3 $SERVICE pass
send_approval_triggered_event $PROJECT combi3 $SERVICE warning
send_approval_triggered_event $PROJECT combi3 $SERVICE fail

send_approval_triggered_event $PROJECT combi4 $SERVICE pass
send_approval_triggered_event $PROJECT combi4 $SERVICE warning
send_approval_triggered_event $PROJECT combi4 $SERVICE fail


sleep 20
# verify the number of open approval events
check_no_open_approvals $PROJECT combi1

keptn get event approval.triggered --project=$PROJECT --stage=combi2

check_number_open_approvals $PROJECT combi2 1

combi2ApprovalId=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi2 | awk '{if(NR>1)print}' | jq -r '.id')
keptn_context_id=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi2 | awk '{if(NR>1)print}' | jq -r '.shkeptncontext')

echo $combi2ApprovalId
keptn send event approval.finished --id=${combi2ApprovalId} --project=delivery-assistant-project --stage=combi2
sleep 10
check_no_open_approvals $PROJECT combi2

# print the response
echo "Resulting approval.finished event by approval:"

response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.combi2.approval.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
echo $response | jq .

# validate the response
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi2"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.result" "pass"

check_number_open_approvals $PROJECT combi3 1

combi3ApprovalId=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}' | jq -r '.id')
keptn_context_id=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}' | jq -r '.shkeptncontext')
keptn send event approval.finished --id=${combi3ApprovalId} --project=delivery-assistant-project --stage=combi3
sleep 5
check_no_open_approvals $PROJECT combi3
combi3EventLength=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi3 | awk '{if(NR>1)print}')

# print the response
echo "Resulting approval.finished event by approval:"

response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.combi3.approval.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
echo $response | jq .

# validate the response
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi3"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.result" "pass"

check_number_open_approvals $PROJECT combi4 2

combi4ApprovalId1=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].id')
keptn_context_id_1=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')
combi4ApprovalId2=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[1].id')
keptn_context_id_2=$(keptn get event approval.triggered --project=delivery-assistant-project --stage=combi4 | awk '{if(NR>1)print}' | jq -r '.[0].shkeptncontext')

keptn send event approval.finished --id=${combi4ApprovalId1} --project=delivery-assistant-project --stage=combi4
sleep 5
check_number_open_approvals $PROJECT combi4 1

# print the response
echo "Resulting approval.finished event by approval:"

response=$(get_keptn_event $PROJECT $keptn_context_id_1 sh.keptn.event.combi4.approval.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
echo $response | jq .

# validate the response
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi4"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.result" "pass"

keptn send event approval.finished --id=${combi4ApprovalId2} --project=delivery-assistant-project --stage=combi4
sleep 5
check_no_open_approvals $PROJECT combi4

# print the response
echo "Resulting approval.finished event by approval:"
response=$(get_keptn_event $PROJECT $keptn_context_id_2 sh.keptn.event.combi4.approval.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

echo $response | jq .

# validate the response
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "combi4"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.result" "pass"
