################################################################
# This is shared library
################################################################

function timestamp() {
  date +"[%m-%d %H:%M:%S]"
}

function print_info() {
  echo -n "[keptn|INFO] "
  echo -n `timestamp`
  echo " $1"
}

function print_debug() {
  echo -n "[keptn|DEBUG] "
  echo -n `timestamp`
  echo " $1"
}

function print_error() {
  echo -n "[keptn|ERROR] "
  echo -n `timestamp`
  echo " $1"
}

function verify_install_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2" && exit 1
  fi
}

function verify_kubectl() {
  if [[ $1 != '0' ]]; then
    print_error "$2" && exit 1
  fi
}

function verify_variable() {
  if [[ -z "$1" ]]; then
    print_error "$2" && exit 1
  fi
}

function wait_for_all_pods_in_namespace {
  RETRY=0; RETRY_MAX=12; NAMESPACE=$1

  CMD="kubectl get pods -n $NAMESPACE && [[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"
  #CMD="[[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    eval $CMD

    if [[ $? == '0' ]]; then
      print_debug "All pods are running in namespace: $NAMESPACE, continue installation."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for pods to start..."
    sleep 10
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Pods in namespace: $NAMESPACE are not running, stop installation." && exit 1
  fi
}

# Waits for all custom resource defintions to be created successfully.
function wait_for_crds() {
  RETRY=0; RETRY_MAX=12
  CRDS=$1 # list of custom resource definitions

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    kubectl get $CRDS

    if [[ $? == '0' ]]; then
      print_debug "All custom resource definitions now available, continue installation."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for custom resource definitions..."
    sleep 10
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Custom resource definitions are missing, stop installation." && exit 1
  fi
}