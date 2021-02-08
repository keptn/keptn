#!/bin/bash

export GOOGLE_APPLICATION_CREDENTIALS=~/gcloud-service-key.json

OS_TYPE=${OS_TYPE:-"linux"}
GCLOUD_VERSION="324.0.0"
GCLOUD_FILENAME="google-cloud-sdk-${GCLOUD_VERSION}-${OS_TYPE}-x86_64.tar.gz"

if [[ ! -f ~/downloads/${GCLOUD_FILENAME} ]]; then
  echo "Downloading ${GCLOUD_FILENAME}"
  wget "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/${GCLOUD_FILENAME}" -O ~/downloads/${GCLOUD_FILENAME}
fi

export CLOUDSDK_CORE_DISABLE_PROMPTS=1;

gunzip -c ~/downloads/${GCLOUD_FILENAME} | tar xopf -
./google-cloud-sdk/install.sh
source ./google-cloud-sdk/completion.bash.inc
source ./google-cloud-sdk/path.bash.inc

# update
gcloud --quiet components update
gcloud --quiet components update kubectl
