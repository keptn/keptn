# Keptn Developer Docs

This folder contains docs for developers. If you are looking for the usage documentation of Keptn, or the `keptn` CLI, 
 please take a look at the [keptn.sh](https://keptn.sh/docs/) website.

## Requirements

* Kubernetes CLI tool [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* Docker
* You have installed the Keptn CLI and have a working installation of Keptn on Kubernetes (see [Quickstart](https://keptn.sh/docs/quickstart/)).
* Docker Hub Account (any other container registry works too)
* Go (Version 1.16.x)
* GitHub Account (required for making Pull Requests)
* If you want to use in-cluster debugging, please take a look at our [debugging guide](debugging.md).

### IDE / Code Editor

While this is not a requirement, we recommend you to use any of the following

* Visual Studio Code (with several Go Plugins)
* JetBrains GoLand (with Google Cloud Code)

## Where to go

Keptn consists of multiple services. We recommend to take a look at the 
[architecture of keptn](https://keptn.sh/docs/concepts/architecture/).

The Keptn core implementation as well as the *batteries-included* services (helm-service, jmeter-service) are stored 
 within this repository ([keptn/keptn](https://github.com/keptn/keptn)). 

In addition, the `go-utils` package is available in [keptn/go-utils](https://github.com/keptn/go-utils/) and contains
 several utility functions that we use in many services.
 
Similarly, the `kubernetes-utils` package is available in [keptn/kubernetes-utils](https://github.com/keptn/kubernetes-utils/) 
 and contains several utility functions that we use in many services.

If you want to contribute to the website or docs provided on the website, the 
 [keptn/keptn.github.io](https://github.com/keptn/keptn.github.io/) is the way to go.

Last but not least, we have a collection of additional services at 
* [github.com/keptn-contrib](https://github.com/keptn-contrib):
    * [dynatrace-service](https://github.com/keptn-contrib/dynatrace-service)
    * [dynatrace-sli-service](https://github.com/keptn-contrib/dynatrace-sli-service)
    * [prometheus-service](https://github.com/keptn-contrib/prometheus-service)
    * [unleash-service](https://github.com/keptn-contrib/unleash-service)
    * *archived* [prometheus-sli-service](https://github.com/keptn-contrib/prometheus-sli-service)
* [github.com/keptn-sandbox](https://github.com/keptn-sandbox):
    * [keptn-service-template-go](https://github.com/keptn-sandbox/keptn-service-template-go)
    * [litmus-service](https://github.com/keptn-sandbox/litmus-service)
    * [locust-service](https://github.com/keptn-sandbox/locust-service)
    * [job-executor-service](https://github.com/keptn-sandbox/job-executor-service)


## Branch Naming Convention

We aim to follow some rules for branch names within our repositories:

* **master** contains the latest (potentially unstable) changes from development
* **release-x.y.z** contains the release x.y.z
* **feature/###/name**, **bug/###/name**, **patch/###/name**, **hotfix/###/name** (where ### references the github issue number) contains 
  code for branches that are under active development.

We are following the git branching model described in [this blog post](https://nvie.com/posts/a-successful-git-branching-model/).
For instance, if a new feature or bug branch is created, the workflow is to create a new branch of the **master** 
 branch, and  name it according to the convention listed above. Once ready, a new Pull Request should be created with 
 the **master** branch as a target. Hotfixes work similar, but should be branched away from the **master** branch. PRs
 for hotfixes should be created to the **master** and respective **release-** branches, ensuring that the latest release
 and the current development version use the fix.
  
![Branch Workflow](assets/git_branch_workflow.png "Git Branch Workflow")

## CI Pipeline: Releases, Nightlies, etc...

We have automated builds for several services and containers using Github Actions. This automatically creates new builds for

* every [GitHub Release](https://github.com/keptn/keptn/releases) tagged with the version number (e.g., `0.8.4`),
* every change in the [master branch](https://github.com/keptn/keptn/tree/master) (unstable) tagged as `x.y.z-dev` (e.g., `0.8.4-dev`) as 
  well as the build-datetime,
* every pull request (unstable) tagged as `x.y.z-dev-PR-$PRID` (e.g., `0.8.4-dev-PR-1234`).


You can check the resulting images out by looking at our [DockerHub container registry](https://hub.docker.com/u/keptn)
 and at the respective containers.
