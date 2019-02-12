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

## Automatic rollback of faulty blue/green deployments. 

## Automated runbook automation with Ansible

