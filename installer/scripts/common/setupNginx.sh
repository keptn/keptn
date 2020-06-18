#!/bin/bash

if [[ "$INGRESS_INSTALL_OPTION" == "Reuse" ]]; then
    print_info "NGINX ingress controller is reused but its compatibility is not checked"
    print_info "Checking if nginx-ingress-controller is available in namespace ingress-nginx"
    wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
    wait_for_all_pods_in_namespace "ingress-nginx"

else
    if [[ "$PLATFORM" == "kubernetes" ]]; then
      # Install nginx service mesh
      print_info "Installing nginx on Kubernetes (this might take a while)"
      kubectl apply -f ../manifests/nginx/nginx.yaml
      verify_install_step $? "Installing nginx deployment failed."
      wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
      verify_install_step $? "Installing nginx failed because deployment not available"

      if [[ "$GATEWAY_TYPE" == "NodePort" ]]; then
        print_info "Install NGINX with a NodePort"
        kubectl apply -f ../manifests/nginx/nginx-svc-nodeport.yaml
        verify_install_step $? "Installing nginx service failed."
      else
        print_info "Install NGINX with a LoadBalancer"
        kubectl apply -f ../manifests/nginx/nginx-svc.yaml
        verify_install_step $? "Installing nginx service failed."
      fi
    else
      print_info "Installing nginx (this might take a while)"
      kubectl apply -f ../manifests/nginx/nginx.yaml
      verify_install_step $? "Installing nginx deployment failed."
      wait_for_deployment_in_namespace "nginx-ingress-controller" "ingress-nginx"
      verify_install_step $? "Installing nginx failed because deployment not available"
      kubectl apply -f ../manifests/nginx/nginx-svc.yaml
      verify_install_step $? "Installing nginx service failed."

      if [[ "$PLATFORM" == "eks" ]]; then
        print_info "Install NGINX on EKS"
        kubectl apply -f ../manifests/nginx/nginx-svc-eks.yaml
        kubectl apply -f ../manifests/nginx/nginx-configmap-eks.yaml
        verify_install_step $? "Installing nginx configmap for EKS failed."
      fi
    fi
fi
