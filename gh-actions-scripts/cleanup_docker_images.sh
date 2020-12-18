#!/bin/bash

##################################################################
# This script deletes outdated images from DockerHub             #
# Required secrets/params:                                       #
# - REGISTRY_USER                                                #
# - REGISTRY_PASSWORD                                            #
##################################################################

##################################################################
# Configuration                                                  #
##################################################################

# list all images that should be checked
IMAGES=("api" "bridge2" "configuration-service" "openshift-route-service" "distributor" "gatekeeper-service" "helm-service" "jmeter-service" "lighthouse-service" "mongodb-datastore" "remediation-service" "shipyard-controller" "shipyard-service")
# max age that images should have before they are marked as outdated
MAX_AGE=30


##################################################################
# Actual Job                                                     #
##################################################################

# get all outdated images (e.g., for repo=keptn/bridge2, tag_filter=patch, max_age_days=30)
function get_outdated_images() {
  REPO=$1
  TAG_FILTER=$2
  MAX_AGE_DAYS=$3

  # Target-Date = Current Date Minus $MAX_AGE_DAYS
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  # Authenticate against DockerHub API
  TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${TOKEN}" "https://hub.docker.com/v2/repositories/keptn/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${TOKEN}" "https://hub.docker.com/v2/repositories/keptn/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" | jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|.name')
  echo $response
}

# delete a tag in a repo
function delete_tag() {
  REPO=$1
  TAG=$2
  TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

  curl -i -X DELETE \
  -H "Accept: application/json" \
  -H "Authorization: JWT ${TOKEN}" \
  https://hub.docker.com/v2/repositories/keptn/$REPO/tags/$TAG/

}


for s in ${IMAGES[@]}; do
  echo "deleting outdated images for service ${s}"
  # get all outdated tag where tag contains "feature"
  outdated_feature_tags=$(get_outdated_images $s "feature" $MAX_AGE)
  # get all outdated tag where tag contains "bug"
  outdated_bug_tags=$(get_outdated_images $s "bug" $MAX_AGE)
  # get all outdated tag where tag contains "patch"
  outdated_patch_tags=$(get_outdated_images $s "patch" $MAX_AGE)
  # get all outdated tag where tag contains "dirty"
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
