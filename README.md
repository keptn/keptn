![keptn](./assets/keptn.png)

# Keptn
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/keptn)
![GitHub Downloads](https://img.shields.io/github/downloads/keptn/keptn/total?logo=github&logoColor=white)
![CI](https://github.com/keptn/keptn/workflows/CI/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/keptn/keptn/branch/master/graph/badge.svg)](https://codecov.io/gh/keptn/keptn)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn/keptn)](https://goreportcard.com/report/github.com/keptn/keptn)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3588/badge)](https://bestpractices.coreinfrastructure.org/projects/3588)

Keptn is an event-based control plane for continuous delivery and automated operations for cloud-native applications. 
Please find the documentation on our [website](https://keptn.sh), and read the motivation about Keptn on our 
[Why Keptn?](https://keptn.sh/why-keptn/) page.

In addition, you can find the roadmap of the Keptn project [here](https://github.com/orgs/keptn/projects/1). It provides 
an overview of user stories that are currently in the focus of development for the next release.

## Quickstart

Keptn runs on Kubernetes. To get started, you can follow our [Quickstart guide](https://keptn.sh/docs/quickstart).

### Developing Keptn

The easiest way to develop is to spin up a Kubernetes cluster locally by using [K3d](https://k3d.io) via the following commands:

```console
curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | TAG=v4.4.4 bash
k3d cluster create mykeptn -p "8082:80@agent[0]" --k3s-server-arg '--no-deploy=traefik' --agents 1
```

Afterwards, install Keptn:
```console
curl -sL https://get.keptn.sh | bash
keptn install --use-case=continuous-delivery
```

You can find more information in our [docs](docs/).

## Community

Please find details on regular hosted community events as well as our Slack workspace in the 
[keptn/community repo](https://github.com/keptn/community).

## Keptn Versions compatibilities

We manage the Keptn *core components* in versions. The respective images in their versions are stored on [DockerHub](https://hub.docker.com/?namespace=keptn).
The versions of the Keptn *core components* and the services are compatible with each other. However, contributed services
as well as services that are not considered *core components* might not follow the same versioning schema.

We are tracking compatibility of those services [on our website](https://keptn.sh/docs/integrations/).

## Contributions

You are welcome to contribute using Pull Requests to the respective repositories. Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md) and our [Code of Conduct](CODE_OF_CONDUCT.md).
Please also check out our list of [good first issues](https://github.com/keptn/keptn/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).

## License

Keptn is an Open Source Project. Please see [LICENSE](LICENSE) for more information.

## Adopters

For a list of users, please refer to [ADOPTERS.md](ADOPTERS.md).

## Further information

* The [Keptn`s website](https://keptn.sh) has the documentation of Keptn and its use cases.
* Please join the [Keptn community](https://keptn.sh/community/).
