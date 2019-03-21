#!/bin/bash

# Create a Bearer token for authenticating against the Kubernetes API
kubectl apply -f kubernetes-monitoring-service-account.yaml
sleep 60

# env variables
export DT_TENANT_ID=$(cat creds.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds.json | jq -r '.dynatraceApiToken')
export kubeURL=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
export kube_API_TOKEN=$(kubectl get secret $(kubectl get sa dynatrace-monitoring -o jsonpath='{.secrets[0].name}' -n dynatrace) -o jsonpath='{.data.token}' -n dynatrace | base64 --decode)
export CLUSTERVERSION=$(curl -s https://$DT_TENANT_ID.live.dynatrace.com/api/v1/config/clusterversion?api-token=$DT_API_TOKEN | jq -r .version[2:5])

# check tenant is at least 1.164
if (( $(echo "$CLUSTERVERSION > 1.163" | bc -l) ))
then
export CLUSTERVERSION=$(curl -s https://$DT_TENANT_ID.live.dynatrace.com/api/v1/config/clusterversion?api-token=$DT_API_TOKEN | jq -r .version)
curl -X POST \
  "https://$DT_TENANT_ID.live.dynatrace.com/api/config/v1/virtualization/kubernetesConfigurations?api-token=$DT_API_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "metadata": {
      "configurationVersions": [0],
      "clusterVersion": "'$CLUSTERVERSION'"

    },
    "label": "dt-acm-keptn",
    "endpointUrl": "'$kubeURL'",
    "authToken": "'$kube_API_TOKEN'"
  }'

fi
