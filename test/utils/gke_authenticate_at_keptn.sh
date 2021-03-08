#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

# authenticate at Keptn API
KEPTN_API_URL="http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/api"
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -o jsonpath='{.data.keptn-api-token}' | base64 --decode)

auth_at_keptn "$KEPTN_API_URL" "$KEPTN_API_TOKEN"
