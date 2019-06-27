[![Build Status](https://travis-ci.org/keptn/keptn.svg?branch=develop)](https://travis-ci.org/keptn/keptn)

![keptn](./assets/keptn.png)

# keptn
keptn is a fabric for cloud-native lifecycle automation at enterprise scale. In its first version it provides an automated setup of the keptn core components as well as a demo application. Also included are three preconfigured use cases for the demo application: automated quality gates, runbook automation, and automated evaluation of blue/green deployments.

## Usage
Please find the documentation of how to get started with keptn in [our official documentation](https://keptn.sh/docs) to get resources on how to use keptn. We recommend to use the [latest stable release](https://github.com/keptn/keptn/releases).

Furthermore, please use the [release section](https://github.com/keptn/keptn/releases) to learn about our current releases, release candidates and pre-releases to get the latest version of keptn.


## Versions compatibilities
We mangage the keptn core components as well as all services (e.g. gitHub-service, helm-service) in versions. The respective images in their versions are stored in [DockerHub](https://hub.docker.com/?namespace=keptn).
The versions of the keptn core components and the services have to be compatible to each other.
Therefore, this section shows the compatibility between these versions.

keptn in [version 0.3.0](https://github.com/keptn/keptn/releases/tag/0.2.1) requires:

Keptn core
- keptn/authenticator:0.2.2
- keptn/bridge:0.1.2
- keptn/control:0.2.4
- keptn/eventbroker:0.2.3
- keptn/eventbroker-ext:0.2.3

Uniform
- keptn/gatekeeper-service:0.1.0
- keptn/github-service:0.2.2
- keptn/helm-service:0.1.0
- keptn/jmeter-service:0.1.0
- keptn/pitometer-service:0.1.3
- keptn/servicenow-service:0.1.2

For Openshift
- keptn/openshift-route-service:0.1.0

keptn in [version 0.2.2](https://github.com/keptn/keptn/releases/tag/0.2.2) requires:
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
  
keptn in [version 0.2.1](https://github.com/keptn/keptn/releases/tag/0.2.1) requires:
- keptn/keptn-authenticator:0.2.1
- keptn/keptn-control:0.2.1
- keptn/keptn-event-broker:0.2.1
- keptn/keptn-event-broker-ext:0.2.1
- keptn/pitometer-service:0.1.1 
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.1 
- keptn/jenkins-service:0.2.0
  - keptn/jenkins-0.5.0

keptn in [version 0.2.0](https://github.com/keptn/keptn/releases/tag/0.2.0) requires:
- keptn/keptn-authenticator:0.2.0
- keptn/keptn-control:0.2.0
- keptn/keptn-event-broker:0.2.0
- keptn/keptn-event-broker-ext:0.2.0
- keptn/pitometer-service:0.1.0
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.0
- keptn/jenkins-service:0.1.0
    - keptn/jenkins:0.4.0

keptn in [version 0.1.3](https://github.com/keptn/keptn/tree/0.1.3) requires:
- keptn/jenkins:0.2
- dynatraceacm/ansibletower:3.3.1-1-2

## Further information
* The [keptn website](https://keptn.sh) has the documentation of keptn and its usecases.

