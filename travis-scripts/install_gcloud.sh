#!/bin/bash

export GOOGLE_APPLICATION_CREDENTIALS=~/gcloud-service-key.json

# check if gcloud is already installed
if [ ! -d "$HOME/google-cloud-sdk/bin" ]; then
  rm -rf $HOME/google-cloud-sdk;
  export CLOUDSDK_CORE_DISABLE_PROMPTS=1;
  
  # download
  curl https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-289.0.0-${OS_TYPE}-x86_64.tar.gz > gcloud.tar.gz
  gunzip -c gcloud.tar.gz | tar xopf -
  ./google-cloud-sdk/install.sh
  source ./google-cloud-sdk/completion.bash.inc
  source ./google-cloud-sdk/path.bash.inc
fi

# update
gcloud --quiet components update
gcloud --quiet components update kubectl
