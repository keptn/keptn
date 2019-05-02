#!/bin/bash
REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

source ./utils.sh

# Creating cluster role binding
kubectl apply -f ../manifests/keptn/keptn-rbac.yaml
verify_kubectl $? "Creating cluster role for keptn failed."

# Creating config map to store mapping
kubectl apply -f ../manifests/keptn/keptn-org-configmap.yaml
verify_kubectl $? "Creating config map for keptn failed."

# Mark internal docker registry as insecure registry for knative controller
VAL=$(kubectl -n knative-serving get cm config-controller -o=json | jq -r .data.registriesSkippingTagResolving | awk '{print $1",'$REGISTRY_URL':5000"}')
kubectl -n knative-serving get cm config-controller -o=yaml | yq w - data.registriesSkippingTagResolving $VAL | kubectl apply -f -
verify_kubectl $? "Marking internal docker registry as insecure failed."

# Deploy knative eventing channels
kubectl apply -f ../../core/eventbroker/config/channel.yaml
verify_kubectl $? "Creating keptn-channel channel failed."
wait_for_channel_in_namespace "keptn-channel" "keptn"

kubectl apply -f ../../core/eventbroker/config/new-artefact-channel.yaml
verify_kubectl $? "Creating new-artefact channel failed."
wait_for_channel_in_namespace "new-artefact" "keptn"

kubectl apply -f ../../core/eventbroker/config/configuration-changed-channel.yaml
verify_kubectl $? "Creating configuration-changed channel failed."
wait_for_channel_in_namespace "configuration-changed" "keptn"

kubectl apply -f ../../core/eventbroker/config/deployment-finished-channel.yaml
verify_kubectl $? "Creating deployment-finished channel failed."
wait_for_channel_in_namespace "deployment-finished" "keptn"

kubectl apply -f ../../core/eventbroker/config/tests-finished-channel.yaml
verify_kubectl $? "Creating tests-finished channel failed."
wait_for_channel_in_namespace "tests-finished" "keptn"

kubectl apply -f ../../core/eventbroker/config/evaluation-done-channel.yaml
verify_kubectl $? "Creating evaluation-done channel failed."
wait_for_channel_in_namespace "evaluation-done" "keptn"

kubectl apply -f ../../core/eventbroker/config/problem-channel.yaml
verify_kubectl $? "Creating problem channel failed."
wait_for_channel_in_namespace "problem" "keptn"

KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)
verify_variable "$KEPTN_API_TOKEN" "KEPTN_API_TOKEN could not be derived." 

kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"
verify_kubectl $? "Creating secret for keptn api token failed."

KEPTN_CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
verify_variable "$KEPTN_CHANNEL_URI" "KEPTN_CHANNEL_URI could not be derived from keptn-channel description." 

# Deploy keptn core components: eventbroker, eventbroker-ext, auth, control
cd ../../core/eventbroker
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $KEPTN_CHANNEL_URI
verify_install_step $? "Deploying keptn event-broker failed."
cd ../../install/scripts

cd ../../core/eventbroker-ext
chmod +x deploy.sh
./deploy.sh
verify_install_step $? "Deploying keptn event-broker-ext failed."
cd ../../install/scripts

cd ../../core/auth
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL
verify_install_step $? "Deploying keptn auth component failed."
cd ../../install/scripts

cd ../../core/control
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $KEPTN_CHANNEL_URI
verify_install_step $? "Deploying keptn control component failed."
cd ../../install/scripts

# Set up SSL
ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
verify_variable "$ISTIO_INGRESS_IP" "ISTIO_INGRESS_IP is empty and could not be derived from the Istio ingress gateway." 

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$ISTIO_INGRESS_IP.xip.io"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
verify_kubectl $? "Creating secret for istio-ingressgateway-certs failed."

kubectl get gateway knative-ingress-gateway --namespace knative-serving -o=yaml | yq w - spec.servers[1].tls.mode SIMPLE | yq w - spec.servers[1].tls.privateKey /etc/istio/ingressgateway-certs/tls.key | yq w - spec.servers[1].tls.serverCertificate /etc/istio/ingressgateway-certs/tls.crt | kubectl apply -f -
verify_kubectl $? "Updating knative ingress gateway with private key failed."

rm key.pem
rm certificate.pem

##############################################
## Start validation of keptn installation   ##
##############################################
wait_for_all_pods_in_namespace "keptn"

wait_for_deployment_in_namespace "event-broker" "keptn" # Wait function also waits for event-broker-ext
wait_for_deployment_in_namespace "auth" "keptn"
wait_for_deployment_in_namespace "control" "keptn"
