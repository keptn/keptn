#!/bin/bash
REGISTRY_URL=$1
CLUSTER_NAME=$2
CLUSTER_ZONE=$3

kubectl create namespace keptn 2> /dev/null

# Create container registry
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml

# Install knative serving, eventing, build
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.5.0/serving.yaml
kubectl apply --filename https://github.com/knative/build/releases/download/v0.5.0/build.yaml
kubectl apply --filename https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml
kubectl apply --filename https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.5.0/monitoring.yaml
kubectl apply --filename https://raw.githubusercontent.com/knative/serving/v0.5.0/third_party/config/build/clusterrole.yaml

# Configure knative serving default domain
rm -f ../manifests/gen/config-domain.yaml

export ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
cat ../manifests/knative/config-domain.yaml | \
  sed 's~ISTIO_INGRESS_IP_PLACEHOLDER~'"$ISTIO_INGRESS_IP"'~' >> ../manifests/gen/config-domain.yaml

kubectl apply -f ../manifests/gen/config-domain.yaml

# Determine the IP scope of the cluster (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
# Gcloud:
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
fi 

kubectl get configmap config-network -n knative-serving -o=yaml | yq w - data['istio.sidecar.includeOutboundIPRanges'] "$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR" | kubectl apply -f - 

echo "Wait 30s for changes to apply... "
sleep 30

kubectl apply -f ../manifests/keptn/keptn-rbac.yaml
kubectl apply -f ../manifests/keptn/keptn-org-configmap.yaml

# Install kaniko build template
kubectl apply -f ../manifests/knative/build/kaniko.yaml -n keptn

# Create build-bot service account
kubectl apply -f ../manifests/knative/build/service-account.yaml

REGISTRY_URL=$(kubectl describe svc docker-registry -n keptn | grep "IP:" | sed 's~IP:[ \t]*~~')

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

export KEPTN_CHANNEL_URI=$(kubectl describe channel keptn-channel -n keptn | grep "Hostname:" | sed 's~[ \t]*Hostname:[ \t]*~~')

export KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)
kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"

# Deploy event broker
cd ../../core/eventbroker
chmod +x deploy.sh
./deploy.sh $REGISTRY_URL $KEPTN_CHANNEL_URI $NEW_ARTEFACT_CHANNEL $START_DEPLOYMENT_CHANNEL $DEPLOYMENT_FINISHED_CHANNEL $START_TESTS_CHANNEL $TESTS_FINISHED_CHANNEL $START_EVALUATION_CHANNEL $EVALUATION_DONE_CHANNEL
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
openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$ISTIO_INGRESS_IP.xip.io"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem

kubectl get gateway knative-ingress-gateway --namespace knative-serving -o=yaml | yq w - spec.servers[1].tls.mode SIMPLE | yq w - spec.servers[1].tls.privateKey /etc/istio/ingressgateway-certs/tls.key | yq w - spec.servers[1].tls.serverCertificate /etc/istio/ingressgateway-certs/tls.crt | kubectl apply -f -

rm key.pem
rm certificate.pem
