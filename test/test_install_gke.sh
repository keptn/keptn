#!/bin/bash

source test/utils.sh

KEPTN_INSTALLER_REPO=${KEPTN_INSTALLER_REPO:-https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz}
PROJECT_NAME=${PROJECT_NAME:-sockshop}

# prepare creds.json file
cd ./test/assets

export CLN=$CLUSTER_NAME_NIGHTLY
export CLZ=$CLOUDSDK_COMPUTE_ZONE	
export PROJ=$PROJECT_NAME

echo "{}" > creds.json # empty credentials file

echo "Installing Keptn on GKE cluster"

# install Keptn using the develop version, which refers to the :latest docker images
keptn install --chart-repo="${KEPTN_INSTALLER_REPO}" --creds=creds.json --verbose --use-case=continuous-delivery --endpoint-service-type=LoadBalancer --hide-sensitive-data
verify_test_step $? "keptn install --chart-repo=${KEPTN_INSTALLER_REPO} - failed"

# authenticate at Keptn API
KEPTN_API_URL=http://$(kubectl -n keptn get service api-gateway-nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)

auth_at_keptn $KEPTN_API_URL $KEPTN_API_TOKEN

# install public-gateway.istio-system
kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: public-gateway
  namespace: istio-system
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
EOF

# set ingress-hostname params
INGRESS_IP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
kubectl create configmap -n keptn ingress-config --from-literal=ingress_hostname_suffix=${INGRESS_IP}.xip.io --from-literal=ingress_port=80 --from-literal=ingress_protocol=http --from-literal=ingress_gateway=public-gateway.istio-system -oyaml --dry-run | kubectl replace -f -

kubectl delete pod -n keptn -lapp.kubernetes.io/name=helm-service
sleep 15

echo "Verifying that services and namespaces have been created"

# verify the deployments within the keptn namespace
verify_deployment_in_namespace "api-gateway-nginx" "keptn"
verify_deployment_in_namespace "api-service" "keptn"
verify_deployment_in_namespace "bridge" "keptn"
verify_deployment_in_namespace "configuration-service" "keptn"
verify_deployment_in_namespace "gatekeeper-service" "keptn"
verify_deployment_in_namespace "jmeter-service" "keptn"
verify_deployment_in_namespace "lighthouse-service" "keptn"

# verify the datastore deployments
verify_deployment_in_namespace "mongodb" "keptn"
verify_deployment_in_namespace "mongodb-datastore" "keptn"

# verify the pods within istio-system
verify_deployment_in_namespace "istio-ingressgateway" "istio-system"
verify_deployment_in_namespace "istio-pilot" "istio-system"
verify_deployment_in_namespace "istio-citadel" "istio-system"
verify_deployment_in_namespace "istio-sidecar-injector" "istio-system"

cd ../..

echo "Installing Keptn on GKE cluster done ✓"

exit 0
