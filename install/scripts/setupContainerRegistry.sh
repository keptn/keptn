#!/bin/bash

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-configmap.yml
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Configmap for container registry could not be created." && exit 1
fi

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-pvc.yml
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Persistent volume claim for container registry could not be created." && exit 1
fi

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-deployment.yml
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Deployment for container registry could not be created." && exit 1
fi

kubectl apply -f ../manifests/container-registry/k8s-docker-registry-service.yml
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Service for container registry could not be created." && exit 1
fi
