#!/bin/bash

source test/utils.sh

# install istio
cd istio-*
export PATH=$PWD/bin:$PATH

istioctl manifest generate --set profile=demo | kubectl delete -f -
kubectl delete namespace istio-system
