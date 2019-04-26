#!/bin/bash

LOG_LOCATION=./logs
exec > >(tee -i $LOG_LOCATION/installKeptn.log)
exec 2>&1

echo "[keptn|INFO] [Fri Sep 09 10:42:29.902022 2011] Starting installation of keptn"

# Environment variables for install Istio and Knative
if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
  echo "[keptn|DEBUG] CLUSTER_IPV4_CIDR not set, retrieve it using gcloud."
  CLUSTER_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - clusterIpv4Cidr)

  if [[ -z "${CLUSTER_IPV4_CIDR}" ]]; then
    echo "[keptn|ERROR] CLUSTER_IPV4_CIDR is undefined, stop installation." && exit 1
  fi
fi

if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
  echo "[keptn|DEBUG] SERVICES_IPV4_CIDR not set, retrieve it using gcloud"
  SERVICES_IPV4_CIDR=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} | yq r - servicesIpv4Cidr)
  
  if [[ -z "${SERVICES_IPV4_CIDR}" ]]; then
    echo "[keptn|ERROR] SERVICES_IPV4_CIDR is undefined, stop installation." && exit 1
  fi
fi

if [[ -z "${GCLOUD_USER}" ]]; then
  echo "[keptn|DEBUG] GCLOUD_USER not set, retrieve it using gcloud"
  GCLOUD_USER=$(gcloud config get-value account)

  if [[ -z "${GCLOUD_USER}" ]]; then
    echo "[keptn|ERROR] GCLOUD_USER is undefined, stop installation." && exit 1
  fi
fi

# Test connection to cluster
echo "[keptn|INFO] Test connection to cluster"
./testConnection.sh

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Connection failed." && exit 1
fi

# Grant cluster admin rights to gcloud user
kubectl create clusterrolebinding keptn-cluster-admin-binding --clusterrole=cluster-admin --user=$GCLOUD_USER
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Cluster role binding could not be created." && exit 1
fi

# Create K8s namespaces
kubectl apply -f ../manifests/k8s-namespaces.yml
if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Namespace could not be created." && exit 1
fi

# Create container registry
echo "[keptn|INFO] Creating container registry"
./setupContainerRegistry.sh

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Creating container registry failed." && exit 1
fi

echo "[keptn|INFO] Creating container registry done"

# Install Istio service mesh
echo "[keptn|INFO] Installing Istio"
./setupIstio.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Installing Istio failed." && exit 1
fi

echo "[keptn|INFO] Installing Istio done"

# Install knative core components
echo "[keptn|INFO] Installing Knative"
./setupKnative.sh $CLUSTER_IPV4_CIDR $SERVICES_IPV4_CIDR

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Installing Knative failed." && exit 1
fi

echo "[keptn|INFO] Installing Knative done"

# Install keptn core services - Install keptn channels
echo "[keptn|INFO] Install keptn"
./setupKeptn.sh

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Installing keptn failed." && exit 1
fi

echo "[keptn|INFO] Install keptn done"

# Install keptn services
echo "[keptn|INFO] Wear uniform"
./wearUniform.sh

if [[ $? != '0' ]]; then
  echo "[keptn|ERROR] Installing keptn's uniform failed." && exit 1
fi

echo "[keptn|INFO] Keptn wears uniform"

# Install done
echo "[keptn|INFO] Installation of keptn complete."

echo "[keptn|INFO] To retrieve the Keptn API Token, please execute the following command:"
echo "[keptn|INFO] kubectl get secret keptn-api-token -n keptn -o=yaml | yq - r data.keptn-api-token | base64 --decode"
