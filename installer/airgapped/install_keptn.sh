#!/bin/bash
# shellcheck disable=SC2181
if [[ $# -ne 3 ]]; then
    echo "Please provide the target registry, organization and helm charts as parameters, e.g., "
    echo "$1 \"docker.io/your-username/\" \"keptn\" \"keptn-0.9.0.tgz\""
    exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}
DOCKER_ORG=${2}
KEPTN_HELM_CHART=${3}

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-"keptn"}
KEPTN_SERVICE_TYPE=${KEPTN_SERVICE_TYPE:-"ClusterIP"}

echo "-----------------------------------------------------------------------"
echo "Installing Keptn Core Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

helm template keptn "${KEPTN_HELM_CHART}" -n "${KEPTN_NAMESPACE}" \
--set="apiGatewayNginx.type=${KEPTN_SERVICE_TYPE},\
global.keptn.registry=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG},\
global.initContainers.image=${TARGET_INTERNAL_DOCKER_REGISTRY}curlimages/curl,\
mongo.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY%/},\
nats.nats.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nats,\
apiGatewayNginx.image.registry="",\
apiGatewayNginx.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nginxinc/nginx-unprivileged"

helm upgrade keptn "${KEPTN_HELM_CHART}" --install -n "${KEPTN_NAMESPACE}" --create-namespace --wait --timeout 12m \
--set="apiGatewayNginx.type=${KEPTN_SERVICE_TYPE},\
global.keptn.registry=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG},\
global.initContainers.image=${TARGET_INTERNAL_DOCKER_REGISTRY}curlimages/curl,\
mongo.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY%/},\
nats.nats.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nats,\
apiGatewayNginx.image.registry="",\
apiGatewayNginx.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nginxinc/nginx-unprivileged"

if [[ $? -ne 0 ]]; then
  echo "Installing Keptn failed."
  exit 1
fi

# add keptn.sh/managed-by annotation to the namespace
kubectl patch namespace "${KEPTN_NAMESPACE}" \
-p "{\"metadata\": {\"annotations\": {\"keptn.sh/managed-by\": \"keptn\"}, \"labels\": {\"keptn.sh/managed-by\": \"keptn\"}}}"

if [[ $? -ne 0 ]]; then
  echo "Patching the namespace failed"
  exit 1
fi
