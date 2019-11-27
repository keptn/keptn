#!/bin/bash
source ./common/utils.sh

# Set up NATS
kubectl apply -f ../manifests/nats/nats-operator-prereqs.yaml
verify_kubectl $? "Creating NATS Operator failed."

kubectl apply -f ../manifests/nats/nats-operator-deploy.yaml
verify_kubectl $? "Creating NATS Operator failed."

wait_for_deployment_in_namespace "nats-operator" "keptn"

kubectl apply -f ../manifests/nats/nats-cluster.yaml
verify_kubectl $? "Creating NATS Cluster failed."

# Add config map in keptn namespace that contains the domain - this will be used by other services as well
rm -f ../manifests/gen/keptn-domain-configmap.yaml

cat ../manifests/keptn/keptn-domain-configmap.yaml | \
  sed 's~DOMAIN_PLACEHOLDER~'"$DOMAIN"'~' >> ../manifests/gen/keptn-domain-configmap.yaml

kubectl apply -f ../manifests/gen/keptn-domain-configmap.yaml
verify_kubectl $? "Creating configmap keptn-domain in keptn namespace failed."

# Creating cluster role binding
kubectl apply -f ../manifests/keptn/rbac.yaml
verify_kubectl $? "Creating cluster role for keptn failed."

# Create keptn secret
KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)
verify_variable "$KEPTN_API_TOKEN" "KEPTN_API_TOKEN could not be derived." 
kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"

# Install logging
print_info "Installing Logging"
kubectl apply -f ../manifests/logging/namespace.yaml
verify_kubectl $? "Creating logging namespace failed."
kubectl apply -f ../manifests/logging/mongodb-k8s/pvc.yaml
verify_kubectl $? "Creating mongodb PVC failed."
kubectl apply -f ../manifests/logging/mongodb-k8s/deployment.yaml
verify_kubectl $? "Creating mongodb deployment failed."
kubectl apply -f ../manifests/logging/mongodb-k8s/svc.yaml
verify_kubectl $? "Creating mongodb service failed."
kubectl apply -f ../manifests/logging/mongodb-datastore/k8s/mongodb-datastore.yaml
verify_kubectl $? "Creating mongodb-datastore service failed."
wait_for_deployment_in_namespace "mongodb-datastore" "keptn-datastore"

kubectl apply -f ../manifests/logging/mongodb-datastore/mongodb-datastore-distributor.yaml
verify_kubectl $? "Creating mongodb-datastore service failed."

kubectl apply -f ../manifests/keptn/core.yaml 
verify_kubectl $? "Deploying keptn core components failed."

kubectl apply -f ../manifests/keptn/core-distributors.yaml 
verify_kubectl $? "Deploying keptn core distributors failed."

##############################################
## Start validation of keptn installation   ##
##############################################
wait_for_all_pods_in_namespace "keptn"

wait_for_deployment_in_namespace "eventbroker-go" "keptn"
wait_for_deployment_in_namespace "api" "keptn"
wait_for_deployment_in_namespace "bridge" "keptn"
wait_for_deployment_in_namespace "gatekeeper-service" "keptn"
wait_for_deployment_in_namespace "helm-service" "keptn"
wait_for_deployment_in_namespace "jmeter-service" "keptn"
wait_for_deployment_in_namespace "shipyard-service" "keptn"
wait_for_deployment_in_namespace "lighthouse-service" "keptn"
wait_for_deployment_in_namespace "configuration-service" "keptn"

wait_for_deployment_in_namespace "helm-service-service-create-distributor" "keptn"
wait_for_deployment_in_namespace "helm-service-configuration-change-distributor" "keptn"
wait_for_deployment_in_namespace "jmeter-service-deployment-distributor" "keptn"
wait_for_deployment_in_namespace "lighthouse-service-tests-finished-distributor" "keptn"
wait_for_deployment_in_namespace "lighthouse-service-start-evaluation-distributor" "keptn"
wait_for_deployment_in_namespace "gatekeeper-service-evaluation-done-distributor" "keptn"
wait_for_deployment_in_namespace "shipyard-service-create-project-distributor" "keptn"

if [ "$INGRESS" = "istio" ]; then
  kubectl apply -f ../manifests/keptn/keptn-gateway.yaml
  verify_kubectl $? "Deploying keptn gateway failed."

  rm -f ../manifests/keptn/gen/keptn-api-virtualservice.yaml
  cat ../manifests/keptn/keptn-api-virtualservice.yaml | \
    sed 's~DOMAIN_PLACEHOLDER~'"$INGRESS_HOST"'~' > ../manifests/keptn/gen/keptn-api-virtualservice.yaml

  kubectl apply -f ../manifests/keptn/gen/keptn-api-virtualservice.yaml
  verify_kubectl $? "Deploying keptn api virtualservice failed."
fi
