# Release Notes 0.2.2

This release improves the usability and stability of keptn. More precisely, it significantly improves the experience of installing keptn, executing commands using the keptn CLI, and deploying new applications.

## New Features

- keptn CLI: New command for installing keptn #272
- keptn CLI: Receives constant status updates during the execution of CLI commands #203
- keptn CLI: New commands for sending new artifact events and arbitrary keptn events #326

- The name of the namespace has the project name as prefix, e.g., sockshop-dev, sockshop-production #229
- Service configuration files follow the naming scheme: `servicename-servicetype.yaml` #295

- Hard-coded URLs for Dynatrace SaaS tenants are removed to support Dynatrace-managed tenants #255
- Deploys dynatrace-service when Dynatrace monitoring is activated #268 #354
- The *deploy* pipeline does not send the Dynatrace deployment event since this is handled by the dynatrace-service #268

## Fixed Issues

- Fix no healthy upstream-problem during blue green deployment #332
- Bugfix in roll-back step of *evalution_done* and *run_tests* pipelines #292
- Set visibility of service to cluster-local #396
- Correct typo in eventbroker and control service #268
- Log reason for failed service onboarding #274
- Mark evaluation as failed when no data can be retrieved from Dynatrace #380

## Known Limitations

- Installation is only supported for GKE (more platforms to come)
- Only one GitHub organization can be configured with the keptn server (will be addressed in [#210](https://github.com/keptn/keptn/issues/210))
