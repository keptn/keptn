#!/bin/bash

source test/utils.sh

KEPTN_EXAMPLES_BRANCH=${KEPTN_EXAMPLES_BRANCH:-"master"}

function cleanup() {
  kubectl delete namespace loadgen
}

trap cleanup EXIT


KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)


# test configuration
PROJECT="sockshop"
SERVICE="carts"
STAGE="production"

PROMETHEUS_SERVICE_VERSION=${PROMETHEUS_SERVICE_VERSION:-master}

kubectl delete namespace $PROJECT-production
keptn delete project $PROJECT
keptn create project $PROJECT --shipyard=./test/assets/shipyard_self_healing_scale.yaml

########################################################################################################################
# Pre-requisites
########################################################################################################################

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/$PROMETHEUS_SERVICE_VERSION/deploy/service.yaml

wait_for_deployment_in_namespace prometheus-service keptn
wait_for_deployment_in_namespace prometheus-service-monitoring-configure-distributor keptn
echo "Prometheus service deployed successfully"

rm -rf examples
git clone --branch ${KEPTN_EXAMPLES_BRANCH} https://github.com/keptn/examples --single-branch

cd examples/onboarding-$SERVICE

###########################################
# onboard carts                           #
###########################################
keptn onboard service $SERVICE --project=$PROJECT --chart=./$SERVICE

###########################################
# onboard carts-db                        #
###########################################
keptn onboard service $SERVICE-db --project=$PROJECT --chart=./$SERVICE-db --deployment-strategy=direct
keptn send event new-artifact --project=$PROJECT --service=$SERVICE-db --image=mongo --sequence=artifact-delivery-db

# add health check in production
keptn add-resource --project=$PROJECT --service=$SERVICE --stage=production --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx

cd ../..

# add remediation.yaml
keptn add-resource --project=$PROJECT --stage=production --service=$SERVICE --resource=./test/assets/self_healing_scaling_remediation.yaml --resourceUri=remediation.yaml
# add slo file
keptn add-resource --project=$PROJECT --service=$SERVICE --stage=production --resource=./test/assets/self_healing_slo.yaml --resourceUri=slo.yaml

# deploy the service
keptn send event new-artifact --project=$PROJECT --service=$SERVICE --image=docker.io/keptnexamples/$SERVICE --tag=0.11.1 --sequence=artifact-delivery

echo "It might take a while for the service to be available on production - waiting 50sec"
sleep 50
echo "Still waiting 50sec ..."
sleep 50

wait_for_deployment_in_namespace $SERVICE-primary $PROJECT-$STAGE

###########################################
# set up prometheus monitoring            #
###########################################

keptn configure monitoring prometheus --project=$PROJECT --service=$SERVICE

wait_for_deployment_in_namespace prometheus-deployment monitoring
echo "Prometheus deployed successfully"

###########################################
# generate load on the service            #
###########################################
cd examples/load-generation/cartsloadgen

kubectl apply -f deploy/cartsloadgen-faulty.yaml
wait_for_deployment_in_namespace cartsloadgen loadgen
echo "Loadgen deployed successfully - waiting for problem notification"

sleep 120
echo "Still waiting 120sec ..."
sleep 120
echo "Still waiting 120sec ..."
sleep 120

event=$(wait_for_problem_open_event ${PROJECT} ${SERVICE} ${STAGE})

echo $event
verify_using_jq "$event" ".source" "prometheus"
verify_using_jq "$event" ".data.project" $PROJECT
verify_using_jq "$event" ".data.stage" "$STAGE"
verify_using_jq "$event" ".data.service" "$SERVICE"

keptn_context_id=$(echo $event | jq -r '.shkeptncontext')

sleep 20

## check remediation.triggered event ##
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


## check remediation.status.changed event ##
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

## check action.triggered event ##
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

## check action.started event ##
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

## check action.finished event ##
response=$(curl -X GET "${KEPTN_ENDPOINT}/mongodb-datastore/event?project=${PROJECT}&type=sh.keptn.event.action.finished&keptnContext=${keptn_context_id}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.events[0]')

# print the response
echo $response | jq .

# validate the response
verify_using_jq "$response" ".source" "helm-service"
verify_using_jq "$response" ".data.project" "$PROJECT"
verify_using_jq "$response" ".data.stage" "$STAGE"
verify_using_jq "$response" ".data.service" "$SERVICE"

replicacount=$(kubectl get deployment -n $PROJECT-$STAGE $SERVICE-primary -ojsonpath='{.spec.replicas}')
if [[ "$replicacount" != "2" ]]; then
  echo "number of replicas has not been increased"
  exit 2
else
  echo "Verified that number of replicas has been increased"
fi

echo "Self healing tests for scaling done ✓"
