#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

function cleanup() {
  sleep 2
  echo "Deleting Project ${PROJECT}"
  keptn delete project "${PROJECT}"

  echo "Deleting echo-service deployment"
  kubectl delete deployments/echo-service -n "${KEPTN_NAMESPACE}"

  echo "Deleting echo-service service2"
  kubectl delete services/echo-service -n "${KEPTN_NAMESPACE}"

  echo "<END>"
  return 0
}

trap cleanup EXIT

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

# get keptn API details
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n "${KEPTN_NAMESPACE}" -o jsonpath='{.data.keptn-api-token}' | base64 --decode)

#test configuration
PROJECT="delayed-task-project"
SERVICE="my-service"

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

echo "Installing keptn-sandbox/echo-service"
kubectl -n "${KEPTN_NAMESPACE}" apply -f https://raw.githubusercontent.com/keptn-sandbox/echo-service/release-0.1.1/deploy/service.yaml

echo "Testing delayed task execution..."

echo "Creating a new project without Git upstream"
keptn create project ${PROJECT} --shipyard=./test/assets/shipyard_delayed_tasks.yaml
verify_test_step $? "keptn create project ${PROJECT} failed."
sleep 2

echo "Creating a new service"
keptn create service ${SERVICE} --project ${PROJECT}
verify_test_step $? "keptn create service ${SERVICE} failed"
sleep 2

echo "Sending trigger echosequence event"
keptn_context_id=$(send_event_json ./test/assets/trigger_echosequence_event_2.json)

echo "Waiting for sequence triggered event"
sequenceTriggeredEvent=$(get_event_with_retry sh.keptn.event.firststage.echosequence.triggered "${keptn_context_id}" "${PROJECT}")
sequenceTriggeredEventTime=$(echo "$sequenceTriggeredEvent" | jq -r '.time')
sequenceTriggeredEventTimeUNIX=$(date -d "$sequenceTriggeredEventTime" +%s)

echo "Waiting for task triggered event"
taskTriggeredEvent=$(get_event_with_retry sh.keptn.event.echo.triggered "${keptn_context_id}" "${PROJECT}")
taskTriggeredEventTime=$(echo "$taskTriggeredEvent" | jq -r '.time')
taskTriggeredEventTimeUNIX=$(date -d "$taskTriggeredEventTime" +%s)

diff="$(($taskTriggeredEventTimeUNIX-$sequenceTriggeredEventTimeUNIX))"

if (( $diff < 20 )); then
  print_error "Test failed. Task was triggered too early"
  exit 2
fi
