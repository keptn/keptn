#!/bin/bash

KEPTN_TAG=${KEPTN_TAG:-"0.13.1"}

if [[ $# -ne 1 ]]; then
  echo "Please provide the target registry as a param, e.g., $1 gcr.io/your-registry/"
  exit 1
fi

TARGET_INTERNAL_DOCKER_REGISTRY=${1}
CONTAINER_ORG="ghcr.io/keptn"

IMAGES_CONTROL_PLANE_THIRD_PARTY=(
  "bitnami/mongodb:6.0.5-debian-11-r4"
  "nats:2.9.15-alpine"
  "nginxinc/nginx-unprivileged:1.23.3-alpine"
  "curlimages/curl:7.85.0"
)
IMAGES_CONTROL_PLANE=(
  "${CONTAINER_ORG}/api:${KEPTN_TAG}"
  "${CONTAINER_ORG}/bridge2:${KEPTN_TAG}"
  "${CONTAINER_ORG}/distributor:${KEPTN_TAG}"
  "${CONTAINER_ORG}/secret-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/shipyard-controller:${KEPTN_TAG}"
  "${CONTAINER_ORG}/remediation-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/mongodb-datastore:${KEPTN_TAG}"
  "${CONTAINER_ORG}/statistics-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/lighthouse-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/approval-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/webhook-service:${KEPTN_TAG}"
  "${CONTAINER_ORG}/resource-service:${KEPTN_TAG}"
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
