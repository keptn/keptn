#!/bin/bash

export DT_TENANT_ID=$(cat creds_dt.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')
export DT_PAAS_TOKEN=$(cat creds_dt.json | jq -r '.dynatracePaaSToken')
export DT_TENANT_URL="$DT_TENANT_ID.live.dynatrace.com"

# Deploy Dynatrace operator
kubectl create namespace dynatrace
kubectl label namespace dynatrace istio-injection=disabled
export LATEST_RELEASE=$(curl -s https://api.github.com/repos/dynatrace/dynatrace-oneagent-operator/releases/latest | grep tag_name | cut -d '"' -f 4)
echo "Installing Dynatrace Operator $LATEST_RELEASE"
kubectl create -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/kubernetes.yaml
sleep 60
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"

rm -f ../manifests/gen/cr.yml

curl -o ../manifests/dynatrace/cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/cr.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/gen/cr.yml

kubectl create -f ../manifests/gen/cr.yml

# Apply auto tagging rules in Dynatrace
echo "--------------------------"
echo "Apply auto tagging rules in Dynatrace "
echo "--------------------------"

./applyAutoTaggingRules.sh $DT_TENANT_ID $DT_API_TOKEN

echo "--------------------------"
echo "End applying auto tagging rules in Dynatrace "
echo "--------------------------"

./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN