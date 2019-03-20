#!/bin/bash

export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds.json | jq -r '.dynatraceApiToken')
export kubeID=$(curl "https://$DT_TENANT_ID.live.dynatrace.com/api/config/v1/virtualization/kubernetesConfigurations?api-token=$DT_API_TOKEN" | jq -r '.values | .[] | select(.name=="dt-acm-keptn") | .id')
export CLUSTERVERSION=$(curl -s https://$DT_TENANT_ID.live.dynatrace.com/api/v1/config/clusterversion?api-token=$DT_API_TOKEN | jq -r .version[2:5])

# check tenant is at least 1.164
if (( $(echo "$CLUSTERVERSION > 1.163" | bc -l) ))
then
curl -X DELETE \
  "https://$DT_TENANT_ID.live.dynatrace.com/api/config/v1/virtualization/kubernetesConfigurations/$kubeID?api-token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  }'
fi
