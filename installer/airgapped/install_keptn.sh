#!/bin/bash
# shellcheck disable=SC2181
if [[ $# -ne 5 ]]; then
    echo "Please provide the target registry, organization and helm charts as parameters, e.g., "
    echo "$1 \"docker.io/your-username/\" \"keptn\" \"keptn-0.9.0.tgz\" \"helm-service-0.9.0.tgz\" \"jmeter-service-0.9.0.tgz\""
    exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}
DOCKER_ORG=${2}
KEPTN_HELM_CHART=${3}
KEPTN_HELM_SERVICE_HELM_CHART=${4}
KEPTN_JMETER_SERVICE_HELM_CHART=${5}

KEPTN_NAMESPACE=${KEPTN_NAMESPACE:-"keptn"}
KEPTN_SERVICE_TYPE=${KEPTN_SERVICE_TYPE:-"ClusterIP"}

echo "-----------------------------------------------------------------------"
echo "Installing Keptn Core Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

kubectl create namespace "${KEPTN_NAMESPACE}"

helm upgrade keptn "${KEPTN_HELM_CHART}" --install --create-namespace -n "${KEPTN_NAMESPACE}" --wait \
--set="control-plane.apiGatewayNginx.type=${KEPTN_SERVICE_TYPE},continuous-delivery.enabled=true,\
control-plane.mongo.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY%/},\
control-plane.nats.nats.image=${TARGET_INTERNAL_DOCKER_REGISTRY}nats:2.7.2-alpine,\
control-plane.nats.reloader.image=${TARGET_INTERNAL_DOCKER_REGISTRY}natsio/nats-server-config-reloader:0.6.2,\
control-plane.nats.exporter.image=${TARGET_INTERNAL_DOCKER_REGISTRY}natsio/prometheus-nats-exporter:0.9.1,\
control-plane.apiGatewayNginx.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}nginxinc/nginx-unprivileged,\
control-plane.apiGatewayNginx.image.tag=1.22.0-alpine,\
control-plane.remediationService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/remediation-service,\
control-plane.apiService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/api,\
control-plane.bridge.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/bridge2,\
control-plane.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/distributor,\
control-plane.shipyardController.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/shipyard-controller,\
control-plane.configurationService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/configuration-service,\
control-plane.resourceService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/resource-service,\
control-plane.mongodbDatastore.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/mongodb-datastore,\
control-plane.statisticsService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/statistics-service,\
control-plane.lighthouseService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/lighthouse-service,\
control-plane.secretService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/secret-service,\
control-plane.approvalService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/approval-service,\
control-plane.webhookService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/webhook-service,\
continuous-delivery.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY}${DOCKER_ORG}/distributor"

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
