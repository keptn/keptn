#!/bin/bash
# shellcheck disable=SC2181

# For official releases we want to upload a helm-chart to google cloud bucket
# This step involves downloading index.yaml and updating it
# as well as uploading the helm-chart with index.yaml

LOCAL_DIRECTORY=${1:-"keptn-charts"}
TARGET_BUCKET=${2:-"keptn-installer"}

# download index.yaml
gsutil cp gs://keptn-installer/index.yaml "${LOCAL_DIRECTORY}/index.yaml"

helm repo index "${LOCAL_DIRECTORY}" --url https://storage.googleapis.com/keptn-installer/ --merge "${LOCAL_DIRECTORY}/index.yaml"
if [ $? -ne 0 ]; then
  echo "Error generating index.yaml, exiting..."
  exit 1
fi

# upload to gcloud
gsutil cp "${LOCAL_DIRECTORY}/index.yaml" "gs://${TARGET_BUCKET}/index.yaml"
gsutil cp "${LOCAL_DIRECTORY}"/* "gs://${TARGET_BUCKET}/"
