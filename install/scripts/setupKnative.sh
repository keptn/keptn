#!/bin/bash
CLUSTER_IPV4_CIDR=$1
SERVICES_IPV4_CIDR=$2

# Will be removed
kubectl create namespace keptn 2> /dev/null

# Create container registry
# Will be removed
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml

# Install knative serving, eventing, build
kubectl apply -f https://github.com/knative/serving/releases/download/v0.5.0/serving.yaml
#TODO: return value 0?
#TODO: check namespaces (serving, eventing, build)
kubectl apply -f https://github.com/knative/build/releases/download/v0.5.0/build.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml
kubectl apply -f https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml
kubectl apply -f https://github.com/knative/serving/releases/download/v0.5.0/monitoring.yaml
kubectl apply -f https://raw.githubusercontent.com/knative/serving/v0.5.0/third_party/config/build/clusterrole.yaml

# Wait max 1min for custom resource virtualservice to be available
RETRY=0
while [ $RETRY -lt 6 ]
do
  # TODO check all crds
  kubectl get virtualservice,handler,......
  if [[ $? == '0' ]]
  then
    echo "[keptn|INFO] CRD virtualservice now available, can continue... "
    break
  fi
  RETRY=$[$RETRY+1]
  echo "[keptn|INFO] Retry: ${RETRY}/6 - Wait 10s for changes to apply... "
  sleep 10
done

# Configure knative serving default domain
rm -f ../manifests/gen/config-domain.yaml

ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
cat ../manifests/knative/config-domain.yaml | \
  sed 's~ISTIO_INGRESS_IP_PLACEHOLDER~'"$ISTIO_INGRESS_IP"'~' >> ../manifests/gen/config-domain.yaml

kubectl apply -f ../manifests/gen/config-domain.yaml

# Determine the IP scope of the cluster (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  echo "[keptn|DEBUG] CLUSTER_IPV4_CIDR not set, retrieve it using gcloud"
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  echo "[keptn|DEBUG] SERVICES_IPV4_CIDR not set, retrieve it using gcloud"
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
fi

kubectl get configmap config-network -n knative-serving -o=yaml | yq w - data['istio.sidecar.includeOutboundIPRanges'] "$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR" | kubectl apply -f - 

echo "[keptn|INFO] Wait 30s for changes to apply... "
sleep 30

kubectl apply -f ../manifests/keptn/keptn-rbac.yaml
kubectl apply -f ../manifests/keptn/keptn-org-configmap.yaml

# Create build-bot service account
kubectl apply -f ../manifests/knative/build/service-account.yaml

##############################################
## Start validation of Knative installation ##
##############################################

# Wait max 2min for webhook deployment
RETRY=0
while [ $RETRY -lt 12 ]
do
  kubectl rollout status deployment webhook -n knative-eventing --timeout=10s
  if [[ $? == '0' ]]
  then
    echo "[keptn|INFO] Deployment webhook in knative-eventing namespace available, can continue... "
    break
  fi
  RETRY=$[$RETRY+1]
done

# Install kaniko build template
kubectl apply -f ../manifests/knative/build/kaniko.yaml -n keptn
