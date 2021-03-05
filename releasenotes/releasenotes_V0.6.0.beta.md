# Release Notes 0.6.0.beta

## New Features
- Keptn support for PKS (Pivotal/VMWare) [#5](https://github.com/keptn/keptn/issues/5)
- Keptn support for Rancher [#462](https://github.com/keptn/keptn/issues/462)
- Keptn installation using a NodePort for Istio ingressgateway [#462](https://github.com/keptn/keptn/issues/462)
- Enable Istio injection for namespaces [#715](https://github.com/keptn/keptn/issues/715)
- Keptn API can be reached without Istio [#1073](https://github.com/keptn/keptn/issues/1073)
- REST API endpoints for: project and service [#893](https://github.com/keptn/keptn/issues/893)
- REST API support for start-evaluation and evaluation-done [#949](https://github.com/keptn/keptn/issues/949)
- Allow filtering by project, stage, and service in keptn-datastore [#980](https://github.com/keptn/keptn/issues/980)
- *helm-service*: Allow changing any file in Helm Chart [#995](https://github.com/keptn/keptn/issues/995)
- *remediation-service*: Remediation action for slowing down requests[#1006](https://github.com/keptn/keptn/issues/1006)
- *lighthouse-service*: Default behaviour if no slo.yml found -> evaluation pass [#1081](https://github.com/keptn/keptn/issues/1081)
- *mongodb-datastore*: Allow event-retrieval filtering by source [#1061](https://github.com/keptn/keptn/issues/1061)
- *bridge*: Improve bridge for lighthouse and SLI events [#1058](https://github.com/keptn/keptn/issues/1058)
- *jmeter/wait-service*: Send test-finished event with start/end time rather than startedat [#1078](https://github.com/keptn/keptn/issues/1078)
- Double check specification and implementation of events [#997](https://github.com/keptn/keptn/issues/997)
- Provide Helm helper functions for Keptn's project, stage, and service [#1109](https://github.com/keptn/keptn/issues/1109)
- Adapt DT integration to create rules based on DT_CUSTOM_PROP [#1110](https://github.com/keptn/keptn/issues/1110)
- Slim down image size of Keptn installer image [#1034](https://github.com/keptn/keptn/issues/1034)

## CLI Enhancements
- Only check image availability on Docker if image contains docker.io [#991](https://github.com/keptn/keptn/issues/991)
- Keptn uninstall commands shows that not everything is uninstalled [#971](https://github.com/keptn/keptn/issues/971)
- Support for start-evaluation and get evaluation-done [#948](https://github.com/keptn/keptn/issues/948)
- Quality gates standalone via CLI needs to set test-strategy [#1069](https://github.com/keptn/keptn/issues/1069)

## New Services
- Lighthouse service [#950](https://github.com/keptn/keptn/issues/950)
- Removed pitometer-service, install lighthouse-service [#1057](https://github.com/keptn/keptn/issues/1057)

## Fixed Issues
- Fixed issues with git clone and secrets [#929](https://github.com/keptn/keptn/issues/929)
- Fixed version in Swagger for Keptn API [#983](https://github.com/keptn/keptn/issues/983)
- Check valid project and stage names [#745](https://github.com/keptn/keptn/issues/745)
- Fixed error message when Keptn installation [#932](https://github.com/keptn/keptn/issues/932)
- Centralize ResolveXipIoWithContext to go-utils for all API calls [#1005](https://github.com/keptn/keptn/issues/1005)
- Fixed issues with Keptn events that are transferred into wrong format when stored by MongoDB [#1021](https://github.com/keptn/keptn/issues/1021)
- Fixed issue for paging the resource list [#1048](https://github.com/keptn/keptn/issues/1048)
- Fixed issue when creating a project with a git upstream repo with . (dots) [#1095](https://github.com/keptn/keptn/issues/1095)
- Forward testStrategy and deployment strategy in evaluation-events[#1098](https://github.com/keptn/keptn/issues/1098)

## Development process
- Travis CI Pipeline runs for all pull requests for external contributors[#957](https://github.com/keptn/keptn/issues/957)
- Changed go-utils from go-dep to go-modules [#959](https://github.com/keptn/keptn/issues/959)
- Provide nightly builds for the CLI and containers [#946](https://github.com/keptn/keptn/issues/946)
- Refactored mongodb-datastore to use go-utils develop version [#1011](https://github.com/keptn/keptn/issues/1011)
- Update of *Contribution Guide* and *Getting started* [#943](https://github.com/keptn/keptn/issues/943)
- Added *Developer Documentation* [#938](https://github.com/keptn/keptn/issues/938)
- Unified Dockerfiles and add Skaffold YAMLs for easy in-cluster debugging [#1026](https://github.com/keptn/keptn/issues/1026)

## Known Limitations
