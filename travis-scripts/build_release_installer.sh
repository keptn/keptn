#!/bin/bash

VERSION=$1
if [ -n "$VERSION" ]; then
  echo 'No Version set, exiting...'
  exit 1
fi

helm package keptn --app-version $VERSION --version $VERSION
if [ $? -ne 0 ]; then
  echo 'Error packing installer, exiting...'
  exit 1
fi

mkdir keptn-charts/
mv keptn-${VERSION}.tgz keptn-charts/

gsutil cp gs://keptn-installer/index.yaml keptn-charts/index.yaml

helm repo index keptn-charts --url https://storage.googleapis.com/keptn-installer/ --merge keptn-charts/index.yaml
if [ $? -ne 0 ]; then
  echo 'Error generating index.yaml, exiting...'
  exit 1
fi

# upload to gcloud
gsutil cp keptn-charts/index.yaml gs://keptn-installer/index.yaml
gsutil cp keptn-charts/keptn-${VERSION}.tgz gs://keptn-installer/keptn-${VERSION}.tgz