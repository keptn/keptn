# Release Notes 0.6.2

This release provides improvements in working with the Keptn Bridge and API as well as troubleshooting support by creating a support-archive through the Keptn CLI. 

:rocket: *Improved automation support with API extensions and deep linking into Keptn Bridge:* 
- *Keptn API:* Internally used endpoints for retrieving Keptn status information are now accessible through the public API, e.g.:
  * GET `/v1/project` - Returns all projects managed by Keptn.
  * GET `/v1/project/{projectName}` - Returns meta-information about a project.
  * GET `/v1/project/{projectName}/stage` - Returns stages from a project.
  * GET `/v1/project/{projectName}/stage/{stageName}/service` - Returns services from a stage.

- *Keptn Bridge:* To provide a convenient and secure way of working with the Keptn Bridge, the Keptn CLI command: `keptn configure bridge --action=expose` has been introduced. This command allows exposure of the Bridge via Istio or Nginx ingress. In addition, basic authentication with username and password can be activated. For improved automation support, deep links into the Keptn Bridge are provided that point to certain UI components.

:squid: *Argo CD for deployment:* With this release, Keptn can be used in combination with Argo CD / Argo Rollout as explained by a [tutorial](https://tutorials.keptn.sh/tutorials/keptn-argo-cd-deployment/#0). While Argo CD is used for deploying an *Argo Rollout*, Keptn is leveraged for testing, evaluating, and promoting this rollout.

:star2: *Easier bug reporting with support-archives:* The CLI now offers the command: `keptn generate support-archive` that fetches all log files from a Keptn deployment and puts them into an archive. This archive can then be used for troubleshooting without connection to the Kubernetes cluster.

Last but not least, this release addresses limitations and issues in regard to create a Keptn project with a not-initialized Git repo. 

## Keptn Specification

Implemented **Keptn spec** version: [0.1.3](https://github.com/keptn/spec/tree/0.1.3)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Allow specifying a domain when installing Keptn (e.g., `keptn install --domain=127.0.0.1.nip.io`) [#1482](https://github.com/keptn/keptn/issues/1482)
- Allow to re-use existing nginx-ingress installation [#1712](https://github.com/keptn/keptn/issues/1712)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *configuration-service:*
  * Improve troubleshooting for git related problems [#1637](https://github.com/keptn/keptn/issues/1637)

</p>
</details>

<details><summary>CLI Enhancements</summary>
<p>

- Create a support-archive for troubleshooting [#1549](https://github.com/keptn/keptn/issues/1549)
- Provide a CLI command for exposing Keptn Bridge [#1560](https://github.com/keptn/keptn/issues/1560)

</p>
</details>

<details><summary>API</summary>
<p>

- Introduce an API-gateway that proxies requests to configuration-service [#1510](https://github.com/keptn/keptn/issues/1510)
- Query a list of projects [#1559](https://github.com/keptn/keptn/issues/1559)
- Provide an endpoint for exposing Keptn's Bridge via Istio or nginx ingress [#1153](https://github.com/keptn/keptn/issues/1153)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Use icons for event types [#1352](https://github.com/keptn/keptn/issues/1352)
- Deep links into Bridge components [#1316](https://github.com/keptn/keptn/issues/1316)
- Format SLOs as floats [#1681](https://github.com/keptn/keptn/issues/1681)

</p>
</details>

## Fixed Issues

- *API:*
  - Do not overwrite `source` property of cloud events [#1643](https://github.com/keptn/keptn/issues/1643)
- *Configuration-service & Shipyard-service:*
  - Catch a not initialized Git repo by creating an initial commit [#1545](https://github.com/keptn/keptn/issues/1545)
  - Fixed error handling (issue with quality-gates multi-stage setups) [#1695](https://github.com/keptn/keptn/issues/1695)
- *Installer:*
  - Check for ImagePullBackOff errors for the installer job [#1521](https://github.com/keptn/keptn/issues/1521)
  - Do not overwrite an existing Keptn installation [#1376](https://github.com/keptn/keptn/issues/1376)
- *Bridge:*
  - Provide proper deep-link functionality for "Problem detected" events [#1557](https://github.com/keptn/keptn/issues/1557)
  - Bridge preselects wrong evaluation event in heatmap view [#1518](https://github.com/keptn/keptn/issues/1518) 
  - Heatmap shows undefined color for test results of type `fail` [#1580](https://github.com/keptn/keptn/issues/1580)
- *Lighthouse:*
  - Evaluating "<=5%" was interpreted as "<=5" (missing percent sign) [#1498](https://github.com/keptn/keptn/issues/1498)
- *Helm:*
  - Helm Service should not require an outbound Internet connection [#1532](https://github.com/keptn/keptn/issues/1532)

## Refactoring

- Added multiple unit tests to improve code coverage
- Refactor api-service and configuration-service [#1510](https://github.com/keptn/keptn/issues/1510)
- Refactor go-utils [#1492](https://github.com/keptn/keptn/issues/1492)
- Change APIVersion from apps/v1beta1 to apps/v1 [#1529](https://github.com/keptn/keptn/issues/1529)

## Development Workflow

- Improve Travis-CI workflow
- Added GitHub actions for linting
- Updated contribution guide

## Good to know / known limitations

- Cluster-internal non-http traffic does not use VirtualServices for Blue/Green deployments [#1715](https://github.com/keptn/keptn/issues/1715)
- For old limitations, please see [Release 0.6.1](https://github.com/keptn/keptn/releases/tag/0.6.1). 
- After executing `keptn configure domain`, an already exposed Keptn Bridge is no longer accessible [#1752](https://github.com/keptn/keptn/issues/1752)
- The installation option --gateway=NodePort currently uses the internal node IP and, hence, a NodePort installation only works if the node can be directly accessed [#1708](https://github.com/keptn/keptn/issues/1708)
