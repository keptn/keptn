#!/bin/bash
source ./common/utils.sh

kubectl apply -f ../manifests/keptn/api-ingress.yaml
verify_install_step $? "Installing Keptn api-ingress failed."


wait_for_k8s_ingress "hostname"
export DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].hostname)
if [[ $? != 0 ]]; then
    print_error "Failed to get k8s ingress gateway information." && exit 1
fi

if [[ "$DOMAIN" == "null" ]]; then
    print_info "Could not get ingress gateway domain name. Trying to retrieve IP address instead."

    wait_for_k8s_ingress "ip"

    export DOMAIN=$(kubectl get ingress api-ingress -n keptn -o json | jq -r .status.loadBalancer.ingress[0].ip)
    if [[ "$DOMAIN" == "null" ]]; then
        print_error "IP of k8s ingress gateway could not be derived."
        exit 1
    fi
    export DOMAIN="$DOMAIN.xip.io"
fi

echo $DOMAIN

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$DOMAIN"
kubectl create secret tls sslcerts --key key.pem --cert certificate.pem -n keptn
rm key.pem
rm certificate.pem  

# Update ingress with updates hosts
cat ../manifests/keptn/api-ingress.yaml | \
    sed 's~domain.placeholder~'"$DOMAIN"'~' > ../manifests/keptn/gen/api-ingress.yaml

kubectl apply -f ../manifests/keptn/gen/api-ingress.yaml
verify_kubectl $? "Deploying Keptn ingress failed."