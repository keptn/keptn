# Release Notes 0.5.0.beta

## New Features
- github.com/keptn is now a mono-repo meaning that it contains all projects relevant for a Keptn deployment [#537](https://github.com/keptn/keptn/issues/537)
- Use new logging format [#535](https://github.com/keptn/keptn/issues/535)
- Remove GitHub email from install process [#555](https://github.com/keptn/keptn/issues/555)
- Delete unused scripts in installer [#614](https://github.com/keptn/keptn/issues/614)
- Use custom dialer for xip io resolving in the websocket communication [#634](https://github.com/keptn/keptn/issues/634)
- Use internal datastore for logging events and service logs [#536](https://github.com/keptn/keptn/issues/536)
- Updated bridge to use internal datastore to show events [#621](https://github.com/keptn/keptn/issues/621)
- Allow to use complete Helm charts when onboarding a new service [#611](https://github.com/keptn/keptn/issues/611)

## CLI Enhancements
- Add EKS support [#608](https://github.com/keptn/keptn/issues/608)
- Provide uninstall command in CLI [#562](https://github.com/keptn/keptn/issues/562)
- After a successful Keptn installation, installer job will be deleted [#663](https://github.com/keptn/keptn/issues/663)
- Add option `--insecure-skip-tls-verify` in installer [#567](https://github.com/keptn/keptn/issues/567)
- Allow update domain from CLI [#553](https://github.com/keptn/keptn/issues/553)
- Refactor installer in keptn CLI [#638](https://github.com/keptn/keptn/issues/638)
- Allow to upload resources from the keptn CLI [#673](https://github.com/keptn/keptn/issues/673)

## New Services
- New **api** for communicating with keptn [#506](https://github.com/keptn/keptn/issues/506)
- New **configuration-service** for managing resources for Keptn project-related entities, i.e., project, stage, and service [#451](https://github.com/keptn/keptn/issues/451)
- New **shipyard-service** to process a shipyard file to create a project and stages [#610](https://github.com/keptn/keptn/issues/610)
- New **wait-service** to wait a certain time before sending an event [#725](https://github.com/keptn/keptn/issues/725)

## Fixed Issues
- Update domain does not create duplicate entries in configmap [#570](https://github.com/keptn/keptn/issues/570)
- Fix xip.io resolving in websocket communication [#634](https://github.com/keptn/keptn/issues/634)

## Known Limitations
