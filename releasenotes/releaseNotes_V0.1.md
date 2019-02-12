# keptn Release Notes V 0.1

## Release Goal

The goal of the V 0.1 release for keptn is to provide an MVP (minimum viable product) for an autonomous cloud fabric supporting
major use cases based on a demo application. 

## Automated Setup

This release of keptn provides an automated setup for the keptn core components as well  as a demo applications. The setup only requires a Kubernetes cluster *TODO* add version and specs. 

If you want to onboard your own application you will have to modify the keptn setup files. Future versions will provide a smoother onboarding experience. 

## Automated Multistage Pipeline

Keptn supports a three stage continuous delivery pipeline with the following stages:

* Development - for integration testing
* Staging/Continuous Performance - for production-like performance testing
* Production - production deployment with automated load generation

## Automated Quality Gates

## Automatic Rollback of Faulty Blue/Green Deployments

## Automated Runbook Automation with Ansible

Keptn provides runbook automation as an auto-remediation approach in response to detected issues in a production environment. Therefore, Keptn automatically sets up and configures an Ansible Tower instance during setup. The example ships with predefined playbooks that are capable of updating the configuration of a service in production, defining configuration change events, and react on them in case of issues are detected. 

The use case provided in this release is as follows.

1. A faulty configuration change is applied to a service in the production environment
1. An issue is detected in production and a problem ticket is opened
1. The faulty configuration change is automatically reverted.
