#!/bin/bash

source test/utils.sh

# test configuration
DYNATRACE_SLI_SERVICE_VERSION=${DYNATRACE_SLI_SERVICE_VERSION:-master}
KEPTN_EXAMPLES_BRANCH=${KEPTN_EXAMPLES_BRANCH:-master}
PROJECT=${PROJECT:-easytravel}
KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

# get keptn API details
if [[ "$PLATFORM" == "openshift" ]]; then
  KEPTN_ENDPOINT=http://api.${KEPTN_NAMESPACE}.127.0.0.1.nip.io/api
  KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)
else
  API_PORT=$(kubectl get svc api-gateway-nginx -n ${KEPTN_NAMESPACE} -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
  INTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
  KEPTN_ENDPOINT="http://${INTERNAL_NODE_IP}:${API_PORT}"/api
  KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)
fi

########################################################################################################################
# Pre-requisites
########################################################################################################################

# ensure dynatrace-sli-service is not installed yet
kubectl -n ${KEPTN_NAMESPACE} get deployment dynatrace-sli-service 2> /dev/null

if [[ $? -eq 0 ]]; then
  echo "Found dynatrace-sli-service. Please uninstall it using:"
  echo "kubectl -n ${KEPTN_NAMESPACE} delete deployment dynatrace-sli-service"
  exit 1
fi

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

# verify that the lighthouse configmap for the project does not exist yet
kubectl -n ${KEPTN_NAMESPACE} get cm lighthouse-config-${PROJECT} 2> /dev/null

if [[ $? -eq 0 ]]; then
  echo "Found configmap lighthouse-config-${PROJECT}. Please remove it using:"
  echo "kubectl -n ${KEPTN_NAMESPACE} delete configmap lighthouse-config-${PROJECT}"
  exit 1
fi

# verify that the Dynatrace credential secret does not exist yet
kubectl -n ${KEPTN_NAMESPACE} get secret dynatrace-credentials-${PROJECT} 2> /dev/null

if [[ $? -eq 0 ]]; then
  echo "Found secret dynatrace-credentials-${PROJECT}. Please remove it using:"
  echo "kubectl -n ${KEPTN_NAMESPACE} delete secret dynatrace-credentials-${PROJECT}"
  exit 1
fi

echo "Testing quality gates standalone for project $PROJECT"

# test keptn create project and create service
rm -rf examples
git clone --branch ${KEPTN_EXAMPLES_BRANCH} https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

echo "Creating a new project without Git upstream"
keptn create project $PROJECT --shipyard=./shipyard-quality-gates.yaml
verify_test_step $? "keptn create project {$PROJECT} - failed"
sleep 10

cd ../..

# verify that the project has been created via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" != "${PROJECT}" ]]; then
  echo "Failed to check that the project exists via the API"
  echo "${response}"
  exit 2
else
  echo "Verified that project exists via the API"
fi

###########################################
# create service frontend                 #
###########################################
SERVICE=frontend
keptn create service $SERVICE --project=$PROJECT
verify_test_step $? "keptn create service ${SERVICE} - failed"
sleep 10

########################################################################################################################
# Testcase 0.a: Send a start-evaluation event for a service that does not exist
########################################################################################################################

echo "Sending start-evaluation event for service 'wrong-service' in stage hardening"

response=$(send_start_evaluation_request $PROJECT wrong-stage wrong-service)

# check if the error response tells us that the service does not exist
if [[ $response != *"Service not found"* ]]; then
  echo "Did not receive expected response from Keptn API"
  exit 1
fi


########################################################################################################################
# Testcase 0.b: Send a start-evaluation event for a stage that does not exist
########################################################################################################################

echo "Sending start-evaluation event for service 'wrong-service' in stage 'wrong-stage'"

response=$(send_start_evaluation_request $PROJECT wrong-stage wrong-service)

# check if the error response tells us that the stage does not exist
if [[ $response != *"Stage not found"* ]]; then
  echo "Did not receive expected response from Keptn API"
  exit 1
fi

########################################################################################################################
# Testcase 0.c: Send a start-evaluation event for a project that does not exist
########################################################################################################################

echo "Sending start-evaluation event for service 'wrong-service' in stage 'wrong-service' in project 'wrong-project'"

response=$(send_start_evaluation_request $PROJECT wrong-stage wrong-service)

# check if the error response tells us that the project does not exist
if [[ $response != *"Project not found"* ]]; then
  echo "Did not receive expected response from Keptn API"
  exit 1
fi

########################################################################################################################
# Testcase 1:
# Project and service should have been created, but no SLO file available and no SLI provider configured
# Sending a start-evaluation event now should pass with an appropriate message
########################################################################################################################

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 10

# try to fetch a evaluation-done event
echo "Getting evaluation-done event with context-id: ${keptn_context_id}"
response=$(get_evaluation_done_event ${keptn_context_id})

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"
verify_using_jq "$response" ".data.evaluationdetails.result" "no evaluation performed by lighthouse because no SLI-provider configured for project ${PROJECT}"
verify_using_jq "$response" ".data.evaluationdetails.score" "0"
verify_using_jq "$response" ".data.evaluationdetails.sloFileContent" ""

########################################################################################################################
# Testcase 2:
# Project and service should have been created, default SLI provider available, but no SLO file available
# Should send a get-sli event
########################################################################################################################

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

# Create a config map containing the default sli-provider for the lighthouse service
kubectl create configmap -n ${KEPTN_NAMESPACE} lighthouse-config --from-literal=sli-provider=dynatrace

keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 10

# try to fetch a get-sli event
echo "Getting get-sli event with context-id: ${keptn_context_id}"
response=$(get_event sh.keptn.internal.event.get-sli ${keptn_context_id} ${PROJECT})

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.internal.event.get-sli"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.sliProvider" "dynatrace"

kubectl delete configmap -n ${KEPTN_NAMESPACE} lighthouse-config

########################################################################################################################
# Testcase 3: Send a start-evaluation event with an SLO file specified, but without an SLI provider configured
########################################################################################################################

# add SLO file for service
echo "Adding SLO File: test/assets/quality_gates_standalone_slo_step1.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_slo_step1.yaml --resourceUri=slo.yaml
verify_test_step $? "keptn add-resource slo.yaml - failed"

# add SLI file for service
echo "Adding SLI File: test/assets/quality_gates_standalone_sli_dynatrace_step1.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_sli_dynatrace_step1.yaml --resourceUri=dynatrace/sli.yaml
verify_test_step $? "keptn add-resource sli.yaml - failed"

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 10

# try to fetch a evaluation-done event
echo "Getting evaluation-done event with context-id: ${keptn_context_id}"
response=$(get_evaluation_done_event ${keptn_context_id})

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"
verify_using_jq "$response" ".data.evaluationdetails.result" "no evaluation performed by lighthouse because no SLI-provider configured for project ${PROJECT}"
verify_using_jq "$response" ".data.evaluationdetails.score" "0"
verify_using_jq "$response" ".data.evaluationdetails.sloFileContent" ""


########################################################################################################################
# Testcase 4: Send a start-evaluation event with an SLO file specified and with an SLI provider set, but no Dynatrace
#             Tenant/API Token configured
########################################################################################################################
echo "Install dynatrace-sli-service from: ${DYNATRACE_SLI_SERVICE_VERSION}"
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/${DYNATRACE_SLI_SERVICE_VERSION}/deploy/service.yaml -n ${KEPTN_NAMESPACE}
sleep 10
wait_for_deployment_in_namespace "dynatrace-sli-service" ${KEPTN_NAMESPACE}

# configure monitoring for Dynatrace
echo "Calling keptn configure monitoring dynatrace --project=$PROJECT"
keptn configure monitoring dynatrace --project=$PROJECT --suppress-websocket
sleep 5
# this should set the configmap 'lighthouse-config-$PROJECT' - verify that it exists

kubectl -n ${KEPTN_NAMESPACE} get configmap "lighthouse-config-${PROJECT}" -oyaml 2> /dev/null
verify_test_step $? "ERROR: Could not find ConfigMap lighthouse-config-$PROJECT (this is expected to be created by: keptn configure monitoring dynatrace --project=$PROJECT)"

# send the start evaluation command again
echo "Sending start-evaluation event for service $SERVICE in stage hardening"
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id: ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continue ..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  kubectl -n ${KEPTN_NAMESPACE} logs svc/dynatrace-sli-service -c dynatrace-sli-service
  print_error "evaluation-done event could not be retrieved"
  # exit 1 - Todo - see below
fi

# okay, evaluation-done event retrieved, parse it
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "fail"


########################################################################################################################
# Testcase 5: Run tests with Dynatrace credentials (Tenant and API token)
########################################################################################################################

# create secret from file
kubectl -n ${KEPTN_NAMESPACE} create secret generic dynatrace-credentials-${PROJECT} --from-literal="DT_TENANT=$QG_INTEGRATION_TEST_DT_TENANT" --from-literal="DT_API_TOKEN=$QG_INTEGRATION_TEST_DT_API_TOKEN"

# send the start evaluation command again
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id: ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continue ..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  echo "Logs from: dynatrace-sli-service"
  kubectl -n ${KEPTN_NAMESPACE} logs svc/dynatrace-sli-service -c dynatrace-sli-service
  echo "Logs from: lighthouse-service"
  kubectl -n ${KEPTN_NAMESPACE} logs svc/lighthouse-service -c lighthouse-service
  print_error "evaluation-done event could not be retrieved"
  exit 1
fi

# okay, evaluation-done event retrieved, parse it
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"

# verify .data.evaluationdetails: There should be 3 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# Verify .data.evaluationdetails: There should be 2 results with status: pass, and 1 with status: info
number_of_pass_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "pass")
number_of_warning_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "warning")
number_of_info_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "info")

if [[ $number_of_pass_results -ne 2 ]]; then
  echo "Expected 2 results with status: pass, but found $number_of_pass_results"
fi

if [[ $number_of_warning_results -ne 0 ]]; then
  echo "Expected 0 results with status: warning, but found $number_of_warning_results"
fi

if [[ $number_of_info_results -ne 1 ]]; then
  echo "Expected 1 results with status: info, but found $number_of_info_results"
fi


########################################################################################################################
# Testcase 6: Run the test again
########################################################################################################################
sleep 30

# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id: ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continue ..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n ${KEPTN_NAMESPACE} logs svc/dynatrace-sli-service -c dynatrace-sli-service
  print_error "evaluation-done event could not be retrieved"
  exit 1
fi

# okay, evaluation-done event retrieved, parse it
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"

# verify .data.evaluationdetails: There should be 3 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# verify .data.evaluationsdetails: There should be 2 results with status: pass, and 1 with status: info
number_of_pass_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "pass")
number_of_warning_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "warning")
number_of_info_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "info")

if [[ $number_of_pass_results -ne 2 ]]; then
  echo "Expected 2 results with status: pass, but found $number_of_pass_results"
fi

if [[ $number_of_warning_results -ne 0 ]]; then
  echo "Expected 0 results with status: warning, but found $number_of_warning_results"
fi

if [[ $number_of_info_results -ne 1 ]]; then
  echo "Expected 1 results with status: info, but found $number_of_info_results"
fi


########################################################################################################################
# Testcase 7: Add slo step2 which contains values that are not handled by dynatrace-sli-service
########################################################################################################################
echo "Adding SLO file: test/assets/quality_gates_standalone_slo_step2.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_slo_step2.yaml --resourceUri=slo.yaml
verify_test_step $? "keptn add-resource slo.yaml - failed"

# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id: ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continue ..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n ${KEPTN_NAMESPACE} logs svc/dynatrace-sli-service -c dynatrace-sli-service
  print_error "evaluation-done event could not be retrieved"
  exit 1
fi

# okay, evaluation-done event retrieved, parse it
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"

# verify .data.evaluationdetails: There should be 3 results that are true, and 2 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 2 ]]; then
  echo "Expected 2 results with success: false, but found $number_of_false_results"
fi

# verify .data.evaluationsdetails: There should be 2 results with status: pass, and 3 with status: info
number_of_pass_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "pass")
number_of_warning_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "warning")
number_of_info_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "info")

if [[ $number_of_pass_results -ne 2 ]]; then
  echo "Expected 2 results with status: pass, but found $number_of_pass_results"
fi

if [[ $number_of_warning_results -ne 0 ]]; then
  echo "Expected 0 results with status: warning, but found $number_of_warning_results"
fi

if [[ $number_of_info_results -ne 3 ]]; then
  echo "Expected 3 results with status: info, but found $number_of_info_results"
fi


########################################################################################################################
# Testcase 8: Also add sli step2 such that dynatrace-sli-service finally has the correct sli configs
########################################################################################################################

# add SLI file for service
echo "Adding SLI File: test/assets/quality_gates_standalone_sli_dynatrace_step2.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_sli_dynatrace_step2.yaml --resourceUri=dynatrace/sli.yaml
verify_test_step $? "keptn add-resource sli.yaml - failed"

# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id: ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continue ..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n ${KEPTN_NAMESPACE} logs svc/dynatrace-sli-service -c dynatrace-sli-service
  print_error "evaluation-done event could not be retrieved"
  exit 1
fi

# okay, evaluation-done event retrieved, parse it
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "pass"

# verify .data.evaluationdetails: There should be 5 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 5 ]]; then
  echo "Expected 5 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# verify .data.evaluationsdetails: There should be 2 results with status: pass, and 3 with status: info
number_of_pass_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "pass")
number_of_warning_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "warning")
number_of_info_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].status' | grep -c "info")

if [[ $number_of_pass_results -ne 2 ]]; then
  echo "Expected 2 results with status: pass, but found $number_of_pass_results"
fi

if [[ $number_of_warning_results -ne 0 ]]; then
  echo "Expected 0 results with status: warning, but found $number_of_warning_results"
fi

if [[ $number_of_info_results -ne 3 ]]; then
  echo "Expected 3 results with status: info, but found $number_of_info_results"
fi


########################################################################################################################
# cleanup
########################################################################################################################

echo "Deleting project ${PROJECT}"
keptn delete project $PROJECT

echo "Uninstalling dynatrace-sli-service"
kubectl -n ${KEPTN_NAMESPACE} delete deployment dynatrace-sli-service

echo "Removing secret dynatrace-credentials-${PROJECT}"
kubectl -n ${KEPTN_NAMESPACE} delete secret dynatrace-credentials-${PROJECT}

echo "Quality gates standalone tests done ✓"

exit 0
