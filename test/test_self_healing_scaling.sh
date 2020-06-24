#!/bin/bash

source test/utils.sh

function cleanup() {
  kubectl delete namespace loadgen
}

trap cleanup EXIT

# get keptn api details
KEPTN_ENDPOINT=https://api.keptn.$(kubectl get cm keptn-domain -n keptn -ojsonpath={.data.app_domain})
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

# test configuration
PROJECT="sockshop"
SERVICE="carts"
STAGE="production"

PROMETHEUS_SERVICE_VERSION=${UNLEASH_SERVICE_VERSION:-0.3.4}

kubectl delete namespace sockshop-dev
kubectl delete namespace sockshop-staging
kubectl delete namespace sockshop-production
keptn delete project sockshop

keptn create project sockshop --shipyard=./test/assets/shipyard_self_healing_scale.yaml

# Prerequisites
rm -rf examples
git clone --branch master https://github.com/keptn/examples --single-branch

cd examples/onboarding-carts

keptn onboard service carts --project=sockshop --chart=./carts
keptn onboard service carts-db --project=sockshop --chart=./carts-db --deployment-strategy=direct
keptn send event new-artifact --project=sockshop --service=carts-db --image=mongo

keptn add-resource --project=sockshop --service=carts --stage=dev --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
keptn add-resource --project=sockshop --service=carts --stage=dev --resource=jmeter/load.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=sockshop --service=carts --stage=staging --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
keptn add-resource --project=sockshop --service=carts --stage=staging --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/load.jmx

keptn add-resource --project=sockshop --service=carts --stage=production --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx

cd ../..

# add remediation.yaml
keptn add-resource --project=sockshop --stage=production --service=carts --resource=./test/assets/self_healing_scaling_remediation.yaml --resourceUri=remediation.yaml
# add slo file
keptn add-resource --project=sockshop --service=carts --stage=production --resource=./test/assets/self_healing_slo.yaml --resourceUri=slo.yaml

keptn send event new-artifact --project=sockshop --service=carts --image=docker.io/keptnexamples/carts --tag=0.11.1

sleep 200

wait_for_deployment_in_namespace carts-primary sockshop-production

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-$PROMETHEUS_SERVICE_VERSION/deploy/service.yaml

wait_for_deployment_in_namespace prometheus-service keptn
echo "Prometheus service deployed successfully"

keptn configure monitoring prometheus --project=sockshop --service=carts

wait_for_deployment_in_namespace prometheus-deployment monitoring
echo "Prometheus deployed successfully"

cd examples/load-generation/cartsloadgen

kubectl apply -f deploy/cartsloadgen-faulty.yaml

wait_for_deployment_in_namespace cartsloadgen loadgen

echo "loadgen deployed successfully waiting for problem notification"

event=$(wait_for_problem_open_event ${PROJECT} ${SERVICE} ${STAGE})

echo $event
verify_using_jq "$event" ".source" "prometheus"
verify_using_jq "$event" ".data.project" $PROJECT
verify_using_jq "$event" ".data.stage" "production"
verify_using_jq "$event" ".data.service" "$SERVICE"

keptn_context_id=$(echo $event | jq -r '.shkeptncontext')

sleep 20

## Check remediation.triggered event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.triggered&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "remediation-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.problem.ProblemTitle" "response_time_p90"
verify_using_jq "$response" ".data.problem.State" "OPEN"


## Check remediation.status.changed event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.remediation.status.changed&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "remediation-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.remediation.result.actionIndex" "0"
verify_using_jq "$response" ".data.remediation.result.actionName" "scaling"


## Check action.triggered event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.triggered&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "remediation-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"
verify_using_jq "$response" ".data.action.action" "scaling"
verify_using_jq "$response" ".data.action.value" "1"

## Check action.started event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.started&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "helm-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"

# wait for the remediation action to be finished
sleep 160

## Check action.finished event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "helm-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"

replicacount=$(kubectl get deployment -n sockshop-production carts-primary -ojsonpath='{.spec.replicas}')
if [[ "$replicacount" != "2" ]]; then
  echo "number of replicas has not been increased"
  exit 2
else
  echo "Verified that number of replicas has been increased"
fi
