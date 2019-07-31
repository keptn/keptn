#!/bin/bash
CLUSTER_NAME=$1
CLUSTER_ZONE=$2

source ./common/utils.sh

# Variables for test connection to cluster
if [[ -z "${GKE_PROJECT}" ]]; then
  print_debug "GKE_PROJECT not set, take it from creds.json"
  GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')
  verify_variable "$GKE_PROJECT" "GKE_PROJECT is not defined in environment variable nor in creds.json file." 
fi

gcloud --quiet config set project $GKE_PROJECT
gcloud --quiet config set container/cluster $CLUSTER_NAME
gcloud --quiet config set compute/zone $CLUSTER_ZONE
gcloud container clusters get-credentials $CLUSTER_NAME --zone $CLUSTER_ZONE --project $GKE_PROJECT

if [[ $? != '0' ]]
then
  exit 1
fi
