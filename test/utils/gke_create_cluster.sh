#!/bin/bash

# configure gcloud
gcloud --quiet config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME_NIGHTLY
gcloud --quiet config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}

# clean up any nightly clusters
clusters=$(gcloud container clusters list --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME)
if echo "$clusters" | grep $CLUSTER_NAME_NIGHTLY; then 
    echo "Deleting nightly cluster..."
    gcloud container clusters delete $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME --quiet
    echo "Finished deleting nigtly cluster"
else 
    echo "No nightly cluster need to be deleted"
fi

# create a new cluster
gcloud container --project $PROJECT_NAME clusters create $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --username "admin" --cluster-version $GKE_VERSION --machine-type "n1-standard-4" --image-type "UBUNTU" --disk-type "pd-standard" --disk-size "100" --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "1" --enable-cloud-logging --enable-cloud-monitoring --no-enable-ip-alias --network "projects/sai-research/global/networks/default" --subnetwork "projects/sai-research/regions/$CLOUDSDK_REGION/subnetworks/default" --addons HorizontalPodAutoscaling,HttpLoadBalancing --no-enable-autoupgrade --no-enable-autorepair
if [[ $? != '0' ]]; then
    echo "gcloud cluster create failed."
    exit 1
fi

# get cluster credentials (this will set kubectl context)
gcloud container clusters get-credentials $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME
if [[ $? != '0' ]]; then
    echo "gcloud get credentials failed."
    exit 1
fi
