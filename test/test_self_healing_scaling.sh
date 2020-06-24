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

PROMETHEUS_SERVICE_VERSION=${UNLEASH_SERVICE_VERSION:-0.3.4}

# Prerequisites

# add remediation.yaml
keptn add-resource --project=sockshop --stage=production --service=carts --resource=./test/assets/self_healing_scaling_remediation.yaml --resourceUri=remediation.yaml

rm -rf examples
git clone --branch master https://github.com/keptn/examples --single-branch

cd examples/onboarding-carts
# add slo file
keptn add-resource --project=sockshop --service=carts --stage=production --resource=slo-self-healing.yaml --resourceUri=slo.yaml

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-$PROMETHEUS_SERVICE_VERSION/deploy/service.yaml

wait_for_deployment_in_namespace prometheus-service keptn
echo "Prometheus service deployed successfully"

keptn configure monitoring prometheus --project=sockshop --service=carts

wait_for_deployment_in_namespace prometheus-deployment monitoring
echo "Prometheus deployed successfully"

cd ../load-generation/cartsloadgen

kubectl apply -f deploy/cartsloadgen-base.yaml

wait_for_deployment_in_namespace cartsloadgen loadgen

echo "loadgen deployed successfully waiting for problem notification"



