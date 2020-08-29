#!/bin/bash

# configure gcloud
gcloud --quiet config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME_NIGHTLY
gcloud --quiet config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}

# clean up any nightly clusters
clusters=$(gcloud container clusters list --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME)
if echo "$clusters" | grep $CLUSTER_NAME_NIGHTLY; then 
    echo "Deleting nightly cluster ${CLUSTER_NAME_NIGHTLY} ..."
    gcloud container clusters delete $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME --quiet
    echo "Finished deleting nightly cluster"
else 
    echo "No nightly cluster need to be deleted"
fi

ISTIO_CONFIG="--istio-config=auth=MTLS_PERMISSIVE"
ADDONS="Istio,HorizontalPodAutoscaling,HttpLoadBalancing"

echo "Creating nightly cluster ${CLUSTER_NAME_NIGHTLY}"

# create a new cluster (Note: disk-size reduced to 25 GB to save resources; pre-emptible nodes used as well)
gcloud beta container --project $PROJECT_NAME clusters create $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --username "admin" --cluster-version $GKE_VERSION \
 --machine-type "n1-standard-8" --image-type "UBUNTU" --preemptible --disk-type "pd-standard" --disk-size "25" \
 --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "1" --enable-cloud-logging --enable-cloud-monitoring --enable-ip-alias --network "projects/sai-research/global/networks/default" --subnetwork "projects/sai-research/regions/$CLOUDSDK_REGION/subnetworks/default" \
 --addons $ADDONS $ISTIO_CONFIG --no-enable-autoupgrade --no-enable-autorepair \
 --labels owner=travis,expiry=auto-delete

if [[ $? != '0' ]]; then
    echo "gcloud cluster create failed"
    exit 1
fi
