#!/bin/bash

echo "--------------------------"
echo "Configure Dynatrace Request Attributes"
echo "--------------------------"

export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds.json | jq -r '.dynatraceApiToken')

curl -X POST \
  "https://$DT_TENANT_ID/api/config/v1/requestAttributes?api-token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "name":"LSN",
    "enabled":true,
    "dataType":"STRING",
    "dataSources":[
        {
          "enabled":true,
          "source":"REQUEST_HEADER",
          "valueProcessing":{
            "splitAt":"",
            "trim":false,
            "extractSubstring":{
              "position":"BETWEEN",
              "delimiter":"LSN=",
              "endDelimiter":";"
            }
        },
        "parameterName":"x-dynatrace-test",
        "capturingAndStorageLocation":"CAPTURE_AND_STORE_ON_SERVER"
      }
    ],
    "normalization":"ORIGINAL",
    "aggregation":"FIRST",
    "confidential":false,
    "skipPersonalDataMasking":false
  }'

  curl -X POST \
    "https://$DT_TENANT_ID/api/config/v1/requestAttributes?Api-Token=$DT_API_TOKEN" \
    -H 'Content-Type: application/json' \
    -H 'cache-control: no-cache' \
    -d '{
      "name":"LTN",
      "enabled":true,
      "dataType":"STRING",
      "dataSources":[
          {
            "enabled":true,
            "source":"REQUEST_HEADER",
            "valueProcessing":{
              "splitAt":"",
              "trim":false,
              "extractSubstring":{
                "position":"BETWEEN",
                "delimiter":"LTN=",
                "endDelimiter":";"
              }
          },
          "parameterName":"x-dynatrace-test",
          "capturingAndStorageLocation":"CAPTURE_AND_STORE_ON_SERVER"
        }
      ],
      "normalization":"ORIGINAL",
      "aggregation":"FIRST",
      "confidential":false,
      "skipPersonalDataMasking":false
    }'

  curl -X POST \
    "https://$DT_TENANT_ID/api/config/v1/requestAttributes?Api-Token=$DT_API_TOKEN" \
    -H 'Content-Type: application/json' \
    -H 'cache-control: no-cache' \
    -d '{
      "name":"TSN",
      "enabled":true,
      "dataType":"STRING",
      "dataSources":[
          {
            "enabled":true,
            "source":"REQUEST_HEADER",
            "valueProcessing":{
              "splitAt":"",
              "trim":false,
              "extractSubstring":{
                "position":"BETWEEN",
                "delimiter":"TSN=",
                "endDelimiter":";"
              }
          },
          "parameterName":"x-dynatrace-test",
          "capturingAndStorageLocation":"CAPTURE_AND_STORE_ON_SERVER"
        }
      ],
      "normalization":"ORIGINAL",
      "aggregation":"FIRST",
      "confidential":false,
      "skipPersonalDataMasking":false
    }'
