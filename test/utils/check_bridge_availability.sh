#!/bin/bash

echo "Exposing Keptn Bridge ..."
# shellcheck disable=SC2155
export BRIDGE_URL=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/bridge

# verify that bridge is available
if curl "${BRIDGE_URL}" -k; then
    echo "Accessing Keptn Bridge failed"
    exit 1
fi
