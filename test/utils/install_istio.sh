#!/bin/bash

source test/utils.sh

# install istio
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.8.2 sh -
cd istio-*
export PATH=$PWD/bin:$PATH
istioctl install --set profile=demo

# verify the pods within istio-system
verify_deployment_in_namespace "istio-ingressgateway" "istio-system"
