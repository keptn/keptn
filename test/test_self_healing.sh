#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n "$KEPTN_NAMESPACE" -o jsonpath='{.data.keptn-api-token}' | base64 --decode)

# test configuration
UNLEASH_SERVICE_VERSION=${UNLEASH_SERVICE_VERSION:-master}
PROJECT="self-healing-project"
SERVICE="frontend"

function print_logs {
  echo "Logs from: remediation-service"
  kubectl -n "$KEPTN_NAMESPACE" logs svc/remediation-service -c remediation-service
}

trap print_logs EXIT

########################################################################################################################
# Pre-requisites
########################################################################################################################

# ensure unleash-service is not installed yet
if kubectl -n "$KEPTN_NAMESPACE" get deployment unleash-service 2> /dev/null; then
  echo "Found unleash-service. Please uninstall it using:"
  echo "kubectl -n ${KEPTN_NAMESPACE} delete deployment unleash-service"
  exit 1
fi

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

echo "Testing self-healing for project $PROJECT ..."

echo "Creating a new project without Git upstream"
keptn create project $PROJECT --shipyard=./test/assets/self_healing_shipyard.yaml
verify_test_step $? "keptn create project ${PROJECT} - failed"
sleep 10

# verify that the project has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" != "${PROJECT}" ]]; then
  echo "Failed to check that the project exists via the API."
  echo "${response}"
  exit 2
else
  echo "Verified that project exists via API"
fi


####################################################################################################################################
# Testcase 1:
# Project exists, but service has not been onboarded yet
# Sending a remediation.triggered event now should result in message: Could not execute remediation action because service is not available
####################################################################################################################################

echo "Sending remediation.triggered event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_remediation_triggered_event.json)
sleep 15

#response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')
response=$(get_keptn_event "$PROJECT" "$keptn_context_id" sh.keptn.event.production.remediation.finished "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN")

# print the response
echo "$response" | jq .

# validate the response
verify_using_jq "$response" ".source" "shipyard-controller"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"
verify_using_jq "$response" ".data.result" "fail"


####################################################################################################################################
# Testcase 2:
# Project exists, service has been onboarded, but no remediation file could be found
# Sending a remediation.triggered event should result in message: Could not execute remediation action because no remediation file available
####################################################################################################################################

###########################################
# create service frontend
###########################################
keptn create service $SERVICE --project=$PROJECT
verify_test_step $? "keptn create service ${SERVICE} --project=${PROJECT} - failed"
sleep 10

# verify that the service has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}/stage/production/service/${SERVICE}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.serviceName')

if [[ "$response" != "${SERVICE}" ]]; then
  echo "Failed to check that the service exists via the API"
  echo "${response}"
  exit 2
else
  echo "Verified that service exists via API"
fi

echo "Sending remediation.triggered event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_remediation_triggered_event.json)
sleep 10

response=$(get_keptn_event "$PROJECT" "$keptn_context_id" sh.keptn.event.production.remediation.finished "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN")
# print the response
echo "$response" | jq .

# validate the response
verify_using_jq "$response" ".source" "shipyard-controller"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"
verify_using_jq "$response" ".data.result" "fail"


##########################################################################################################################################
# Testcase 3:
# Project exists, service has been onboarded, remediation file available, but no service executor available
# Sending a remediation.triggered event should case an action.triggered event to be sent
##########################################################################################################################################

echo "Uploading remediation.yaml to $PROJECT/production/$SERVICE"
keptn add-resource --project=$PROJECT --service=$SERVICE --stage=production --resource=./test/assets/self_healing_remediation.yaml --resourceUri=remediation.yaml

echo "Sending remediation.triggered event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_remediation_triggered_event.json)
sleep 10

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.production.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

if [[ "$response" != "0" ]]; then
  echo "Received unexpected remediation.finished event"
  echo "${response}"
  exit 2
else
  echo "Verified that no remediation.finished event has been sent"
fi

response=$(get_keptn_event "$PROJECT" "$keptn_context_id" sh.keptn.event.action.triggered "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN")

# print the response
echo "$response" | jq .

# validate the response
verify_using_jq "$response" ".source" "shipyard-controller"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.action.name" "toggle-feature"
verify_using_jq "$response" ".data.action.value.EnablePromotion" "off"

##########################################################################################################################################
# Testcase 3:
# Project exists, service has been onboarded, remediation file available, action executor is available, but will fail
# Sending a problem.open event now should result in message: Action toggle-feature triggered but not executed after waiting for 2 minutes.
##########################################################################################################################################

# install unleash service
echo "Installing unleash-service version ${UNLEASH_SERVICE_VERSION}"
kubectl apply -f "https://raw.githubusercontent.com/keptn-contrib/unleash-service/${UNLEASH_SERVICE_VERSION}/deploy/service.yaml" -n "${KEPTN_NAMESPACE}"

sleep 10

wait_for_deployment_in_namespace "unleash-service" "${KEPTN_NAMESPACE}"

kubectl get deployment -n "$KEPTN_NAMESPACE" unleash-service -oyaml

echo "Sending remediation.triggered event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_remediation_triggered_event.json)

sleep 10

response=$(get_keptn_event "$PROJECT" "$keptn_context_id" sh.keptn.event.production.remediation.finished "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN")
# print the response
echo "$response" | jq .

# validate the response
verify_using_jq "$response" ".source" "shipyard-controller"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo "$response" | jq .

# validate the response
verify_using_jq "$response" ".source" "unleash-service"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"

echo "Remediation Service logs:"
kubectl logs -n "$KEPTN_NAMESPACE" svc/remediation-service -c remediation-service

echo "Unleash service logs:"
kubectl logs -n "$KEPTN_NAMESPACE" svc/unleash-service -c unleash-service

echo "Unleash service distributor logs:"
kubectl logs -n "$KEPTN_NAMESPACE" svc/unleash-service -c distributor

echo "Self healing tests done ✓"
