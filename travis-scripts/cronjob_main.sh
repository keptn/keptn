#!/usr/bin/env bash

gcloud --quiet config set project $PROJECT_NAME
gcloud --quiet config set container/cluster $CLUSTER_NAME_NIGHTLY
gcloud --quiet config set compute/zone ${CLOUDSDK_COMPUTE_ZONE}

clusters=$(gcloud container clusters list --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME)
if echo "$clusters" | grep $CLUSTER_NAME_NIGHTLY; then 
    echo "First delete old keptn installation"
    cd ./installer/scripts
    ./common/uninstallKeptn.sh
    cd ../..

    echo "Start deleting nightly cluster"
    gcloud container clusters delete $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME --quiet
    echo "Finished deleting nigtly cluster"
else 
    echo "No nightly cluster available"
fi

gcloud container --project $PROJECT_NAME clusters create $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --username "admin" --cluster-version "1.12.8-gke.10" --machine-type "n1-standard-4" --image-type "UBUNTU" --disk-type "pd-standard" --disk-size "100" --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --num-nodes "1" --enable-cloud-logging --enable-cloud-monitoring --no-enable-ip-alias --network "projects/sai-research/global/networks/default" --subnetwork "projects/sai-research/regions/$CLOUDSDK_REGION/subnetworks/default" --addons HorizontalPodAutoscaling,HttpLoadBalancing --no-enable-autoupgrade --no-enable-autorepair
if [[ $? != '0' ]]; then
    print_error "gcloud cluster create failed."
    exit 1
fi

gcloud container clusters get-credentials $CLUSTER_NAME_NIGHTLY --zone $CLOUDSDK_COMPUTE_ZONE --project $PROJECT_NAME
if [[ $? != '0' ]]; then
    print_error "gcloud get credentials failed."
    exit 1
fi

# Install hub
HUB_VERSION="v2.11.2"
HUB_INSTALLER="hub-darwin-amd64-2.11.2"

curl -L -s https://github.com/github/hub/releases/download/${HUB_VERSION}/${HUB_INSTALLER}.tgz  -o ${HUB_INSTALLER}.tgz
tar xopf ${HUB_INSTALLER}.tgz
sudo mv ${HUB_INSTALLER}/bin/hub /usr/local/bin/hub

# Build and install keptn CLI
cd cli/
dep ensure
go build -o keptn
sudo mv keptn /usr/local/bin/keptn
cd ..

# Prepare creds.json file
cd ./installer/scripts

export GITU=$GITHUB_USER_NAME_NIGHTLY	
export GITAT=$GITHUB_TOKEN_NIGHTLY	
export CLN=$CLUSTER_NAME_NIGHTLY	
export CLZ=$CLOUDSDK_COMPUTE_ZONE	
export PROJ=$PROJECT_NAME	
export GITO=$GITHUB_ORG_NIGHTLY	

source ./gke/defineCredentialsHelper.sh
replaceCreds

# Install keptn
keptn install --keptn-version=develop --creds=creds.json --verbose
cd ../..

# Execute end-to-end test
cd test
source ./testOnboarding.sh