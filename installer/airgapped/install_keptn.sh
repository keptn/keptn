#!/bin/bash

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
control-plane.mongodb.image.repositoryImageName=centos/mongodb-36-centos7,\
control-plane.nats.nats.image=${TARGET_INTERNAL_DOCKER_REGISTRY}nats:2.1.9-alpine3.14,\
control-plane.nats.reloader.image=${TARGET_INTERNAL_DOCKER_REGISTRY}connecteverything/nats-server-config-reloader:0.6.0,\
control-plane.nats.exporter.image=${TARGET_INTERNAL_DOCKER_REGISTRY}synadia/prometheus-nats-exporter:0.5.0,\
control-plane.apiGatewayNginx.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.apiGatewayNginx.image.repositoryImageName=nginxinc/nginx-unprivileged,\
control-plane.apiGatewayNginx.image.tag=1.19.4-alpine,\
control-plane.remediationService.image.registry=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.remediationService.image.repositoryImageName=keptn/remediation-service,\
control-plane.apiService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.apiService.image.repositoryImageName=keptn/api,\
control-plane.bridge.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.bridge.image.repositoryImageName=keptn/bridge2,\
control-plane.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.distributor.image.repositoryImageName=keptn/distributor,\
control-plane.shipyardController.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.shipyardController.image.repositoryImageName=keptn/shipyard-controller,\
control-plane.configurationService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.configurationService.image.repositoryImageName=keptn/configuration-service,\
control-plane.mongodbDatastore.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.mongodbDatastore.image.repositoryImageName=keptn/mongodb-datastore,\
control-plane.statisticsService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.statisticsService.image.repositoryImageName=keptn/statistics-service,\
control-plane.lighthouseService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.lighthouseService.image.repositoryImageName=keptn/lighthouse-service,\
control-plane.secretService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.secretService.image.repositoryImageName=keptn/secret-service,\
control-plane.approvalService.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
control-plane.approvalService.image.repositoryImageName=keptn/approval-service,\
continuous-delivery.distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
continuous-delivery.distributor.image.repositoryImageName=keptn/distributor"

echo ""

echo "-----------------------------------------------------------------------"
echo "Installing Keptn Helm-Service Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

helm upgrade helm-service "${KEPTN_HELM_SERVICE_HELM_CHART}" --install -n "${KEPTN_NAMESPACE}" \
--set="helmservice.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
helmservice.image.repositoryImageName=keptn/helm-service,\
distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
distributor.image.repositoryImageName=${TARGET_INTERNAL_DOCKER_REGISTRY}keptn/distributor"

echo ""

echo "-----------------------------------------------------------------------"
echo "Installing Keptn JMeter-Service Helm Chart in Namespace ${KEPTN_NAMESPACE}"
echo "-----------------------------------------------------------------------"

helm upgrade jmeter-service "${KEPTN_JMETER_SERVICE_HELM_CHART}" --install -n "${KEPTN_NAMESPACE}" \
--set="jmeterservice.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
jmeterservice.image.repositoryImageName=keptn/jmeter-service,\
distributor.image.repository=${TARGET_INTERNAL_DOCKER_REGISTRY},\
distributor.image.repositoryImageName=keptn/distributor"


# add keptn.sh/managed-by annotation to the namespace
kubectl patch namespace "${KEPTN_NAMESPACE}" \
-p "{\"metadata\": {\"annotations\": {\"keptn.sh/managed-by\": \"keptn\"}, \"labels\": {\"keptn.sh/managed-by\": \"keptn\"}}}"
