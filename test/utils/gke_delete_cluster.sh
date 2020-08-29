#!/bin/bash

# clean up any nightly clusters
clusters=$(gcloud container clusters list --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME)
if echo "$clusters" | grep $CLUSTER_NAME_NIGHTLY; then 
    echo "Deleting nightly cluster ..."
    gcloud container clusters delete $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME --quiet
    echo "Finished deleting nightly cluster"
else 
    echo "No nightly cluster need to be deleted"
fi
