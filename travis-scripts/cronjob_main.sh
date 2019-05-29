#!/usr/bin/env bash

source ./travis-scripts/setup_functions.sh

# prints the full command before output of the command.
set -x

install_hub
install_yq
install_helm

setup_gcloud_nightly

uninstall_keptn

delete_nightly_cluster
create_nightly_cluster

build_and_install_cli

keptn install --keptn-version=develop

export ISTIO_INGRESS=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
export_names

# Execute unit tests
execute_cli_tests

# Execute end-to-end test
cd test
source ./testOnboarding.sh

