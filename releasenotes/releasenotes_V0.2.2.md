# Release Notes 0.2.2

This release improves the usability and stability of keptn. More precisely, it significantly improves the experience of installing keptn, executing commands in the CLI, as well as deploying new applictions.

## New Features

- CLI: Introduce a new command for installing keptn [#272](https://github.com/keptn/keptn/issues/272)
- CLI: Receive constant status updates during the execution of a CLI command [#203](https://github.com/keptn/keptn/issues/203)
- CLI: Introduce new commands for sending new artifact events as well as arbitrary keptn events [#326](https://github.com/keptn/keptn/issues/326)

- The name of the namespace has the project name as prefix, e.g., sockshop-dev, sockshop-production [#229](https://github.com/keptn/keptn/issues/229)
- Service configuration files follow the naming scheme: `servicename-servicetype.yaml` [#295](https://github.com/keptn/keptn/issues/295)

- Hard-coded URLs for Dynatrace SaaS tenants are removed to also support Dynatrace-managed tenants [#255](https://github.com/keptn/keptn/issues/255)
- The deployDynatrace.sh script installs the dynatrace-service [#268](https://github.com/keptn/keptn/issues/268)
- The *deploy* pipeline does not send the Dynatrace deployment event since this is handled by the dynatrace-service now [#268](https://github.com/keptn/keptn/issues/268)


## Fixed Issues
- Fix no healthy upstream-problem during b/g deployment [#332](https://github.com/keptn/keptn/issues/332)
- Bugfix in roll-back step of *evalution_done* and *run_tests* pipelines [#292](https://github.com/keptn/keptn/issues/292)
- Set visibility of service to cluster-local [#396](https://github.com/keptn/keptn/issues/255)
- Correct typo in eventbroker and control service [#268](https://github.com/keptn/keptn/issues/324) 
- Log reason for failed service onboarding [#274](https://github.com/keptn/keptn/issues/274)
- Mark evaluation as failed when no data can be retrieved from Dynatrace [#380](https://github.com/keptn/keptn/issues/380)

## Known Limitations

- Installation is only supported for GKE (more platforms to come)
- Only one GitHub organization can be configured with the keptn server (will be adressed in [#210](https://github.com/keptn/keptn/issues/210))
