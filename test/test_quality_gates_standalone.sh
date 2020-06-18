#!/bin/bash

source test/utils.sh

# test configuration
DYNATRACE_SLI_SERVICE_VERSION=${DYNATRACE_SLI_SERVICE_VERSION:-0.4.1}
KEPTN_EXAMPLES_BRANCH=${KEPTN_EXAMPLES_BRANCH:-master}
PROJECT=${PROJECT:-easytravel}

# get keptn api details
KEPTN_ENDPOINT=https://api.keptn.$(kubectl get cm keptn-domain -n keptn -ojsonpath={.data.app_domain})
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

########################################################################################################################
# Pre-requesits
########################################################################################################################

# ensure dynatrace-sli-service is not installed yet
kubectl -n keptn get deployment dynatrace-sli-service

if [[ $? -eq 0 ]]; then
  echo "Found dynatrace-sli-service. Please uninstall it using"
  echo "kubectl -n keptn delete deployment dynatrace-sli-service"
  exit 1
fi


# verify that the project does not exiset yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/configuration-service/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

# verify that the lighthouse configmap for the project does not exiset yet
kubectl -n keptn get cm lighthouse-config-${PROJECT}

if [[ $? -eq 0 ]]; then
  echo "Found configmap lighthouse-config-${PROJECT}. Please remove it using"
  echo "kubectl -n keptn delete configmap lighthouse-config-${PROJECT}"
  exit 1
fi

# verify that the Dynatrace credential secret does not exist yet
kubectl -n keptn get secret dynatrace-credentials-${PROJECT}

if [[ $? -eq 0 ]]; then
  echo "Found secret dynatrace-credentials-${PROJECT}. Please remove it using"
  echo "kubectl -n keptn delete secret dynatrace-credentials-${PROJECT}"
  exit 1
fi


echo "Testing quality gates standalone for project $PROJECT ..."

# Test keptn create-project and create service
rm -rf examples
git clone --branch ${KEPTN_EXAMPLES_BRANCH} https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts

echo "Creating a new project without git upstream"
keptn create project $PROJECT --shipyard=./shipyard-quality-gates.yaml
verify_test_step $? "keptn create project command failed."
sleep 10

cd ../..

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
# create service catalogue                #
###########################################
SERVICE=frontend
keptn create service $SERVICE --project=$PROJECT
verify_test_step $? "keptn create service ${SERVICE} failed."
sleep 10

########################################################################################################################
# Testcase 1:
# Project and service should have been created, but no SLO file added and no SLI provider configured
# Sending a start-evaluation event now should pass with an appropriate message
########################################################################################################################

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 10

echo "Getting evaluation-done event with context-id ${keptn_context_id}"

# try to fetch a evaluation-done event
response=$(get_evaluation_done_event ${keptn_context_id})

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" ""
verify_using_jq "$response" ".data.evaluationdetails.result" "no evaluation performed by lighthouse because no SLO found for service ${SERVICE}"
verify_using_jq "$response" ".data.evaluationdetails.score" "0"
verify_using_jq "$response" ".data.evaluationdetails.sloFileContent" ""


########################################################################################################################
# Testcase 2: Send a start-evaluation event with an SLO file specified, but without an SLI provider configured
########################################################################################################################

# add SLO file for service
echo "Adding SLO File: test/assets/quality_gates_standalone_slo_step1.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_slo_step1.yaml --resourceUri=slo.yaml
verify_test_step $? "keptn add-resource slo.yaml failed."

# add SLI file for service
echo "Adding SLI File: test/assets/quality_gates_standalone_sli_dynatrace_step1.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_sli_dynatrace_step1.yaml --resourceUri=dynatrace/sli.yaml
verify_test_step $? "keptn add-resource sli.yaml failed."

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 10

echo "Getting evaluation-done event with context-id ${keptn_context_id}"

# try to fetch a evaluation-done event
response=$(get_evaluation_done_event ${keptn_context_id})

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "lighthouse-service"
verify_using_jq "$response" ".type" "sh.keptn.events.evaluation-done"
verify_using_jq "$response" ".data.project" "${PROJECT}"
verify_using_jq "$response" ".data.stage" "hardening"
verify_using_jq "$response" ".data.service" "${SERVICE}"
verify_using_jq "$response" ".data.result" "failed"
verify_using_jq "$response" ".data.evaluationdetails.result" "no evaluation performed by lighthouse because no SLI-provider configured for project ${PROJECT}"
verify_using_jq "$response" ".data.evaluationdetails.score" "0"
verify_using_jq "$response" ".data.evaluationdetails.sloFileContent" ""

########################################################################################################################
# Testcase 3: Send a start-evaluation event with an SLO file specified and with an SLI provider set, but no Dynatrace
#             Tenant/API Token configured
########################################################################################################################

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/${DYNATRACE_SLI_SERVICE_VERSION}/deploy/service.yaml
sleep 10

wait_for_deployment_in_namespace "dynatrace-sli-service" "keptn"

# now configure monitoring for Dynatrace
echo "Calling keptn configure monitoring dynatrace --project=$PROJECT"
keptn configure monitoring dynatrace --project=$PROJECT --suppress-websocket
sleep 5
# this should set the configmap 'lighthouse-config-$PROJECT' - verify that it exists

kubectl -n keptn get configmap "lighthouse-config-${PROJECT}" -oyaml
verify_test_step $? "ERROR: Could not find configmap lighthouse-config-$PROJECT (this is expected to be created by keptn configure monitoring dynatrace --project=$PROJECT)"

echo "Sending start-evaluation event for service $SERVICE in stage hardening"

# now that this is set, let's send the start evaluation command again
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# Wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continuing..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  kubectl -n keptn logs svc/dynatrace-sli-service
  print_error "evaluation-done event could not be retrieved"
  # exit 1 - Todo - see below
fi

# ToDo: there is no response here right now, but in fact dynatrace-sli-service should send a note...
echo $response | grep "No event returned"

if [[ $? -ne 0 ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n keptn logs svc/dynatrace-sli-service
  echo "Expected an 'No event returned' in the response, but got"
  echo $response
  exit 1
fi

########################################################################################################################
# Testcase 4: Run tests with Dynatrace credentials (tenant and api token) set
########################################################################################################################

# create secret from file
kubectl -n keptn create secret generic dynatrace-credentials-${PROJECT} --from-literal="DT_TENANT=$QG_INTEGRATION_TEST_DT_TENANT" --from-literal="DT_API_TOKEN=$QG_INTEGRATION_TEST_DT_API_TOKEN"

# now that this is set, let's send the start evaluation command again
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# Wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continuing..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n keptn logs svc/dynatrace-sli-service
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

# Verify .data.evaluationdetails: There should be 3 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# Verify .data.evaluationsdetails: There should be 2 results with status: pass, and 1 with status: info
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
# Testcase 5: Run the test again
########################################################################################################################
sleep 30

# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# Wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continuing..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n keptn logs svc/dynatrace-sli-service
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

# Verify .data.evaluationdetails: There should be 3 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# Verify .data.evaluationsdetails: There should be 2 results with status: pass, and 1 with status: info
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
# Testcase 6: Add slo step2 which contains values that are not handled by dynatrace-sli-service
########################################################################################################################
echo "Adding SLO File: test/assets/quality_gates_standalone_slo_step2.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_slo_step2.yaml --resourceUri=slo.yaml
verify_test_step $? "keptn add-resource slo.yaml failed."

# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# Wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continuing..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n keptn logs svc/dynatrace-sli-service
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

# Verify .data.evaluationdetails: There should be 3 results that are true, and 2 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 3 ]]; then
  echo "Expected 3 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 2 ]]; then
  echo "Expected 2 results with success: false, but found $number_of_false_results"
fi

# Verify .data.evaluationsdetails: There should be 2 results with status: pass, and 3 with status: info
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
# Testcase 7: Also add sli step2 such that dynatrace-sli-service finally has the correct sli configs
########################################################################################################################

# add SLI file for service
echo "Adding SLI File: test/assets/quality_gates_standalone_sli_dynatrace_step2.yaml"
keptn add-resource --project=$PROJECT --stage=hardening --service=$SERVICE --resource=test/assets/quality_gates_standalone_sli_dynatrace_step2.yaml --resourceUri=dynatrace/sli.yaml
verify_test_step $? "keptn add-resource sli.yaml failed."


# send start-evaluation event
keptn_context_id=$(send_start_evaluation_event $PROJECT hardening $SERVICE)
sleep 30

# Wait until a evaluation-done event is retrieved
echo "Trying to get evaluation-done event with context-id ${keptn_context_id}"

RETRY=0; RETRY_MAX=30;

while [[ $RETRY -lt $RETRY_MAX ]]; do
  # try to fetch the evaluation-done event
  response=$(get_evaluation_done_event ${keptn_context_id})

  # check if this contains an error
  echo $response | grep "No event returned"

  if [[ $? -ne 0 ]]; then
    echo "Received an evaluation-done event, continuing..."
    break
  else
    RETRY=$[$RETRY+1]
    echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for evaluation-done event"
    sleep 10
  fi
done

if [[ $RETRY == $RETRY_MAX ]]; then
  # print logs of dynatrace-sli-service
  kubectl -n keptn logs svc/dynatrace-sli-service
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

# Verify .data.evaluationdetails: There should be 5 results that are true, and 0 false
number_of_true_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "true")
number_of_false_results=$(echo $response | jq -r '.data.evaluationdetails.indicatorResults[].value.success' | grep -c "false")

if [[ $number_of_true_results -ne 5 ]]; then
  echo "Expected 5 results with success: true, but found $number_of_true_results"
fi

if [[ $number_of_false_results -ne 0 ]]; then
  echo "Expected 0 results with success: false, but found $number_of_false_results"
fi

# Verify .data.evaluationsdetails: There should be 2 results with status: pass, and 3 with status: info
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
kubectl -n keptn delete deployment dynatrace-sli-service

echo "Removing secret dynatrace-credentials-${PROJECT}"
kubectl -n keptn delete secret dynatrace-credentials-${PROJECT}

echo "Done!"

exit 0
