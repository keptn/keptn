![keptn](./assets/keptn.png)

# Keptn
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/keptn)
![Github Downloads](https://img.shields.io/github/downloads/keptn/keptn/total?logo=github&logoColor=white)
[![Build Status](https://travis-ci.org/keptn/keptn.svg?branch=master)](https://travis-ci.org/keptn/keptn)
[![codecov](https://codecov.io/gh/keptn/keptn/branch/master/graph/badge.svg)](https://codecov.io/gh/keptn/keptn)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn/keptn)](https://goreportcard.com/report/github.com/keptn/keptn)

Keptn is an event-based control plane for continuous delivery and automated operations for cloud-native applications. 
Please find the documentation on our [website](https://keptn.sh), and read the motivation about Keptn on our 
[Why Keptn?](https://keptn.sh/why-keptn/) page.

## Usage

Please find the documentation of how to get started with Keptn in the [Quick Start](https://keptn.sh/docs/quickstart/) and the [Installation instructions](https://keptn.sh/docs/0.6.0/installation/setup-keptn/). We recommend using the [latest stable release](https://github.com/keptn/keptn/releases).

Furthermore, you can learn about our current releases, release candidates and pre-releases on the [release section](https://github.com/keptn/keptn/releases).

## Versions compatibilities
We manage the Keptn core components in versions. The respective images in their versions are stored on [DockerHub](https://hub.docker.com/?namespace=keptn).
The versions of the Keptn core components and the services have to be compatible with each other.
Therefore, this section shows the compatibility between these versions.

**Keptn in [version 0.6.1](https://github.com/keptn/keptn/releases/tag/0.6.1) requires:**

*Keptn core:*
- keptn/api:0.6.1
- keptn/bridge:0.6.1
- keptn/configuration-service:0.6.1
- keptn/distributor:0.6.1
- keptn/eventbroker-go:0.6.1
- keptn/gatekeeper-service:0.6.1
- keptn/helm-service:0.6.1
- keptn/jmeter-service:0.6.1
- keptn/lighthouse-service:0.6.1
- keptn/mongodb-datastore:0.6.1
- keptn/shipyard-service:0.6.1
- keptn/wait-service:0.6.1
- keptn/remediation-service:0.6.1

*for Openshift:*
- keptn/openshift-route-service:0.6.1

<details><summary>Keptn version 0.6.0</summary>
<p>

*Keptn core:*
- keptn/api:0.6.0
- keptn/bridge:0.6.0
- keptn/configuration-service:0.6.0
- keptn/distributor:0.6.0
- keptn/eventbroker-go:0.6.0
- keptn/gatekeeper-service:0.6.0
- keptn/helm-service:0.6.0
- keptn/jmeter-service:0.6.0
- keptn/lighthouse-service:0.6.0
- keptn/mongodb-datastore:0.6.0
- keptn/shipyard-service:0.6.0
- keptn/wait-service:0.6.0
- keptn/remediation-service:0.6.0

*for Openshift:*
- keptn/openshift-route-service:0.6.0

</p>
</details>

<details><summary>Keptn version 0.5.2</summary>
<p>

*Keptn core:*
- keptn/api:0.5.0
- keptn/bridge:0.5.0
- keptn/configuration-service:0.5.0
- keptn/distributor:0.5.0
- keptn/eventbroker-go:0.5.0
- keptn/gatekeeper-service:0.5.0
- keptn/helm-service:0.5.1
- keptn/jmeter-service:0.5.0
- keptn/mongodb-datastore:0.5.0
- keptn/pitometer-service:0.5.0
- keptn/shipyard-service:0.5.0
- keptn/wait-service:0.5.0
- keptn/remediation-service:0.5.0


*Keptn uniform:*
- keptn/dynatrace-service:0.2.0
- keptn/prometheus-service:0.2.0
- keptn/servicenow-service:0.1.4

*for Openshift:*
- keptn/openshift-route-service:0.5.0

</p>
</details>

<details><summary>Keptn version 0.5.1</summary>
<p>

*Keptn core:*
- keptn/api:0.5.0
- keptn/bridge:0.5.0
- keptn/configuration-service:0.5.0
- keptn/distributor:0.5.0
- keptn/eventbroker-go:0.5.0
- keptn/gatekeeper-service:0.5.0
- keptn/helm-service:0.5.1
- keptn/jmeter-service:0.5.0
- keptn/mongodb-datastore:0.5.0
- keptn/pitometer-service:0.5.0
- keptn/shipyard-service:0.5.0
- keptn/wait-service:0.5.0
- keptn/remediation-service:0.5.0


*Keptn uniform:*
- keptn/dynatrace-service:0.3.1
- keptn/prometheus-service:0.2.0
- keptn/servicenow-service:0.1.4

*for Openshift:*
- keptn/openshift-route-service:0.5.0

</p>
</details>
<details><summary>Keptn version 0.5.0</summary>
<p>

Keptn in [version 0.5.0](https://github.com/keptn/keptn/releases/tag/0.5.0) requires:

*Keptn core:*
- keptn/api:0.5.0
- keptn/bridge:0.5.0
- keptn/configuration-service:0.5.0
- keptn/distributor:0.5.0
- keptn/eventbroker-go:0.5.0
- keptn/gatekeeper-service:0.5.0
- keptn/helm-service:0.5.0
- keptn/jmeter-service:0.5.0
- keptn/mongodb-datastore:0.5.0
- keptn/pitometer-service:0.5.0
- keptn/shipyard-service:0.5.0
- keptn/wait-service:0.5.0
- keptn/remediation-service:0.5.0


*Keptn uniform:*
- keptn/dynatrace-service:0.3.1
- keptn/prometheus-service:0.2.0
- keptn/servicenow-service:0.1.4

*for Openshift:*
- keptn/openshift-route-service:0.5.0

</p>
</details>

Please check out the [GitHub releases page](https://github.com/keptn/keptn/releases) if you need information for older Keptn versions.

## Roadmap

The roadmap of the Keptn project can be found [here](https://github.com/orgs/keptn/projects/1). It gives you an overview of user stories that are currently in the focus of development for the next release.

## Contributions

You are welcome to contribute using Pull Requests to the respective repositories. Before contributing, please read our [Contributing Guidelines](CONTRIBUTING.md) and our [Code of Conduct](CODE_OF_CONDUCT.md).
Please also check out our list of [good first issues](https://github.com/keptn/keptn/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22).

## License

Keptn is an Open Source Project. Please see [LICENSE](LICENSE) for more information.

## Further information

* The [Keptn`s website](https://keptn.sh) has the documentation of Keptn and its use cases.
* Please join the [Keptn community](https://github.com/keptn/community).
