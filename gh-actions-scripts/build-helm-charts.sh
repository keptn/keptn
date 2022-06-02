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

# ####################
# INSTALLER HELM CHART
# ####################
INSTALLER_BASE_PATH=installer/manifests

helm dependency build ${INSTALLER_BASE_PATH}/keptn/charts/control-plane

helm package ${INSTALLER_BASE_PATH}/keptn --app-version "$IMAGE_TAG" --version "$VERSION"
if [ $? -ne 0 ]; then
  echo "Error packing installer, exiting..."
  exit 1
fi

mv "keptn-${VERSION}.tgz" "keptn-charts/keptn-${VERSION}.tgz"

# verify the chart
echo "::group::Template install of keptn"
helm template --debug "keptn-charts/keptn-${VERSION}.tgz"

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for installer has templating errors - exiting"
  echo "::endgroup::"
  exit 1
fi

echo "::endgroup::"

# verify the chart with install
echo "::group::Dry run install of keptn"
helm install --dryrun "keptn-charts/keptn-${VERSION}.tgz"

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for installer has runtime errors - exiting"
  echo "::endgroup::"
  exit 1
fi
echo "::endgroup::"


# ####################
# HELM-SVC HELM CHART
# ####################
HELM_SVC_BASE_PATH=helm-service

helm dependency build ${HELM_SVC_BASE_PATH}/chart

helm package ${HELM_SVC_BASE_PATH}/chart --app-version "$IMAGE_TAG" --version "$VERSION"
if [ $? -ne 0 ]; then
  echo "Error packaging installer, exiting..."
  exit 1
fi

mv "helm-service-${VERSION}.tgz" "keptn-charts/helm-service-${VERSION}.tgz"

#verify the chart
echo "::group::Template of helm-service"
helm template --debug "keptn-charts/helm-service-${VERSION}.tgz"

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for helm-svc has templating errors -exiting"
  echo "::endgroup::"
  exit 1
fi
echo "::endgroup::"

# verify the chart with install
echo "::group::Dry run install of helm-service"
helm install --dryrun "keptn-charts/helm-service-${VERSION}.tgz" > out.yml

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for installer has runtime errors - exiting"
  echo "::endgroup::"
  exit 1
fi
echo "::endgroup::"

# ####################
# JMETER-SVC HELM CHART
# ####################
JMETER_SVC_BASE_PATH=jmeter-service

helm dependency build ${JMETER_SVC_BASE_PATH}/chart

helm package ${JMETER_SVC_BASE_PATH}/chart --app-version "$IMAGE_TAG" --version "$VERSION"
if [ $? -ne 0 ]; then
  echo "Error packaging installer, exiting..."
  exit 1
fi

mv "jmeter-service-${VERSION}.tgz" "keptn-charts/jmeter-service-${VERSION}.tgz"

echo "::group::Template of jmeter-service"
#verify the chart
helm template --debug "keptn-charts/jmeter-service-${VERSION}.tgz"

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for jmeter-svc has templating errors -exiting"
  echo "::endgroup::"
  exit 1
fi
echo "::endgroup::"


# verify the chart with install
echo "::group::Dry run install of jmeter-service"
helm install --dryrun "keptn-charts/jmeter-service-${VERSION}.tgz"

if [ $? -ne 0 ]; then
  echo "::error::Helm Chart for installer has runtime errors - exiting"
  echo "::endgroup::"
  exit 1
fi
echo "::endgroup::"

echo "Generated files:"
echo " - keptn-charts/keptn-${VERSION}.tgz"
echo " - keptn-charts/helm-service-${VERSION}.tgz"
echo " - keptn-charts/jmeter-service-${VERSION}.tgz"
