#!/bin/bash

source ./utils.sh

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
verify_kubectl $? "Configmap for container registry could not be created, stop installation."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
verify_kubectl $? "Persistent volume claim for container registry could not be created, stop installation."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
verify_kubectl $? "Deployment for container registry could not be created, stop installation."

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml
verify_kubectl $? "Service for container registry could not be created, stop installation."
