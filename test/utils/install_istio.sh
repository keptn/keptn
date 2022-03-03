#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

# install istio
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.9.6 sh -
cd istio-* || exit
export PATH=$PWD/bin:$PATH

# run istio prechecks
istioctl x precheck

# install istio onto cluster
istioctl install -y --verify

verify_test_step $? "Failed to install Istio"

# verify the pods within istio-system
verify_deployment_in_namespace "istio-ingressgateway" "istio-system"
