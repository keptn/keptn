# Release Notes 0.6.1

This release introduces the new and shiny :sparkles: **Keptn's Bridge** :sparkles:

* This new Keptn's Bridge is implemented on an Angular and NodeJS stack. An excerpt of its feature set shows: 
  * Auto-reload Keptn events 
  * Display labels attached to Keptn events
  * Breakdown of an SLO evaluation into SLIs
  * Comparison of evaluation results across multiple evaluations
  * Eye-catching visuals to highlight noticeable elements
  * Deep-link to a Keptn project
  * Link to Keptn-deployed services

* Next to the Keptn's Bridge, this release provides documentation on how to run Keptn on a **Microkube 1.2**. This allows deploying Keptn on a single-node Kubernetes cluster inside a Virtual Machine (VM) that is hosted locally. 

* Last but not least, this release addresses known limitations and issues in using Keptn Quality Gates. Also many thanks for the [Keptn Enhancement Proposals](https://github.com/keptn/enhancement-proposals) that were submitted and implemented!

## Keptn Specification

Implemented **Keptn spec** version: [0.1.3](https://github.com/keptn/spec/tree/0.1.3)

## New Features

<details><summary>Quality Gates</summary>
<p>

- Return an event with `result=failure` when no SLI-provider is available, but an SLO file is found [#1212](https://github.com/keptn/keptn/issues/1212)
- Consider the test result of functional tests to determine the result of the evaluation [#1380](https://github.com/keptn/keptn/issues/1380)
- Configure SLI provider when `keptn configure monitoring` is executed [#1341](https://github.com/keptn/keptn/issues/1341)
- Retrieve SLIs even if tests fail [#1289](https://github.com/keptn/keptn/issues/1289)

</p>
</details>

<details><summary>Platform Support / Installer</summary>
<p>

- Restarting the NATS cluster requires distributors to reconnect [#1209](https://github.com/keptn/keptn/issues/1209)
- OpenShift support for Keptn quality gates standalone [#1197](https://github.com/keptn/keptn/issues/1197)
- Create *DestinationRule* for exposed services when using Istio [#1408](https://github.com/keptn/keptn/issues/1408)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *jmeter-service:*
  - [KEP 0005](https://github.com/keptn/enhancement-proposals/blob/master/text/0005-deployment-finished-with-deployed-endpointurl.md) - Use deploymentURILocal/deploymentURIPublic instead of guessing the service URL [#1403](https://github.com/keptn/keptn/issues/1403)
- *helm-service:* 
  - Create namespaces on demand and not by default [#1417](https://github.com/keptn/keptn/issues/1417)
  - [KEP 0005](https://github.com/keptn/enhancement-proposals/blob/master/text/0005-deployment-finished-with-deployed-endpointurl.md) - Send deploymentURIPublic and deploymentURILocal after successful deployment [#1417](https://github.com/keptn/keptn/issues/1417)
  - Only wait for deployments contained in Helm release [#1225](https://github.com/keptn/keptn/issues/1225)
- *configuration-service:* 
  - Allow reading files without using git pull [#1396](https://github.com/keptn/keptn/issues/1396)
  - Improve mutex per project [#1395](https://github.com/keptn/keptn/issues/1395)
  - Enhance API call GET `/projects` to return more information [#1394](https://github.com/keptn/keptn/issues/1394)
- API: Allow to set `keptnContext` in events sent to API endpoint [#1355](https://github.com/keptn/keptn/issues/1355)

</p>
</details>

<details><summary>CLI Enhancements</summary>
<p>

- Check if a new CLI version is available once a day [#1190](https://github.com/keptn/keptn/issues/1190)
- Check Keptn and Kubernetes compatibility based on K8s cluster version before installing [#1326](https://github.com/keptn/keptn/issues/1326)

</p>
</details>

<details><summary>Integrations / Keptn contrib</summary>
<p>

- *dynatrace-SLI-service:*
  - Ensure compatibility with new metrics v2 api (/api/v2/metrics/query) [#1282](https://github.com/keptn/keptn/issues/1282)
- *dynatrace-service:*
  - Automatically create custom alert rules [#1265](https://github.com/keptn/keptn/issues/1265)
  - Support of new API for Dashboards [#1358](https://github.com/keptn/keptn/issues/1358)
- *servicenow-service:*
  - servicenow-service supports Keptn 0.6.x in [0.2.0](https://github.com/keptn-contrib/servicenow-service/releases/tag/0.2.0)

</p>
</details>

## Fixed Issues
- Check K8s cluster version only for full installation [#1398](https://github.com/keptn/keptn/issues/1398)
- Fixed memory leak in *mongodb-datastore* [#1440](https://github.com/keptn/keptn/issues/1440)
- Use right mount path for persistant volume in *mongodb-datastore* [#1360](https://github.com/keptn/keptn/issues/1360)

## Fixed Limitations
- Distributors can not automatically reconnect to NATS cluster [#1209](https://github.com/keptn/keptn/issues/1209)
- Dynatrace SLI Service / Metrics API will change at the end of Q1 [#1282](https://github.com/keptn/keptn/issues/1282)
- The Quality-gates standalone version is currently not supported on OpenShift [#1197](https://github.com/keptn/keptn/issues/1197)

## Good to know / known limitations
- Installation on AKS with K8s version before 1.15.5 might fail [#1429](https://github.com/keptn/keptn/issues/1429)
- For old limitations, please see [Release 0.6.0](https://github.com/keptn/keptn/releases/tag/0.6.0). 
