#!/bin/bash

#####################################################################################################
# This script deletes outdated images from DockerHub                                                #
# Required secrets/params:                                                                          #
# - REGISTRY_USER: A dockerhub user.                                                                #
# - REGISTRY_PASSWORD: The corresponding password. CAUTION: Personal access tokens don't work here! #
# - DOCKER_ORG: The organization on dockerhub where the images are located                          #
#####################################################################################################

#####################################################################################################
# Configuration                                                                                     #
#####################################################################################################

# max age that images should have before they are marked as outdated
MAX_AGE_DAYS=30
IMAGES=(
  "api"
  "bridge2"
  "configuration-service"
  "resource-service"
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
DOCKER_API_TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

if [[ "$DOCKER_API_TOKEN" == "null" ]]; then
  echo "Failed to authenticate on DockerHub Api."
  exit 1
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
  TARGET_DATE=$(echo "$(date +%s) - (${MAX_AGE_DAYS} * 24 * 60 * 60)" | bc)
else
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)
fi

function get_outdated_commit_hash_tags() {
  REPO=$1

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?page_size=1" | jq -r '.count')
# unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  >&2 echo "Found $COUNT tags for $REPO without filter"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?ordering=-last_updated&page_size=${COUNT}" | \
     jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|select(.name | match("\\b[0-9a-f]{7}\\b"))|.name')
  echo "$response"
}

function get_outdated_datetime_tags() {
  REPO=$1
  TAG_FILTER=$2

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')
  # unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  >&2 echo "Found $COUNT tags for $REPO (filter $TAG_FILTER)"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" | \
    jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|select(.name | match("^\\b[0-9]{8}\\b"))|.name')
  echo "$response"
}

# get all outdated images (e.g., for repo=keptn/bridge2, tag_filter=patch, max_age_days=30)
function get_outdated_images() {
  REPO=$1
  TAG_FILTER=$2

  # Count number of tags based on the tag filter
  COUNT=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&page_size=1" | jq -r '.count')
  # unfortunately, anything above 100 doesn't work for pagination with docker hub api; leaving it in for debug purposes
  >&2 echo "Found $COUNT tags for $REPO (filter $TAG_FILTER)"

  # get all tags, ordered by last_update (get the newest), and filter with jq based on TARGET_DATE
  response=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/?name=${TAG_FILTER}&ordering=-last_updated&page_size=${COUNT}" | jq -r --argjson date "$TARGET_DATE" '.results|.[]|select (.last_updated | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|.name')
  echo "$response"
}

# delete a tag in a repo
function delete_tag() {
  REPO=$1
  TAG=$2

  echo -ne "Deleting ${REPO}:${TAG}"

  response=$(curl -s -o /dev/null -i -X DELETE \
    -w "%{http_code}" \
    -H "Accept: application/json" \
    -H "Authorization: JWT ${DOCKER_API_TOKEN}" \
    "https://hub.docker.com/v2/repositories/$DOCKER_ORG/$REPO/tags/$TAG/")

  if [[ "$response" != "204" ]]; then
    echo " - Delete failed with response $response"
  else
    echo " - Done!"
  fi
}

for s in "${IMAGES[@]}"; do
  echo "Deleting outdated images for service ${s}"

  # get all outdated commit hash tags
  outdated_commit_hash_tags=$(get_outdated_commit_hash_tags "$s")

  outdated_datetime_tags=$(get_outdated_datetime_tags "$s" "2021")

  # get all outdated tag where tag contains "dev-PR"
   outdated_dev_pr_tags=$(get_outdated_images "$s" "dev-PR")

  # ToDo: Also Check for "x.y.z-dev.20" tags (e.g., 0.8.0-dev.20210101)
   outdated_dev_tags=$(get_outdated_images "$s" "dev.20")

  # get all outdated tag where tag contains "dirty"
  outdated_dirty_tags=$(get_outdated_images "$s" "dirty")

  for tag in ${outdated_commit_hash_tags}; do
    delete_tag "$s" "$tag"
  done

  for tag in ${outdated_datetime_tags}; do
    delete_tag "$s" "$tag"
  done

  for tag in ${outdated_dev_pr_tags}; do
    delete_tag "$s" "$tag"
  done

  for tag in ${outdated_dev_tags}; do
    delete_tag "$s" "$tag"
  done

  for tag in ${outdated_dirty_tags}; do
    delete_tag "$s" "$tag"
  done
done
