# Keptn Installer

This repository contains **scripts** and **manifest files** that are needed to install Keptn on a Kubernetes cluster. The scripts help to install the manifests in correct order and to manipulate them for Cloud platform specific requirements.  

The scripts and manifests are finally put into a container, which starts the initial script `installKeptn.sh`. Depending on the platform parameter that is handed over to the `installKeptn.sh`, the script then runs the installation process for one of the following Kubernetes platforms:
* AKS
* EKS
* OpenShift
* GKE
* PKS
* Kubernetes
* Minikube 1.2

## Customise Charts with own values.yaml for local testing
Create `installer/manifests/keptn/values-local.yaml` file for your local values to be stored. The file should look like this:
For local templating use `helm template . -f values-local.yaml --name-template test-control-plane --output-dir ../../temp`

```
global:
  keptn:
    registry: "testregistry/keptn"      # keptn registry/image name
    tag: "0.0.1"                        # keptn version/tag

# only change if version at ./charts/control-plane/values.yaml --> apiGatewayNginx.registry/tag is not satisfying
#  apiGatewayNginx:
#    registry: this.is.a.test           # nginx registry/image name
#    tag: 10.0.0                        # ngnix version/tag
```
