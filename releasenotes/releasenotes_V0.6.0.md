# Release Notes 0.6.0 SLI (SLI Sta'lone)

This release focuses on improvements of the **Keptn Quality Gates** capability that allows the evaluation of **Service Level Objectives** (SLOs), which are determined by **Service Level Indicators** (SLI). For these improvements, the **Lighthouse** service has been introduced, which is a service that is responsible for conducting an evaluation based on data from SLI-providers. The Lighthouse service supersedes Pitometer.
Furthermore, the following highlights should be noticed:

- Because of popular demand, Keptn Quality Gates is now available as a standalone feature. That allows enriching existing CD pipelines with quality gates that ensure that services are only promoted if they meet defined SLOs. To integrate Keptn in an existing pipeline, REST endpoints are provided to trigger an evaluation and to pull the evaluation results.

- The Keptn Quality Gates capability doesn't require deployment and testing features. We've provided an installation option in the Keptn installer that excludes the deployment and testing components, which results in a smaller resource footprint.

- The Keptn's bridge has been enhanced in order to get better visibility into evaluation results. Thus, it is possible to see why an SLO has not been met.

- Due to the flexible architecture of Keptn, it is possible to exchange an SLI-provider to gather data for the quality gate from another monitoring or testing provider.

## Keptn Specification

Implemented [Keptn spec](https://github.com/keptn/spec/tree/0.1.2) version: 0.1.2

## New Features

<details><summary>Quality Gates</summary>
<p>
- REST API support for start-evaluation and evaluation-done [#949](https://github.com/keptn/keptn/issues/949)
- *ligthouse-service*: Default behaviour if no slo.yml found -> evaluation pass [#1081](https://github.com/keptn/keptn/issues/1081)
- *bridge*: Improve bridge for lighthouse and SLI events [#1058](https://github.com/keptn/keptn/issues/1058)
- Lighthouse-service forwards custom data to SLI providers [#1147](https://github.com/keptn/keptn/issues/1147)
- Lighthouse-service forwards environment variable keptn_deployment to SLI providers [#1161](https://github.com/keptn/keptn/issues/1161)
- Service Level Indicators (SLI) are stored in Keptn's git repository [#1192](https://github.com/keptn/keptn/issues/1192)
- Implemented SLI Provider for [Dynatrace](https://github.com/keptn-contrib/dynatrace-sli-service) and [Prometheus](https://github.com/keptn-contrib/prometheus-sli-service)
- Forward testStrategy and deployment strategy in evaluation-events[#1098](https://github.com/keptn/keptn/issues/1098)
- *lighthouse-service*: Prevent previous failed SLI results to be used for comparison-based evaluation [#1263](https://github.com/keptn/keptn/issues/1263)
- *lighthouse-service*: can deal with an SLO file that has no SLI criteria defined [#1213](https://github.com/keptn/keptn/issues/1213)
</p>
</details>


<details><summary>Platform Support / Installer (Uniform)</summary>
<p>
- Keptn support for PKS (Pivotal/VMWare) [#5](https://github.com/keptn/keptn/issues/5)
- Keptn support for Rancher [#462](https://github.com/keptn/keptn/issues/462)
- Keptn installation using a NodePort for Istio ingressgateway [#462](https://github.com/keptn/keptn/issues/462)
- Enable Istio injection for namespaces [#715](https://github.com/keptn/keptn/issues/715)
- Keptn API can be reached without Istio [#1073](https://github.com/keptn/keptn/issues/1073)
- Nginx support for EKS [#1124](https://github.com/keptn/keptn/issues/1124)
- Allow to configure a custom domain when using Nginx [#1167](https://github.com/keptn/keptn/issues/1167)
- Slim down image size of Keptn installer image [#1034](https://github.com/keptn/keptn/issues/1034)
- Fluent-bit removed from installation[#1172](https://github.com/keptn/keptn/issues/1172)
- Keptn Install Quality Gates: Only install relevant services [#1130](https://github.com/keptn/keptn/issues/1130)
- Provide Helm helper functions for Keptn's project, stage, and service [#1109](https://github.com/keptn/keptn/issues/1109)
- Single Istio gateway accepting HTTPS traffic [#1231](https://github.com/keptn/keptn/issues/1231)
- Check if Istio is already installed before installing it [#1208](https://github.com/keptn/keptn/issues/1208)
- Servicenow-service and prometheus-service have been removed from Keptn's uniform [#1302](https://github.com/keptn/keptn/issues/1302)
- Retrieving hostname/ip of ingress improved [#1199](https://github.com/keptn/keptn/issues/1199)

</p>
</details>


<details><summary>Control-plane enhancements/changes</summary>
<p>
- *helm-service*: Allow changing any file in Helm Chart [#995](https://github.com/keptn/keptn/issues/995)
- *remediation-service*: Remediation action for slowing down requests[#1006](https://github.com/keptn/keptn/issues/1006)
- Self-healing use-case with Dynatrace [#1185](https://github.com/keptn/keptn/issues/1185)
- Self-healing use-case based on feature toggle remediation with Unleash [#1104](https://github.com/keptn/keptn/issues/1104)
- *jmeter/wait-service*: Send test-finished event with start/end time rather than startedat [#1078](https://github.com/keptn/keptn/issues/1078)
- REST API endpoints for: project and service [#893](https://github.com/keptn/keptn/issues/893)
- Allow filtering by project, stage, and service in keptn-datastore [#980](https://github.com/keptn/keptn/issues/980)
- *mongodb-datastore*: Allow event-retrieval filtering by source [#1061](https://github.com/keptn/keptn/issues/1061)
- Adapt DT integration to create rules based on DT_CUSTOM_PROP [#1110](https://github.com/keptn/keptn/issues/1110)
- *jmeter-service*: Sends a test-finished with status `fail` if test execution failed [#542](https://github.com/keptn/keptn/issues/542)
</p>
</details>


<details><summary>CLI Enhancements</summary>
<p>
- Only check image availability on Docker if image contains docker.io [#991](https://github.com/keptn/keptn/issues/991)
- Keptn uninstall commands shows that not everything is uninstalled [#971](https://github.com/keptn/keptn/issues/971)
- Support for start-evaluation and get evaluation-done [#948](https://github.com/keptn/keptn/issues/948)
- Quality gates standalone via CLI needs to set test-strategy [#1069](https://github.com/keptn/keptn/issues/1069)
- Avoid installing incompatible Keptn version [#1162](https://github.com/keptn/keptn/issues/1162)
- Allow start and end datetime (instead of timeframe) for start-evaluation [#1131](https://github.com/keptn/keptn/issues/1131)
- Setup Dynatrace monitoring using the Keptn CLI `configure monitoring` command [#443](https://github.com/keptn/keptn/issues/443)
- Configure domain assumes dedicated Keptn gateway [#1295](https://github.com/keptn/keptn/issues/1295)
</p>
</details>


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
- Reset NodePorts of services in generated Helm chart [#1181](https://github.com/keptn/keptn/issues/1181)
- Keptn API: Swagger UI refers to https instead of http [#1012](https://github.com/keptn/keptn/issues/1012)
- Allow project, stage, and service names with hyphens [#1166](https://github.com/keptn/keptn/issues/1166)
- Check whether Keptn-service name matches Kubernetes-service name [#1261](https://github.com/keptn/keptn/issues/1261)
- api.keptn.DOMAIN/v1/events endpoints requires query parameters [#1210](https://github.com/keptn/keptn/issues/1210)
- *mongodb-datastore*: eventContext removed from events that are read from mongodb-datastore [#1300](https://github.com/keptn/keptn/issues/1300)


## Development process
- Travis CI Pipeline runs for all pull requests for external contributors[#957](https://github.com/keptn/keptn/issues/957)
- Changed go-utils from go-dep to go-modules [#959](https://github.com/keptn/keptn/issues/959)
- Provide nightly builds for the CLI and containers [#946](https://github.com/keptn/keptn/issues/946)
- Refactored mongodb-datastore to use go-utils develop version [#1011](https://github.com/keptn/keptn/issues/1011)
- Update of *Contribution Guide* and *Getting started* [#943](https://github.com/keptn/keptn/issues/943)
- Added *Developer Documentation* [#938](https://github.com/keptn/keptn/issues/938)
- Unified Dockerfiles and add Skaffold YAMLs for easy in-cluster debugging [#1026](https://github.com/keptn/keptn/issues/1026)


## Known Limitations

- There is no new-artifact cloud event - `keptn send event new-artifact` causes a configuration change event [#1218](https://github.com/keptn/keptn/issues/1218)
- Currently, Helm v2 is used, which is now entering maintenance mode as Helm v3 has been released. This means that Helm v2 is going to be unmaintained at the end of 2020 (Source: https://helm.sh/blog/2019-10-22-helm-2150-released/#helm-2-support-plan ).
- Currently, CloudEvents v0.2 is used. We aim to implement v1.0 in the future which will cause changes in our CloudEvents handling. [#1178](https://github.com/keptn/keptn/issues/1178)
- Using an Apache Proxy with `ProxyPreserveHost on` in blue-green deployment scenario fails because of Istio Envoy/Proxy
- Distributors can not automatically reconnect to NATS cluster [#1209](https://github.com/keptn/keptn/issues/1209) 
- Dynatrace SLI Service / Metrics API will change at the end of Q1 [#1282](https://github.com/keptn/keptn/issues/1282) 
- Sending multiple new-artifacts events can cause problems which result in a "no healthy upstream" error [#1229](https://github.com/keptn/keptn/issues/1229), [#1219](https://github.com/keptn/keptn/issues/1219)
- Istio can potentially cause the kube-apiserver to experience downtimes (see https://github.com/istio/istio/issues/19481 ) [#1298](https://github.com/keptn/keptn/issues/1298)
- The Quality-gates standalone version is currently not supported on OpenShift and PKS [#1197](https://github.com/keptn/keptn/issues/1197)
