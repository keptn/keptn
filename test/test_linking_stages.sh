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

echo "KEPTN_ENDPOINT $KEPTN_ENDPOINT"

#test configuration
PROJECT="link-project"
SERVICE="my-service"


# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

echo "Installing keptn-sandbox/echo-service"
kubectl -n "${KEPTN_NAMESPACE}" apply -f https://raw.githubusercontent.com/keptn-sandbox/echo-service/release-0.1.0/deploy/service.yaml


echo "Testing link staging..."

echo "Creating a new project without Git upstream"
keptn create project ${PROJECT} --shipyard=./test/assets/linking_stages_shipyard.yaml
verify_test_step $? "keptn create project ${PROJECT} failed."
sleep 2

echo "Creating a new service"
keptn create service ${SERVICE} --project ${PROJECT}
verify_test_step $? "keptn create service ${SERVICE} failed"
sleep 2


####################################################################################################################################
# Testcase:
# 1) sh.keptn.event.firststage.echosequence.triggered
# 2) check if the stages in the provided shipyard file gets started by shipyard controller
####################################################################################################################################

echo "Sending trigger echosequence event"
keptn_context_id=$(send_event_json ./test/assets/trigger_echosequence_event.json)
sleep 20

declare -a list_of_events=("sh.keptn.event.firststage.echosequence.triggered" "sh.keptn.event.firststage.echosequence.finished" "sh.keptn.event.secondstage.echosequence.triggered" "sh.keptn.event.secondstage.echosequence.finished" "sh.keptn.event.thirdstage.echosequence.triggered" "sh.keptn.event.thirdstage.echosequence.finished")
for e in "${list_of_events[@]}"; do
  echo "Verifying that event $e was sent"
  verify_event_not_null "$(get_keptn_event "$PROJECT" "$keptn_context_id" "$e" "$KEPTN_ENDPOINT" "$KEPTN_API_TOKEN")"
  if [ "$?" -eq "-1" ];then
    print_error "Event $e could not be fetched. Exiting test..."
    exit 2
  fi
done
