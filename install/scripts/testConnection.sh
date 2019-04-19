#!/bin/bash
GKE_PROJECT=$1
CLUSTER_NAME=$2
CLUSTER_ZONE=$3

gcloud --quiet config set project $GKE_PROJECT
gcloud --quiet config set container/cluster $CLUSTER_NAME
gcloud --quiet config set compute/zone $CLUSTER_ZONE
gcloud container clusters get-credentials $CLUSTER_NAME --zone $CLUSTER_ZONE --project $GKE_PROJECT

if [[ $? != '0' ]]
then
  echo -e "[keptn|0]Could not connect to cluster. Please ensure you have set the correct values for your Cluster Name, GKE Project, and Cluster Zone during the credentials setup."
  exit 1
fi

if [[ $? = '0' ]]
then
  echo -e "[keptn|0]Connection to cluster successful."
fi
