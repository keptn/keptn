#!/bin/bash

VERSION=$1
IMAGE_TAG=$2

if [ -z "$VERSION" ]; then
  echo "No Version set, exiting..."
  exit 1
fi

if [ -z "$IMAGE_TAG" ]; then
  echo "No Image Tag set, defaulting to version"
  IMAGE_TAG=$VERSION
fi

BASE_PATH=installer/manifests

helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm dependency build ${BASE_PATH}/keptn/charts/control-plane

# replace "appVersion: latest" with "appVersion: $VERSION" in all Chart.yaml files
find -name Chart.yaml -exec sed -i -- "s/appVersion: latest/appVersion: ${IMAGE_TAG}/g" {} \;
find -name Chart.yaml -exec sed -i -- "s/version: latest/version: ${VERSION}/g" {} \;

helm package ${BASE_PATH}/keptn --app-version $VERSION --version $VERSION
if [ $? -ne 0 ]; then
  echo "Error packing installer, exiting..."
  exit 1
fi

mkdir keptn-charts/
mv keptn-${VERSION}.tgz keptn-charts/keptn-installer-${VERSION}.tgz

# verify the chart
helm template --debug keptn-charts/keptn-installer-${VERSION}.tgz

if [ $? -ne 0 ]; then
  echo "::error Helm Chart has templating errors - exiting"
  exit 1
fi

# download index.yaml chart
#gsutil cp gs://keptn-installer/index.yaml keptn-charts/index.yaml
#
#helm repo index keptn-charts --url https://storage.googleapis.com/keptn-installer/ --merge keptn-charts/index.yaml
#if [ $? -ne 0 ]; then
#  echo "Error generating index.yaml, exiting..."
#  exit 1
#fi
#
## upload to gcloud
#gsutil cp keptn-charts/index.yaml gs://keptn-installer/index.yaml
#gsutil cp keptn-charts/keptn-${VERSION}.tgz gs://keptn-installer/keptn-${VERSION}.tgz



echo "Generated files:"
echo " - keptn-charts/keptn-installer-${VERSION}.tgz"


