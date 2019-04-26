#!/bin/bash
CLUSTER_IPV4_CIDR=$1
SERVICES_IPV4_CIDR=$2

# Apply custom resource definitions for Istio
kubectl apply -f ../manifests/istio/istio-crds-knative.yaml
# TODO: check return value

# Wait max 1min for custom resource virtualservice to be available
RETRY=0
while [ $RETRY -lt 6 ]
do
  # TODO check all crds
  kubectl get virtualservice,handler,...
  if [[ $? == '0' ]]
  then
    echo "[keptn|INFO] All custom resource definitions now available, continue installation. "
    break
  fi
  RETRY=$[$RETRY+1]
  echo "[keptn|INFO] Retry: ${RETRY}/6 - Wait 10s for changes to apply... "
  sleep 10
done

# Apply Istio configuration
rm -f ../manifests/gen/istio-knative.yaml
cat ../manifests/istio/istio-knative.yaml | \
  sed 's~INCLUDE_OUTBOUND_IP_RANGES_PLACEHOLDER~'"$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR"'~' >> ../manifests/gen/istio-knative.yaml

kubectl apply -f ../manifests/gen/istio-knative.yaml
# TODO: check return value

# TODO: script from Florian
echo "Wait 4 minutes for changes to apply... "
sleep 240

# Delete all pods in keptn to apply Istio changes
kubectl delete pods --all -n keptn
# TODO: validate result

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
#     echo "[keptn|DEBUG] IP of Istio ingressgateway: ${ISTIO_INGRESS_IP}"
#     echo "[keptn|INFO] IP of Istio ingressgateway available, can continue... "
#     break
#   fi
#   RETRY=$[$RETRY+1]
#   echo "[keptn|INFO] Retry: ${RETRY}/24 - Wait 10s for changes to apply... "
#   sleep 10
# done
