#!/bin/bash

source test/utils.sh

# authenticate at Keptn API
KEPTN_API_URL=$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

auth_at_keptn $KEPTN_API_URL $KEPTN_API_TOKEN
