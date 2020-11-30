#!/bin/bash

function get_outdated_images() {
  REPO=$1
  TAG_FILTER=$2
  MAX_AGE_DAYS=$3

  DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

  COUNT=$(curl -s -H "Authorization: JWT ${TOKEN}" "https://hub.docker.com/v2/repositories/keptn/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')

  response=$(curl -s -H "Authorization: JWT ${TOKEN}" "https://hub.docker.com/v2/repositories/keptn/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" | jq -r --argjson date "$DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|.name')
  echo $response
}

function delete_tag() {
  REPO=$1
  TAG=$2
  TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

  curl -i -X DELETE \
  -H "Accept: application/json" \
  -H "Authorization: JWT ${TOKEN}" \
  https://hub.docker.com/v2/repositories/keptn/$REPO/tags/$TAG/

}

allServices=("api" "bridge2" "configuration-service" "openshift-route-service" "distributor" "gatekeeper-service" "helm-service" "jmeter-service" "lighthouse-service" "mongodb-datastore" "remediation-service" "shipyard-controller" "shipyard-service")
MAX_AGE=30

for s in ${allServices[@]}; do
  echo "deleting outdated images for service ${s}"
  outdated_feature_tags=$(get_outdated_images $s "feature" $MAX_AGE)
  outdated_bug_tags=$(get_outdated_images $s "bug" $MAX_AGE)
  outdated_patch_tags=$(get_outdated_images $s "patch" $MAX_AGE)
  outdated_dirty_tags=$(get_outdated_images $s "dirty" $MAX_AGE)

  for tag in ${outdated_feature_tags}; do
    echo "Deleting ${s}:${tag}"
    delete_tag $s $tag
  done

  for tag in ${outdated_bug_tags}; do
    echo "Deleting ${s}:${tag}"
    delete_tag $s $tag
  done

  for tag in ${outdated_patch_tags}; do
    echo "Deleting ${s}:${tag}"
    delete_tag $s $tag
  done

  for tag in ${outdated_dirty_tags}; do
    echo "Deleting ${s}:${tag}"
    delete_tag $s $tag
  done
done
