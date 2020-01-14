#!/bin/bash
source ./common/utils.sh

kubectl apply -f ../manifests/keptn/api-ingress.yaml
verify_install_step $? "Installing Keptn api-ingress failed."

if [[ "$GATEWAY_TYPE" == "LoadBalancer" ]]; then
  wait_for_k8s_ingress
  export DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].hostname)
  if [[ $? != 0 ]]; then
      print_error "Failed to get K8s ingress gateway information." && exit 1
  fi
  export INGRESS_HOST=$DOMAIN

  if [[ "$DOMAIN" == "null" ]]; then
      print_info "Could not get ingress gateway domain name. Retrieving IP address instead."

      export DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].ip)
      if [[ "$DOMAIN" == "null" ]]; then
          print_error "IP address of ingress gateway could not be retrieved."
          exit 1
      fi
      export DOMAIN="$DOMAIN.xip.io"
      export INGRESS_HOST=$DOMAIN
  fi
elif [[ "$GATEWAY_TYPE" == "NodePort" ]]; then
    NODE_PORT=$(kubectl -n ingress-nginx get service ingress-nginx -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
    NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
    export DOMAIN="$NODE_IP.xip.io:$NODE_PORT"
    export INGRESS_HOST="$NODE_IP.xip.io"
fi

echo $DOMAIN
echo $INGRESS_HOST

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$INGRESS_HOST"
kubectl create secret tls sslcerts --key key.pem --cert certificate.pem -n keptn
rm key.pem
rm certificate.pem  

# Update ingress with updated hosts
cat ../manifests/keptn/api-ingress.yaml | \
    sed 's~domain.placeholder~'"$INGRESS_HOST"'~' > ../manifests/keptn/gen/api-ingress.yaml

kubectl apply -f ../manifests/keptn/gen/api-ingress.yaml
verify_kubectl $? "Deploying Keptn ingress failed."
