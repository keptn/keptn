#!/bin/bash

source ./utils.sh

DT_TENANT_ID=$(cat creds_dt.json | jq -r '.dynatraceTenant')
DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')
DT_PAAS_TOKEN=$(cat creds_dt.json | jq -r '.dynatracePaaSToken')

# Deploy Dynatrace operator
LATEST_RELEASE=$(curl -s https://api.github.com/repos/dynatrace/dynatrace-oneagent-operator/releases/latest | grep tag_name | cut -d '"' -f 4)
print_info "Installing Dynatrace Operator $LATEST_RELEASE"

kubectl create namespace dynatrace
verify_kubectl $? "Creating namespace dynatrace for oneagent operator failed."

kubectl label namespace dynatrace istio-injection=disabled

kubectl apply -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/kubernetes.yaml
verify_kubectl $? "Applying Dynatrace operator failed."
wait_for_crds "oneagent"

# Create Dynatrace secret
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"
verify_kubectl $? "Creating secret for Dynatrace OneAgent failed."

# Create Dynatrace OneAgent
rm -f ../manifests/gen/cr.yml

curl -o ../manifests/dynatrace/cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/cr.yml | sed 's/ENVIRONMENTID/'"$DT_TENANT_ID"'/' >> ../manifests/gen/cr.yml

kubectl apply -f ../manifests/gen/cr.yml
verify_kubectl $? "Creating Dynatrace OneAgent failed."

# Apply auto tagging rules in Dynatrace
print_info "Applying auto tagging rules in Dynatrace."
./applyAutoTaggingRules.sh $DT_TENANT_ID $DT_API_TOKEN
verify_install_step $? "Applying auto tagging rules in Dynatrace failed."
print_info "Applying auto tagging rules in Dynatrace done."

print_info "Creating service entries for Dynatrace OneAgent."
./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN
verify_install_step $? "Creating service entries for Dynatrace OneAgent failed."
print_info "Creating service entries for Dynatrace OneAgent done."

# Create secrets to be used by keptn services
kubectl -n keptn create secret generic dynatrace --from-literal="DT_API_TOKEN=$DT_API_TOKEN" --from-literal="DT_TENANT_ID=$DT_TENANT_ID"
verify_kubectl $? "Creating dynatrace secret for keptn services failed."
