#!/bin/bash
REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep IP: | sed 's~IP:[ \t]*~~')

# Mark internal docker registry as insecure registry for knative controller
val=$(kubectl -n knative-serving get cm config-controller -o=json | jq -r .data.registriesSkippingTagResolving | awk '{print $1",'$REGISTRY_URL':5000"}')
kubectl -n knative-serving get cm config-controller -o=yaml | yq w - data.registriesSkippingTagResolving $val | kubectl apply -f -

# Deploy knative eventing channels (keptn-channel)
kubectl apply -f ../../core/eventbroker/config/channel.yaml
kubectl apply -f ../../core/eventbroker/config/new-artefact-channel.yaml
kubectl apply -f ../../core/eventbroker/config/configuration-changed-channel.yaml
kubectl apply -f ../../core/eventbroker/config/deployment-finished-channel.yaml
kubectl apply -f ../../core/eventbroker/config/tests-finished-channel.yaml
kubectl apply -f ../../core/eventbroker/config/evaluation-done-channel.yaml
kubectl apply -f ../../core/eventbroker/config/problem-channel.yaml

KEPTN_CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')
KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)

kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"

# Deploy event broker
cd ../../core/eventbroker
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $KEPTN_CHANNEL_URI
cd ../../install/scripts

cd ../../core/eventbroker-ext
chmod +x deploy.sh
./deploy.sh
cd ../../install/scripts

cd ../../core/auth
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL
cd ../../install/scripts

cd ../../core/control
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $KEPTN_CHANNEL_URI
cd ../../install/scripts

# Set up SSL
ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')

openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$ISTIO_INGRESS_IP.xip.io"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
kubectl get gateway knative-ingress-gateway --namespace knative-serving -o=yaml | yq w - spec.servers[1].tls.mode SIMPLE | yq w - spec.servers[1].tls.privateKey /etc/istio/ingressgateway-certs/tls.key | yq w - spec.servers[1].tls.serverCertificate /etc/istio/ingressgateway-certs/tls.crt | kubectl apply -f -

rm key.pem
rm certificate.pem
