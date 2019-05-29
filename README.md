![keptn](./assets/keptn.png)

# keptn
keptn is a fabric for cloud-native lifecycle automation at enterprise scale. In its first version it provides an automated setup of the keptn core components as well as a demo application. Also included are three preconfigured use cases for the demo application: automated quality gates, runbook automation, and automated evaluation of blue/green deployments.

## Usage

Here is the best way to getting started with keptn:
- If you want to try out the latest stable release with your own services and application, please head over to the [release section](https://github.com/keptn/keptn/releases) of keptn and follow the official [documentation of keptn](https://keptn.sh/docs). We recommend to work with this version.
- If you want to try out the latest version of keptn with your own services and application, please use the [master branch](https://github.com/keptn/keptn/tree/master) and follow the [documentation](https://keptn.sh/docs/) on the https://keptn.sh website. 
- If you want to work with the latest version of keptn that is currently under development, please use the [develop branch](https://github.com/keptn/keptn/tree/develop). (:warning: this is the development branch, so it might not be stable all the time)
- Please use the [docs on the keptn website](https://keptn.sh/docs) to get resources on how to use keptn.
- Please use the [release section](https://github.com/keptn/keptn/releases) to learn about our current releases, release candidates and pre-releases to get the latest version of keptn.

## Versions compatibilities
We mangage the keptn core components as well as all services (e.g., jenkins-service, github-service) in versions. The respective images in their versions are stored in [DockerHub](https://hub.docker.com/u/keptn).
The versions of the keptn core components and the services have to be compatible to each other.
Therefore, this section shows the compatibility between these versions.

keptn in [version 0.2.2](https://github.com/keptn/keptn/releases/tag/0.2.2) requires:
- keptn/authenticator:0.2.2
- keptn/bridge:0.1.0
- keptn/control:0.2.2
- keptn/eventbroker:0.2.2
- keptn/eventbroker-ext:0.2.2
- keptn/pitometer-service:0.1.2 
- keptn/servicenow-service:0.1.1
- keptn/github-service:0.2.0
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

