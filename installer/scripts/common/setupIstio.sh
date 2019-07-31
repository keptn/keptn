#!/bin/bash
source ./common/utils.sh

kubectl create namespace istio-system

# Apply custom resource definitions for Istio
kubectl apply -f ../manifests/istio/crd-10.yaml
verify_kubectl $? "Creating istio custom resource definitions failed."
kubectl apply -f ../manifests/istio/crd-11.yaml
verify_kubectl $? "Creating istio custom resource definitions failed."
kubectl apply -f ../manifests/istio/crd-12.yaml
verify_kubectl $? "Creating istio custom resource definitions failed."
kubectl apply -f ../manifests/istio/crd-certmanager-10.yaml
verify_kubectl $? "Creating istio custom resource definitions failed."
kubectl apply -f ../manifests/istio/crd-certmanager-11.yaml
verify_kubectl $? "Creating istio custom resource definitions failed."
wait_for_crds "virtualservices,destinationrules,serviceentries,gateways,envoyfilters,policies,meshpolicies,httpapispecbindings,httpapispecs,quotaspecbindings,quotaspecs,rules,attributemanifests"

kubectl apply -f ../manifests/istio/istio-lean.yaml
wait_for_deployment_in_namespace "istio-ingressgateway" "istio-system"
wait_for_deployment_in_namespace "istio-pilot" "istio-system"
verify_kubectl $? "Creating all istio components failed."
wait_for_all_pods_in_namespace "istio-system"