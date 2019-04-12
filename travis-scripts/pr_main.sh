#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# prints the full command before output of the command.
set -x

setup_gcloud
setup_glcoud_pr
install_yq
setup_knative_pr
export_names
execute_core_component_tests
execute_cli_tests