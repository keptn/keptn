# Running Integration Tests

Several tests are specified in [.travis.yml](../.travis.yml). They usually follow this layout:

```yaml
    stage: Some Integration Test (--platform=kubernetes)
    os: linux
    before_script:
      # download and install kubectl and setup cluster
    script:
      - kubectl get nodes || travis_terminate 1
      # finally install keptn quality gates
      - test/test_install_kubernetes_quality_gates.sh
      - keptn status
      - export PROJECT=musicshop
      - test/test_quality_gates_standalone.sh
    after_success:
      # clean up cluster, etc...
      - echo "Tests were successful, cleaning up the cluster now..."
    after_failure:
      # print some debug info
      - echo "Keptn Installation Log:"
      - cat ~/.keptn/keptn-installer.log
```

* `before_script`: Set up `kubectl`, `keptn` CLI and create a Kubernetes cluster
* `script`: Verify connection to the kubernetes cluster, install Keptn (e.g., `test/test_install_kubernetes_quality_gates.sh`),
  verify connection to Keptn (`keptn status`) and run the actual test (e.g., `test/test_quality_gates_standalone.sh`)
* `after_success`: Perform some cleanup (e.g., delete the cluster if necessary)
* `after_failure`: Print some debug output

## Run Integration Tests in Travis-CI

There are some caveats when creating a Kubernetes cluster on Travis-CI:

* It is possible to create and connect to an external cluster (e.g., GKE, AWS, ...). However, this requires secrets
  for the cloud provider to be set on Travis-CI, and accessing those secrets only works for commits that are pushed
  within the current repository. It does **not work** for Pull Requests from forks.
* It is possible to run docker on Travis-CI.
* It is not possible to create a virtual-machine on Travis-CI. Therefore, setting up K3s, Minishift or Minikube on Travis-CI
  only works with some specific settings (e.g., `minishift start --vm-driver=generic` or `minikube start --vm-driver=none`).
  With this, `minishift` and `minikube` use the local docker installation of Travis-CI. Please note, with this setup
  we are limited to 2 vCPUs and 4 GB of memory (see [https://docs.travis-ci.com/user/reference/overview/#virtualisation-environment-vs-operating-system](https://docs.travis-ci.com/user/reference/overview/#virtualisation-environment-vs-operating-system) for details)

Above limitations lead to the following setup with Keptn's `.travis.yml`:

* We test the full installation (`--use-case=continuous-delivery`) on `GKE` once per day on the master branch via a Travis-CI cron job.
* We test the default installation on Minikube, MicroK8s and Minishift for every push on the master branch (e.g., after merging a pull-request).
  * This type of installation is limited to `--gateway=NodePort`

## Prepare your local environment to run integration tests

When running integration tests locally, we recommend using either K3s, Minikube or Minishift.

### Set-up steps for K3s (recommended on Linux)

Starting and setting up K3s is easy:

1. Download, install and run K3s (tested with versions 1.16 to 1.18):
    ```console
    curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.18.3+k3s1 K3S_KUBECONFIG_MODE="644" sh -s - --no-deploy=traefik
    ```
1. Export the Kubernetes config using:
    ```console
    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
    ```  
    Please read more about it [here](https://rancher.com/docs/k3s/latest/en/).
1. Verify that everything has worked using `kubectl get nodes`

### Set-up steps for Minikube

In case you plan to use Minikube, the following steps are necessary to get Keptn running on Minikube:

1. Download and install Minikube (versions 1.4 to 1.10 should work) from [https://github.com/kubernetes/minikube/releases](https://github.com/kubernetes/minikube/releases) for your operating system.
1. Make sure you create a fresh Minikube cluster (using a VM):
    ```console
    minikube stop # optional
    minikube delete # optional
    minikube start [--vm-driver=...] --cpus 2 --memory 4096 # or: minikube start --cpus 6 --memory 12200
    ```
    **Note**: In some cases you have to specify `--vm-driver=...` to select virtualbox, kvm, hyperv, etc... - please find out more that the [official MiniKube docs](https://kubernetes.io/docs/setup/learning-environment/minikube/#specifying-the-vm-driver)
    
    **Note**: On Linux `--vm-driver=docker` *should* work in most cases
    
    **Note**: For the quality gates use case 2 CPUs and 4096 MB memory is enough, for the full installation we recommend at least 6 CPUs and 12200 MB memory

1. Verify that everything has worked using `kubectl get nodes`

### Set-up steps for Minishift

In case you plan to use Minishift, the following steps are necessary to get Keptn running on Minishift (1.34.2):

1. Download and install Minishift 1.34.2 (from [https://github.com/minishift/minishift/releases](https://github.com/minishift/minishift/releases))
   for your operating system
1. Setup Minishift profile, cpu and memory limits:
   ```console
   # make sure you have a profile is set correctly
   minishift profile set keptn-dev
   # minimum memory required for the minishift VM
   minishift config set memory 4GB # Note: use 12 GB for full installation
   # the minimum required vCpus for the minishift vm
   minishift config set cpus 2 # Note: use 6 for full installation
   # Add new user called admin with password admin having role cluster-admin
   minishift addons enable admin-user
   # Allow the containers to be run with uid 0
   minishift addons enable anyuid
   ```
   
   **Note**: For the quality-gates use-case 2 vCPUs and 4 GB memory is enough. For the full installation you should use
   more CPUs (e.g., 6) and more memory (e.g., 12 GB)
   
1. Start Minishift:
   ```console
   minishift start
   ```
1. Enable admission WebHooks on your OpenShift master node:
   ```console
   minishift openshift config set --target=kube --patch '{
       "admissionConfig": {
           "pluginConfig": {
               "ValidatingAdmissionWebhook": {
                   "configuration": {
                       "apiVersion": "apiserver.config.k8s.io/v1alpha1",
                       "kind": "WebhookAdmission",
                       "kubeConfigFile": "/dev/null"
                   }
               },
               "MutatingAdmissionWebhook": {
                   "configuration": {
                       "apiVersion": "apiserver.config.k8s.io/v1alpha1",
                       "kind": "WebhookAdmission",
                       "kubeConfigFile": "/dev/null"
                   }
               }
           }
       }
   }'
   ```
1. Login via `oc` cli (you might need to try a couple of times):
   ```console
   oc login -u admin -p admin
   ```
   **Note**: It takes a couple of minutes before your Minishift cluster is ready. Expect error messages like
   ```
   Error from server (InternalError): Internal error occurred: unexpected response: 400
   ```
   and retry again.
1. Set policies/permissions:
   ```console
   oc adm policy --as system:admin add-cluster-role-to-user cluster-admin admin
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:default:default
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:kube-system:default
   ```
1. Note down the Minishift server URL printed by `oc status` (e.g., `https://192.168.99.101:8443`)


## Run the integration tests locally

### Run the quality-gates integration test (with dynatrace-sli-service) locally

Pre-requesits:

* A K8s cluster (see above)
  * 4 GB of free memory (less is possible but you will notice a significant slowdown)
  * at least 2 CPU cores on your system (less is possible but you will notice a significant slowdown)
* `kubectl`
* `keptn` CLI
* A local copy of this repository (keptn/keptn)
* Linux, Mac OS or Linux Subsystem for Windows
* Bash
* Dynatrace Tenant with a service that has the tags `FrontEnd` and `testdeploy`
* Dynatrace API Token that is allowed to read metrics from the Dynatrace Tenant
* working internet connection


1. Setup your Kubernetes cluster (see above)
1. `keptn` CLI available
1. Verify `kubectl` is configured to the right cluster:
   ```console
   kubectl get nodes
   ```
1. Install Keptn
   * K3s
       ```console
       keptn install --verbose
       ```
   * Minikube
       ```console
       keptn install --endpoint-service-type=NodePort --verbose
       ```
       **Note**: You can use `--endpoint-service-type=LoadBalancer` if you run `minikube tunnel`
   * Minishift
      ```console
      keptn install --platform=openshift --verbose
      ```
   **Note**: Use `--chart-repo=https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz` to install the latest master branch, see [master installation docs](install_master.md).
1. Verify the installation has worked
   ```console
   keptn status
   ```
1. Verify which images have been deployed
   ```console
   kubectl -n keptn get deployments -owide
   ```
1. Expose Keptn Bridge for better troubleshooting and debugging
   ```console
   keptn configure bridge --action=expose
   ```
1. Verify accessing Keptn Bridge works
1. Configuration for the test:
   ```console
   export PROJECT=easytravel
   export SERVICE=frontend
   export STAGE=hardening
   export QG_INTEGRATION_TEST_DT_TENANT=<INSERT_YOUR_DT_TENANT_HERE>
   export QG_INTEGRATION_TEST_DT_API_TOKEN=<INSERT_YOUR_DT_API_TOKEN_HERE>
   ```
1. Run the test
   ```console
   bash ./test/test_quality_gates_standalone.sh
   ```

### Run the full installation integration test locally

Pre-requesits:

* A K8s cluster (see above)
  * 12 GB of free memory (less is possible but you will notice a significant slowdown)
  * at least 6 CPU cores on your system (less is possible but you will notice a significant slowdown)
* `kubectl`
* `keptn` CLI
* A local copy of this repository (keptn/keptn)
* Linux, Mac OS or Linux Subsystem for Windows
* Bash
* working internet connection

1. Setup your Kubernetes cluster (see above)
1. Verify `kubectl` is configured to the right cluster:
   ```console
   kubectl get nodes
   ```
1. `keptn` CLI available
1. Install Keptn
   * K3s
       ```console
       keptn install --use-case=continuous-delivery --verbose
       ```
   * Minikube
       ```console
       keptn install --use-case=continuous-delivery --endpoint-service-type=NodePort --verbose
       ```
       **Note**: You can use `--endpoint-service-type=LoadBalancer` if you run `minikube tunnel`
   * Minishift
      ```console
      keptn install --use-case=continuous-delivery --platform=openshift --verbose
      ```
   **Note**: Use `--chart-repo=https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz` to install the latest master branch, see [master installation docs](install_master.md).
1. Verify the installation has worked
   ```console
   keptn status
   ```
1. Verify which images have been deployed
   ```console
   kubectl -n keptn get deployments -owide
   ```
1. Configuration for the test:
   ```console
   export PROJECT=sockshop
   ```
1. Run the test for onboarding a service
   ```console
   bash ./test/test_onboard_service.sh
   ```
1. Run the test for sending a new artifact
   ```console
   bash ./test/test_new_artifact.sh
   ```
