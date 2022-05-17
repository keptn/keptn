# Integration Tests


<details>
<summary>Table of Contents</summary>

<!-- toc -->

- [Integration Tests](#integration-tests)
  - [Structure of Integration Tests](#structure-of-integration-tests)
    - [Adding a new Integration Test](#adding-a-new-integration-test)
  - [Running Integration Tests](#running-integration-tests)
  - [Run Integration Tests remotely on Github](#run-integration-tests-remotely-on-github)
  - [Run Integration Tests locally](#run-integration-tests-locally)
    - [Prepare your local environment to run integration tests](#prepare-your-local-environment-to-run-integration-tests)
      - [**Setup steps for K3d (recommended on Linux)**](#setup-steps-for-k3d-recommended-on-linux)
      - [**Setup steps for Minishift (not recommended)**](#setup-steps-for-minishift-not-recommended)
    - [Run the full installation of Integration Tests locally](#run-the-full-installation-of-integration-tests-locally)

<!-- tocstop -->
</details>

## Structure of Integration Tests

The Integration Tests and their resources are located under the `/test` directory in this repository. For running the Integration Tests, there are two main directories, we will focus on:
* `/test/assets` -> directory containing resources and scripts, which are used during the run of the Integration Tests
* `/test/go-tests` -> Integration Tests and testsuites

Integration Tests are organized into four main testsuites (testsuite files have `testsuite_` prefixes), where every testsuite consists of tests that are run on a specific Kubernetes platform:
* GKE
* K3D
* K3S
* Minishift

These testsuites are run in parallel during pipeline execution on Github.

Each testsuite consists of one or more tests, which are actually Go functions executed in a testing context. These functions (tests) are stored in files with `test_` prefixes. Also, each test can be part of one or more testsuites.

### Adding a new Integration Test

Adding a new Integration Test consists of two steps:
1. Write an Integration Test (please take other Integrations Tests as reference/inspirations) and put it into the `/test/go-tests` directory. Please note that the naming convention is important (e.g. test_myNewTest.go).
2. Add the test to one or more testsuites (files with `testsuite_` prefix).

## Running Integration Tests

There are two possibilities to run Integration Tests:
* running Integration Tests remotely on Github
* running Integration Tests locally

## Run Integration Tests remotely on Github

The possibility to run Integration Tests remotely is restricted to users, who are part of the Keptn project. There are two possibilities how to run Integration Tests:
* Running Integration Tests with the default context for a specific branch (code changes outside of the `/test` directory)
* Running Integration Tests with a context from a specific branch (code changes inside of the `/test` directory)

These two options can be also combined and currently, only executions of all Integration Tests for all testsuites are supported. The execution of the Tests is fairly easy:
1. Navigate to the `Actions` tab in the `keptn/keptn` repository (https://github.com/keptn/keptn)
2. Choose `Integration Tests` from the left side menu
3. Click on `Run Workflow`, where a dialog window will appear. 
   Here, you need to choose the context (`Use Workflow from`) of the tests you wish to use (`master` is the default). 
   You should use this `master` context unless you have not made any changes in the Integration Tests pipeline. 
   Secondly, you choose a branch, from which the CI build artifacts (docker images) should be used. 
   Here, you mostly use the branch of the code you are currently working on and want to run Integration Tests for your code changes. Please be aware, that you need to wait for the docker images to be built before you can execute the Integartion Tests.

## Run Integration Tests locally

### Prepare your local environment to run integration tests

When running integration tests locally, we recommend using either [K3d](https://k3d.io/) or [Minishift](https://github.com/minishift/minishift). Please use the setup steps below to set up your local environment before installing Keptn and running the Integration Tests.

#### **Setup steps for K3d (recommended on Linux)**

Starting and setting up K3d is easy:

1. Download and install K3d (version 5.2.2 is recommended) (**Note**: please be aware you need to have Docker installed, more info here: https://k3d.io/):
    ```console
    curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | TAG=v5.2.2 bash
    ```
2. Create Kubernetes cluster:
    ```console
    k3d cluster create mykeptn -p "8082:80@loadbalancer" --k3s-arg "--no-deploy=traefik@server:*" --k3s-arg "--no-deploy=servicelb@server:*" --k3s-arg "--kube-proxy-arg=conntrack-max-per-core=0@server:*"  --agents 1
    ```  
3. Verify that everything has worked using `kubectl get nodes`

#### **Setup steps for Minishift (not recommended)**

In the case you plan to use Minishift, the following steps are necessary to get Keptn running on Minishift (1.34.3):

1. Download and install Minishift 1.34.3 (from [https://github.com/minishift/minishift/releases](https://github.com/minishift/minishift/releases))
   for your operating system
   * extract the archive `tar -zxvf minishift-1.34.3-linux-amd64.tgz && cd minishift-1.34.3-linux-amd64`
   * export the minishift path `export PATH=$PATH:$HOME/your_path/minishift-1.34.3-linux-amd64`

1. Setup Minishift profile, cpu and memory limits:
   ```console
   # make sure you have a profile is set correctly
   minishift profile set keptn-dev
   # minimum memory required for the minishift VM
   minishift config set memory 12GB
   # the minimum required vCpus for the minishift vm
   minishift config set cpus 6
   # Add new user called admin with password admin having role cluster-admin
   minishift addons enable admin-user
   # Allow the containers to be run with uid 0
   minishift addons enable anyuid
   ```
   
2. Start Minishift:
   ```console
   minishift start
   ```
   **Note**: Please make sure you have your Virtualization Environment properly set up before executing `minishift start`: https://docs.okd.io/3.11/minishift/getting-started/setting-up-virtualization-environment.html
3. Enable admission WebHooks on your OpenShift master node:
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
4. Login via `oc` cli (you might need to try a couple of times):
   ```console
   oc login -u admin -p admin
   ```
   **Note**: It takes a couple of minutes before your Minishift cluster is ready. Expect error messages like
   ```
   Error from server (InternalError): Internal error occurred: unexpected response: 400
   ```
   and retry again.
5. Set policies/permissions:
   ```console
   oc adm policy --as system:admin add-cluster-role-to-user cluster-admin admin
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:default:default
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:kube-system:default
   ```
6. Note down the Minishift server URL printed by `oc status` (e.g., `https://192.168.99.101:8443`)

### Run the full installation of Integration Tests locally

Pre-requisites:

* A K8s cluster (see above)
  * 12 GB of free memory (less is possible but you will notice a significant slowdown)
  * at least 6 CPU cores on your system (less is possible but you will notice a significant slowdown)
* `kubectl`
* A local copy of this repository (keptn/keptn)
* Linux, Mac OS or Linux Subsystem for Windows
* Bash
* working internet connection

1. Setup your Kubernetes cluster (see above)
2. Install `keptn` CLI:
   ```console
   curl -sL https://get.keptn.sh | bash
   ```
   **Note**: Please use the newest available Keptn version. Available versions can be found here: https://github.com/keptn/keptn/tags
3. Install Keptn
   * K3d
       ```console
       keptn install --use-case=continuous-delivery
       ```
   * Minishift
      ```console
      keptn install --use-case=continuous-delivery --platform=openshift --verbose
      ```
   **Note**: If you want to upgrade to the latest developer version, please use `helm upgrade` with `--reuse-values` option after installation.
4. Store Keptn namespace to env variable
   ```console
   KEPTN_NAMESPACE=keptn
   ```
5. Install Mockserver
   ```console
   helm repo add mockserver https://www.mock-server.com
   helm upgrade --install --namespace $KEPTN_NAMESPACE --version 5.13.0 mockserver mockserver/mockserver --set service.type=ClusterIP
   ```
6. Install Gitea
   ```console
   curl -SL https://raw.githubusercontent.com/keptn/keptn/master/docs/developer/install_gitea.sh | bash
   ```
7. Expose Keptn
   ```console
   curl -SL https://raw.githubusercontent.com/keptn/examples/master/quickstart/expose-keptn.sh | bash
   ```
8. Open a new terminal and type:
   ```console
   kubectl -n $KEPTN_NAMESPACE port-forward service/api-gateway-nginx 8080:80
   ```
   After executing, return back to the original terminal.
9.  Set Up env variables
   ```console
   KEPTN_ENDPOINT=http://$(kubectl -n $KEPTN_NAMESPACE get ingress api-keptn-ingress -ojsonpath='{.spec.rules[0].host}')/api
   KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n $KEPTN_NAMESPACE -ojsonpath='{.data.keptn-api-token}' | base64 --decode)
   KEPTN_BRIDGE_URL=http://$(kubectl -n $KEPTN_NAMESPACE get ingress api-keptn-ingress -ojsonpath='{.spec.rules[0].host}')/bridge
   ```
11. Authenticate Keptn:
   ```console
   keptn auth --endpoint=$KEPTN_ENDPOINT --api-token=$KEPTN_API_TOKEN
   ```
11. Verify the installation has worked
   ```console
   keptn status
   ```
11. Verify which images have been deployed
   ```console
   kubectl -n $KEPTN_NAMESPACE get deployments
   ```
11. Run tests (e.g., UniformRegistration):
   ```console
   cd test/go-tests && go test ./...
   ```
   **Note**: If you want to run a single test, (e.g. BackupTestore_Test), please add `_test` suffix to the test file name, so it becomes executable. Otherwise, you cano run only `testsuite_*_test.go` files. For running a single test use:
   ```console
   cd test/go-tests && go test ./... -v -run <NameOfTheTest>
   ```
