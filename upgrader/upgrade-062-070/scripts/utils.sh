#!/usr/bin/env bash
# shellcheck disable=SC2181

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
    print_error "Stopping keptn installation. Already created resources are not deleted; run keptn uninstall to clean-up."
    exit 1
  fi
}

function verify_kubectl() {
  if [[ $1 != '0' ]]; then
    print_error "$2"
    print_error "Stopping keptn installation. Already created resources are not deleted; run keptn uninstall to clean-up."
    exit 1
  fi
}

function verify_variable() {
  if [[ -z "$1" ]]; then
    print_error "$2"
    print_error "Stopping keptn installation. Already created resources are not deleted; run keptn uninstall to clean-up."
    exit 1
  fi
}

# Waits for a deployment in a given namespace to be available.
function wait_for_deployment_in_namespace() {
  DEPLOYMENT=$1; NAMESPACE=$2;
  RETRY=0; RETRY_MAX=24;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    DEPLOYMENT_LIST=$(eval "kubectl get deployments -n ${NAMESPACE} | awk '/$DEPLOYMENT/'" | awk '{print $1}') # list of multiple deployments when starting with the same name
    if [[ -z "$DEPLOYMENT_LIST" ]]; then
      RETRY=$((RETRY+1))
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 10
    else
      break
    fi
  done

  if [[ "$RETRY" == "$RETRY_MAX" ]]; then
    print_error "Deployment ${DEPLOYMENT} in namespace ${NAMESPACE} is not available"
    exit 1
  fi

  RETRY=0

  verify_variable "$DEPLOYMENT_LIST" "List of deployments in namespace $NAMESPACE could not be derived."

  # shellcheck disable=SC2206
  array=(${DEPLOYMENT_LIST// / })

  for DEPLOYED_DEPLOYMENT in "${array[@]}"
  do
    while [[ $RETRY -lt $RETRY_MAX ]]; do
      kubectl rollout status deployment "$DEPLOYMENT" -n "$NAMESPACE"

      if [[ $? == '0' ]]
      then
        print_debug "Deployment ${DEPLOYED_DEPLOYMENT} in ${NAMESPACE} namespace available."
        break
      fi
      RETRY=$((RETRY+1))
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYED_DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 10
    done

    if [[ "$RETRY" == "$RETRY_MAX" ]]; then
      print_error "Deployment ${DEPLOYED_DEPLOYMENT} in namespace ${NAMESPACE} is not available"
      exit 1
    fi
  done
}

# Waits for all pods in a given namespace to be up and running.
function wait_for_all_pods_in_namespace() {
  NAMESPACE=$1;
  RETRY=0; RETRY_MAX=24;

  CMD="kubectl get pods -n $NAMESPACE && [[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"
  #CMD="[[ \$(kubectl get pods -n $NAMESPACE 2>&1 | grep -c -v -E '(Running|Completed|Terminating|STATUS)') -eq 0 ]]"

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    eval "$CMD"

    if [[ $? == '0' ]]; then
      print_debug "All pods are running in namespace ${NAMESPACE}."
      break
    fi
    RETRY=$((RETRY+1))
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for pods to start in namespace ${NAMESPACE} ..."
    sleep 10
  done

  if [[ $RETRY == "$RETRY_MAX" ]]; then
    print_error "Pods in namespace ${NAMESPACE} are not running."
    # show the pods that have problems
    kubectl get pods --field-selector=status.phase!=Running -n "${NAMESPACE}"
    exit 1
  fi
}

# Waits for all custom resource definitions to be created successfully.
function wait_for_crds() {
  CRDS=$1; # list of custom resource definitions
  RETRY=0; RETRY_MAX=24;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    kubectl get "$CRDS"

    if [[ $? == '0' ]]; then
      print_debug "All custom resource definitions are available."
      break
    fi
    RETRY=$((RETRY+1))
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for custom resource definitions ..."
    sleep 10
  done

  if [[ $RETRY == "$RETRY_MAX" ]]; then
    print_error "Custom resource definitions are missing."
    exit 1
  fi
}

# Waits for ip of Istio ingress gateway (max wait time 30sec)
function wait_for_istio_ingressgateway() {
  PROPERTY=$1;
  RETRY=0; RETRY_MAX=6;
  DOMAIN="";

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    DOMAIN=$(kubectl get svc istio-ingressgateway -o json -n istio-system | jq -r ".status.loadBalancer.ingress[0].${PROPERTY}")
    if [[ $DOMAIN = "null" ]]; then
      DOMAIN=""
    fi

    if [[ "$DOMAIN" != "" ]]; then
      print_debug "${PROPERTY} of Istio ingress gateway is available."
      break
    fi
    RETRY=$((RETRY+1))
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 5s for ${PROPERTY} of Istio ingress gateway to be available ..."
    sleep 5
  done
}

# Waits for hostname or ip of ingress gateway (max wait time 120sec)
function wait_for_k8s_ingress() {
  RETRY=0; RETRY_MAX=24;
  DOMAIN="";

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    # Check availability of hostname
    DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].hostname)
    if [[ $DOMAIN = "null" ]]; then
      DOMAIN=""
    fi

    if [[ "$DOMAIN" != "" ]]; then
      print_debug "Domain name of ingress gateway is available."
      break
    fi

    # Check availability of IP
    DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].ip)
    if [[ $DOMAIN = "null" ]]; then
      DOMAIN=""
    fi

    if [[ "$DOMAIN" != "" ]]; then
      print_debug "IP address of ingress gateway is available."
      break
    fi

    RETRY=$((RETRY+1))
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 5s for domain name or IP address of ingress gateway to be available ..."
    sleep 5

  done
}
