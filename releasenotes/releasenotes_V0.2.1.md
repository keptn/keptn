# Release Notes 0.2.1

This release is a stability improvement release. It does not add any new use cases, but significantly improves the installation experience of keptn.

## New Features

- Improved installation script that verifies each steps of the installation to provide accurate feedback of the installation procedure #250 #285

## Fixed Issues

- Invalid characters in project and service names now detected by keptn CLI #292 
- Error handling in core services #287

## Version dependencies:

keptn is installed by using these images from the [keptn Dockerhub registry](https://hub.docker.com/u/keptn):

- keptn/keptn-authenticator:0.2.1
- keptn/keptn-control:0.2.1
- keptn/keptn-event-broker:0.2.1
- keptn/keptn-event-broker-ext:0.2.1
- keptn/pitometer-service:0.1.1 (fixed issue #291)
- keptn/servicenow-service:0.1.0
- keptn/github-service:0.1.1 (fixed issue #265)
- keptn/jenkins-service:0.2.0 (fixed issues #250 #268, known limitations #292)
  - keptn/jenkins-0.5.0

## Known Limitations

- Installation currently only on GKE (more platforms to come)
- Only one GitHub organization can be configured with the keptn server (will be adressed in #210)
- For use cases that require Dynatrace: support for Dynatrace SaaS tenants only (will be adressed in #255)
- keptn CLI output not reliably reflecting success/error of keptn services: the CLI only reflects the successful acknowledgment of the CLI command but not its successful execution (will be adressed in #203 #143)
