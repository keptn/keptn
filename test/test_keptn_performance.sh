#!/bin/bash

source test/utils.sh

function cleanup() {
  echo "Executing cleanup..."

  echo "Delete lighthouse-config configmap"
  kubectl delete configmap -n ${KEPTN_NAMESPACE} lighthouse-config

  echo "Deleting project ${SELF_MONITORING_PROJECT}"
  keptn delete project $SELF_MONITORING_PROJECT

  echo "Uninstalling dynatrace-sli-service"
  kubectl -n ${KEPTN_NAMESPACE} delete deployment dynatrace-sli-service

  echo "Removing secret dynatrace-credentials-${SELF_MONITORING_PROJECT}"
  kubectl -n ${KEPTN_NAMESPACE} delete secret dynatrace-credentials-${SELF_MONITORING_PROJECT}

  for project_nr in $(seq 1 ${NR_PROJECTS})
  do
    keptn delete project project-${project_nr}
  done
}
#trap cleanup EXIT SIGINT

# test configuration
DYNATRACE_SLI_SERVICE_VERSION=${DYNATRACE_SLI_SERVICE_VERSION:-master}
SELF_MONITORING_PROJECT=${SELF_MONITORING_PROJECT:-keptn-selfmonitoring}
KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}

NR_PROJECTS=${NR_PROJECTS:-20}
NR_SERVICES_PER_PROJECT=${NR_SERVICES_PER_PROJECT:-15}
NR_EVALUATIONS_PER_SERVICE=${NR_EVALUATIONS_PER_SERVICE:-100}

if [[ $QG_INTEGRATION_TEST_DT_TENANT == "" ]]; then
  echo "No DT Tenant env var provided. Exiting."
  exit 1
fi

if [[ $QG_INTEGRATION_TEST_DT_API_TOKEN == "" ]]; then
  echo "No DZ API Token env var provided. Exiting."
  exit 1
fi

cleanup

# get keptn API details
if [[ "$PLATFORM" == "openshift" ]]; then
  KEPTN_ENDPOINT=http://api.${KEPTN_NAMESPACE}.127.0.0.1.nip.io/api
else
  if [[ "$KEPTN_SERVICE_TYPE" == "NodePort" ]]; then
    API_PORT=$(kubectl get svc api-gateway-nginx -n ${KEPTN_NAMESPACE} -o jsonpath='{.spec.ports[?(@.name=="http")].nodePort}')
    INTERNAL_NODE_IP=$(kubectl get nodes -o jsonpath='{ $.items[0].status.addresses[?(@.type=="InternalIP")].address }')
    KEPTN_ENDPOINT="http://${INTERNAL_NODE_IP}:${API_PORT}"/api
  else
    KEPTN_ENDPOINT=http://$(kubectl -n ${KEPTN_NAMESPACE} get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')/api
  fi
fi

KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n ${KEPTN_NAMESPACE} -ojsonpath={.data.keptn-api-token} | base64 --decode)


# deploy dynatrace-sli service
kubectl -n ${KEPTN_NAMESPACE} create secret generic dynatrace-credentials-${SELF_MONITORING_PROJECT} --from-literal="DT_TENANT=$QG_INTEGRATION_TEST_DT_TENANT" --from-literal="DT_API_TOKEN=$QG_INTEGRATION_TEST_DT_API_TOKEN"

echo "Install dynatrace-sli-service from: ${DYNATRACE_SLI_SERVICE_VERSION}"
kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/dynatrace-sli-service/${DYNATRACE_SLI_SERVICE_VERSION}/deploy/service.yaml -n ${KEPTN_NAMESPACE}
sleep 5
# TODO: delete line below !
kubectl -n ${KEPTN_NAMESPACE} set image deployment/dynatrace-sli-service dynatrace-sli-service=keptncontrib/dynatrace-sli-service:0.6.0-master

sleep 10
wait_for_deployment_in_namespace "dynatrace-sli-service" ${KEPTN_NAMESPACE}

kubectl create configmap -n ${KEPTN_NAMESPACE} lighthouse-config-${SELF_MONITORING_PROJECT} --from-literal=sli-provider=dynatrace

# create the project

keptn create project $SELF_MONITORING_PROJECT --shipyard=./test/assets/shipyard-quality-gates-self-monitoring.yaml

keptn add-resource --project=$SELF_MONITORING_PROJECT --resource=./test/assets/self_monitoring_sli.yaml --resourceUri=dynatrace/sli.yaml

# create services

SERVICES=("bridge" "eventbroker-go" "configuration-service" "mongodb-datastore" "gatekeeper-service" "remediation-service" "lighthouse-service" "statistics-service" "gatekeeper-service" "dynatrace-sli-service" "jmeter-service" "dynatrace-service" "api-service" "api-gateway-nginx")

for SERVICE in "${SERVICES[@]}"
do
    keptn create service $SERVICE --project=$SELF_MONITORING_PROJECT
done

for SERVICE in "${SERVICES[@]}"
do
    keptn add-resource --project=$SELF_MONITORING_PROJECT --service=$SERVICE --stage=hardening --resource=./test/assets/self_monitoring_slo.yaml --resourceUri=slo.yaml
done

# initial evaluation

SELF_MONITORING_SERVICE=mongodb-datastore

keptn add-resource --project=$SELF_MONITORING_PROJECT --service=$SELF_MONITORING_SERVICE --stage=hardening --resource=./test/assets/mongodb-performance.jmx --resourceUri=jmeter/load.jmx

cat << EOF > ./tmp-trigger-evaluation.json
{
  "type": "sh.keptn.event.hardening.evaluation.triggered",
  "specversion": "1.0",
  "source": "travis-ci",
  "contenttype": "application/json",
  "data": {
    "project": "$SELF_MONITORING_PROJECT",
    "stage": "hardening",
    "service": "$SELF_MONITORING_SERVICE",
    "deployment": {
      "deploymentURIsLocal": ["$SELF_MONITORING_SERVICE:8080"]
    },
    "labels": {
      "nr_projects": "0"
    }
  }
}
EOF

cat tmp-trigger-evaluation.json

keptn_context_id=$(send_event_json ./tmp-trigger-evaluation.json)
rm tmp-trigger-evaluation.json

# try to fetch a evaluation-done event
echo "Getting evaluation-done event with context-id: ${keptn_context_id}"
response=$(get_event_with_retry sh.keptn.event.evaluation.finished ${keptn_context_id} ${SELF_MONITORING_PROJECT})

echo $response | jq .


# Create projects, services and evaluations to generate data
for project_nr in $(seq 1 ${NR_PROJECTS})
do
  keptn create project project-${project_nr} --shipyard=./test/assets/shipyard-quality-gates.yaml

  for service_nr in $(seq 1 ${NR_SERVICES_PER_PROJECT})
  do
    keptn create service service-${service_nr} --project=project-${project_nr}

    for evaluation_r in $(seq 1 ${NR_EVALUATIONS_PER_SERVICE})
    do
      send_start_evaluation_request project-${project_nr} hardening service-${service_nr}
    done

    # do the evaluation again
    cat << EOF > ./tmp-trigger-evaluation.json
    {
      "type": "sh.keptn.event.hardening.evaluation.triggered",
      "specversion": "1.0",
      "source": "travis-ci",
      "contenttype": "application/json",
      "data": {
        "project": "$SELF_MONITORING_PROJECT",
        "stage": "hardening",
        "service": "$SELF_MONITORING_SERVICE",
        "deployment": {
          "deploymentURIsLocal": ["$SELF_MONITORING_SERVICE:8080"]
        },
        "labels": {
          "nr_projects": "$SELF_MONITORING_PROJECT_nr",
          "nr_services_per_project": "$service_nr",
          "nr_evaluations_per_service": "$NR_NR_EVALUATIONS_PER_SERVICE"
        }
      }
    }
EOF

    keptn_context_id=$(send_event_json ./tmp-trigger-evaluation.json)
    rm tmp-trigger-evaluation.json

    # try to fetch a evaluation-done event
    echo "Getting evaluation-done event with context-id: ${keptn_context_id}"
    response=$(get_event_with_retry sh.keptn.event.evaluation.finished ${keptn_context_id} ${SELF_MONITORING_PROJECT})

    echo $response | jq .
  done
done
