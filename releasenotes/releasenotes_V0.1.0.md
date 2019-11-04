# keptn Release Notes V 0.1

## Release Goal

The goal of the V 0.1 release for keptn is to provide an MVP (minimum viable product) for an autonomous cloud fabric supporting
major use cases based on a demo application. 

## Automated Setup

This release of keptn provides an automated setup for the keptn core components as well as a demo application. The setup only requires a Kubernetes cluster with minimum master version 1.10.11. Detail specs can be found in the [getting started section](../GettingStarted.md).

If you want to use keptn with your own application, you will have to modify the keptn setup files. Future versions will provide a smoother onboarding experience. 

## Automated Multistage Pipeline

Keptn supports a three-stage continuous delivery pipeline with the following stages:

* Development - for integration testing
* Staging/Continuous Performance - for production-like performance testing
* Production - production deployment with automated load generation

## Automated Quality Gates

Keptn provides a quality gate from the development to the staging and from the staging to the production environment. While the gate in development is a basic check regarding the availability of the service, the quality gate in staging assumes an execution of a performance test that gets validated against a performance signature. This signature defines the thresholds to mark a deployment as unstable and to stop the delivery pipeline.

The use case provided in this release is as follows:

1. The source code of a service is changed, and the service gets deployed to the development environment. 
1. The service passes the quality gates in the development environment.
1. However, the service does not pass the quality gate in staging due to an increase of the response time detected by a performance test.

## Automated Runbook Automation with Ansible

Keptn provides runbook automation as an auto-remediation approach in response to detected issues in a production environment. Therefore, keptn automatically sets up and configures an Ansible Tower instance during setup. The example ships with predefined playbooks that are capable of updating the configuration of a service in production, defining configuration change events, and reacting on them in case of detected issues. 

The use case provided in this release is as follows:

1. A configuration change is applied to a service in the production environment, leading to an increase of the failure rate.
1. An issue is detected in production and a problem ticket is opened.
1. The configuration change is automatically reverted.

## Production Deployments with Canary Releases

Keptn provides the runbooks to release a new version and to automatically switch back to the previous version if an issue is detected. As described above, keptn relies on Ansible Tower for auto-remediation capabilities. Thus, keptn is shipped with pre-defined playbooks that can deploy a new version and take care of re-routing traffic in case of detected problems.

The use case provided in this release is as follows:

1. A faulty service version is deployed to the production environment and traffic is routed to this new version in a canary release manner, starting with only 10 % of the traffic and increasing this percentage over time.
1. An issue is detected in production and a problem ticket is opened.
1. The traffic routing is changed to redirect traffic to the previous (non-faulty) version.

