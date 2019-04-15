#!/bin/bash

kubectl apply -f ../manifests/istio/istio-crds-knative.yaml
sleep 60
kubectl apply -f ../manifests/istio-knative.yaml

echo "Wait 4 minutes for changes to apply... "
sleep 240

echo "Wait 10s for changes to apply... "
sleep 10

kubectl delete pods --all -n keptn

# kubectl delete meshpolicies.authentication.istio.io default # fix for the MySQL connection error caused by Istio
