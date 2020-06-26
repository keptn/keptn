#!/bin/bash

source ./common/utils.sh

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

case $PLATFORM in
  aks|eks|gke|pks|kubernetes)
    ./common/install.sh
    ;;
  openshift)
    echo "Install Keptn on OpenShift"
    ./openshift/installOnOpenshift.sh
    ;;
  *)
    echo "Platform not provided"
    echo "Installation aborted, please provide platform when executing keptn install --platform="
    exit 1
    ;;
esac
