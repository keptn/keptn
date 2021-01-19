#!/bin/bash

source test/utils.sh

# get keptn API details
if [[ "$PLATFORM" == "openshift" ]]; then
  KEPTN_ENDPOINT=http://api.${KEPTN_NAMESPACE}.127.0.0.1.nip.io/api
else
  if [[ "$KEPTN_SERVICE_TYPE" == "NodePort" ]]; then
    API_PORT=$(kubectl get svc api-gateway-nginx -n ${KEPTN_NAMESPACE} -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
    INTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
    KEPTN_ENDPOINT="http://${INTERNAL_NODE_IP}:${API_PORT}"/api
  else
    KEPTN_ENDPOINT=http://$(kubectl -n ${KEPTN_NAMESPACE} get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/api
  fi
fi

KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)

# test configuration
UNLEASH_SERVICE_VERSION=${UNLEASH_SERVICE_VERSION:-master}
PROJECT="self-healing-project"
SERVICE="frontend"

########################################################################################################################
# Pre-requisites
########################################################################################################################

# ensure unleash-service is not installed yet
kubectl -n ${KEPTN_NAMESPACE} get deployment unleash-service 2> /dev/null

if [[ $? -eq 0 ]]; then
  echo "Found unleash-service. Please uninstall it using:"
  echo "kubectl -n ${KEPTN_NAMESPACE} delete deployment unleash-service"
  exit 1
fi

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

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
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

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
# Sending a problem.open event now should result in message: Could not execute remediation action because service is not available
####################################################################################################################################

echo "Sending problem.open event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_problem_open_event.json)
sleep 10

#response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')
response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.remediation.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "remediation-service"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"
verify_using_jq "$response" ".data.result" "fail"


####################################################################################################################################
# Testcase 2:
# Project exists, service has been onboarded, but no remediation file could be found
# Sending a problem.open event should result in message: Could not execute remediation action because no remediation file available
####################################################################################################################################

###########################################
# create service frontend                
###########################################
keptn create service $SERVICE --project=$PROJECT
verify_test_step $? "keptn create service ${SERVICE} --project=${PROJECT} - failed"
sleep 10

# verify that the service has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}/stage/production/service/${SERVICE}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.serviceName')

if [[ "$response" != "${SERVICE}" ]]; then
  echo "Failed to check that the service exists via the API"
  echo "${response}"
  exit 2
else
  echo "Verified that service exists via API"
fi

echo "Sending problem.open event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_problem_open_event.json)
sleep 10

#response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')
response=$(get_keptn_event $PROJECT $keptn_context_id sh.keptn.event.remediation.finished $KEPTN_ENDPOINT $KEPTN_API_TOKEN)
# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "remediation-service"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.status" "errored"
verify_using_jq "$response" ".data.result" "fail"
verify_using_jq "$response" ".data.message" "Remediation disabled for service $SERVICE in project $PROJECT in stage production"


##########################################################################################################################################
# Testcase 3:
# Project exists, service has been onboarded, remediation file available, but no service executor available
# Sending a problem.open event now should result in message: Action toogle-feature triggered but not executed after waiting for 2 minutes.
##########################################################################################################################################

echo "Uploading remediation.yaml to $PROJECT/production/$SERVICE"
keptn add-resource --project=$PROJECT --service=$SERVICE --stage=production --resource=./test/assets/self_healing_remediation.yaml --resourceUri=remediation.yaml

echo "Sending problem.open event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_problem_open_event.json)
sleep 10

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

if [[ "$response" != "0" ]]; then
  echo "Received unexpected remediation.finished event"
  echo "${response}"
  exit 2
else
  echo "Verified that no remediation.finished event has been sent"
fi

sleep 60

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

if [[ "$response" != "0" ]]; then
  echo "Received unexpected remediation.finished event"
  echo "${response}"
  exit 2
else
  echo "Verified that no remediation.finished event has been sent"
fi
# TODO: we need a timeout mechanism for actions in the remediation service
#response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')
#
## print the response
#echo $response | jq .
#
## validate the response
#verify_using_jq "$response" ".source" "remediation-service"
#verify_using_jq "$response" ".data.project" "self-healing-project"
#verify_using_jq "$response" ".data.stage" "production"
#verify_using_jq "$response" ".data.service" "$SERVICE"
#verify_using_jq "$response" ".data.remediation.status" "errored"
#verify_using_jq "$response" ".data.remediation.result" "failed"
#verify_using_jq "$response" ".data.remediation.message" "Action toogle-feature triggered but not executed after waiting for 2 minutes."


##########################################################################################################################################
# Testcase 3:
# Project exists, service has been onboarded, remediation file available, first action executor is available, but not the second
# Sending a problem.open event now should result in message: Action toogle-feature triggered but not executed after waiting for 2 minutes.
##########################################################################################################################################

# install unleash service
echo "Installing unleash-service version ${UNLEASH_SERVICE_VERSION}"
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/unleash-service/${UNLEASH_SERVICE_VERSION}/deploy/service.yaml -n ${KEPTN_NAMESPACE}

kubectl -n ${KEPTN_NAMESPACE} set image deployment/unleash-service unleash-service=keptncontrib/unleash-service:0.0.0-master

sleep 10

wait_for_deployment_in_namespace "unleash-service" "${KEPTN_NAMESPACE}"

kubectl get deployment -n ${KEPTN_NAMESPACE} unleash-service -oyaml

echo "Sending problem.open event"
keptn_context_id=$(send_event_json ./test/assets/self_healing_problem_open_event.json)

sleep 10

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

if [[ "$response" != "0" ]]; then
  echo "Received unexpected remediation.finished event"
  echo "${response}"
  exit 2
else
  echo "Verified that no remediation.finished event has been sent"
fi

sleep 60

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

if [[ "$response" != "0" ]]; then
  echo "Received unexpected remediation.finished event"
  echo "${response}"
  exit 2
else
  echo "Verified that no remediation.finished event has been sent"
fi

# TODO: we need a timeout mechanism for actions in the remediation service
#response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')
#
## print the response
#echo $response | jq .
#
## validate the response
#verify_using_jq "$response" ".source" "remediation-service"
#verify_using_jq "$response" ".data.project" "self-healing-project"
#verify_using_jq "$response" ".data.stage" "production"
#verify_using_jq "$response" ".data.service" "$SERVICE"
#verify_using_jq "$response" ".data.remediation.status" "errored"
#verify_using_jq "$response" ".data.remediation.result" "failed"
#verify_using_jq "$response" ".data.remediation.message" "Action run-snow-wf triggered but not executed after waiting for 2 minutes."

response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

echo "Remediation Service logs:"
kubectl logs -n ${KEPTN_NAMESPACE} svc/remediation-service -c remediation-service

echo "Unleash service logs:"
kubectl logs -n ${KEPTN_NAMESPACE} svc/unleash-service -c unleash-service

echo "Unleash service distributor logs:"
kubectl logs -n ${KEPTN_NAMESPACE} svc/unleash-service -c distributor

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "unleash-service"
verify_using_jq "$response" ".data.project" "self-healing-project"
verify_using_jq "$response" ".data.stage" "production"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.action.status" "errored"
# TODO: we need a message field for that
# verify_using_jq "$response" ".data.action.message" "Action run-snow-wf triggered but not executed after waiting for 2 minutes."

echo "Self healing tests done ✓"
