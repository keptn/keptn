function timestamp() {
  date +"[%Y-%m-%d %H:%M:%S]"
}

function print_error() {
  echo "[keptn|ERROR] $(timestamp) $1"
}

function send_start_evaluation_event() {
  PROJECT=$1
  STAGE=$2
  SERVICE=$3

  response=$(keptn send event start-evaluation --project=$PROJECT --stage=$STAGE --service=$SERVICE --timeframe=5m)
  keptn_context_id=$(echo $response | awk -F'Keptn context:' '{ print $2 }' | xargs)

  echo "$keptn_context_id"
}

function get_evaluation_done_event() {
  keptn_context_id=$1
  keptn get event evaluation-done --keptn-context="${keptn_context_id}" | tail -n +2
}

function send_evaluation_done_event() {
  PROJECT=$1
  STAGE=$2
  SERVICE=$3
  RESULT=$4

  cat ./test/assets/evaluation_done_event_template.json | jq -r --arg project $PROJECT --arg stage $STAGE --arg service $SERVICE --arg result $RESULT '.data.project=$project | .data.stage=$stage | .data.service=$service | .data.result=$result' > tmp_evaluation_done_event.json

  response=$(keptn send event --file=tmp_evaluation_done_event.json)
  rm tmp_evaluation_done_event.json

  keptn_context_id=$(echo $response | awk -F'Keptn context:' '{ print $2 }' | xargs)
  echo "$keptn_context_id"
}

function send_event_json() {
  EVENT_JSON_FILE_URI=$1

  response=$(keptn send event --file=$EVENT_JSON_FILE_URI)
  keptn_context_id=$(echo $response | awk -F'Keptn context:' '{ print $2 }' | xargs)

  echo "$keptn_context_id"
}

function get_keptn_event() {
  PROJECT=$1
  keptn_context_id=$2
  type=$3
  KEPTN_ENDPOINT=$4
  KEPTN_API_TOKEN=$5
  curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=${type}&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]'
}

function check_no_open_approvals() {
  PROJECT=$1
  STAGE=$2

  result=$(keptn get event approval.triggered --project=$PROJECT --stage=$STAGE | awk '{if(NR>1)print}')
  if [[ "$result" != "No approval.triggered events have been found" ]]; then
    echo "Received unexpected number of approval.triggered events"
    echo "${result}"
    exit 2
  else
    echo "Verified number of approval.triggered events"
  fi
}

function check_number_open_approvals() {
  PROJECT=$1
  STAGE=$2
  EXPECTED=$3

  result=$(keptn get event approval.triggered --project=$PROJECT --stage=$STAGE | awk '{if(NR>1)print}' | jq -r '.' | jq -r --slurp 'length')
  if [[ "$result" != "$EXPECTED" ]]; then
    echo "Received unexpected number of approval.triggered events"
    echo "${result}"

    exit 2
  else
    echo "Verified number of approval.triggered events"
  fi
}

function verify_using_jq() {
  payload=$1
  attribute=$2
  expected=$3

  actual=$(echo "${payload}" | jq -r "${attribute}")

  if [[ "${actual}" != "${expected}" ]]; then
    print_error "ERROR: Checking $attribute, expected '${expected}', got '${actual}' ❌"
    exit 1
  else
    echo "Checking $attribute: ${actual} ✓"
  fi

  return 0
}

function verify_test_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Keptn Test failed."
    exit 1
  fi
}

function wait_for_url() {
  URL=$1
  RETRY=0; RETRY_MAX=50;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    curl $URL -k

    if [[ $? -eq 0 ]]; then
      echo "Verified access to ${URL}!"
      break
    else
      RETRY=$[$RETRY+1]
      echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for url ${URL} ..."
      sleep 10
    fi
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "URL ${URL} could not be reached"
    exit 1
  fi
}

function verify_image_of_deployment() {
  DEPLOYMENT=$1; NAMESPACE=$2; IMAGE_NAME=$3;

  CURRENT_IMAGE_NAME=$(kubectl get deployment ${DEPLOYMENT} -n ${NAMESPACE} -o=jsonpath='
  }{$.spec.template.spec.containers[:1].image}')

  if [[ "$CURRENT_IMAGE_NAME" == "$IMAGE_NAME" ]]; then
    echo "Found image ${CURRENT_IMAGE_NAME} in deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
  else
    echo "ERROR: Found image ${CURRENT_IMAGE_NAME} but expected ${IMAGE_NAME}  in deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
    exit -1
  fi
}

function wait_for_deployment_in_namespace() {
  DEPLOYMENT=$1; NAMESPACE=$2;
  RETRY=0; RETRY_MAX=50;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    DEPLOYMENT_LIST=$(eval "kubectl get deployments -n ${NAMESPACE} | awk '/$DEPLOYMENT /'" | awk '{print $1}') # list of multiple deployments when starting with the same name
    if [[ -z "$DEPLOYMENT_LIST" ]]; then
      RETRY=$[$RETRY+1]
      echo "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
      sleep 15
    else
      echo "Found deployment ${DEPLOYMENT} in namespace ${NAMESPACE}: ${DEPLOYMENT_LIST}"
      break
    fi
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    echo "Error: Could not find deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
    exit -1
  fi
}

function verify_deployment_in_namespace() {
  DEPLOYMENT=$1; NAMESPACE=$2;

  DEPLOYMENT_LIST=$(eval "kubectl get deployments -n ${NAMESPACE} | awk '/$DEPLOYMENT /'" | awk '{print $1}') # list of multiple deployments when starting with the same name
  if [[ -z "$DEPLOYMENT_LIST" ]]; then
    echo "Error: Could not find deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
    exit -1
  else
    echo "Found deployment ${DEPLOYMENT} in namespace ${NAMESPACE}: ${DEPLOYMENT_LIST}"
  fi
}

function verify_pod_in_namespace() {
  POD=$1; NAMESPACE=$2;

  POD_LIST=$(eval "kubectl get pod -n ${NAMESPACE} | awk '/$POD/'" | awk '{print $1}') # list of multiple deployments when starting with the same name
  if [[ -z "$POD_LIST" ]]; then
    echo "Error: Could not find pod ${POD} in namespace ${NAMESPACE}"
    exit -1
  else
    echo "Found pod ${POD} in namespace ${NAMESPACE}: ${POD_LIST}"
  fi
}

function verify_namespace_exists() {
  NAMESPACE=$1;

  NAMESPACE_LIST=$(eval "kubectl get namespaces -L istio-injection | grep ${NAMESPACE} | awk '/$NAMESPACE/'" | awk '{print $1}')

  if [[ -z "$NAMESPACE_LIST" ]]; then
    echo "Error: Could not find namespace ${NAMESPACE}"
    exit -2
  else
    echo "Found namespace ${NAMESPACE}"
  fi
}
