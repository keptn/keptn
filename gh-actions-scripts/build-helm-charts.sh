#!/bin/bash
# shellcheck disable=SC2181

VERSION=$1
IMAGE_TAG=$2
KEPTN_SPEC_VERSION=$3
DOCKER_ORG=$4
SIGN_CHART=$5

if [ $# -lt 4 ]; then
  echo "Usage: $0 VERSION IMAGE_TAG KEPTN_SPEC_VERSION SIGN_CHART"
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

if [ -z "$SIGN_CHART" ]; then
  echo "No flag set for signing charts, defaulting to no signing"
  SIGN_CHART=false
fi

if [[ "$SIGN_CHART" == 'true' ]]; then
  if [ -z "$SIGNING_KEY_BASE64" ] || [ -z "$SIGNING_KEY_PASSPHRASE_BASE64" ] || [ -z "$SIGNING_KEY_NAME" ]; then
    echo 'The following variable need to be set to enable Helm chart signing: SIGNING_KEY_BASE64 SIGNING_KEY_PASSPHRASE_BASE64 SIGNING_KEY_NAME'
    exit 2
  fi

  echo "Creating necessary files for chart signing..."

  KEY_PATH='.gpg-dir'
  SIGNING_KEY_PATH="$KEY_PATH/secring.gpg"
  SIGNING_KEY_PASSPHRASE_PATH="$KEY_PATH/passphrase"
  mkdir "$KEY_PATH"
  base64 -d <<< "$SIGNING_KEY_BASE64" > "$SIGNING_KEY_PATH"
  base64 -d <<< "$SIGNING_KEY_PASSPHRASE_BASE64" > "$SIGNING_KEY_PASSPHRASE_PATH"

  echo "Done..."
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

if [[ "$SIGN_CHART" == 'true' ]]; then
  echo "Packaging chart with signage..."
  # shellcheck disable=SC2002
  cat "$SIGNING_KEY_PASSPHRASE_PATH" | helm package ${COMMON_CHART_BASE_PATH} --version "$VERSION" --sign --key "$SIGNING_KEY_NAME" --keyring "$SIGNING_KEY_PATH" --passphrase-file -
  mv "common-${VERSION}.tgz.prov" 'keptn-charts/'
else
  echo "Packaging chart without signage..."
  helm package ${COMMON_CHART_BASE_PATH} --version "$VERSION"
fi

if [ $? -ne 0 ]; then
  echo "Error packaging common chart, exiting..."
  exit 1
fi

mv "common-${VERSION}.tgz" 'keptn-charts/'

# ####################
# INSTALLER HELM CHART
# ####################

declare -A charts
charts[keptn]=installer/manifests/keptn

for i in "${!charts[@]}"
do
  echo "=== Building $i ==="
  BASE_NAME=$i
  BASE_PATH=${charts[$i]}

  echo "::group::Helm dependency build"
  helm dependency build ${BASE_PATH}
  echo "::endgroup::"

  if [[ "$SIGN_CHART" == 'true' ]]; then
    # shellcheck disable=SC2002
    cat "$SIGNING_KEY_PASSPHRASE_PATH" | helm package ${BASE_PATH} --app-version "$IMAGE_TAG" --version "$VERSION" --sign --key "$SIGNING_KEY_NAME" --keyring "$SIGNING_KEY_PATH" --passphrase-file -
    mv "${BASE_NAME}-${VERSION}.tgz.prov" 'keptn-charts/'
  else
    helm package ${BASE_PATH} --app-version "$IMAGE_TAG" --version "$VERSION"
  fi

  if [ $? -ne 0 ]; then
    echo "Error packing ${BASE_NAME}, exiting..."
    exit 1
  fi

  mv "${BASE_NAME}-${VERSION}.tgz" 'keptn-charts/'

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
ls -l 'keptn-charts/'
