#!/bin/bash
source ./common/utils.sh

# determine whether istio is already installed
kubectl get ns istio-system
ISTIO_AVAILABLE=$?

if [[ "$ISTIO_AVAILABLE" == 0 ]] && [[ "$INGRESS_INSTALL_OPTION" == "Reuse" ]]; then
    # An istio-version is already installed
    print_info "Istio installation is reused but its full compatibility is not checked"
    print_info "Checking if istio-ingressgateway is available in namespace istio-system"
    wait_for_deployment_in_namespace "istio-ingressgateway" "istio-system"
    wait_for_all_pods_in_namespace "istio-system"

elif [[ "$ISTIO_AVAILABLE" == 0 ]] && ([[ "$INGRESS_INSTALL_OPTION" == "StopIfInstalled" ]] || [[ "$INGRESS_INSTALL_OPTION" == "" ]] || [[ "$INGRESS_INSTALL_OPTION" == "INGRESS_INSTALL_PLACEHOLDER" ]]); then
    print_error "Istio is already installed but is not used due to unknown compatibility"
    exit 1
else
    if [[ "$ISTIO_AVAILABLE" == 0 ]] && [[ "$INGRESS_INSTALL_OPTION" == "Overwrite" ]]; then
        print_info "Istio installation is overwritten"
    fi

    # Istio installation
    print_info "Install Istio"
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
fi

kubectl apply -f ../manifests/istio/public-gateway.yaml
verify_kubectl $? "Deploying public-gateway failed."
