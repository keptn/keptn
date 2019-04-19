#!/bin/bash

export DT_TENANT_ID=$(cat creds_dt.json | jq -r '.dynatraceTenant')
export DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')
export DT_PAAS_TOKEN=$(cat creds_dt.json | jq -r '.dynatracePaaSToken')

# Deploy Dynatrace operator
kubectl create namespace dynatrace
kubectl label namespace dynatrace istio-injection=disabled

export LATEST_RELEASE=$(curl -s https://api.github.com/repos/dynatrace/dynatrace-oneagent-operator/releases/latest | grep tag_name | cut -d '"' -f 4)
echo "Installing Dynatrace Operator $LATEST_RELEASE"

kubectl apply -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/kubernetes.yaml

# Wait for custom resource OneAgent to be available
SLEEP_TIME=1
SLEEP_ROUND=1

while [ $SLEEP_TIME -lt 100 ]
do
  kubectl get OneAgent
  if [[ $? == '0' ]]
  then
    echo "CRD OneAgent now available, can continue... "
    break
  fi
  SLEEP_TIME=$[$SLEEP_ROUND*$SLEEP_ROUND]
  SLEEP_ROUND=$[$SLEEP_ROUND+1]
  echo "Wait ${SLEEP_TIME}s for changes to apply... "
  sleep ${SLEEP_TIME}
done

# Create Dynatrace OneAgent
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"

rm -f ../manifests/gen/cr.yml

curl -o ../manifests/dynatrace/cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/cr.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/gen/cr.yml

kubectl apply -f ../manifests/gen/cr.yml

# Apply auto tagging rules in Dynatrace
echo "Apply auto tagging rules in Dynatrace "
./applyAutoTaggingRules.sh $DT_TENANT_ID $DT_API_TOKEN
echo "End applying auto tagging rules in Dynatrace "

./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN

# Create secrets to be used by keptn services
kubectl -n keptn create secret generic dynatrace --from-literal="DT_API_TOKEN=$DT_API_TOKEN" --from-literal="DT_TENANT_ID=$DT_TENANT_ID"
