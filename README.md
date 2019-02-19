[![Build Status](https://travis-ci.org/keptn/keptn.svg?branch=master)](https://travis-ci.org/keptn/keptn)

![keptn](./assets/keptn.png)

# keptn
keptn is a fabric for cloud-native lifecycle automation at enterprise scale. In its first version it provides an automated setup of the keptn core components as well as a demo application. Also included are three preconfigured use cases for the demo application: automated quality gates, runbook automation, and automated evaluation of blue/green deployments.

<!-- 
##### Table of Contents
 * [Introduction](#intro)
 * [Repositories](#repos)
 * [To start using keptn](#using-keptn)
 * [To start developing keptn](#developing-keptn)

## Introduction <a id="intro"></a>

In a nutshell, keptn provides following advantages:
* **Best practices out-of-the box:** Keptn supports best practices and the latest delivery platforms utilized by modern applications while ensuring that you can easily get started in minutes.
* **Future-proof and pluggable:** One-off implementations often lead to eventual maintenance problems. Keptn provides a core framework that you can use to build cohesive, standardized cloud native-fabric for your organization. Keptn also enables you to replace tools as you wish and avoid vendor lock-in.
* **Smart and flexible:** Concepts like GitOps, self-healing, and unbreakable deployments can be implemented in different ways, varying from one organization to the next. While keptn has built-in intelligence that assists you in taking advantage of these industry best practices, it’s flexible enough to address your specific needs.
-->

## Repositories <a id="repos"></a>
* [keptn/keptn](README.md). This is the main repository that you are currently looking at. It hosts keptn's core components and documents that govern the keptn open source project. It includes:

    * [cloudevents](./cloudevents/): Events are everywhere. However, event producers tend to describe events differently. To provide a definition of all events keptn understands, this directory maintains a list of events that follow the [CloudEvent specification](https://github.com/cloudevents/spec). 
    * [designDocs](./designDocs/): We're already designing the architecture and use cases for the next releases - you can find the current design docs [here](./designDocs). Please feel free to review and comment the designs, after all, we're encouraging all of you to collaborate on keptn.
    * [install](./install/): This directory contains all artifacts that are required to install keptn.
    * [releasenotes](./releasenotes/): You can find the current release notes in this directory.
    * [onboard](./onboard/): This directory contains the srcipts to onboard an application (or single services) so that keptn takes care of it.

* [keptn/examples](https://github.com/keptn/examples). This repository contains examples to explore keptn and to learn more about the cloud-native lifecycle automation based on different use cases.

<!-- 
## To start using keptn <a id="using-keptn"></a>

## To start developing keptn <a id="developing-keptn"></a>
-->
