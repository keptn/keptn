# Keptn Installer

This repository contains **scripts** and **manifest files** that are needed to install Keptn on a Kubernets cluster. The scripts help to install the manifests in correct order and to manipulate them for Cloud platform specific requirements.  

The scripts and manifests are finally put into a container, which starts the initial script `installKeptn.sh`. Depending on the platform parameter that is handed over to the `installKeptn.sh`, the script then runs the installation process for one of the following Kubernetes platforms:
* AKS
* EKS
* OpenShift
* GKE
* PKS (not released yet)
* Kubernetes