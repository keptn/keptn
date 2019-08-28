#!/bin/bash
source ./common/utils.sh

kubectl create namespace istio-system

helm template ../manifests/istio --name istio --namespace istio-system --values ../manifests/istio/values.yaml | kubectl apply -f -
verify_kubectl $? "Installing Istio failed."
wait_for_deployment_in_namespace "istio-ingressgateway" "istio-system"
wait_for_deployment_in_namespace "istio-pilot" "istio-system"
wait_for_deployment_in_namespace "istio-citadel" "istio-system"
wait_for_deployment_in_namespace "istio-sidecar-injector" "istio-system"
wait_for_all_pods_in_namespace "istio-system"