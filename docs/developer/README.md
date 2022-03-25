# Keptn Developer Docs

This folder contains docs for developers. If you are looking for the usage documentation of Keptn, or the `keptn` CLI,
please take a look at the [keptn.sh](https://keptn.sh/docs/) website.
We recommend to read our [CONTRIBUTING](../CONTRIBUTING.md) guidelines before reading this document.

<details>
<summary>Table of Contents</summary>

<!-- toc -->

- [Requirements](#requirements)
  * [IDE / Code Editor](#ide--code-editor)
- [Developer resources](#developer-resources)
- [Where to go](#where-to-go)
- [CI Pipeline: Releases, Nightlies, etc...](#ci-pipeline-releases-nightlies-etc)

<!-- tocstop -->

</details>

## Requirements

* Kubernetes CLI tool [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
* Docker
* You have installed the Keptn CLI and have a working installation of Keptn on Kubernetes (see [Quickstart](https://keptn.sh/docs/quickstart/)).
* Docker Hub Account (any other container registry works too)
* Go (Version 1.17.x)
* GitHub Account (required for making Pull Requests)
* If you want to use in-cluster debugging, please take a look at our [debugging guide](debugging.md).

### IDE / Code Editor

While this is not a requirement, we recommend you to use any of the following

* Visual Studio Code (with several Go Plugins)
* JetBrains GoLand (with Google Cloud Code)

## Developer resources

- [Bridge](./bridge.md)
- [Core](./core.md)
- [Debugging](./debugging.md)
- [Deploy /master](./install_master.md)
- [Fork](./fork.md)
- [Integration Tests](./integration_tests.md)
- [Pipelines](./pipelines.md)

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

Last but not least, we have a collection of additional services at [github.com/keptn-contrib](https://github.com/keptn-contrib)
and [github.com/keptn-sandbox](https://github.com/keptn-sandbox).

## CI Pipeline: Releases, Nightlies, etc...

We have automated builds for several services and containers using Github Actions. This automatically creates new builds for

* every [GitHub Release](https://github.com/keptn/keptn/releases) tagged with the version number (e.g., `0.8.4`),
* every change in the [master branch](https://github.com/keptn/keptn/tree/master) (unstable) tagged as `x.y.z-dev` (e.g., `0.8.4-dev`) as
  well as the build-datetime,
* every pull request (unstable) tagged as `x.y.z-dev-PR-$PRID` (e.g., `0.8.4-dev-PR-1234`).


You can check the resulting images out by looking at our [DockerHub container registry](https://hub.docker.com/u/keptn)
 and at the respective containers.
