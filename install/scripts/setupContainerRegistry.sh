#!/bin/bash

source ./utils.sh

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
verify_kubectl $? "Creating config map for container registry failed."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
verify_kubectl $? "Creating persistent volume claim for container registry failed."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
verify_kubectl $? "Creating deployment for container registry failed."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml
verify_kubectl $? "Creating service for container registry failed."
