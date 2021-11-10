#!/bin/bash

KEPTN_TAG=${KEPTN_TAG:-"0.10.0"}

if [[ $# -ne 1 ]]; then
    echo "Please provide the target registry as a param, e.g., $1 gcr.io/your-registry/"
    exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}

IMAGES_CONTROL_PLANE_THIRD_PARTY="bitnami/mongodb:4.4.9-debian-10-r0 nats:2.1.9-alpine3.12 connecteverything/nats-server-config-reloader:0.6.0 synadia/prometheus-nats-exporter:0.5.0 nginxinc/nginx-unprivileged:1.21.3-alpine"
IMAGES_CONTROL_PLANE="keptn/api:${KEPTN_TAG} keptn/bridge2:${KEPTN_TAG} keptn/configuration-service:${KEPTN_TAG} keptn/distributor:${KEPTN_TAG} keptn/secret-service:${KEPTN_TAG} keptn/shipyard-controller:${KEPTN_TAG} keptn/remediation-service:${KEPTN_TAG} keptn/mongodb-datastore:${KEPTN_TAG} keptn/statistics-service:${KEPTN_TAG} keptn/lighthouse-service:${KEPTN_TAG} keptn/approval-service:${KEPTN_TAG} keptn/webhook-service:${KEPTN_TAG}"
IMAGES_CONTINUOUS_DELIVERY="keptn/helm-service:${KEPTN_TAG} keptn/jmeter-service:${KEPTN_TAG}"

IMAGES="$IMAGES_CONTROL_PLANE_THIRD_PARTY $IMAGES_CONTROL_PLANE $IMAGES_CONTINUOUS_DELIVERY"


for img in $IMAGES
do
    echo "Processing $img..."
    docker pull "$img"
    docker tag "$img" "${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
    docker push "${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
    echo "$img -> ${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
done
