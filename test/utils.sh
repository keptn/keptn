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

function verify_deployment_in_namespace() {
  DEPLOYMENT=$1; NAMESPACE=$2;

  DEPLOYMENT_LIST=$(eval "kubectl get deployments -n ${NAMESPACE} | awk '/$DEPLOYMENT/'" | awk '{print $1}') # list of multiple deployments when starting with the same name
  if [[ -z "$DEPLOYMENT_LIST" ]]; then
    echo "Error: Could not find deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
    exit -1
  else
    echo "Found deployment ${DEPLOYMENT} in namespace ${NAMESPACE}"
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