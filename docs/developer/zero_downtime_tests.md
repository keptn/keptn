# Zero Downtime Tests


<details>
<summary>Table of Contents</summary>

<!-- toc -->

- [Zero Downtime Tests](#zero-downtime-tests)
  - [Structure of Zero Downtime Tests](#structure-of-zero-downtime-tests)
    - [Adding a new Zero Downtime Test](#adding-a-new-zero-downtime-test)
  - [Running Zero Downtime Tests](#running-zero-downtime-tests)
  - [Run Zero Downtime Tests remotely on Github](#run-zero-downtime-tests-remotely-on-github)
  - [Run Zero Downtime Tests locally](#run-zero-downtime-tests-locally)
    - [Prepare your local environment to run Zero Downtime tests](#prepare-your-local-environment-to-run-zero-downtime-tests)
    - [Run the full installation of Zero Downtime Tests locally](#run-the-full-installation-of-zero-downtime-tests-locally)

<!-- tocstop -->
</details>

## Structure of Zero Downtime Tests

The Zero Downtime Tests and their resources are located under the `/test` directory in this repository. There are two main directories:
* `/test/zero-downtime/assets` -> directory containing resources and scripts, which are used during the run of the Zero Downtime Tests 
* `/test/zero-downtime` -> Zero Downtime Tests and testsuites

Zero Downtime Tests are organized into testsuites, a main suite to run all tests, and a separate file for each specific scenario:
* testsuite_zerodowntime_test contains the Test_ZeroDowntime function that runs all test suites, each test is based on the ZDTestTemplate that runs a test function in parallel to rolling upgrade and API probing
* sequence_test contains test for default sequences: evaluation and deployment 
* API_test contains API probes of "indempotent" endpoints to be run nonstop to verify keptn API availability during upgrades 
* webhook_test contains tests for sequences calling webhooks

API tests are run in parallel during pipeline execution on Github, webhook and sequences are passed to the  ZDTemplate and are run sequentially.

### Adding a new Zero Downtime Test

Adding a new Zero DowntimeTest consists of two steps:
1. Copy the webhook test suite file in `/test/zero-downtime` directory as a reference. Each `Test_` function will be run sequentially, `SetupSuite` will run once at suite creation. To change the behavior of the test suite during upgrades edit
    ```
    func Webhook(t *testing.T, env *ZeroDowntimeEnv) 
    ```
   This is currently set to run a new suite continuously until the upgrades are terminated
2. Add the new created function to the `testsuite_zerodowntime_test` like so 
   ```
   func (suite *TestSuiteDowntime) TestWebhook() {
   ZDTestTemplate(suite.T(), Webhook, "Webhook")
   }.
   ```

## Running Zero Downtime Tests

There are two possibilities to run Zero Downtime Tests:
* running Zero Downtime Tests remotely on Github
* running Zero Downtime Tests locally

## Run Zero Downtime Tests remotely on Github

The possibility to run Zero Downtime Tests remotely is restricted to users, who are part of the Keptn project. There are two possibilities how to run Zero Downtime Tests:
* Running Zero Downtime Tests with the default context for a specific branch (code changes outside of the `/test` directory)
* Running Zero Downtime Tests with a context from a specific branch (code changes inside of the `/test` directory)

These two options can be also combined and currently, only executions of all Zero Downtime Tests for all testsuites are supported. The execution of the Tests is fairly easy:
1. Navigate to the `Actions` tab in the `keptn/keptn` repository (https://github.com/keptn/keptn)
2. Choose `Zero Downtime Tests` from the left side menu
3. Click on `Run Workflow`, where a dialog window will appear. 
   Here, you need to choose the context (`Use Workflow from`) of the tests you wish to use (`master` is the default). 
   You should use this `master` context unless you have not made any changes in the Zero Downtime Tests pipeline. 
   Secondly, you choose a branch, from which the CI build artifacts (docker images) should be used.
   Here, you mostly use the branch of the code you are currently working on and want to run Zero Downtime Tests for your code changes. Please be aware, that you need to wait for the docker images to be built before you can execute the Integartion Tests.
   You have then to pass the image tag for the upgrade command e.g. 0.15.1-dev or a PR one like 0.15.0-PR-6897.
   If you want to perform UI test tick the first box. If you prefer to delete the generated cluster immediately tick the second box.

## Run Zero Downtime Tests locally

### Prepare your local environment to run Zero Downtime Tests

When running Zero Downtime Tests locally, we recommend using either [K3d](https://k3d.io/) or [Minishift](https://github.com/minishift/minishift). Please follow the instruction for Integration Test to setup the cluster correctly.

### Run the full installation of Zero Downtime Tests locally

Pre-requisites:

* A K8s cluster (see above)
  * 12 GB of free memory (less is possible but you will notice a significant slowdown)
  * at least 6 CPU cores on your system (less is possible but you will notice a significant slowdown)
* `kubectl`
* A local copy of this repository (keptn/keptn)
* Linux, Mac OS or Linux Subsystem for Windows
* Bash
* working internet connection
* `keptn` CLI
* Authenticated and exposed Keptn installation
* Mockserver
* Gitea

Env variables: 
   ```console
   export KEPTN_NAMESPACE=keptn
   export KEPTN_ENDPOINT=http://$(kubectl -n $KEPTN_NAMESPACE get ingress api-keptn-ingress -ojsonpath='{.spec.rules[0].host}')/api
   export KEPTN_API_TOKEN=$(kubectl get secret keptn-api-token -n $KEPTN_NAMESPACE -ojsonpath='{.data.keptn-api-token}' | base64 --decode)
   export KEPTN_BRIDGE_URL=http://$(kubectl -n $KEPTN_NAMESPACE get ingress api-keptn-ingress -ojsonpath='{.spec.rules[0].host}')/bridge
   ```

Mandatory variables:

* `INSTALL_HELM_CHART` specifies the current installed version e.g. "https://github.com/keptn/helm-charts-dev/raw/69eea439a26a99ecc163e296860dbb5d43e41600/packages/keptn-0.15.1-dev.tgz"
* `UPGRADE_HELM_CHART` specifies what version we want to upgrade to e.g. "https://github.com/keptn/helm-charts-dev/raw/gh-pages/packages/keptn-0.15.0-dev.tgz"

A few variables allows to setup the behaviour of the test locally, those can also be edited in the default field of the ZeroDowntimeEnv struct or by setting up local env variables.

* `NUMBER_OF_UPGRADES` default:2 specifies the amount of times we upgrade back and forth from the original installed version to the upgrade version
* `API_PROBES_INTERVAL` default:5s  sets frequencies between API tests runs
* `SEQUENCES_INTERVAL` default:15s  sets frequences between sequences/webhook tests suite runs


Run a single test ex. an API probe or a single run of Test_Sequences
   ```console
   cd test/go-tests && go test ./... -v -run <Test_Function_Name>
   ```
Run the whole zero downtime suite
   ```console
   cd test/go-tests && go test ./... -v -run Test_ZeroDowntime
   ```
