#!/bin/bash

CLUSTER_NAME=$1
CLUSTER_ZONE=$2

# Determine the IP scope of the cluster (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
# Gcloud:
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
fi 

rm -f ../manifests/gen/istio-knative.yaml

cat ../manifests/istio/istio-knative.yaml | \
  sed 's~INCLUDE_OUTBOUND_IP_RANGES_PLACEHOLDER~'"$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR"'~' >> ../manifests/gen/istio-knative.yaml

kubectl apply -f ../manifests/istio/istio-crds-knative.yaml

echo "Wait 4 minutes for changes to apply... "
sleep 240

kubectl apply -f ../manifests/gen/istio-knative.yaml

echo "Wait 10s for changes to apply... "
sleep 10

kubectl delete pods --all -n keptn
