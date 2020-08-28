#!/bin/bash

echo "Exposing Keptn Bridge..."
export BRIDGE_URL=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/bridge

# verify that bridge is available
curl "${BRIDGE_URL}" -k

if [[ $? != '0' ]]; then
    echo "Accessing Keptn Bridge failed"
    exit 1
fi
