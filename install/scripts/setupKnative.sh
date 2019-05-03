#!/bin/bash
CLUSTER_IPV4_CIDR=$1
SERVICES_IPV4_CIDR=$2

source ./utils.sh

# Create namespace and container registry
# Needed for pull request Travis Build - will be removed
kubectl create namespace keptn #2> /dev/null

./setupContainerRegistry.sh
verify_install_step $? "Creating container registry failed."

# Install knative serving, building, eventing
for i in 1 2
do
  kubectl apply -f https://github.com/knative/serving/releases/download/v0.4.0/serving.yaml
  verify_kubectl $? "Applying knative serving components failed."
  sleep 5
  wait_for_crds "clusteringresses,configurations,images,podautoscalers,revisions,routes,services"
  wait_for_all_pods_in_namespace "knative-serving"

  kubectl apply -f https://github.com/knative/build/releases/download/v0.4.0/build.yaml
  verify_kubectl $? "Applying knative building components failed."
  sleep 5
  wait_for_crds "builds,buildtemplates,clusterbuildtemplates,images"
  wait_for_all_pods_in_namespace "knative-build"

  kubectl apply -f https://github.com/knative/eventing/releases/download/v0.4.0/release.yaml
  ## KNOWN ISSUE: Race condition regarding custom resource defintion ClusterChannelProvisioner when applying release.yaml - https://github.com/knative/eventing/issues/680
  print_info "KNOWN ISSUE: Applying the knative v.0.4.0 release.yaml runs into a race condition for the custom resource definition ClusterChannelProvisioner." 
  verify_kubectl $? "Applying knative eventing components failed."
  sleep 5
  wait_for_crds "channels,clusterchannelprovisioners,subscriptions"
  wait_for_all_pods_in_namespace "knative-eventing"

  kubectl apply -f https://github.com/knative/eventing-sources/releases/download/v0.4.0/release.yaml
  verify_kubectl $? "Applying knative sources failed."
  sleep 5
  wait_for_crds "awssqssources,containersources,cronjobsources,githubsources,kuberneteseventsources"
  wait_for_all_pods_in_namespace "knative-sources"

  kubectl apply -f https://github.com/knative/serving/releases/download/v0.4.0/monitoring.yaml
  verify_kubectl $? "Applying knative monitoring components failed."
  sleep 5
  wait_for_all_pods_in_namespace "knative-monitoring"

  sleep 30
done

# Creating cluster role binding for knative
kubectl apply -f https://raw.githubusercontent.com/knative/serving/v0.4.0/third_party/config/build/clusterrole.yaml
verify_kubectl $? "Creating cluster role for knative failed."

# Configure knative serving default domain
rm -f ../manifests/gen/config-domain.yaml

ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
verify_variable "$ISTIO_INGRESS_IP" "ISTIO_INGRESS_IP is empty and could not be derived from the Istio ingress gateway." 

cat ../manifests/knative/config-domain.yaml | \
  sed 's~ISTIO_INGRESS_IP_PLACEHOLDER~'"$ISTIO_INGRESS_IP"'~' >> ../manifests/gen/config-domain.yaml

kubectl apply -f ../manifests/gen/config-domain.yaml
verify_kubectl $? "Creating configmap config-domain in knative-serving namespace failed."

kubectl get configmap config-network -n knative-serving -o=yaml | yq w - data['istio.sidecar.includeOutboundIPRanges'] "$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR" | kubectl apply -f - 
verify_kubectl $? "Updating configmap config-network in knative-serving namespace failed."

# Create build-bot service account
kubectl apply -f ../manifests/knative/build/service-account.yaml
verify_kubectl $? "Creating service account for build bot in keptn namespace failed."

# Install kaniko build template
kubectl apply -f ../manifests/knative/build/kaniko.yaml -n keptn
verify_kubectl $? "Creating the kaniko build template failed."

##############################################
## Start validation of Knative installation ##
##############################################
