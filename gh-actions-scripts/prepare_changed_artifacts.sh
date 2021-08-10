#!/bin/bash

CHANGED_FILES=$1

if [ $# -ne 1 ]; then
  echo "Usage: $0 CHANGED_FILES"
  exit
fi

# initialize variables with false (make sure they are also set in needs.prepare_ci_run.outputs !!!)
BUILD_INSTALLER=false
BUILD_API=false
BUILD_CLI=false
BUILD_OS_ROUTE_SVC=false
BUILD_BRIDGE=false
BUILD_JMETER=false
BUILD_HELM_SVC=false
BUILD_APPROVAL_SVC=false
BUILD_DISTRIBUTOR=false
BUILD_SHIPYARD_CONTROLLER=false
BUILD_SECRET_SVC=false
BUILD_CONFIGURATION_SVC=false
BUILD_REMEDIATION_SVC=false
BUILD_LIGHTHOUSE_SVC=false
BUILD_MONGODB_DS=false
BUILD_STATISTICS_SVC=false

artifacts=(
  "$BRIDGE_ARTIFACT_PREFIX"
  "$API_ARTIFACT_PREFIX"
  "$OS_ROUTE_SVC_ARTIFACT_PREFIX"
  "$JMETER_SVC_ARTIFACT_PREFIX"
  "$HELM_SVC_ARTIFACT_PREFIX"
  "$APPROVAL_SVC_ARTIFACT_PREFIX"
  "$DISTRIBUTOR_ARTIFACT_PREFIX"
  "$SHIPYARD_CONTROLLER_ARTIFACT_PREFIX"
  "$SECRET_SVC_ARTIFACT_PREFIX"
  "$CONFIGURATION_SVC_ARTIFACT_PREFIX"
  "$REMEDIATION_SVC_ARTIFACT_PREFIX"
  "$LIGHTHOUSE_SVC_ARTIFACT_PREFIX"
  "$MONGODB_DS_ARTIFACT_PREFIX"
  "$STATISTICS_SVC_ARTIFACT_PREFIX"
)

echo "changed files:"
echo "$CHANGED_FILES"
matrix_config='{"config":['
# shellcheck disable=SC2016
build_artifact_template='{"artifact": $artifact, "working-dir": $working_dir, "should-run": $should_run, "test-folders": $test_folders, "go-flags": $go_flags }'

for changed_file in $CHANGED_FILES; do
  echo "Checking if $changed_file leads to a build..."

  if [[ $changed_file == "${INSTALLER_FOLDER}"* ]]; then
    echo "Found changes in Installer"
    BUILD_INSTALLER=true
    continue
  fi

  if [[ $changed_file == "${CLI_FOLDER}"* ]]; then
    echo "Found changes in CLI"
    BUILD_CLI=true
    continue
  fi

  for artifact in "${artifacts[@]}"; do
    # Prepare variables
    artifact_fullname="${artifact}_ARTIFACT"
    artifact_folder="${artifact}_FOLDER"
    should_build_artifact="BUILD_${artifact}"
    artifact_go_flags="${artifact}_GO_FLAGS"
    artifact_test_folders="${artifact}_TEST_FOLDERS"

    if [[ ( $changed_file == ${!artifact_folder}* ) && ( "${!should_build_artifact}" != 'true' ) && ( $BUILD_EVERYTHING != 'true' ) ]]; then
      echo "Found changes in $artifact"
      IFS= read -r "${should_build_artifact}" <<< "true"
      artifact_config=$(jq -n \
        --arg artifact "${!artifact_fullname}" \
        --arg working_dir "${!artifact_folder}" \
        --arg should_run "${!should_build_artifact}" \
        --arg test_folders "${!artifact_test_folders}" \
        --arg go_flags "${!artifact_go_flags}" \
        "$build_artifact_template"
      )
      matrix_config="$matrix_config $artifact_config,"
    elif [[ ( $BUILD_EVERYTHING == 'true' ) && ( "${!should_build_artifact}" != 'true' ) ]]; then
      echo "No changes in $artifact but build is set to build everything"
      artifact_config=$(jq -n \
        --arg artifact "${!artifact_fullname}" \
        --arg working_dir "${!artifact_folder}" \
        --arg should_run "false" \
        --arg test_folders "${!artifact_test_folders}" \
        --arg go_flags "${!artifact_go_flags}" \
        "$build_artifact_template"
      )
      matrix_config="$matrix_config $artifact_config,"
    fi
  done
done

# Terminate matrix JSON config and remove trailing comma
matrix_config="${matrix_config%,}]}"

# Escape newlines for multiline string support in GH actions
# Reference: https://github.community/t/set-output-truncates-multiline-strings/16852
matrix_config="${matrix_config//'%'/'%25'}"
matrix_config="${matrix_config//$'\n'/'%0A'}"
matrix_config="${matrix_config//$'\r'/'%0D'}"

# print job outputs (make sure they are also set in needs.prepare_ci_run.outputs !!!)
echo "::set-output name=BUILD_INSTALLER::$BUILD_INSTALLER"
echo "::set-output name=BUILD_API::$BUILD_API"
echo "::set-output name=BUILD_CLI::$BUILD_CLI"
echo "::set-output name=BUILD_OS_ROUTE_SVC::$BUILD_OS_ROUTE_SVC"
echo "::set-output name=BUILD_BRIDGE::$BUILD_BRIDGE"
echo "::set-output name=BUILD_JMETER::$BUILD_JMETER"
echo "::set-output name=BUILD_HELM_SVC::$BUILD_HELM_SVC"
echo "::set-output name=BUILD_APPROVAL_SVC::$BUILD_APPROVAL_SVC"
echo "::set-output name=BUILD_DISTRIBUTOR::$BUILD_DISTRIBUTOR"
echo "::set-output name=BUILD_SHIPYARD_CONTROLLER::$BUILD_SHIPYARD_CONTROLLER"
echo "::set-output name=BUILD_SECRET_SVC::$BUILD_SECRET_SVC"
echo "::set-output name=BUILD_CONFIGURATION_SVC::$BUILD_CONFIGURATION_SVC"
echo "::set-output name=BUILD_REMEDIATION_SVC::$BUILD_REMEDIATION_SVC"
echo "::set-output name=BUILD_LIGHTHOUSE_SVC::$BUILD_LIGHTHOUSE_SVC"
echo "::set-output name=BUILD_MONGODB_DS::$BUILD_MONGODB_DS"
echo "::set-output name=BUILD_STATISTICS_SVC::$BUILD_STATISTICS_SVC"
echo "::set-output name=BUILD_MATRIX::$matrix_config"

echo "The following artifacts will be tested and built:"
echo "BUILD_INSTALLER: $BUILD_INSTALLER"
echo "BUILD_API: $BUILD_API"
echo "BUILD_CLI: $BUILD_CLI"
echo "BUILD_OS_ROUTE_SVC: $BUILD_OS_ROUTE_SVC"
echo "BUILD_BRIDGE: $BUILD_BRIDGE"
echo "BUILD_JMETER: $BUILD_JMETER"
echo "BUILD_HELM_SVC: $BUILD_HELM_SVC"
echo "BUILD_APPROVAL_SVC: $BUILD_APPROVAL_SVC"
echo "BUILD_DISTRIBUTOR: $BUILD_DISTRIBUTOR"
echo "BUILD_SHIPYARD_CONTROLLER: $BUILD_SHIPYARD_CONTROLLER"
echo "BUILD_SECRET_SVC: $BUILD_SECRET_SVC"
echo "BUILD_CONFIGURATION_SVC: $BUILD_CONFIGURATION_SVC"
echo "BUILD_REMEDIATION_SVC: $BUILD_REMEDIATION_SVC"
echo "BUILD_LIGHTHOUSE_SVC: $BUILD_LIGHTHOUSE_SVC"
echo "BUILD_MONGODB_DS: $BUILD_MONGODB_DS"
echo "BUILD_STATISTICS_SVC: $BUILD_STATISTICS_SVC"

if [[ "$matrix_config" == '{"config":[]}' ]]; then
  echo "Build matrix is emtpy, setting output..."
  echo "::set-output name=BUILD_MATRIX_EMPTY::true"
else
  echo "Build matrix is NOT emtpy, setting output..."
  echo "::set-output name=BUILD_MATRIX_EMPTY::false"
fi
