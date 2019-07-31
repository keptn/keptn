#!/bin/bash
CLUSTER_NAME=$1

LOG_LOCATION=.
exec > >(tee -i $LOG_LOCATION/setupPRstatuscheckCluster.log)
exec 2>&1

gcloud beta container --project "sai-research" clusters create "$CLUSTER_NAME" --zone "us-central1-a" --username "admin" --cluster-version "1.11.6-gke.2" --machine-type "n1-standard-1" --image-type "UBUNTU" --disk-type "pd-standard" --disk-size "100" --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "3" --enable-cloud-logging --enable-cloud-monitoring --no-enable-ip-alias --network "projects/sai-research/global/networks/default" --subnetwork "projects/sai-research/regions/us-central1/subnetworks/default" --enable-autoscaling --min-nodes "2" --max-nodes "7" --addons HorizontalPodAutoscaling,HttpLoadBalancing --no-enable-autoupgrade --no-enable-autorepair

gcloud container clusters get-credentials $CLUSTER_NAME --zone us-central1-a --project sai-research

export GCLOUD_USER=$(gcloud config get-value account)
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER

kubectl create -f ../manifests/k8s-namespaces.yml 

kubectl apply -f ../manifests/istio/istio-crds.yml
kubectl apply -f ../manifests/istio/istio-demo.yml
