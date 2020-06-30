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

# Creating cluster role binding
kubectl apply -f ../manifests/keptn/rbac.yaml
verify_kubectl $? "Creating cluster role for keptn failed."

# Create keptn secret
KEPTN_API_TOKEN=$(head -c 16 /dev/urandom | base64)
verify_variable "$KEPTN_API_TOKEN" "KEPTN_API_TOKEN could not be derived." 
kubectl create secret generic -n keptn keptn-api-token --from-literal=keptn-api-token="$KEPTN_API_TOKEN"

# Install Keptn Datastore
print_info "Installing Keptn Datastore"
kubectl apply -f ../manifests/logging/rbac.yaml
verify_kubectl $? "Creating rbac failed."
kubectl apply -f ../manifests/logging/mongodb/pvc.yaml
verify_kubectl $? "Creating mongodb PVC failed."
kubectl apply -f ../manifests/logging/mongodb/secret.yaml
verify_kubectl $? "Creating mongodb secret failed."
kubectl apply -f ../manifests/logging/mongodb/deployment.yaml
verify_kubectl $? "Creating mongodb deployment failed."
kubectl apply -f ../manifests/logging/mongodb/svc.yaml
verify_kubectl $? "Creating mongodb service failed."
kubectl apply -f ../manifests/logging/mongodb-datastore/mongodb-datastore.yaml
verify_kubectl $? "Creating mongodb-datastore service failed."
wait_for_deployment_in_namespace "mongodb-datastore" "keptn-datastore"

kubectl apply -f ../manifests/logging/mongodb-datastore/mongodb-datastore-distributor.yaml
verify_kubectl $? "Creating mongodb-datastore service failed."

# Install Keptn core and use case dependent components 
print_debug "Deploying Keptn core"
kubectl apply -f ../manifests/keptn/core.yaml 
verify_kubectl $? "Deploying Keptn core components failed."

##############################################
## Start validation of Keptn core           ##
##############################################
wait_for_all_pods_in_namespace "keptn"
wait_for_deployment_in_namespace "api-service" "keptn"
wait_for_deployment_in_namespace "bridge" "keptn"
wait_for_deployment_in_namespace "eventbroker-go" "keptn"
wait_for_deployment_in_namespace "helm-service" "keptn"
wait_for_deployment_in_namespace "shipyard-service" "keptn"
wait_for_deployment_in_namespace "configuration-service" "keptn"
wait_for_deployment_in_namespace "helm-service-service-create-distributor" "keptn"
wait_for_deployment_in_namespace "shipyard-service-create-project-distributor" "keptn"
wait_for_deployment_in_namespace "shipyard-service-delete-project-distributor" "keptn"

# Install API Gateway NGINX
print_debug "Deploying API Gateway NGINX"
kubectl apply -f ../manifests/keptn/api-gateway-nginx.yaml
verify_kubectl $? "Deploying API Gateway NGINX failed."
wait_for_all_pods_in_namespace "keptn"
wait_for_deployment_in_namespace "api-gateway-nginx" "keptn"

case $USE_CASE in
  "")
    print_debug "Deploying Keptn quality gates"
    kubectl apply -f ../manifests/keptn/quality-gates.yaml 
    verify_kubectl $? "Deploying Keptn quality gates components failed."

    print_debug "Deploying Keptn continuous operations"
    kubectl apply -f ../manifests/keptn/continuous-operations.yaml
    verify_kubectl $? "Deploying Keptn continuous operations components failed."

    ################################################
    ## Start validation of Keptn all capabilities ##
    ################################################
    wait_for_all_pods_in_namespace "keptn"
    wait_for_deployment_in_namespace "lighthouse-service" "keptn"
    wait_for_deployment_in_namespace "lighthouse-service-distributor" "keptn"
    wait_for_deployment_in_namespace "remediation-service-distributor" "keptn"
    wait_for_deployment_in_namespace "wait-service-deployment-distributor" "keptn"
    ;;
  continuous-delivery)
    print_debug "Deploying Keptn continuous deployment"
    kubectl apply -f ../manifests/keptn/continuous-deployment.yaml 
    verify_kubectl $? "Deploying Keptn continuous deployment components failed."

    print_debug "Deploying Keptn quality gates"
    kubectl apply -f ../manifests/keptn/quality-gates.yaml 
    verify_kubectl $? "Deploying Keptn quality gates components failed."

    print_debug "Deploying Keptn continuous operations"
    kubectl apply -f ../manifests/keptn/continuous-operations.yaml 
    verify_kubectl $? "Deploying Keptn continuous operations components failed."

    ################################################
    ## Start validation of Keptn all capabilities ##
    ################################################
    wait_for_all_pods_in_namespace "keptn"
    wait_for_deployment_in_namespace "gatekeeper-service" "keptn"
    wait_for_deployment_in_namespace "jmeter-service" "keptn"
    wait_for_deployment_in_namespace "lighthouse-service" "keptn"
    wait_for_deployment_in_namespace "remediation-service" "keptn"
    wait_for_deployment_in_namespace "wait-service" "keptn"
    wait_for_deployment_in_namespace "lighthouse-service-distributor" "keptn"
    wait_for_deployment_in_namespace "gatekeeper-service-evaluation-done-distributor" "keptn"
    wait_for_deployment_in_namespace "helm-service-distributor" "keptn"
    wait_for_deployment_in_namespace "jmeter-service-deployment-distributor" "keptn"
    wait_for_deployment_in_namespace "remediation-service-distributor" "keptn"
    wait_for_deployment_in_namespace "wait-service-deployment-distributor" "keptn"
    ;;
  *)
    echo "Use case not provided"
    ;;
esac
