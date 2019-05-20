#!/bin/bash

source ./utils.sh

# API documentation
# https://www.dynatrace.com/support/help/extend-dynatrace/dynatrace-api/configuration/auto-tag-api/

DT_TENANT=$1
DT_API_TOKEN=$2


DT_RULE_NAME=service
# check if rule already exists in Dynatrace tenant
export DT_ID=
export DT_ID=$(curl -X GET \
  "https://$DT_TENANT/api/config/v1/autoTags?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  | jq -r '.values[] | select(.name == "'$DT_RULE_NAME'") | .id')

# if exists, then delete it
if [ "$DT_ID" != "" ]
then
  echo "Removing $DT_RULE_NAME since exists. Replacing with new definition."
  curl -f -X DELETE \
  "https://$DT_TENANT/api/config/v1/autoTags/$DT_ID?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'

  if [[ $? -ne 0 ]]; then
    print_error "Tagging rule: $DT_RULE_NAME could not be deleted in tenant $DT_TENANT_ID."
    exit 1
  fi
fi




curl -f -X POST \
  "https://$DT_TENANT/api/config/v1/autoTags?api-token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "name": "'$DT_RULE_NAME'",
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

if [[ $? != '0' ]]; then
  echo ""
  print_error "Tagging rule for service could not be created in tenant $DT_TENANT."
  exit 1
fi


DT_RULE_NAME=environment
# check if rule already exists in Dynatrace tenant
export DT_ID=
export DT_ID=$(curl -X GET \
  "https://$DT_TENANT/api/config/v1/autoTags?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  | jq -r '.values[] | select(.name == "'$DT_RULE_NAME'") | .id')

# if exists, then delete it
if [ "$DT_ID" != "" ]
then
  echo "Removing $DT_RULE_NAME since exists. Replacing with new definition."
  curl -f -X DELETE \
  "https://$DT_TENANT/api/config/v1/autoTags/$DT_ID?Api-Token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'

  if [[ $? -ne 0 ]]; then
    print_error "Tagging rule: $DT_RULE_NAME could not be deleted in tenant $DT_TENANT_ID."
    exit 1
  fi
fi

curl -f -X POST \
  "https://$DT_TENANT/api/config/v1/autoTags?api-token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"name": "'$DT_RULE_NAME'",
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

if [[ $? != '0' ]]; then
  echo ""
  print_error "Tagging rule for environment could not be created in tenant $DT_TENANT."
  exit 1
fi
