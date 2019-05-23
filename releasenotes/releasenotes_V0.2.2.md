# Release Notes 0.2.2

This release is a stability improvement release. It does not add any new use cases, but significantly improves the installation experience of keptn.

## New Features

- The deployDynatrace.sh script installs the dynatrace-service [#268](https://github.com/keptn/keptn/issues/268)
- CLI: Introduce a new command for installing keptn [#272](https://github.com/keptn/keptn/issues/272)
- CLI: Introduce new commands for sending new artifact events as well as arbitrary keptn events [#326](https://github.com/keptn/keptn/issues/326)

## Fixed Issues
- Correct typo in eventbroker and control service [#268](https://github.com/keptn/keptn/issues/324) 

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
