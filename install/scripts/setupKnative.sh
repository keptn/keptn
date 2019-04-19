#!/bin/bash
CLUSTER_NAME=$1
CLUSTER_ZONE=$2

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

ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
cat ../manifests/knative/config-domain.yaml | \
  sed 's~ISTIO_INGRESS_IP_PLACEHOLDER~'"$ISTIO_INGRESS_IP"'~' >> ../manifests/gen/config-domain.yaml

kubectl apply -f ../manifests/gen/config-domain.yaml

# Determine the IP scope of the cluster (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  echo "[keptn|1]CLUSTER_IPV4_CIDR not set, retrieve it using gcloud"
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  echo "[keptn|1]SERVICES_IPV4_CIDR not set, retrieve it using gcloud"
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
fi

kubectl get configmap config-network -n knative-serving -o=yaml | yq w - data['istio.sidecar.includeOutboundIPRanges'] "$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR" | kubectl apply -f - 

echo "[keptn|0]Wait 30s for changes to apply... "
sleep 30

kubectl apply -f ../manifests/keptn/keptn-rbac.yaml
kubectl apply -f ../manifests/keptn/keptn-org-configmap.yaml

# Install kaniko build template
kubectl apply -f ../manifests/knative/build/kaniko.yaml -n keptn

# Create build-bot service account
kubectl apply -f ../manifests/knative/build/service-account.yaml
