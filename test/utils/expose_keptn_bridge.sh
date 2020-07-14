#!/bin/bash

echo "Exposing Keptn Bridge..."
kubectl -n keptn delete secret bridge-credentials
# expose Keptn's Bridge (for easier troubleshooting/debugging afterwards)
kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=${NIGHTLY_BRIDGE_USERNAME}" --from-literal="BASIC_AUTH_PASSWORD=${NIGHTLY_BRIDGE_PASSWORD}"
# restart bridge pod to make use of the secret
kubectl -n keptn delete pods --selector=run=bridge
sleep 5
kubectl patch svc bridge -n keptn -p '{"spec": {"type": "LoadBalancer"}}'
sleep 10
export BRIDGE_URL=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/bridge
# verify that bridge is available
curl "${BRIDGE_URL}" -k

if [[ $? != '0' ]]; then
    echo "Accessing Keptn Bridge failed"
    exit 1
fi
