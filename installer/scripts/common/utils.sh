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
      RETRY=$[$RETRY+1]
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 10
    else
      break
    fi
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Deployment ${DEPLOYMENT} in namespace ${NAMESPACE} is not available"
    exit 1
  fi

  RETRY=0

  verify_variable "$DEPLOYMENT_LIST" "List of deployments in namespace $NAMESPACE could not be derived."

  array=(${DEPLOYMENT_LIST// / })

  for DEPLOYED_DEPLOYMENT in "${array[@]}" 
  do
    while [[ $RETRY -lt $RETRY_MAX ]]; do
      kubectl rollout status deployment $DEPLOYMENT -n $NAMESPACE

      if [[ $? == '0' ]]
      then
        print_debug "Deployment ${DEPLOYED_DEPLOYMENT} in ${NAMESPACE} namespace available."
        break
      fi
      RETRY=$[$RETRY+1]
      print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for deployment ${DEPLOYED_DEPLOYMENT} in namespace ${NAMESPACE} ..."
      sleep 10
    done

    if [[ $RETRY == $RETRY_MAX ]]; then
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
    eval $CMD

    if [[ $? == '0' ]]; then
      print_debug "All pods are running in namespace ${NAMESPACE}."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for pods to start in namespace ${NAMESPACE} ..."
    sleep 10
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Pods in namespace ${NAMESPACE} are not running."
    # show the pods that have problems
    kubectl get pods --field-selector=status.phase!=Running -n ${NAMESPACE}
    exit 1
  fi
}

# Waits for all custom resource defintions to be created successfully.
function wait_for_crds() {
  CRDS=$1; # list of custom resource definitions
  RETRY=0; RETRY_MAX=24;

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    kubectl get $CRDS

    if [[ $? == '0' ]]; then
      print_debug "All custom resource definitions are available."
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 10s for custom resource definitions ..."
    sleep 10
  done

  if [[ $RETRY == $RETRY_MAX ]]; then
    print_error "Custom resource definitions are missing."
    exit 1
  fi
}

# Waits for ip of Istio ingress gateway (max wait time 30sec)
function wait_for_ingressgateway() {
  PROPERTY=$1;SVC=$2;NAMESPACE=$3;
  RETRY=0; RETRY_MAX=6;
  DOMAIN="";

  while [[ $RETRY -lt $RETRY_MAX ]]; do
    DOMAIN=$(kubectl get svc ${SVC} -o json -n ${NAMESPACE} | jq -r .status.loadBalancer.ingress[0].${PROPERTY})
    if [[ $DOMAIN = "null" ]]; then
      DOMAIN=""
    fi

    if [[ "$DOMAIN" != "" ]]; then
      print_debug "${PROPERTY} of Istio ingress gateway is available: ${DOMAIN}"
      break
    fi
    RETRY=$[$RETRY+1]
    print_debug "Retry: ${RETRY}/${RETRY_MAX} - Wait 5s for ${PROPERTY} of Istio ingress gateway to be available ..."
    sleep 5
  done
}

function setupKeptnDomain() {
  PROVIDER=$1;SVC=$2;NAMESPACE=$3;

  print_info "Determining ingress hostname/ip for Keptn (using ${PROVIDER})"
  # Domain used for routing to keptn services
  if [[ "$GATEWAY_TYPE" == "LoadBalancer" ]]; then
    wait_for_ingressgateway "hostname" $SVC $NAMESPACE
    export DOMAIN=$(kubectl get svc $SVC -o json -n $NAMESPACE | jq -r .status.loadBalancer.ingress[0].hostname)
    if [[ $? != 0 ]]; then
        print_error "Failed to get ingress gateway information." && exit 1
    fi
    export INGRESS_HOST=$DOMAIN

    if [[ "$DOMAIN" == "null" ]]; then
        print_info "Could not get domain name. Trying to retrieve IP address instead."

        wait_for_ingressgateway "ip" $SVC $NAMESPACE

        export DOMAIN=$(kubectl get svc $SVC -o json -n $NAMESPACE | jq -r .status.loadBalancer.ingress[0].ip)
        if [[ "$DOMAIN" == "null" ]]; then
            print_error "Could not get IP."
            exit 1
        fi
        export DOMAIN="$DOMAIN.xip.io"
        export INGRESS_HOST=$DOMAIN
    fi
  elif [[ "$GATEWAY_TYPE" == "NodePort" ]]; then
      NODE_PORT=$(kubectl -n $NAMESPACE get service $SVC -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
      NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
      export DOMAIN="$NODE_IP.xip.io:$NODE_PORT"
      export INGRESS_HOST="$NODE_IP.xip.io"
  fi

  print_info "Determined ${DOMAIN} and ${INGRESS_HOST}"

  if [[ "$PLATFORM" == "eks" ]]; then
      print_info "For EKS: No SSL certificate created. Please use keptn configure domain at the end of the installation."
  else
      # Set up SSL (self-signed certificate)
      print_info "Setting up self-signed SSL certificate."
      openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$INGRESS_HOST"

      if [[ "$PROIVDER" == "istio" ]]; then
        kubectl create --namespace $NAMESPACE secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
      elif [[ "$PROVIDER" == "nginx" ]]; then
        kubectl create secret tls sslcerts --key key.pem --cert certificate.pem -n keptn
      fi
      #verify_kubectl $? "Creating secret for istio-ingressgateway-certs failed."

      rm key.pem
      rm certificate.pem
  fi

  # Add config map in keptn namespace that contains the domain - this will be used by other services as well
  cat ../manifests/keptn/keptn-domain-configmap.yaml | \
    sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' | kubectl apply -f -
  verify_kubectl $? "Creating configmap keptn-domain in keptn namespace failed."
}