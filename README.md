![keptn](./assets/keptn.png)

# Keptn

![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/keptn)
![GitHub Downloads](https://img.shields.io/github/downloads/keptn/keptn/total?logo=github&logoColor=white)
![CI](https://github.com/keptn/keptn/workflows/CI/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/keptn/keptn/branch/master/graph/badge.svg)](https://codecov.io/gh/keptn/keptn)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn/keptn)](https://goreportcard.com/report/github.com/keptn/keptn)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3588/badge)](https://bestpractices.coreinfrastructure.org/projects/3588)



Keptn is an event-based control plane for continuous delivery and automated operations for cloud-native applications. 
Please find the documentation on our [website](https://keptn.sh/), and read the motivation about Keptn on our 
[Why Keptn?](https://keptn.sh/why-keptn/) page.

In addition, you can find the roadmap of the Keptn project [here](https://github.com/orgs/keptn/projects/10). It provides 
an overview of user stories that are currently in the focus of development for the next release.

## Keptn Today! Keptn Lifecycle Toolkit Tomorrow!

### Keptn: Moving towards our 1.0 milestone!

3 years of hard work will soon reach a long awaited milestone: [Keptn 1.0 with LTS (Long Time Support)](https://docs.google.com/document/d/1RdFegnZrxjWxJAem9auaeVQ5_mKl5wFlwd6MgF1ot0s/edit#heading=h.qoctq8iujkhs) brings you automated release validation based on SLOs that you can easily integrate into your existing DevOps Tools (deployment, test and observability).

If you want to explore Keptn then keep scrolling down to get all information!

### Keptn Lifecycle Toolkit: The kuber-native road ahead!

At KubeCon 2022 in Detroit we announced the direction we are heading: Keptn Lifecycle Toolkit!
Keptn Lifecycle Toolkit brings application-aware deployment lifecycle management to your k8s cluster: 
* kubernetes-native: no external dependencies, everything in your CRDs and 
* pipeline-less: works with any delivery tool (ArgoCD, Flux, Jenkins, GitHub, GitLab, Harness ...) without having to integrate it with Keptn

To decide whether you want to stick with Keptn 1.0 or focus on Keptn Lifecycle Toolkit do this:
1. Watch our [Keptn Lifecycle Toolkit in a Nutshell](https://www.youtube.com/watch?v=K-cvnZ8EtGc) video including live demo
2. Try Keptn Lifecycle Toolkit yourself: [GitHub Repo](https://github.com/keptn/lifecycle-toolkit/)

## Quickstart

Keptn runs on Kubernetes. To get started, you can follow our [Quickstart guide](https://keptn.sh/docs/quickstart).

### Developing Keptn

The easiest way to develop is to spin up a Kubernetes cluster locally by using [K3d](https://k3d.io) (requires `docker`) via the following commands:

```console
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | TAG=v5.3.0 bash
k3d cluster create mykeptn -p "8082:80@loadbalancer" --k3s-arg "--no-deploy=traefik@server:*"
```

Afterwards, install Keptn CLI:
```console
curl -sL https://get.keptn.sh | bash
```

And install Keptn using Helm:
```console
helm repo add keptn https://charts.keptn.sh && helm repo update
helm install keptn keptn/keptn \
-n keptn --create-namespace \
--wait \
--set=apiGatewayNginx.type=LoadBalancer
```

In case you want to install `helm-service` and `jmeter-service`, execute:

```console
helm install jmeter-service keptn/jmeter-service -n keptn
helm install helm-service keptn/helm-service -n keptn
```

Please follow the instructions printed by the CLI to connect to your Keptn installation.

### Installing Keptn from Master branch

Note: This will install a potentially unstable version of Keptn.

If you want to install the latest master version of Keptn onto your cluster you can do that by using the development helm charts repository located at https://charts-dev.keptn.sh .
By default, the global registry used is ´ghcr.io/keptn´, so you will need to override it.

```console
helm repo add keptn-dev https://charts-dev.keptn.sh    # Add the keptn-dev helm repo
helm repo update                                       # Update all repo contents
helm search repo keptn-dev --devel --versions          # List all versions present in the keptn-dev repo

# Select a chart version from the previous command that you want to install

helm install -n keptn-dev keptn keptn-dev/keptn --set=global.keptn.registry=ghcr.io/keptn --create-namespace --version "<the-version-you-selected-previously>"
```

You can find more information in our [docs](docs/).

## Community

Please find details on regular hosted community events as well as our Slack workspace in the 
[keptn/community repo](https://github.com/keptn/community).

## Keptn Versions compatibilities

We manage the Keptn *core components* in versions.
The versions of the Keptn *core components* and the services are compatible with each other. However, contributed services
as well as services that are not considered *core components* might not follow the same versioning schema.

We are tracking compatibility of those services [on our website](https://keptn.sh/docs/integrations/).

## Container Images

Keptn provides container images of all *core components*.
The respective images in their versions are stored on the following container registries:

* [GitHub Container Registry](https://github.com/orgs/keptn/packages?repo_name=keptn)
* [Quay.io Container Registry](https://quay.io/organization/keptn)

From version 0.19.0, all released container images are signed using [cosign](https://github.com/sigstore/cosign)
with a keyless signing mechanism.
That means that Keptn uses short-lived code signing certificates and keys together with OIDC and a transparency log
to sign all its container images.
More info on keyless signed container images can be found [here](https://github.com/sigstore/cosign/blob/main/KEYLESS.md).


## Helm Chart

Keptn provides Helm charts for easy installation of all control plane components.
From version 0.19.0, the charts are signed and can be verified with the public key that can be found in [assets/pubring.gpg](assets/pubring.gpg)
and attached to every release.
More info on signed Helm charts can be found [here](https://helm.sh/docs/topics/provenance/).

## Contributions

You are welcome to contribute using Pull Requests to the respective repositories. Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md) and our [Code of Conduct](CODE_OF_CONDUCT.md).
Please also check out our list of [good first issues](https://github.com/keptn/keptn/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).

## License

Keptn is an Open Source Project. Please see [LICENSE](LICENSE) for more information.

## Adopters

For a list of users, please refer to [ADOPTERS.md](https://github.com/keptn/community/blob/main/ADOPTERS).

## Further information

* The [Keptn`s website](https://keptn.sh) has the documentation of Keptn and its use cases.
* Please join the [Keptn community](https://keptn.sh/community/).
