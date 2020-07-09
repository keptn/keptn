#!/bin/bash

VERSION=${1:-latest}

# Note: Currently, the Helm chart always has version 0.1.0 and
# $VERSION cannot be used as it is an invalid Semantic Version

cd keptn/charts/control-plane
helm dependency build
cd ../../..

helm package keptn --app-version $VERSION
if [ $? -ne 0 ]; then
  echo 'Error packing installer, exiting...'
  exit 1
fi

mkdir keptn-charts/
mv keptn-0.1.0.tgz keptn-charts/

helm repo index keptn-charts --url https://storage.googleapis.com/keptn-installer/${VERSION}
if [ $? -ne 0 ]; then
  echo 'Error generating index.yaml, exiting...'
  exit 1
fi

# upload to gcloud
gsutil cp keptn-charts/index.yaml gs://keptn-installer/${VERSION}/index.yaml
gsutil cp keptn-charts/keptn-0.1.0.tgz gs://keptn-installer/${VERSION}/keptn-0.1.0.tgz