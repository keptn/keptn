#!/bin/bash
source ./common/utils.sh

case $PLATFORM in
  aks)
    if [ "$INGRESS" = "istio" ]; then        
        # Install Istio service mesh
        print_info "Installing Istio on AKS (this might take a while)"
        source ./common/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"    
        
    elif [ "$INGRESS" = "nginx" ]; then
        # Install nginx service mesh
        print_info "Installing nginx on AKS (this might take a while)"
        kubectl apply -f ../manifests/nginx/nginx.yaml
        verify_install_step $? "Installing nginx deployment failed."
        wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
        verify_install_step $? "Installing nginx failed because deployment not available"
        kubectl apply -f ../manifests/nginx/nginx-svc.yaml
        verify_install_step $? "Installing nginx service failed."

        source ./installIngressForApi.sh
    fi    
    ;;
  eks)
    if [ "$INGRESS" = "istio" ]; then
        # Install Istio service mesh
        print_info "Installing Istio on EKS (this might take a while)"
        source ./common/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"    
        
    elif [ "$INGRESS" = "nginx" ]; then
        # Install nginx service mesh
        print_info "Installing nginx on EKS (this might take a while)"
        kubectl apply -f ../manifests/nginx/nginx.yaml
        verify_install_step $? "Installing nginx deployment failed."
        wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
        verify_install_step $? "Installing nginx failed because deployment not available"
        kubectl apply -f ../manifests/nginx/nginx-svc.yaml
        verify_install_step $? "Installing nginx service failed."
        kubectl apply -f ../manifests/nginx/nginx-svc-eks.yaml
        kubectl apply -f ../manifests/nginx/nginx-configmap-eks.yaml
        verify_install_step $? "Installing nginx configmap for EKS failed."

        source ./installIngressForApi.sh
    fi
    ;;
  openshift)
    echo "Install Keptn on OpenShift"
    if [ "$INGRESS" = "istio" ]; then
        # Install Istio service mesh
        print_info "Installing Istio on OpenShift (this might take a while)"
        source ./openshift/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"
    fi 
    ;;
  gke)
    if [ "$INGRESS" = "istio" ]; then        
        # Install Istio service mesh
        print_info "Installing Istio on GKE (this might take a while)"
        source ./common/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"    
        
    elif [ "$INGRESS" = "nginx" ]; then
        # Install nginx service mesh
        print_info "Installing nginx on GKE (this might take a while)"
        kubectl apply -f ../manifests/nginx/nginx.yaml
        verify_install_step $? "Installing nginx deployment failed."
        wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
        verify_install_step $? "Installing nginx failed because deployment not available"
        kubectl apply -f ../manifests/nginx/nginx-svc.yaml
        verify_install_step $? "Installing nginx service failed."     

        source ./installIngressForApi.sh
    fi 
    ;;
  pks)
    if [ "$INGRESS" = "istio" ]; then        
        # Install Istio service mesh
        print_info "Installing Istio on PKS (this might take a while)"
        source ./common/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"    
        
    elif [ "$INGRESS" = "nginx" ]; then
        # Install nginx service mesh
        print_info "Installing nginx on PKS (this might take a while)"
        kubectl apply -f ../manifests/nginx/nginx.yaml
        verify_install_step $? "Installing nginx deployment failed."
        wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
        verify_install_step $? "Installing nginx failed because deployment not available" 
        kubectl apply -f ../manifests/nginx/nginx-svc.yaml
        verify_install_step $? "Installing nginx service failed."   

        source ./installIngressForApi.sh     
    fi 
    ;;
  kubernetes)
    if [ "$INGRESS" = "istio" ]; then        
        # Install Istio service mesh
        print_info "Installing Istio on Kubernetes (this might take a while)"
        source ./common/setupIstio.sh
        verify_install_step $? "Installing Istio failed."
        print_info "Installing Istio done"    
        
    elif [ "$INGRESS" = "nginx" ]; then
        # Install nginx service mesh
        print_info "Installing nginx on Kubernetes (this might take a while)"
        kubectl apply -f ../manifests/nginx/nginx.yaml
        verify_install_step $? "Installing nginx deployment failed."
        wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
        verify_install_step $? "Installing nginx failed because deployment not available"

        if [ "$GATEWAY_TYPE" = "NodePort" ]; then
          kubectl apply -f ../manifests/nginx/nginx-svc-nodeport.yaml
          verify_install_step $? "Installing nginx service failed."  
        else 
          kubectl apply -f ../manifests/nginx/nginx-svc.yaml
          verify_install_step $? "Installing nginx service failed."  
        fi

        source ./installIngressForApi.sh
    fi 
    ;;
  *)
    echo "Platform not provided"
    echo "Installation aborted, please provide platform when executing keptn install --platform="
    exit 1
    ;;
esac
