#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

case $PLATFORM in
  aks)
    echo "Install on AKS"
    ./common/install.sh
    ;;
  eks)
    echo "Install on EKS"
    ./common/install.sh
    ;;
  openshift)
    echo "Install on OpenShift"
    ./openshift/installOnOpenshift.sh
    ;;
  gke)
    echo "Install on GKE"
    ./common/install.sh
    ;;
  pks)
    echo "Install on PKS"
    ./common/install.sh
    ;;
  kubernetes) 
    echo "Install on Kubernetes"
    ./common/install.sh
    ;;
  *)
    echo "Platform not provided"
    echo "Installation aborted, please provide platform when executing keptn install --platform="
    exit 1
    ;;
esac
