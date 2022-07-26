#!/bin/bash
# shellcheck disable=SC2181

VERSION=$1
IMAGE_TAG=$2
KEPTN_SPEC_VERSION=$3
DOCKER_ORG=$4

if [ $# -lt 3 ]; then
  echo "Usage: $0 VERSION IMAGE_TAG KEPTN_SPEC_VERSION"
  exit
fi

if [ -z "$VERSION" ]; then
  echo "No Version set, exiting..."
  exit 1
fi

if [ -z "$IMAGE_TAG" ]; then
  echo "No Image Tag set, defaulting to version"
  IMAGE_TAG=$VERSION
fi

if [ -z "$DOCKER_ORG" ]; then
  echo "No Docker organisation set, defaulting to keptn"
  DOCKER_ORG='keptn'
fi


# replace "appVersion: latest" with "appVersion: $VERSION" in all Chart.yaml files
find . -name Chart.yaml -exec sed -i -- "s/appVersion: latest/appVersion: ${IMAGE_TAG}/g" {} \;
find . -name Chart.yaml -exec sed -i -- "s/version: latest/version: ${VERSION}/g" {} \;
# replace "keptnSpecVersion: latest" with "keptnSpecVersion: $KEPTN_SPEC_VERSION" in all values.yaml files
find . -name values.yaml -exec sed -i -- "s/keptnSpecVersion: latest/keptnSpecVersion: ${KEPTN_SPEC_VERSION}/g" {} \;
find . -name values.yaml -exec sed -i -- "s/docker.io\/keptn\//docker.io\/${DOCKER_ORG}\//g" {} \;

mkdir keptn-charts/

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add nats https://nats-io.github.io/k8s/helm/charts/

# ####################
# COMMON HELM CHART
# ####################
COMMON_CHART_BASE_PATH=installer/manifests/common

helm package ${COMMON_CHART_BASE_PATH} --version "$VERSION"
if [ $? -ne 0 ]; then
  echo "Error packaging common chart, exiting..."
  exit 1
fi

mv "common-${VERSION}.tgz" "keptn-charts/common-${VERSION}.tgz"

# ##################################################################
# INSTALLER HELM CHART # HELM-SVC HELM CHART # JMETER-SVC HELM CHART
# ##################################################################

declare -A charts
charts[keptn]=installer/manifests/keptn
charts[helm-service]=helm-service/chart
charts[jmeter-service]=jmeter-service/chart

for i in "${!charts[@]}"
do
  echo "=== Building $i ==="
  BASE_NAME=$i
  BASE_PATH=${charts[$i]}

  echo "::group::Helm dependency build"
  helm dependency build ${BASE_PATH}
  echo "::endgroup::"

  helm package ${BASE_PATH} --app-version "$IMAGE_TAG" --version "$VERSION"
  if [ $? -ne 0 ]; then
    echo "Error packing ${BASE_NAME}, exiting..."
    exit 1
  fi

  mv "${BASE_NAME}-${VERSION}.tgz" "keptn-charts/${BASE_NAME}-${VERSION}.tgz"

  # verify the chart
  echo "::group::Template install of ${BASE_NAME}"
  helm template --debug "keptn-charts/${BASE_NAME}-${VERSION}.tgz"

  if [ $? -ne 0 ]; then
    echo "::error::Helm Chart for ${BASE_NAME} has templating errors - exiting"
    echo "::endgroup::"
    exit 1
  fi

  echo "::endgroup::"

done

echo "Generated files:"
echo " - keptn-charts/keptn-${VERSION}.tgz"
echo " - keptn-charts/helm-service-${VERSION}.tgz"
echo " - keptn-charts/jmeter-service-${VERSION}.tgz"
