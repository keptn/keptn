#!/bin/bash

# Environment variables for test connection to cluster
if [[ -z "${GKE_PROJECT}" ]]; then
  echo "[keptn|DEBUG] GKE_PROJECT not set, take it from creds.json"
  GKE_PROJECT=$(cat creds.json | jq -r '.gkeProject')
  # TODO: break installation when GKE_PROJECT empty
fi

if [[ -z "${CLUSTER_NAME}" ]]; then
  echo "[keptn|DEBUG] CLUSTER_NAME not set, take it from creds.json"
  CLUSTER_NAME=$(cat creds.json | jq -r '.clusterName')
  # TODO: break installation when GKE_PROJECT empty
fi

if [[ -z "${CLUSTER_ZONE}" ]]; then
  echo "[keptn|DEBUG] CLUSTER_ZONE not set, take it from creds.json"
  CLUSTER_ZONE=$(cat creds.json | jq -r '.clusterZone')
  # TODO: break installation when GKE_PROJECT empty
fi

gcloud --quiet config set project $GKE_PROJECT
gcloud --quiet config set container/cluster $CLUSTER_NAME
gcloud --quiet config set compute/zone $CLUSTER_ZONE
gcloud container clusters get-credentials $CLUSTER_NAME --zone $CLUSTER_ZONE --project $GKE_PROJECT

if [[ $? != '0' ]]
then
  echo "[keptn|ERROR] Could not connect to cluster. Please ensure you have set the correct values for your Cluster Name, GKE Project, and Cluster Zone during the credentials setup."
  exit 1
else
  echo "[keptn|INFO] Connection to cluster successful."
fi
