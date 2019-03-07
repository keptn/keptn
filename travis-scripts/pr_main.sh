#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# Causes the shell to exit immediately if a simple command exits with a nonzero exit value as well as
# prints the full command before output of the command.
set -e -x

setup_gcloud
setup_glcoud_pr
install_yq
setup_knative
export_names
execute_core_component_tests
execute_cli_tests