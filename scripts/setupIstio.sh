#!/bin/bash

DT_TENANT_ID=$1
DT_PAAS_TOKEN=$2

kubectl apply -f ../manifests/istio/crd-10.yaml
kubectl apply -f ../manifests/istio/crd-11.yaml
kubectl apply -f ../manifests/istio/crd-certmanager-10.yaml
kubectl apply -f ../manifests/istio/crd-certmanager-11.yaml

sleep 30

kubectl apply -f ../manifests/istio/istio-demo.yaml

echo "wait a couple of minutes for changes to apply... "
sleep 250
echo "wait even longer..."
sleep 250
echo "continue..."

kubectl label namespace production istio-injection=enabled

# Istio configuration for production namespace
kubectl create -f ../repositories/k8s-deploy-production/istio/gateway.yml
kubectl create -f ../repositories/k8s-deploy-production/istio/destination_rule.yml
kubectl create -f ../repositories/k8s-deploy-production/istio/virtual_service.yml

# ./createServiceEntry.sh $DT_TENANT_ID $DT_PAAS_TOKEN

sleep 10

kubectl delete pods --all -n production
kubectl delete pods --all -n staging
kubectl delete pods --all -n dev

kubectl delete meshpolicies.authentication.istio.io default # fix for the MySQL connection error caused by Istio
