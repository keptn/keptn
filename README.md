![keptn](./assets/keptn.png)

# Keptn
Keptn is a fabric for cloud-native lifecycle automation at enterprise scale. In the current version it provides an automated setup of the Keptn core components as well as a demo application. Also included are three pre-configured use cases for the demo application: automated quality gates, runbook automation, and automated evaluation of blue/green deployments.

## Usage
Please find the documentation of how to get started with Keptn in [our official documentation](https://keptn.sh/docs) to get resources on how to use Keptn. We recommend to use the [latest stable release](https://github.com/keptn/keptn/releases).

Furthermore, please use the [release section](https://github.com/keptn/keptn/releases) to learn about our current releases, release candidates and pre-releases to get the latest version of Keptn.

## Versions compatibilities
We mangage the Keptn core components as well as all services (e.g., github-service, helm-service, etc.) in versions. The respective images in their versions are stored on [DockerHub](https://hub.docker.com/?namespace=keptn).
The versions of the Keptn core components and the services have to be compatible to each other.
Therefore, this section shows the compatibility between these versions.

**Keptn in [version 0.5.0](https://github.com/keptn/keptn/releases/tag/0.5.0) requires:**

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

<details><summary>Keptn version 0.4.0</summary>
<p>

Keptn in [version 0.4.0](https://github.com/keptn/keptn/releases/tag/0.4.0) requires:

*Keptn core:*
- keptn/authenticator:0.2.3
- keptn/bridge:0.1.3
- keptn/control:0.3.0
- keptn/eventbroker-go:0.1.0
- keptn/eventbroker-ext:0.3.0

*Keptn uniform:*
- keptn/gatekeeper-service:0.1.1
- keptn/github-service:0.3.0
- keptn/helm-service:0.1.1
- keptn/jmeter-service:0.1.1
- keptn/pitometer-service:0.2.0
- keptn/servicenow-service:0.1.3

*for Openshift:*
- keptn/openshift-route-service:0.1.1

</p>
</details>

<details><summary>Keptn version 0.3.0</summary>
<p>

Keptn in [version 0.3.0](https://github.com/keptn/keptn/releases/tag/0.3.0) requires:

*Keptn core:*
- keptn/authenticator:0.2.2
- keptn/bridge:0.1.2
- keptn/control:0.2.4
- keptn/eventbroker:0.2.3
- keptn/eventbroker-ext:0.2.3

*Keptn uniform:*
- keptn/gatekeeper-service:0.1.0
- keptn/github-service:0.2.2
- keptn/helm-service:0.1.0
- keptn/jmeter-service:0.1.0
- keptn/pitometer-service:0.1.3
- keptn/servicenow-service:0.1.2

*for Openshift:*
- keptn/openshift-route-service:0.1.0

</p>
</details>

<details><summary>Keptn version 0.2.2</summary>
<p>

Keptn in [version 0.2.2](https://github.com/keptn/keptn/releases/tag/0.2.2) requires:
- keptn/authenticator:0.2.2
- keptn/bridge:0.1.2
- keptn/control:0.2.3
- keptn/eventbroker:0.2.2
- keptn/eventbroker-ext:0.2.2
- keptn/pitometer-service:0.1.2
- keptn/servicenow-service:0.1.1
- keptn/github-service:0.2.1 
- keptn/jenkins-service:0.3.0
  - keptn/jenkins-0.6.0

</p>
</details>

<details><summary>Keptn version 0.2.1</summary>
<p>

Keptn in [version 0.2.1](https://github.com/keptn/keptn/releases/tag/0.2.1) requires:
- keptn/keptn-authenticator:0.2.1
- keptn/keptn-control:0.2.1
- keptn/keptn-event-broker:0.2.1
- keptn/keptn-event-broker-ext:0.2.1
- keptn/pitometer-service:0.1.1 
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.1 
- keptn/jenkins-service:0.2.0
  - keptn/jenkins-0.5.0

</p>
</details>

<details><summary>Keptn version 0.2.0</summary>
<p>

Keptn in [version 0.2.0](https://github.com/keptn/keptn/releases/tag/0.2.0) requires:
- keptn/keptn-authenticator:0.2.0
- keptn/keptn-control:0.2.0
- keptn/keptn-event-broker:0.2.0
- keptn/keptn-event-broker-ext:0.2.0
- keptn/pitometer-service:0.1.0
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.0
- keptn/jenkins-service:0.1.0
    - keptn/jenkins:0.4.0

</p>
</details>

<details><summary>Keptn version 0.1.3</summary>
<p>

Keptn in [version 0.1.3](https://github.com/keptn/keptn/tree/0.1.3) requires:

- keptn/jenkins:0.2
- dynatraceacm/ansibletower:3.3.1-1-2

</p>
</details>

## Further information
* The [Keptn`s website](https://keptn.sh) has the documentation of Keptn and its use cases.

