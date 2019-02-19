#!/bin/bash

DT_TENANT_ID=$1
DT_PAAS_TOKEN=$2

kubectl apply -f ../manifests/istio/istio-crds.yml
kubectl apply -f ../manifests/istio/istio-demo.yml

echo "Wait 4 minutes for changes to apply... "
sleep 240
echo "Wait additional 4 minutes for changes to apply... "
sleep 240

kubectl label namespace production istio-injection=enabled

./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN

sleep 10

kubectl delete pods --all -n production
kubectl delete pods --all -n staging
kubectl delete pods --all -n dev

kubectl delete meshpolicies.authentication.istio.io default # fix for the MySQL connection error caused by Istio
