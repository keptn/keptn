#!/bin/bash
CLUSTER_NAME=$1
CLUSTER_ZONE=$2

# Determine the IP scope of the cluster (https://github.com/knative/docs/blob/master/serving/outbound-network-access.md)
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  echo "[keptn|1]CLUSTER_IPV4_CIDR not set, retrieve it using gcloud"
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  echo "[keptn|1]SERVICES_IPV4_CIDR not set, retrieve it using gcloud"
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
fi 

# Apply custom resource definitions for Istio
kubectl apply -f ../manifests/istio/istio-crds-knative.yaml

# Wait max 1min for custom resource virtualservice to be available
sleep 2
RETRY=0
while [ $RETRY -lt 6 ]
do
  kubectl get virtualservice
  if [[ $? == '0' ]]
  then
    echo "[keptn|0]CRD virtualservice now available, can continue... "
    break
  fi
  RETRY=$[$RETRY+1]
  echo "[keptn|0]Retry: ${RETRY}/6 - Wait 10s for changes to apply... "
  sleep 10
done

# Wait max 1min for custom resource handler to be available
sleep 2
RETRY=0
while [ $RETRY -lt 6 ]
do
  kubectl get handler
  if [[ $? == '0' ]]
  then
    echo "[keptn|0]CRD handler now available, can continue... "
    break
  fi
  RETRY=$[$RETRY+1]
  echo "[keptn|0]Retry: ${RETRY}/6 - Wait 10s for changes to apply... "
  sleep 10
done

# Apply Istio configuration
rm -f ../manifests/gen/istio-knative.yaml
cat ../manifests/istio/istio-knative.yaml | \
  sed 's~INCLUDE_OUTBOUND_IP_RANGES_PLACEHOLDER~'"$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR"'~' >> ../manifests/gen/istio-knative.yaml

kubectl apply -f ../manifests/gen/istio-knative.yaml

echo "Wait 4 minutes for changes to apply... "
sleep 240

# Delete all pods in keptn to apply Istio changes
kubectl delete pods --all -n keptn

##############################################
## Start validation of Istio installation   ##
##############################################

# Wait max 4min for IP of Istio ingressgateway
# sleep 2
# RETRY=0
# while [ $RETRY -lt 24 ]
# do
#   ISTIO_INGRESS_IP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
#   if [[ -z "${ISTIO_INGRESS_IP}" ]]
#   then
#     echo "[keptn|1]IP of Istio ingressgateway: ${ISTIO_INGRESS_IP}"
#     echo "[keptn|0]IP of Istio ingressgateway available, can continue... "
#     break
#   fi
#   RETRY=$[$RETRY+1]
#   echo "[keptn|0]Retry: ${RETRY}/24 - Wait 10s for changes to apply... "
#   sleep 10
# done
