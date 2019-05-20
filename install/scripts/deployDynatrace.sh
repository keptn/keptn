#!/bin/bash

source ./utils.sh

DT_TENANT=$(cat creds_dt.json | jq -r '.dynatraceTenant')
DT_API_TOKEN=$(cat creds_dt.json | jq -r '.dynatraceApiToken')
DT_PAAS_TOKEN=$(cat creds_dt.json | jq -r '.dynatracePaaSToken')

# Deploy Dynatrace operator
DT_OPERATOR_LATEST_RELEASE=$(curl -s https://api.github.com/repos/dynatrace/dynatrace-oneagent-operator/releases/latest | grep tag_name | cut -d '"' -f 4)
print_info "Installing Dynatrace Operator $DT_OPERATOR_LATEST_RELEASE"

kubectl create namespace dynatrace
verify_kubectl $? "Creating namespace dynatrace for oneagent operator failed."

kubectl label namespace dynatrace istio-injection=disabled

kubectl apply -f https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$DT_OPERATOR_LATEST_RELEASE/deploy/kubernetes.yaml
verify_kubectl $? "Applying Dynatrace operator failed."
wait_for_crds "oneagent"

# Create Dynatrace secret
kubectl -n dynatrace create secret generic oneagent --from-literal="apiToken=$DT_API_TOKEN" --from-literal="paasToken=$DT_PAAS_TOKEN"
verify_kubectl $? "Creating secret for Dynatrace OneAgent failed."

# Create Dynatrace OneAgent
rm -f ../manifests/gen/cr.yml

curl -o ../manifests/dynatrace/cr.yml https://raw.githubusercontent.com/Dynatrace/dynatrace-oneagent-operator/$DT_OPERATOR_LATEST_RELEASE/deploy/cr.yaml
cat ../manifests/dynatrace/cr.yml | sed 's~ENVIRONMENTID.live.dynatrace.com~'"$DT_TENANT"'~' >> ../manifests/gen/cr.yml

kubectl apply -f ../manifests/gen/cr.yml
verify_kubectl $? "Creating Dynatrace OneAgent failed."

# Apply auto tagging rules in Dynatrace
print_info "Applying auto tagging rules in Dynatrace."
./applyAutoTaggingRules.sh $DT_TENANT $DT_API_TOKEN
verify_install_step $? "Applying auto tagging rules in Dynatrace failed."
print_info "Applying auto tagging rules in Dynatrace done."

# Create secrets to be used by dynatrace-service
kubectl -n keptn create secret generic dynatrace --from-literal="DT_API_TOKEN=$DT_API_TOKEN" --from-literal="DT_TENANT=$DT_TENANT"
verify_kubectl $? "Creating dynatrace secret for keptn services failed."

# Create dynatrace-service
DT_SERVICE_RELEASE="0.1.0"
print_info "Deploying dynatrace-service $DT_SERVICE_RELEASE"
kubectl apply -f https://raw.githubusercontent.com/keptn/dynatrace-service/$DT_SERVICE_RELEASE/dynatrace-service.yaml
verify_kubectl $? "Deploying dynatrace-service failed."
