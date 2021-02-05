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

# list all images that should be checked
MAX_AGE_DAYS=30

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


# list of images to be checked
IMAGES=("api" "bridge2" "configuration-service" "openshift-route-service" "distributor" "gatekeeper-service" "helm-service" "jmeter-service" "lighthouse-service" "mongodb-datastore" "remediation-service" "shipyard-controller")
# additional old images that we want to keep
ADDITIONAL_OLD_IMAGES=("installer" "bridge" "upgrader" "shipyard-service" "wait-service" "pitometer-service")
# Older than version 0.5
ADDITIONAL_VERY_OLD_IMAGES=("keptn-authenticator" "keptn-control" "keptn-event-broker" "keptn-event-broker-ext" "slack-service" "control" "eventbroker" "eventbroker-ext" "github-service")

# merge IMAGES and ADDITIONAL_OLD_IMAGES
IMAGES=("${IMAGES[@]}" "${ADDITIONAL_OLD_IMAGES[@]}")
IMAGES=("${IMAGES[@]}" "${ADDITIONAL_VERY_OLD_IMAGES[@]}")


# Authenticate against DockerHub API
DOCKER_API_TOKEN=$(curl -s -H "Content-Type: application/json" -X POST -d '{"username": "'${REGISTRY_USER}'", "password": "'${REGISTRY_PASSWORD}'"}' https://hub.docker.com/v2/users/login/ | jq -r .token)

if [[ "$DOCKER_API_TOKEN" == "null" ]]; then
  echo "Failed to authenticate on DockerHub Api."
  exit 1
fi


# get all github release tags
function get_releases() {
  REPO=$1
  curl --silent https://api.github.com/repos/$1/releases | jq -r '.[].tag_name'
}

# check if a repo + tag is stale
function check_if_stale() {
  REPO=$1
  TAGS=$2

  # Target-Date = Current Date Minus $MAX_AGE_DAYS
  TARGET_DATE=$(date -d "-${MAX_AGE_DAYS} days" +%s)

  # for each tag, check if the tag is stale
  for TAG in ${TAGS[@]}; do
    HTTP_RESPONSE=$(curl -s -H "Authorization: JWT ${DOCKER_API_TOKEN}" --write-out "HTTPSTATUS:%{http_code}" "https://hub.docker.com/v2/repositories/${DOCKER_ORG}/${REPO}/tags/${TAG}/")

    # extract body and status
    HTTP_BODY=$(echo "$HTTP_RESPONSE" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    HTTP_STATUS=$(echo "$HTTP_RESPONSE" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')

    if [ "$HTTP_STATUS" == "200" ]; then
      # extract TAG_LAST_PULLED
      TAG_LAST_PULLED=$(echo $HTTP_BODY | jq -r '.tag_last_pulled')
      # and check if it is null (e.g., not tracked yet)
      if [ "$TAG_LAST_PULLED" == "null" ]; then
        echo "$TAG"
      else
        # if it's != null, we need to compare it to TARGET_DATE
        echo "$HTTP_BODY" | jq -r --argjson date "$TARGET_DATE" 'select (.tag_last_pulled | sub(".[0-9]+Z$"; "Z") | fromdate < $date)|.name'
      fi
    fi
  done
}

# query all releases from current repo
RELEASE_TAGS=$(get_releases "keptn/keptn")


for IMAGE in "${IMAGES[@]}"; do
  echo "Detecting stale images for ${DOCKER_ORG}/${IMAGE} for all release tags"
  STALE_TAGS=$(check_if_stale "${IMAGE}" "${RELEASE_TAGS}")

  # pull each stale tag
  if [[ -n "$STALE_TAGS" ]]; then
    for TAG in ${STALE_TAGS[@]}; do
      echo "Pulling ${DOCKER_ORG}/${IMAGE}:${TAG} to ensure images are not stale..."
      docker pull "${DOCKER_ORG}/${IMAGE}:${TAG}"
    done
  else
    echo "All images are fine."
  fi
done
