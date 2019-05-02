################################################################
# This is shared library for the installation                  #
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
    print_error "Stopping keptn update. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
    exit 1
  fi
}

function verify_kubectl() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Stopping keptn update. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
    exit 1
  fi
}

function verify_variable() {
  if [[ -z "$1" ]]; then
    print_error "$2"
    print_error "Stopping keptn update. Already created resources are not deleted; execute the uninstallKeptn.sh script to clean-up."
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
        print_debug "Deployment ${DEPLOYMENT} in ${NAMESPACE} namespace available, continue update."
        break
      fi
      RETRY=$[$RETRY+1]
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 10
    done

    if [[ $RETRY == $RETRY_MAX ]]; then
      print_error "Deployment ${DEPLOYMENT} in namespace ${NAMESPACE} is not available"
      exit 1
    fi
  done
}