#!/bin/bash

# configure gcloud
gcloud --quiet config set project $GCLOUD_PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME_NIGHTLY
gcloud --quiet config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}

# get cluster credentials (this will set kubectl context)
gcloud container clusters get-credentials "$CLUSTER_NAME_NIGHTLY" --zone "$CLOUDSDK_COMPUTE_ZONE" --project "$GCLOUD_PROJECT_NAME"
if [[ $? != '0' ]]; then
    echo "gcloud get credentials failed"
    exit 1
fi
