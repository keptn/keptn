#!/bin/bash

DT_TENANT_ID=$1
DT_PAAS_TOKEN=$2

kubectl apply --filename https://github.com/knative/serving/releases/download/v0.4.0/istio-crds.yaml && \
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.4.0/istio.yaml

echo "Wait 4 minutes for changes to apply... "
sleep 240
echo "Wait 4 additional minutes for changes to apply... "
sleep 240

kubectl label namespace production istio-injection=enabled

./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN

echo "Wait 10s for changes to apply... "
sleep 10

kubectl delete pods --all -n production
kubectl delete pods --all -n staging
kubectl delete pods --all -n dev

kubectl delete meshpolicies.authentication.istio.io default # fix for the MySQL connection error caused by Istio
