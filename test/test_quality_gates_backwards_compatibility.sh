#!/bin/bash

 source test/utils.sh

# test configuration
KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)

########################################################################################################################
# Testcase 1: backwards compatibility for 0.7.x evaluation-done events
########################################################################################################################

# feed a legacy evaluation-done event to Keptn
echo "Testing backwards compatibility with Kepn 0.7.x evaluation-done events"

echo "Sending legacy evaluation-done event"
legacy_event_context_id=$(send_event_json ./test/assets/07x_evaluation_done_event.json)
sleep 5

echo "Trying to retrieve event using the sh.keptn.event.evaluation.finished type filter"
# check if the event is returned when checking for sh.keptn.event.evluation.finished
response=$(get_event sh.keptn.event.evaluation.finished ${legacy_event_context_id} "legacy-project")

# print the response
echo $response | jq .

verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".specversion" "1.0"
verify_using_jq "$response" ".type" "sh.keptn.event.evaluation.finished"
verify_using_jq "$response" ".data.project" "legacy-project"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "legacy-service"
verify_using_jq "$response" ".data.evaluation.result" "fail"


echo "Retrieving the event directly via the API endpoint"
number_of_events=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter=shkeptncontext:${legacy_event_context_id}%20AND%20AND%20AND%20data.project:legacy-project&excludeInvalidated=true" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

verify_value "number of services" $number_of_events 1

echo "Sending evaluation.invalidated event"
send_evaluation_invalidated_event "legacy-project" "hardening" "legacy-service" "evaluation-done-id" $legacy_event_context_id

echo "Retrieving the event directly via the API endpoint - now it should be excluded"
number_of_events=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished?filter='shkeptncontext:${legacy_event_context_id}%20AND%20AND%20AND%20data.project:legacy-project'&excludeInvalidated=true" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events | length')

verify_value "number of events" $number_of_events 0

exit 0
