[![Build Status](https://travis-ci.org/keptn/keptn.svg?branch=master)](https://travis-ci.org/keptn/keptn)

![keptn](./assets/keptn.png)

# keptn
keptn is a fabric for cloud-native lifecycle automation at enterprise scale. In its first version it provides an automated setup of the keptn core components as well as a demo application. Also included are three preconfigured use cases for the demo application: automated quality gates, runbook automation, and automated evaluation of blue/green deployments.

## Usage

Here is the best way to getting started with keptn:
- If you want to try out keptn with a predefined application and work through the main use cases that keptn provides, please use the latest 0.1.x release, currently [0.1.3](https://github.com/keptn/keptn/tree/0.1.3). (:warning: 0.1.x version will be obsolete as soon as version 0.2 has been released).
- If you want to try out the latest version of keptn with your own services and application, please use the [0.2.x prerelease branch](https://github.com/keptn/keptn/tree/prerelease-0.2.x) and follow the [pre-release documentation](https://keptn.sh/docs/0.2.0/) on the https://keptn.sh website. We recommend to work with this version.
- If you want to work with the latest version of keptn that is currently under development, please use the master branch. (:warning: this is the development branch, so it might not be stable all the time)
- Please use the [docs on the keptn website](https://keptn.sh/docs) to get resources on how to use keptn.
- Please use the [release section](https://github.com/keptn/keptn/releases) to learn about our current releases, release candidates and pre-releases to get the latest version of keptn.

## Repositories <a id="repos"></a>
* [keptn/keptn](README.md). This is the main repository that you are currently looking at. It hosts keptn's core components and documents that govern the keptn open source project. It includes:
    * [designDocs](./designDocs/): We're already designing the architecture and use cases for the next releases - you can find the current design docs [here](./designDocs). Please feel free to review and comment the designs, after all, we're encouraging all of you to collaborate on keptn.
    * [install](./install/): This directory contains all artifacts that are required to install keptn. Please note that the latest docs can be found in the [docs section of the keptn website](https://keptn.sh/docs).
    * [releasenotes](./releasenotes/): You can find the current release notes in this directory.

## Versions compatibilities
We mangage the keptn core components as well as all services (e.g. Jenkins-service, GitHub-service) in versions. The respective images in their versions are stored in [DockerHub](https://hub.docker.com/?namespace=keptn).
The versions of the keptn core components and the services have to be compatible to each other.
Therefore, this section shows the compatibility between these versions.

keptn in [version 0.2](https://github.com/keptn/keptn/tree/prerelease-0.2.x) requires:
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

