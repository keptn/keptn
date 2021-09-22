#!/bin/bash
# shellcheck disable=SC2181
if [[ $# -ne 4 ]]; then
    echo "Please provide the target registry and helm charts as parameters, e.g., "
    echo "$1 \"docker.io/your-username/\" \"keptn-0.9.0.tgz\" \"helm-service-0.9.0.tgz\" \"jmeter-service-0.9.0.tgz\""
    exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}
KEPTN_HELM_CHART=${2}
KEPTN_HELM_SERVICE_HELM_CHART=${3}
KEPTN_JMETER_SERVICE_HELM_CHART=${4}

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-"keptn"}
KEPTN_SERVICE_TYPE=${KEPTN_SERVICE_TYPE:-"ClusterIP"}

echo "-----------------------------------------------------------------------"
echo "Installing Keptn Core Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

kubectl create namespace "${KEPTN_NAMESPACE}"

helm upgrade keptn "${KEPTN_HELM_CHART}" --install --create-namespace -n "${KEPTN_NAMESPACE}" --wait \
--set="control-plane.apiGatewayNginx.type=${KEPTN_SERVICE_TYPE},continuous-delivery.enabled=true,\
control-plane.mongodb.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.nats.nats.image=${TARGET_INTERNAL_DOCKER_REGISTRY}nats:2.1.9-alpine3.12,\
control-plane.nats.reloader.image=${TARGET_INTERNAL_DOCKER_REGISTRY}connecteverything/nats-server-config-reloader:0.6.0,\
control-plane.nats.exporter.image=${TARGET_INTERNAL_DOCKER_REGISTRY}synadia/prometheus-nats-exporter:0.5.0,\
control-plane.apiGatewayNginx.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nginxinc/nginx-unprivileged,\
control-plane.apiGatewayNginx.image.tag=1.21.3-alpine,\
control-plane.remediationService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/remediation-service,\
control-plane.apiService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/api,\
control-plane.bridge.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/bridge2,\
control-plane.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/distributor,\
control-plane.shipyardController.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/shipyard-controller,\
control-plane.configurationService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/configuration-service,\
control-plane.mongodbDatastore.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/mongodb-datastore,\
control-plane.statisticsService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/statistics-service,\
control-plane.lighthouseService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/lighthouse-service,\
control-plane.secretService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/secret-service,\
control-plane.approvalService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/approval-service,\
control-plane.webhookService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/webhook-service,\
continuous-delivery.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/distributor"

if [[ $? -ne 0 ]]; then
  echo "Installing Keptn failed."
  exit 1
fi

echo ""

echo "-----------------------------------------------------------------------"
echo "Installing Keptn Helm-Service Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

helm upgrade helm-service "${KEPTN_HELM_SERVICE_HELM_CHART}" --install -n "${KEPTN_NAMESPACE}" \
--set="helmservice.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/helm-service,\
distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/distributor"

if [[ $? -ne 0 ]]; then
  echo "Installing helm-service failed."
  exit 1
fi

echo ""

echo "-----------------------------------------------------------------------"
echo "Installing Keptn JMeter-Service Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

helm upgrade jmeter-service "${KEPTN_JMETER_SERVICE_HELM_CHART}" --install -n "${KEPTN_NAMESPACE}" \
--set="jmeterservice.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/jmeter-service,\
distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/distributor"

if [[ $? -ne 0 ]]; then
  echo "Installing jmeter-service failed."
  exit 1
fi

# add keptn.sh/managed-by annotation to the namespace
kubectl patch namespace "${KEPTN_NAMESPACE}" \
-p "{\"metadata\": {\"annotations\": {\"keptn.sh/managed-by\": \"keptn\"}, \"labels\": {\"keptn.sh/managed-by\": \"keptn\"}}}"

if [[ $? -ne 0 ]]; then
  echo "Patching the namespace failed"
  exit 1
fi
