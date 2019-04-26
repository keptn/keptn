#!/bin/bash
CLUSTER_IPV4_CIDR=$1
SERVICES_IPV4_CIDR=$2

source ./utils.sh

# Create namespace and container registry
# Needed for Pull Request Travis Build - will be removed
kubectl create namespace keptn #2> /dev/null

./setupContainerRegistry.sh
verify_install_step $? "Creating container registry failed, stop installation."

# Install knative serving, building, eventing
kubectl apply -f https://github.com/knative/serving/releases/download/v0.5.0/serving.yaml
## (!) KOWN ISSUE: Custom resource definition Image not created in time ##
print_info "KNOWN ISSUE: Custom resource definition Image not created in time, skip verification of kubectl cmd."
# verify_kubectl $? "Applying knative serving components failed, stop installation."
wait_for_crds "certificates,clusteringresses,configurations,images,podautoscalers,revisions,routes,services,serverlessservices"
wait_for_all_pods_in_namespace "knative-serving"

kubectl apply -f https://github.com/knative/build/releases/download/v0.5.0/build.yaml
verify_kubectl $? "Applying knative building components failed, stop installation."
wait_for_crds "builds,buildtemplates,clusterbuildtemplates,images"
wait_for_all_pods_in_namespace "knative-build"

kubectl apply -f https://github.com/knative/eventing/releases/download/v0.5.0/release.yaml
## (!) KOWN ISSUE: Custom resource defintion ClusterChannelProvisioner not created in time ##
print_info "KNOWN ISSUE: Custom resource definition ClusterChannelProvisioner not created in time, skip verification of kubectl cmd."
#verify_kubectl $? "Applying knative eventing components failed, stop installation."
wait_for_crds "brokers,channels,clusterchannelprovisioners,subscriptions,triggers"
wait_for_all_pods_in_namespace "knative-eventing"

kubectl apply -f https://github.com/knative/eventing-sources/releases/download/v0.5.0/eventing-sources.yaml
verify_kubectl $? "Applying knative sources failed, stop installation."
wait_for_crds "awssqssources,containersources,cronjobsources,githubsources,kuberneteseventsources"
wait_for_all_pods_in_namespace "knative-sources"

kubectl apply -f https://github.com/knative/serving/releases/download/v0.5.0/monitoring.yaml
verify_kubectl $? "Applying knative monitoring components failed, stop installation."
## (!) KOWN ISSUE: Two pods in knative-monitoring don't run ##
print_info "KNOWN ISSUE: Two pods in knative-monitoring don't run, skip verification of pod status."
#wait_for_all_pods_in_namespace "knative-monitoring"

kubectl apply -f https://raw.githubusercontent.com/knative/serving/v0.5.0/third_party/config/build/clusterrole.yaml
verify_kubectl $? "Creating cluster role for knative failed, stop installation."

# Configure knative serving default domain
rm -f ../manifests/gen/config-domain.yaml

ISTIO_INGRESS_IP=$(kubectl describe svc istio-ingressgateway -n istio-system | grep "LoadBalancer Ingress:" | sed 's~LoadBalancer Ingress:[ \t]*~~')
verify_variable $ISTIO_INGRESS_IP "ISTIO_INGRESS_IP is empty, stop installation." 

cat ../manifests/knative/config-domain.yaml | \
  sed 's~ISTIO_INGRESS_IP_PLACEHOLDER~'"$ISTIO_INGRESS_IP"'~' >> ../manifests/gen/config-domain.yaml

kubectl apply -f ../manifests/gen/config-domain.yaml
verify_kubectl $? "Creating config map failed, stop installation."

kubectl get configmap config-network -n knative-serving -o=yaml | yq w - data['istio.sidecar.includeOutboundIPRanges'] "$CLUSTER_IPV4_CIDR,$SERVICES_IPV4_CIDR" | kubectl apply -f - 
verify_kubectl $? "Updating configmap config-network failed, stop installation."

# Create build-bot service account
kubectl apply -f ../manifests/knative/build/service-account.yaml
verify_kubectl $? "Creating service account for build bot in keptn namespace failed, stop installation."

# Install kaniko build template
kubectl apply -f ../manifests/knative/build/kaniko.yaml -n keptn
verify_kubectl $? "Creating the kaniko build template failed, stop installation."

# ##############################################
# ## Start validation of Knative installation ##
# ##############################################
