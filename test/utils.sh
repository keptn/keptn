function timestamp() {
  date +"[%Y-%m-%d %H:%M:%S]"
}

function print_error() {
  echo "[keptn|ERROR] $(timestamp) $1"
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

  CURRENT_IMAGE_NAME=$(kubectl get deployment ${DEPLOYMENT} -n ${NAMESPACE} -o=jsonpath='{$.spec.template.spec.containers[:1].image}')

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