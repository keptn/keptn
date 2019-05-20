################################################################
# This is shared library for the keptn installation            #
################################################################

function timestamp() {
  date +"[%Y-%m-%d %H:%M:%S]"
}

function print_info() {
  echo "[keptn|INFO] $(timestamp) $1"
}

function print_debug() {
  echo "[keptn|DEBUG] $(timestamp) $1"
}

function print_error() {
  echo "[keptn|ERROR] $(timestamp) $1"
}

function verify_install_step() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Stopping keptn installation. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
    exit 1
  fi
}

function verify_kubectl() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Stopping keptn installation. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
    exit 1
  fi
}

function verify_variable() {
  if [[ -z "$1" ]]; then
    print_error "$2"
    print_error "Stopping keptn installation. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
    exit 1
  fi
}

# Waits for a deployment in a given namespace to be available.
function wait_for_deployment_in_namespace() {
  DEPL=$1; NAMESPACE=$2;
  RETRY=0; RETRY_MAX=12; 

  DEPLOYMENT_LIST=$(eval "kubectl get deployments -n $NAMESPACE | awk '/$DEPL/'" | awk '{print $1}') # list of multiple deployments when starting with the same name, e.g.: event-broker, event-broker-ext
  verify_variable "$DEPLOYMENT_LIST" "DEPLOYMENT_LIST could not be derived from deployments list of namespace $NAMESPACE."

  array=(${DEPLOYMENT_LIST// / })

  for DEPLOYMENT in "${array[@]}" 
  do
    while [[ $RETRY -lt $RETRY_MAX ]]; do
      kubectl rollout status deployment $DEPLOYMENT -n $NAMESPACE

      if [[ $? == '0' ]]
      then
        print_debug "Deployment ${DEPLOYMENT} in ${NAMESPACE} namespace available."
        break
      fi
      RETRY=$[$RETRY+1]
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 20s for deployment ${DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 20
    done

    if [[ $RETRY == $RETRY_MAX ]]; then
      print_error "Deployment ${DEPLOYMENT} in namespace ${NAMESPACE} is not available"
      exit 1
    fi
  done
}

# Waits for all pods in a given namespace to be up and running.
function wait_for_channel_in_namespace() {
  CHANNEL=$1; NAMESPACE=$2;
  RETRY=0; RETRY_MAX=12; 

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    kubectl get channel $CHANNEL -n $NAMESPACE

    if [[ $? == '0' ]]; then
      print_debug "Channel ${CHANNEL} in namespace ${NAMESPACE} available."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 20s for channel ${CHANNEL} in namespace ${NAMESPACE} to be available ..."
    sleep 20
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Channel in namespace ${NAMESPACE} not available."
    exit 1
  fi
}

# Waits for all pods in a given namespace to be up and running.
function wait_for_all_pods_in_namespace() {
  NAMESPACE=$1;
  RETRY=0; RETRY_MAX=12; 

  CMD="kubectl get pods -n $NAMESPACE && [[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"
  #CMD="[[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    eval $CMD

    if [[ $? == '0' ]]; then
      print_debug "All pods are running in namespace ${NAMESPACE}."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 20s for pods to start in namespace ${NAMESPACE} ..."
    sleep 20
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Pods in namespace ${NAMESPACE} are not running."
    exit 1
  fi
}

# Waits for all custom resource defintions to be created successfully.
function wait_for_crds() {
  CRDS=$1; # list of custom resource definitions
  RETRY=0; RETRY_MAX=12;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    kubectl get $CRDS

    if [[ $? == '0' ]]; then
      print_debug "All custom resource definitions are available."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 20s for custom resource definitions ..."
    sleep 20
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Custom resource definitions are missing."
    exit 1
  fi
}
