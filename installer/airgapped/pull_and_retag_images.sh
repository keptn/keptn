#!/bin/bash

KEPTN_TAG=${KEPTN_TAG:-"0.13.1"}

if [[ $# -ne 1 ]]; then
  echo "Please provide the target registry as a param, e.g., $1 gcr.io/your-registry/"
  exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}
DOCKER_ORG="keptn"

if [[ "$KEPTN_TAG" == *"dev"* ]]; then
  DOCKER_ORG="keptndev"
fi

IMAGES_CONTROL_PLANE_THIRD_PARTY=(
  "bitnami/mongodb:6.0.3-debian-11-r9"
  "nats:2.9.11-alpine"
  "nginxinc/nginx-unprivileged:1.23.3-alpine"
  "curlimages/curl:7.85.0"
)
IMAGES_CONTROL_PLANE=(
  "${DOCKER_ORG}/api:${KEPTN_TAG}"
  "${DOCKER_ORG}/bridge2:${KEPTN_TAG}"
  "${DOCKER_ORG}/distributor:${KEPTN_TAG}"
  "${DOCKER_ORG}/secret-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/shipyard-controller:${KEPTN_TAG}"
  "${DOCKER_ORG}/remediation-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/mongodb-datastore:${KEPTN_TAG}"
  "${DOCKER_ORG}/statistics-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/lighthouse-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/approval-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/webhook-service:${KEPTN_TAG}"
  "${DOCKER_ORG}/resource-service:${KEPTN_TAG}"
)

IMAGES=()
IMAGES+=("${IMAGES_CONTROL_PLANE_THIRD_PARTY[@]}")
IMAGES+=("${IMAGES_CONTROL_PLANE[@]}")

for img in "${IMAGES[@]}"; do
  echo "Processing $img..."
  docker pull "$img"
  docker tag "$img" "${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
  docker push "${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
  echo "$img -> ${TARGET_INTERNAL_DOCKER_REGISTRY}${img}"
done
