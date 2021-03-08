#!/bin/bash

# shellcheck disable=SC1091
source test/utils.sh

# install istio
cd istio-* || exit
export PATH=$PWD/bin:$PATH

istioctl manifest generate --set profile=demo | kubectl delete -f -
kubectl delete namespace istio-system
