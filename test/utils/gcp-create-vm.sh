#!/bin/bash

echo "VM Instance Name: $VM_INSTANCE_NAME"

# configure gcloud
gcloud --quiet config set project "${GCLOUD_PROJECT_NAME}"
gcloud --quiet config set compute/zone "${CLOUDSDK_COMPUTE_ZONE}"

echo "Creating VM Instance..."
gcloud beta compute instances create "$VM_INSTANCE_NAME" \
--zone=us-east1-b \
--machine-type=e2-standard-2 \
--subnet=default \
--network-tier=PREMIUM \
--no-restart-on-failure \
--maintenance-policy=TERMINATE \
--preemptible --no-service-account \
--no-scopes \
--image=ubuntu-2004-focal-v20210825 \
--image-project=ubuntu-os-cloud \
--boot-disk-size=10GB \
--boot-disk-type=pd-balanced \
--boot-disk-device-name=gh-nightly-tutorial-runner \
--shielded-secure-boot \
--shielded-vtpm \
--shielded-integrity-monitoring \
--reservation-affinity=any
