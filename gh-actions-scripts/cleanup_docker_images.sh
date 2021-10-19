#!/bin/bash

##################################################################
# This script deletes outdated images from DockerHub             #
# Required secrets/params:                                       #
# - REGISTRY_USER                                                #
# - REGISTRY_PASSWORD                                            #
# - DOCKER_ORG                                                   #
##################################################################

##################################################################
# Configuration                                                  #
##################################################################

# max age that images should have before they are marked as outdated
MAX_AGE=30
IMAGES=(
  "api"
  "bridge2"
  "configuration-service"
  "openshift-route-service"
  "distributor"
  "approval-service"
  "helm-service"
  "jmeter-service"
  "lighthouse-service"
  "mongodb-datastore"
  "remediation-service"
  "secret-service"
  "shipyard-controller"
  "statistics-service"
  "webhook-service"
)

# ensure the params/variables are set
if [ -z "$REGISTRY_USER" ]; then
  echo "REGISTRY_USER is not set. Please set REGISTRY_USER to the username of your container registry."
  exit 1
fi

if [ -z "$REGISTRY_PASSWORD" ]; then
  echo "REGISTRY_PASSWORD is not set. Please set REGISTRY_PASSWORD to the password of your container registry."
  exit 1
fi

if [ -z "$DOCKER_ORG" ]; then
  echo "DOCKER_ORG is not set. Please set DOCKER_ORG to the organization that you want to check stale images for."
  exit 1
fi

##################################################################
# Actual Job                                                     #
##################################################################

# Authenticate at Docker Hub
# Authenticate against DockerHub API

# shellcheck disable=SC2086
DOCKER_HUB_TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

if [[ "$DOCKER_HUB_TOKEN" == "null" ]]; then
  echo "Failed to authenticate on DockerHub API."
  exit 1
fi

function get_outdated_commit_hash_tags() {
  REPO=$1
  MAX_AGE_DAYS=$2

  # Target-Date = Current Date Minus $MAX_AGE_DAYS
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?page_size=1" | jq -r '.count')
  # unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  echo >&2 "Found $COUNT tags for $REPO without filter"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?ordering=-last_updated&page_size=${COUNT}" |
    jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|select(.name | match("\\b[0-9a-f]{7}\\b"))|.name')
  echo "$response"
}

function get_outdated_datetime_tags() {
  REPO=$1
  TAG_FILTER=$2
  MAX_AGE_DAYS=$3

  # Target-Date = Current Date Minus $MAX_AGE_DAYS
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')
  # unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  echo >&2 "Found $COUNT tags for $REPO (filter $TAG_FILTER)"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" |
    jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|select(.name | match("^\\b[0-9]{8}\\b"))|.name')
  echo "$response"
}

# get all outdated images (e.g., for repo=keptn/bridge2, tag_filter=patch, max_age_days=30)
function get_outdated_images() {
  REPO=$1
  TAG_FILTER=$2
  MAX_AGE_DAYS=$3

  # Target-Date = Current Date Minus $MAX_AGE_DAYS
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')
  # unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  echo >&2 "Found $COUNT tags for $REPO (filter $TAG_FILTER)"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_HUB_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" | jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|.name')
  echo "$response"
}

# delete a tag in a repo
function delete_tag() {
  REPO=$1
  TAG=$2
  DOCKER_REGISTRY_TOKEN=$3

  echo -e "Fetching image manifest for ${REPO}:${TAG}"

  image_digest=$(curl -s -I \
    -H "Authorization: Bearer ${DOCKER_REGISTRY_TOKEN}" \
    -H "Accept: application/vnd.docker.distribution.manifest.v2+json" \
    "https://index.docker.io/v2/${DOCKER_ORG}/${REPO}/manifests/${TAG}" \
    | awk '$1 == "docker-content-digest:" { print $2 }' \
    | tr -d $'\r' \
  )

  echo -e "Deleting ${REPO}:${TAG} with digest ${image_digest}"
  response=$(curl -s \
    -X DELETE \
    -H "Authorization: Bearer ${DOCKER_REGISTRY_TOKEN}" \
    "https://registry-1.docker.io/v2/$DOCKER_ORG/$REPO/manifests/${image_digest}" \
  )

  if [[ "$response" != "204" ]]; then
    echo " - Delete failed with response $response"
    exit 1
  else
    echo " - Done!"
  fi
}

for service in "${IMAGES[@]}"; do
  echo "Deleting outdated images for service ${service}"

  # get all outdated commit hash tags
  outdated_commit_hash_tags=$(get_outdated_commit_hash_tags "$service" "$MAX_AGE")

  outdated_datetime_tags=$(get_outdated_datetime_tags "$service" "2020" "$MAX_AGE")

  # get all outdated tag where tag contains "dev-PR"
  # outdated_dev_pr_tags=$(get_outdated_images "$service" "dev-PR" "$MAX_AGE")

  # ToDo: Also Check for "x.y.z-dev.20" tags (e.g., 0.8.0-dev.20210101)
  # outdated_dev_tags=$(get_outdated_images "$service" "dev.20" "$MAX_AGE")

  # get all outdated tag where tag contains "feature"
  outdated_feature_tags=$(get_outdated_images "$service" "feature" "$MAX_AGE")
  # get all outdated tag where tag contains "bug"
  outdated_bug_tags=$(get_outdated_images "$service" "bug" "$MAX_AGE")
  # get all outdated tag where tag contains "patch"
  outdated_patch_tags=$(get_outdated_images "$service" "patch" "$MAX_AGE")

  # get all outdated tag where tag contains "dirty"
  outdated_dirty_tags=$(get_outdated_images "$service" "dirty" "$MAX_AGE")

  docker_registry_token=$(curl -s \
    -H "Content-Type: application/json" \
    "https://auth.docker.io/token?service=registry.docker.io&scope=repository:${DOCKER_ORG}/${service}:*" \
    | jq -r .token\
  )

  for tag in ${outdated_commit_hash_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done

  for tag in ${outdated_datetime_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done

  # variable definitions commented out above

  # for tag in ${outdated_dev_pr_tags}; do
  #   echo "dummy" "$service" "$tag"
  # done

  # for tag in ${outdated_dev_tags}; do
  #   echo "dummy" "$service" "$tag"
  # done

  for tag in ${outdated_feature_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done

  for tag in ${outdated_bug_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done

  for tag in ${outdated_patch_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done

  for tag in ${outdated_dirty_tags}; do
    delete_tag "$service" "$tag" "$docker_registry_token"
  done
done
