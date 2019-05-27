#!/bin/bash
REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

CONTROL_RELEASE="develop"
AUTHENTICATOR_RELEASE="develop"
EVENTBROKER_RELEASE="develop"
EVENTBROKER_EXT_RELEASE="develop"

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
kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/keptn-channel.yaml
verify_kubectl $? "Creating keptn-channel channel failed."
wait_for_channel_in_namespace "keptn-channel" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/new-artefact-channel.yaml
verify_kubectl $? "Creating new-artefact channel failed."
wait_for_channel_in_namespace "new-artefact" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/configuration-changed-channel.yaml
verify_kubectl $? "Creating configuration-changed channel failed."
wait_for_channel_in_namespace "configuration-changed" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/deployment-finished-channel.yaml
verify_kubectl $? "Creating deployment-finished channel failed."
wait_for_channel_in_namespace "deployment-finished" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/tests-finished-channel.yaml
verify_kubectl $? "Creating tests-finished channel failed."
wait_for_channel_in_namespace "tests-finished" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/evaluation-done-channel.yaml
verify_kubectl $? "Creating evaluation-done channel failed."
wait_for_channel_in_namespace "evaluation-done" "keptn"

kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/problem-channel.yaml
verify_kubectl $? "Creating problem channel failed."
wait_for_channel_in_namespace "problem" "keptn"

KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)
verify_variable "$KEPTN_API_TOKEN" "KEPTN_API_TOKEN could not be derived." 

kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"
#verify_kubectl $? "Creating secret for keptn api token failed."

KEPTN_CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
verify_variable "$KEPTN_CHANNEL_URI" "KEPTN_CHANNEL_URI could not be derived from keptn-channel description." 

# Deploy eventbroker component
kubectl delete -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/eventbroker.yaml --ignore-not-found
kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker/$EVENTBROKER_RELEASE/config/eventbroker.yaml
verify_kubectl $? "Deploying keptn eventbroker component failed."

# Deploy eventbroker-ext component
kubectl delete -f https://raw.githubusercontent.com/keptn/eventbroker-ext/$EVENTBROKER_EXT_RELEASE/config/eventbroker-ext.yaml --ignore-not-found
kubectl apply -f https://raw.githubusercontent.com/keptn/eventbroker-ext/$EVENTBROKER_EXT_RELEASE/config/eventbroker-ext.yaml
verify_kubectl $? "Deploying keptn eventbroker-ext component failed."

# Deploy authenticator component
kubectl delete -f https://raw.githubusercontent.com/keptn/authenticator/$AUTHENTICATOR_RELEASE/config/authenticator.yaml --ignore-not-found
kubectl apply -f https://raw.githubusercontent.com/keptn/authenticator/$AUTHENTICATOR_RELEASE/config/authenticator.yaml
verify_kubectl $? "Deploying keptn authenticator component failed."

# Deploy control component
curl -o ../manifests/keptn/control.yaml https://raw.githubusercontent.com/keptn/control/$CONTROL_RELEASE/config/control.yaml

rm -f ../manifests/keptn/gen/control.yaml
cat ../manifests/keptn/control.yaml | \
  sed 's~CHANNEL_URI_PLACEHOLDER~'"$KEPTN_CHANNEL_URI"'~' >> ../manifests/keptn/gen/control.yaml
  
kubectl delete -f ../manifests/keptn/gen/control.yaml --ignore-not-found
kubectl apply -f ../manifests/keptn/gen/control.yaml
verify_kubectl $? "Deploying keptn control component failed."

# Set up SSL
ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
verify_variable "$ISTIO_INGRESS_IP" "ISTIO_INGRESS_IP is empty and could not be derived from the Istio ingress gateway." 

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$ISTIO_INGRESS_IP.xip.io"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
#verify_kubectl $? "Creating secret for istio-ingressgateway-certs failed."

kubectl get gateway knative-ingress-gateway --namespace knative-serving -o=yaml | yq w - spec.servers[1].tls.mode SIMPLE | yq w - spec.servers[1].tls.privateKey /etc/istio/ingressgateway-certs/tls.key | yq w - spec.servers[1].tls.serverCertificate /etc/istio/ingressgateway-certs/tls.crt | kubectl apply -f -
verify_kubectl $? "Updating knative ingress gateway with private key failed."

rm key.pem
rm certificate.pem

##############################################
## Start validation of keptn installation   ##
##############################################
wait_for_all_pods_in_namespace "keptn"

wait_for_deployment_in_namespace "event-broker" "keptn" # Wait function also waits for eventbroker-ext
wait_for_deployment_in_namespace "auth" "keptn"
wait_for_deployment_in_namespace "control" "keptn"
