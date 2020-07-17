# Release Notes 0.7.0

:rocket: *Delivery Assistant:* - [SPEC 26](https://github.com/keptn/spec/pull/26)

:sparkles: *Closed-loop Remediation with custom Integrations:* - [KEP 09](https://github.com/keptn/enhancement-proposals/pull/9) | [SPEC 31](https://github.com/keptn/spec/pull/31)

:rocket: *Improved automation support with API extensions:* - [KEP 10](https://github.com/keptn/enhancement-proposals/pull/10)

:star2: *CD upgrade from Helm 2.0 to 3.0:*

:hammer: *Hardening Keptn:* RBAC

:star2: *Removed Istio and NGINX:* 

## Keptn Specification

Implemented **Keptn spec** version: [0.1.4](https://github.com/keptn/spec/tree/0.1.4)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Kubernetes 1.14 - 1.18 support [#1777](https://github.com/keptn/keptn/issues/1777)
- Keptn on K3s support [#1896](https://github.com/keptn/keptn/issues/1896)
- *Hardening:* Use K8s service account with a restricted set of permissions instead of default [#1862](https://github.com/keptn/keptn/issues/1862)
- *Hardening:* Added Kubernetes recommended labels to the Keptn installation [#1996](https://github.com/keptn/keptn/issues/1996)
- *Installer*: Removed Istio and NGNIX from installer [#1960](https://github.com/keptn/keptn/issues/1960)
- *OpenShift:* `keptn uninstall` command mistakenly recommended to delete several OpenShift namespaces [#1781](https://github.com/keptn/keptn/issues/1781)

</p>
</details>

<details><summary>API</summary>
<p>

- Expose `/event` endpoint from mongodb-datastore to the public Keptn API [#1791](https://github.com/keptn/keptn/issues/1791)
- Change Keptn API and Keptn Bridge path on ingress from subdomain to suffix [#1994](https://github.com/keptn/keptn/issues/1994)
- Retrieve metadata of Keptn installation [#1843](https://github.com/keptn/keptn/issues/1843)
- *Keptn Configure Bridge:* Do not expose the service, nor apply Istio/NGINX manifests [#1962](https://github.com/keptn/keptn/issues/1962) 

</p>
</details>


<details><summary>CLI</summary>
<p>

- Polished the user output messages and double check links [#2042](https://github.com/keptn/keptn/issues/2042)
- Removed `--scheme=http` when using Keptn CLI with HTTP instead of HTTPs [#1948](https://github.com/keptn/keptn/issues/1948)
- `keptn onboard service` is aborted when continuous.delivery is not installed [#2047](https://github.com/keptn/keptn/issues/2047)
- `keptn install` removed anything related to Istio and NGINX [#1961](https://github.com/keptn/keptn/issues/1961)
- `keptn install` removed `--platform` flag [#1967](https://github.com/keptn/keptn/issues/1967)
- Keptn generate support-archive should have a separate check for ingress options [#1941](https://github.com/keptn/keptn/issues/1941)
- Show warning when creating a project without Git upstream [#1840](https://github.com/keptn/keptn/issues/1840)
- Allow specify an upstream Git for existing projects [#1517](https://github.com/keptn/keptn/issues/1517)
- Allow user to send an approval event to the provided stage and to approve a deployment using the CLI [#1749](https://github.com/keptn/keptn/issues/1749)
- Removed fixed host header `api.keptn` in CLI commands [#1797](https://github.com/keptn/keptn/issues/1797)
- Implemented delivery assistant for approving a deployment [#1835](https://github.com/keptn/keptn/issues/1835)
- Implemented get projects, services, stages, and metadata [#1624](https://github.com/keptn/keptn/issues/1624)
- Enforce username and password when configuring Keptn Bridge [#1893](https://github.com/keptn/keptn/issues/1893)
- Improved the output of Keptn CLI for troubleshooting [#1928](https://github.com/keptn/keptn/issues/1928)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *configuration-service:*
  * Manage open remediation workflows in the materialized view [#1848](https://github.com/keptn/keptn/issues/1848)
  * Allow to retrieve all open approval events for a specific project, stage, and service [#1757](https://github.com/keptn/keptn/issues/1757)

- *gatekeeper-service:*
  * React on an approval.finished event to send configuration.changed event for the current stage [#1737](https://github.com/keptn/keptn/issues/1737)
  * Read approval_strategy and send event based on configured strategy and evaluation result [#1658](https://github.com/keptn/keptn/issues/1658)

- *helm-service:*
  * Introduce a new ConfigMap for INGRESS_HOSTNAME_SUFFIX [#1963](https://github.com/keptn/keptn/issues/1963)
  * Gateway in generated VirtualServices is configurable via environment variable [#1986](https://github.com/keptn/keptn/issues/1986)

- *jmeter-service:*
  * Properly handle errors from configuration-service [#1480](https://github.com/keptn/keptn/issues/1480)

- *mongodb-service:*
  * Manage open approval events in a collection [#1756](https://github.com/keptn/keptn/issues/1756)
  * Moved MongoDB credentials into a Kubernetes secret [#1528](https://github.com/keptn/keptn/issues/1528) 
  * Increased MongoDB datastore volume size [#1900](https://github.com/keptn/keptn/issues/1900)

- *remediation-service:*
  * Extracted featuretoggle action from remediation-service into *unleash-service* [#1816](https://github.com/keptn/keptn/issues/1816)
  * Moved functionality of scaler to *helm-service* [#1817](https://github.com/keptn/keptn/issues/1817)
  * Moved posting Dynatrace problem comments to *dynatrace-service* [#1818](https://github.com/keptn/keptn/issues/1818)
  * React on problem.open and process pre-defined workflow: trigger action, wait, evaluate, continue remediation or send a remediation.finished [#1849](https://github.com/keptn/keptn/issues/1849)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Update UI look-and-feel [#1974](https://github.com/keptn/keptn/issues/1974)
- Splitted UI into *Environment* and *Services* view [#1698](https://github.com/keptn/keptn/issues/1698)
- *Environment view:* Click on stage shows stage information and currently deployed services in panel on the right side [#1699](https://github.com/keptn/keptn/issues/1699)
- *Environment view:* Displays that a service is *out-of-sync* in stage overview and detail info [#1700](https://github.com/keptn/keptn/issues/1700)
- *Environment view:* Introduced buttons to approve/decline a deployment of a service that is *out-of-sync* [#1701](https://github.com/keptn/keptn/issues/1701)
- *Environment view:* Shows status information in stages when stage is empty (no service deployed) [#1860](https://github.com/keptn/keptn/issues/1860)
- Extended horizontal axis of the bar chart from a timeline to fixed distances [#1668](https://github.com/keptn/keptn/issues/1668)
- Get HeatMap of evaluation-done event including deep link into Bridge [#1677](https://github.com/keptn/keptn/issues/1677)
- Provide a "COPY JSON" button on the Bridge [#1794](https://github.com/keptn/keptn/issues/1794)
- Improved JSON payload visualization [#1420](https://github.com/keptn/keptn/issues/1420)
- Use the public API for query list of projects, stages, and services instead of connecting directly to configuration-service [#1657](https://github.com/keptn/keptn/issues/1657)
- Notify user of new available Keptn Bridge in UI [#1547](https://github.com/keptn/keptn/issues/1547)
- Filter events in list of root events [#1342](https://github.com/keptn/keptn/issues/1342)
- Unit tests for Bridge [#1486](https://github.com/keptn/keptn/issues/1486)

</p>
</details>

## Fixed Issues
- Project with two stages broke after lighthouse run at second stage [#1695](https://github.com/keptn/keptn/issues/1695)
- After `keptn configure domain` an already exposed bridge was no longer accessible [#1752](https://github.com/keptn/keptn/issues/1752)
- Eventbroker-go crashesd with out-of-memory [#1901](https://github.com/keptn/keptn/issues/1901)
- MongoDB performed incomplete read of message header [#1907](https://github.com/keptn/keptn/issues/1907)
- MongoDB failed on start with certain PVs [#1519](https://github.com/keptn/keptn/issues/1519)
- Keptn CLI used wrong Kubernetes Context/Profile [#1942](https://github.com/keptn/keptn/issues/1942)
- *OpenShift*: Start of api-gateway-nginx failed with permission denied error [#1951](https://github.com/keptn/keptn/issues/1951)

## Development Process / Testing

- Integration tests for verifying implementation of KEP 18 [#1965](https://github.com/keptn/keptn/issues/1965)
- Integration tests for the use-case of scaling a ReplicaSet based on a Prometheus alert [#1847](https://github.com/keptn/keptn/issues/1847)
- Platform/integration test for manual approval use-case [#1750](https://github.com/keptn/keptn/issues/1750)
- Platform/integration test for self-healing use-case [#1846](https://github.com/keptn/keptn/issues/1846)
- Improve integration test for quality gates use-case [#1591](https://github.com/keptn/keptn/issues/1591)

## Good to know / Known Limitations

* **Upgrade from 0.6.2 to 0.7:** *Keptn 0.7 uses Helm 3 while previous Keptn releases rely on Helm 2*. To upgrade  your Helm releases from Helm 2 to 3, two options are provided: 
  1. *Job without Helm 3 Upgrade:* This option is needed when the cluster contains Helm releases not managed by Keptn. If this job is executed, it is necessary to manually converted the releases from Helm 2 to 3 as explained on [keptn.sh/docs](https://keptn.sh/docs/0.7.0/operate/upgrade/#job-without-helm-3-0-upgrade).
  1. *Job with Helm 3 Upgrade:* Full automation of Helm upgrade for installations were just Keptn is installed. If this job is executed, **all** Helm releases on the cluster are converted from Helm 2 to 3 and Tiller will be removed.
  