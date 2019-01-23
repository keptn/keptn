#!/bin/bash
DT_TENANT_ID=$1
DT_PAAS_TOKEN=$2

kubectl apply -f ../manifests/istio/istio-crds.yml
kubectl apply -f ../manifests/istio/istio-demo.yml

sleep 500

kubectl label namespace production istio-injection=enabled
kubectl create -f ../manifests/istio/istio-gateway.yml

./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN

sleep 10

kubectl delete pods --all -n production

kubectl delete meshpolicies.authentication.istio.io default # fix for the MySQL connection error caused by Istio
