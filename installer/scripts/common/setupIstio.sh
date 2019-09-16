#!/bin/bash
source ./common/utils.sh

kubectl create namespace istio-system

helm template ../manifests/istio/helm/istio-init --name istio-init --namespace istio-system | kubectl apply -f -
verify_kubectl $? "Creating Istio resources failed"
wait_for_crds "adapters.config.istio.io,attributemanifests.config.istio.io,authorizationpolicies.rbac.istio.io,clusterrbacconfigs.rbac.istio.io,destinationrules.networking.istio.io,envoyfilters.networking.istio.io,gateways.networking.istio.io,handlers.config.istio.io,httpapispecbindings.config.istio.io,httpapispecs.config.istio.io,instances.config.istio.io,meshpolicies.authentication.istio.io,policies.authentication.istio.io,quotaspecbindings.config.istio.io,quotaspecs.config.istio.io,rbacconfigs.rbac.istio.io,rules.config.istio.io,serviceentries.networking.istio.io,servicerolebindings.rbac.istio.io,serviceroles.rbac.istio.io,sidecars.networking.istio.io,templates.config.istio.io,virtualservices.networking.istio.io"
helm template ../manifests/istio/helm/istio --name istio --namespace istio-system --values ../manifests/istio/helm/istio/values-istio-minimal.yaml | kubectl apply -f -
verify_kubectl $? "Installing Istio failed."
wait_for_deployment_in_namespace "istio-ingressgateway" "istio-system"
wait_for_deployment_in_namespace "istio-pilot" "istio-system"
wait_for_deployment_in_namespace "istio-citadel" "istio-system"
wait_for_deployment_in_namespace "istio-sidecar-injector" "istio-system"
wait_for_all_pods_in_namespace "istio-system"