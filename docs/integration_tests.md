# Integration Tests

## Structure of Integration Tests

### Adding new Integration Test
Adding a new Integration Tests consists of two steps:
1. Writing an Integration Test (please inspire yourself with other Integrations Tests) and place it to /test/go-tests folder. Please note that the naming convention is important (for example test_myNewTest.go).
2. Add the test to one or more Testsuites (files with `testsuite_` prefix).

## Running Integration Tests

There are two possibilities to run Integration Tests on Keptn project:
* running Integration Tests remotely on Github
* running Integration Tests locally

## Run Integration Tests remotely on Github

The possibility to run Integration Tests remotely is privileged to users, who are part of the Keptn project. There are two possibilities how to run Integration Tests:
* Running Integration Tests with default context for a specific branch (code changes ouside of /test folder)
* Running Integration Tests with a context from a specific branch (code changes inside of /test folder)

These two options can be also combined and currently only executions of all Integration Tests for all Testsuites is supported. The execution of the Tests is fairly easy:
1. Navigate to the `Actions` tab in `keptn/keptn` repository (https://github.com/keptn/keptn)
2. Choose `Integration Tests` from the left side menu
3. Click on `Run Workflow`, where a dialog window will appear. Here you need to choose the context (`Use Workflow from`) of the tests you wish to use (`master` is default). You should use this `master` context unless you have not made any changes in the Integration Tests. Secondly you choose a branch, from which the CI build artifacts (docker images) should be used from. Here you mostly use the branch of the code you are currently working on and want to run Integration Tests for your code changes.

## Run Integration Tests locally

### Prepare your local environment to run integration tests

When running integration tests locally, we recommend using either K3d or Minishift. Please use the set-up steps below to set up your local enviroment before installing 
Keptn and running the Integration Tests.

#### Set-up steps for K3d (recommended on Linux)

Starting and setting up K3d is easy:

1. Download and install K3d (**Note**: please be aware you need to have Docker installed, more info here: https://k3d.io/v5.2.2/):
    ```console
    curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
    ```
2. Create Kubernetes cluster:
    ```console
    k3d cluster create mykeptn -p "8082:80@loadbalancer" --k3s-arg "--no-deploy=traefik@server:*" --k3s-arg "--no-deploy=servicelb@server:*" --k3s-arg "--kube-proxy-arg=conntrack-max-per-core=0@server:*"  --agents 1
    ```  
3. Verify that everything has worked using `kubectl get nodes`

#### Set-up steps for Minishift

In case you plan to use Minishift, the following steps are necessary to get Keptn running on Minishift (1.34.3):

1. Download and install Minishift 1.34.3 (from [https://github.com/minishift/minishift/releases](https://github.com/minishift/minishift/releases))
   for your operating system
2. Setup Minishift profile, cpu and memory limits:
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
   
3. Start Minishift:
   ```console
   minishift start
   ```
4. Enable admission WebHooks on your OpenShift master node:
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
5. Login via `oc` cli (you might need to try a couple of times):
   ```console
   oc login -u admin -p admin
   ```
   **Note**: It takes a couple of minutes before your Minishift cluster is ready. Expect error messages like
   ```
   Error from server (InternalError): Internal error occurred: unexpected response: 400
   ```
   and retry again.
6. Set policies/permissions:
   ```console
   oc adm policy --as system:admin add-cluster-role-to-user cluster-admin admin
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:default:default
   oc adm policy  add-cluster-role-to-user cluster-admin system:serviceaccount:kube-system:default
   ```
7. Note down the Minishift server URL printed by `oc status` (e.g., `https://192.168.99.101:8443`)

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
2. Verify `kubectl` is configured to the right cluster:
   ```console
   kubectl get nodes
   ```
3. Install `keptn` CLI:
   ```console
   curl -sL https://get.keptn.sh | KEPTN_VERSION=0.12.0 bash
   ```
   **Note**: Please use the newest available Keptn version. Available versions can be found here: https://github.com/keptn/keptn/tags
4. Install Keptn
   * K3d
       ```console
       keptn install --use-case=continuous-delivery
       ```
   * Minishift
      ```console
      keptn install --use-case=continuous-delivery --platform=openshift --verbose
      ```
   **Note**: If you want to upgrade to the latest developer version, please use `helm upgrade` with `--reuse-values` option after installation.
5. Expose Keptn
   ```console
   curl -SL https://raw.githubusercontent.com/keptn/examples/master/quickstart/expose-keptn.sh | bash
   ```
6. Open a new terminal and type:
   ```console
   kubectl -n keptn port-forward service/api-gateway-nginx 8080:80
   ```
   After executing, return back to the original terminal.
7. Authenticate Keptn:
   ```console
   keptn auth --endpoint=http://127.0.0.1:8080/api --api-token=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)
   ```
8. Verify the installation has worked
   ```console
   keptn status
   ```
9. Verify which images have been deployed
   ```console
   kubectl -n keptn get deployments
   ```
10. Run tests (e.g., UniformRegistration):
   ```console
   cd test/go-tests && KEPTN_ENDPOINT="http://127.0.0.1:8080/api" go test ./...
   ```
   **Note**: If you want to run a single test, (for example BackupTestore_Test), please add `_test` suffix to the test file name, so it becomes executable. Otherwise, you will be able to run only the `testsuite_*_test.go` files. For running a single test use:
   ```console
   cd test/go-tests && KEPTN_ENDPOINT="http://127.0.0.1:8080/api" go test ./... -v -run <NameOfTheTest>
   ```
