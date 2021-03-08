#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

KEPTN_EXAMPLES_BRANCH=${KEPTN_EXAMPLES_BRANCH:-"master"}
KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-keptn}
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n "$KEPTN_NAMESPACE" -o jsonpath='{.data.keptn-api-token}' | base64 --decode)

echo "Testing onboarding for project $PROJECT"

echo "Using remote execution plane services: ${REMOTE_EXECUTION_PLANE}"

# check if REMOTE_EXECUTION_PLANE is set to true. If yes, scale down the helm-service and jmeter in the keptn namespace and install the services via their helm charts
if [[ "${REMOTE_EXECUTION_PLANE}" == "true" ]]; then
  kubectl scale deployment/helm-service -n ${KEPTN_NAMESPACE} --replicas=0
  kubectl scale deployment/jmeter-service -n ${KEPTN_NAMESPACE} --replicas=0

  KEPTN_API_HOSTNAME=$(kubectl -n ${KEPTN_NAMESPACE} get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  KEPTN_VERSION_TAG=$(kubectl get deployment -n ${KEPTN_NAMESPACE} helm-service -ogo-template='{{index .metadata.labels "app.kubernetes.io/version"}}')

  helm install helm-service ./helm-service/chart -f ./helm-service/chart/values.yaml -n keptn-helm-service --set helmservice.image.tag=${KEPTN_VERSION_TAG} --set distributor.image.tag=${KEPTN_VERSION_TAG} --set remoteControlPlane.enabled=true --set remoteControlPlane.api.protocol=http --set remoteControlPlane.api.hostname=${KEPTN_API_HOSTNAME} --set remoteControlPlane.api.token=${KEPTN_API_TOKEN} --create-namespace
  helm install jmeter-service ./jmeter-service/chart -f ./jmeter-service/chart/values.yaml -n keptn-jmeter-service --set jmeterservice.image.tag=${KEPTN_VERSION_TAG} --set distributor.image.tag=${KEPTN_VERSION_TAG} --set remoteControlPlane.enabled=true --set remoteControlPlane.api.protocol=http --set remoteControlPlane.api.hostname=${KEPTN_API_HOSTNAME} --set remoteControlPlane.api.token=${KEPTN_API_TOKEN} --create-namespace
fi

# verify that the project does not exist yet via the Keptn API
response=$(curl -X GET "${KEPTN_ENDPOINT}/controlPlane/v1/project/${PROJECT}" -H  "accept: application/json" -H  "x-token: ${KEPTN_API_TOKEN}" -k 2>/dev/null | jq -r '.projectName')

if [[ "$response" == "${PROJECT}" ]]; then
  echo "Project ${PROJECT} already exists. Please delete it using:"
  echo "keptn delete project ${PROJECT}"
  exit 2
fi

# test keptn create-project and onboard
rm -rf examples
git clone --branch "${KEPTN_EXAMPLES_BRANCH}" https://github.com/keptn/examples --single-branch
cd examples/onboarding-carts || exit

echo "Creating a new project without Git upstream"
keptn create project "$PROJECT" --shipyard=../../test/assets/shipyard_onboard_service.yaml
verify_test_step $? "keptn create project ${PROJECT} failed."
sleep 10

###########################################
# onboard carts                           #
###########################################
keptn onboard service carts --project="$PROJECT" --chart=./carts
verify_test_step $? "keptn onboard service carts failed"
sleep 10

# add functional tests
keptn add-resource --project="$PROJECT" --service=carts --stage=dev --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/basiccheck.jmx
# add performance tests
keptn add-resource --project="$PROJECT" --service=carts --stage=staging --resource=jmeter/basiccheck.jmx --resourceUri=jmeter/load.jmx

###########################################
# onboard carts-db                        #
###########################################
keptn onboard service carts-db --project="$PROJECT" --chart=./carts-db
verify_test_step $? "keptn onboard service carts-db failed"

echo "Onboarding done ✓"

exit 0
