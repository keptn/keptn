#!/bin/bash
source ./common/utils.sh

kubectl create namespace istio-system

helm template istio-init ../manifests/istio/helm/istio-init --namespace istio-system | kubectl apply -f -
verify_kubectl $? "Creating Istio resources failed"
wait_for_crds "adapters.config.istio.io,attributemanifests.config.istio.io,authorizationpolicies.rbac.istio.io,clusterrbacconfigs.rbac.istio.io,destinationrules.networking.istio.io,envoyfilters.networking.istio.io,gateways.networking.istio.io,handlers.config.istio.io,httpapispecbindings.config.istio.io,httpapispecs.config.istio.io,instances.config.istio.io,meshpolicies.authentication.istio.io,policies.authentication.istio.io,quotaspecbindings.config.istio.io,quotaspecs.config.istio.io,rbacconfigs.rbac.istio.io,rules.config.istio.io,serviceentries.networking.istio.io,servicerolebindings.rbac.istio.io,serviceroles.rbac.istio.io,sidecars.networking.istio.io,templates.config.istio.io,virtualservices.networking.istio.io"

# We tested it with helm --set according to the descriptions provided in https://istio.io/docs/setup/install/helm/
# However, it did not work out. Therefore, we are using sed
sed 's/LoadBalancer #change to NodePort, ClusterIP or LoadBalancer if need be/'$GATEWAY_TYPE'/g' ../manifests/istio/helm/istio/charts/gateways/values.yaml  > ../manifests/istio/helm/istio/charts/gateways/values_tmp.yaml
mv ../manifests/istio/helm/istio/charts/gateways/values_tmp.yaml ../manifests/istio/helm/istio/charts/gateways/values.yaml
helm template istio ../manifests/istio/helm/istio --namespace istio-system --values ../manifests/istio/helm/istio/values-istio-minimal.yaml | kubectl apply -f -
verify_kubectl $? "Installing Istio failed."
wait_for_deployment_in_namespace "istio-ingressgateway" "istio-system"
wait_for_deployment_in_namespace "istio-pilot" "istio-system"
wait_for_deployment_in_namespace "istio-citadel" "istio-system"
wait_for_deployment_in_namespace "istio-sidecar-injector" "istio-system"
wait_for_all_pods_in_namespace "istio-system"

oc adm policy add-scc-to-user anyuid -z istio-ingress-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z default -n istio-system
oc adm policy add-scc-to-user anyuid -z prometheus -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-egressgateway-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-citadel-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-ingressgateway-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-cleanup-old-ca-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-mixer-post-install-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-mixer-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-pilot-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-sidecar-injector-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-galley-service-account -n istio-system
oc adm policy add-scc-to-user anyuid -z istio-security-post-install-account -n istio-system

oc expose svc istio-ingressgateway -n istio-system

ROUTER_POD=$(oc get pods -n default -l router=router -ojsonpath={.items[0].metadata.name})
# allow wildcard domains
oc project default
oc adm router --replicas=0
verify_kubectl $? "Scaling down router failed"
oc set env dc/router ROUTER_ALLOW_WILDCARD_ROUTES=true
verify_kubectl $? "Configuration of openshift router failed"
oc scale dc/router --replicas=1
verify_kubectl $? "Upscaling of router failed"

oc delete pod $ROUTER_POD -n default --force --grace-period=0 --ignore-not-found

# create wildcard route for istio ingress gateway

BASE_URL=$(oc get route -n istio-system istio-ingressgateway -oyaml | yq r - spec.host | sed 's~istio-ingressgateway-istio-system.~~')
# Domain used for routing to keptn services
export DOMAIN="ingress-gateway.$BASE_URL"

oc create route passthrough istio-wildcard-ingress-secure-keptn --service=istio-ingressgateway --hostname="www.keptn.ingress-gateway.$BASE_URL" --port=https --wildcard-policy=Subdomain --insecure-policy='None' -n istio-system


oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:keptn:default
verify_kubectl $? "Adding cluster-role failed."

# Set up SSL
openssl req -nodes -newkey rsa:2048 -keyout key.pem -out certificate.pem  -x509 -days 365 -subj "/CN=$DOMAIN"

kubectl create --namespace istio-system secret tls istio-ingressgateway-certs --key key.pem --cert certificate.pem
#verify_kubectl $? "Creating secret for istio-ingressgateway-certs failed."

rm key.pem
rm certificate.pem

kubectl apply -f ../manifests/istio/public-gateway.yaml
verify_kubectl $? "Deploying public-gateway failed."
