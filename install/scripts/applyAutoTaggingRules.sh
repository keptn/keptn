#!/bin/bash

source ./utils.sh

DT_TENANT_ID=$1
DT_API_TOKEN=$2

curl -X POST \
  "https://$DT_TENANT_ID.live.dynatrace.com/api/config/v1/autoTags?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "name": "service",
  "rules": [
    {
      "type": "SERVICE",
      "enabled": true,
      "valueFormat": "{ProcessGroup:KubernetesContainerName}",
      "propagationTypes": [],
      "conditions": [
        {
          "key": {
            "attribute": "PROCESS_GROUP_PREDEFINED_METADATA",
            "dynamicKey": "KUBERNETES_CONTAINER_NAME",
            "type": "PROCESS_PREDEFINED_METADATA_KEY"
          },
          "comparisonInfo": {
            "type": "STRING",
            "operator": "EXISTS",
            "value": null,
            "negate": false,
            "caseSensitive": null
          }
        }
      ]
    }
  ]
}'

if [[ $1 != '0' ]]; then
  echo ""
  print_error "Tagging rule for service could not be created in tenant $DT_TENANT_ID."
  exit 1
fi

curl -X POST \
  "https://$DT_TENANT_ID.live.dynatrace.com/api/config/v1/autoTags?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name": "environment",
  "rules": [
    {
      "type": "SERVICE",
      "enabled": true,
      "valueFormat": "{ProcessGroup:KubernetesNamespace}",
      "propagationTypes": [],
      "conditions": [
        {
          "key": {
            "attribute": "PROCESS_GROUP_PREDEFINED_METADATA",
            "dynamicKey": "KUBERNETES_NAMESPACE",
            "type": "PROCESS_PREDEFINED_METADATA_KEY"
          },
          "comparisonInfo": {
            "type": "STRING",
            "operator": "EXISTS",
            "value": null,
            "negate": false,
            "caseSensitive": null
          }
        }
      ]
    }
  ]
}'

if [[ $1 != '0' ]]; then
  echo ""
  print_error "Tagging rule for environment could not be created in tenant $DT_TENANT_ID."
  exit 1
fi
