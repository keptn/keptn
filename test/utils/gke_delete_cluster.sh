#!/bin/bash

# clean up any nightly clusters
# get all clusters
clusters=$(gcloud container clusters list --zone $CLOUDSDK_COMPUTE_ZONE --project $GCLOUD_PROJECT_NAME)

echo "Deleting nightly $CLUSTER_NAME_NIGHTLY cluster ..."

if echo "$clusters" | grep $CLUSTER_NAME_NIGHTLY; then
    gcloud container clusters delete $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $GCLOUD_PROJECT_NAME --quiet
    echo "Finished deleting nightly cluster $CLUSTER_NAME_NIGHTLY"
else 
    echo "No nightly cluster needs to be deleted"
fi
