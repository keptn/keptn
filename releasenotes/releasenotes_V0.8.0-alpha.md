# Release Notes 0.8.0-alpha

Please note: This is an **alpha release** and should **not** be used for production workloads.

You should expect changes for the final release.

*Info*: Update of LICENSE file [2725](https://github.com/keptn/keptn/issues/2725) 

## Keptn Specification

Implemented **Keptn spec** version: [0.2.0-alpha](https://github.com/keptn/spec/tree/0.2.0-alpha)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Lower K8s resource limits for distributors [2649](https://github.com/keptn/keptn/issues/2649) 
- Upgrade NGNIX unprivileged to latest version [2653](https://github.com/keptn/keptn/issues/2653) 
- Test Keptn control-plane for Kubernetes 1.19 using K3s [2411](https://github.com/keptn/keptn/issues/2411) 

</p>
</details>

<details><summary>API</summary>
<p>

- Remove WebSocket communication between CLI and API [2727](https://github.com/keptn/keptn/issues/2727)

</p>
</details>

<details><summary>CLI</summary>
<p>

- Upgrader for migrating from Shipyard v0.1 to Shipyard v0.2 [2500](https://github.com/keptn/keptn/issues/2500)
- Continue working with current Keptn context and remove Keptn context switch from keptn --help [2721](https://github.com/keptn/keptn/issues/2721)
- Improvement to write version mismatch to std::err [2761](https://github.com/keptn/keptn/issues/2761)
- Re-add the version check into the root command [2571](https://github.com/keptn/keptn/issues/2571)
- Adapt CLI command `keptn send event new-artifact` to CloudEvents spec of 0.8.0 [2558](https://github.com/keptn/keptn/issues/2558)
- Improve post-installation steps by including Keptn API endpoint IP [2444](https://github.com/keptn/keptn/issues/2444)
- Adapt CLI commands `create service`, `onboard service` and `delete service` to use endpoint of the shipyard-controller [2557](https://github.com/keptn/keptn/issues/2557) 
- CLI support creating a project using the new shipyard spec [2266](https://github.com/keptn/keptn/issues/2266) 
- Improved `keptn install --help` messages [2584](https://github.com/keptn/keptn/issues/2584) 
- Keptn support for multiple plans [1863](https://github.com/keptn/keptn/issues/1863) 
- YAML input support for URIs [1648](https://github.com/keptn/keptn/issues/1648) 
- Improved error message when no connection to Keptn API could be established [1349](https://github.com/keptn/keptn/issues/1349) 

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *configuration-service:*
  - Keep track of last processed artifact in materialized view [2692](https://github.com/keptn/keptn/issues/2692)
  - HEAD branch not properly set [2735](https://github.com/keptn/keptn/issues/2735)
  - Updating existing upstream not working [2708](https://github.com/keptn/keptn/issues/2708)
  - Include Git commit ID in GET resource responses [2307](https://github.com/keptn/keptn/issues/2307)

- *distributor*:
  - Having a subscription topic does not have to be a requirement [2562](https://github.com/keptn/keptn/issues/2562)
  - Extend distributor to bridge traffic from Keptn service to Keptn API [2220](https://github.com/keptn/keptn/issues/2220)
  - Sidecar for polling open *.triggered events [2166](https://github.com/keptn/keptn/issues/2166)

- *eventbroker*:
  - Remove eventbroker from Keptn core [2254](https://github.com/keptn/keptn/issues/2254)

- *gatekeeper-service*:
  - gatekeeper-service becomes the approval-service for automatic approvals [2533](https://github.com/keptn/keptn/issues/2533)

- *helm-service*:
  - Create a sequence diagram for helm-service [2592](https://github.com/keptn/keptn/issues/2592)
  - Return Git commit ID in finished events sent by helm-service [2531](https://github.com/keptn/keptn/issues/2531)
  - helm-service reacts on `release.triggered` and sends `release.started/finished` events [2265](https://github.com/keptn/keptn/issues/2265)
  - helm-service reacts on `deployment.triggered` and sends `deployment.started/finished` events [2262](https://github.com/keptn/keptn/issues/2262)

- *jmeter-service*:
  - jmeter-service reacts on `test.triggered` and sends `test.started/finished` events [2263](https://github.com/keptn/keptn/issues/2263)

- *lighthouse-service*:
  - Support quality gates use-case with updated services [2724](https://github.com/keptn/keptn/issues/2724)
  - lighthouse-service reacts on `evaluation.triggered` and sends `evaluation.started/finished` events [2264](https://github.com/keptn/keptn/issues/2264)

- *mongodb-datastore*:
  - Quering (root) events via mongodb-datastore is slow when there is many events in the DB [2759](https://github.com/keptn/keptn/issues/2759)
  - Fixed: mongodb-datastore does not contain "triggeredid" in input [2514](https://github.com/keptn/keptn/issues/2514)

- *remediation-service*
  - Include `triggerid` property in `remediation.status.changed/finished` events [1917](https://github.com/keptn/keptn/issues/1917)
  - Support remediation use-case with updated services [2663](https://github.com/keptn/keptn/issues/2663)

- *shipyard-controller*:
  - Fixed: Shipyard-controller does not set result field of next `.triggered` event [2816](https://github.com/keptn/keptn/issues/2816)
  - Shipyard-controller subscribes to trigger-events defined in the shipyard.yaml and provides a built-in task sequence for evaluations [2529](https://github.com/keptn/keptn/issues/2529)
  - Shipyard-controller is integrated into Travis CI build for release branches [2273](https://github.com/keptn/keptn/issues/2273)
  - Controls the task sequences defined in the Shipyard [2193](https://github.com/keptn/keptn/issues/2193)
  - Manages open *.started events in a mongoDB collection per project [2159](https://github.com/keptn/keptn/issues/2159)
  - Manages open *.triggered events in a mongoDB collection per project [2158](https://github.com/keptn/keptn/issues/2158)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Fixed: Keptn Bridge is not showing notification about the new Keptn version [2693](https://github.com/keptn/keptn/issues/2693)
- Fixed: Keptn Bridge ignores deployed service artifact [2543](https://github.com/keptn/keptn/issues/2543)
- Use an HTTP-interceptor to add default headers and implement generic error handling [1987](https://github.com/keptn/keptn/issues/1987) 
- Added COPY button for SLO content [1997](https://github.com/keptn/keptn/issues/1997)

</p>
</details>

## Fixed Issues

- Fixed several spelling mistakes [2849](https://github.com/keptn/keptn/issues/2849)
- Fixed: Helm chart for continuous-delivery has dependencies to control-plane [2840](https://github.com/keptn/keptn/issues/2840)
- Fixed: Required Flags are not validated before PreRunE is called [2729](https://github.com/keptn/keptn/issues/2729)
- Fixed: CLI does not work when using GPG pass [2638](https://github.com/keptn/keptn/issues/2638)
- Fixed: Storing the credentials does not work with Linux/OpenShift combination [2712](https://github.com/keptn/keptn/issues/2712)
- Fixed: Commands are taking too long to return when no connection to a cluster can be established [2505](https://github.com/keptn/keptn/issues/2505) 
- Fixed: Using `keptn create project` with `--shipyard` pointing to an URL does not properly work [2511](https://github.com/keptn/keptn/issues/2511) 

## Development Process / Testing

<details><summary>Moved CI builds and integration test from Travis-CI to GitHub Actions</summary>
<p>

- Travis-CI builds are disabled due to negative credit balance [2715](https://github.com/keptn/keptn/issues/2715)
- Migrate integration tests from Travis-CI to GitHub Actions [2811](https://github.com/keptn/keptn/issues/2811)
- Migrate go-utils and kubernetes-utils from Travis-CI to GitHub Actions [2796](https://github.com/keptn/keptn/issues/2796)
- Migrate CI from travis-ci.org to travis-ci.com (by Dec. 2020) [2356](https://github.com/keptn/keptn/issues/2356)
- Move Docker builds from Travis-CI to GitHub Actions [2752](https://github.com/keptn/keptn/issues/2752)
- Move unit test execution from TravisCI to GitHub Actions [2716](https://github.com/keptn/keptn/issues/2716)
- Remove hard-dependency of MacOS builds in Travis-CI [2719](https://github.com/keptn/keptn/issues/2719)
- Auto-updating go-utils and kubernetes-utils in keptn/keptn needs to be a signed commit (and moved to GitHub Actions) [2750](https://github.com/keptn/keptn/issues/2750)

</p>
</details>

<details><summary>Fixed CI issues</summary>
<p>

- Fixed: Flaky integration tests: Integration tests fail (in unpredictable situations) [2149](https://github.com/keptn/keptn/issues/2149)
- Fixed: Integration Test stalls at the Keptn auth command [2704](https://github.com/keptn/keptn/issues/2704)
- Fixed: Integration Tests: Setup of Keptn fails due to server version check [2701](https://github.com/keptn/keptn/issues/2701)
- Fixed: Unable to do remote debugging of mongodb-datastore due to liveness-probe [2536](https://github.com/keptn/keptn/issues/2536)
- Fixed: GitHub Action Reviewdog Fails: The `add-path` command is disabled [2694](https://github.com/keptn/keptn/issues/2694)

</p>
</details>

- Test the linking of stages based on task sequence events: `sh.keptn.event.stage-name.sequence-name.finished` [2534](https://github.com/keptn/keptn/issues/2534)
- GKE Integration Tests (for 1.16): Legacy monitoring is not supported in this version [2789](https://github.com/keptn/keptn/issues/2789)
- After a feature/bug/patch/hotfix has been merged, the respective (temporary) images are deleted [1037](https://github.com/keptn/keptn/issues/1037)
- DockerHub: Stale images are going to be deleted soon [2710](https://github.com/keptn/keptn/issues/2710)
- Move tests for delivery assistant and self-healing to K3s [2771](https://github.com/keptn/keptn/issues/2771)
- Only run integration tests on Travis for nightlies to save some build credits [2753](https://github.com/keptn/keptn/issues/2753)
- Move check of deprecated K8s versions from Travis-CI to GitHub Actions [2717](https://github.com/keptn/keptn/issues/2717)
- Reduce number of platform/integration tests on Travis-ci [2718](https://github.com/keptn/keptn/issues/2718)
- Switch from CLA Bot to DCO [2690](https://github.com/keptn/keptn/issues/2690)
- Increased test coverage for helm-service [2530](https://github.com/keptn/keptn/issues/2530)
- Include pluto to automatically check for deprecated K8s apiVersions [2382](https://github.com/keptn/keptn/issues/2383)
- Makefile: Fixed - Build-CLI works, but the resulting binary is not [2504](https://github.com/keptn/keptn/issues/2504)
- Makefile: Add a way to build all Dockerfile [2464](https://github.com/keptn/keptn/issues/2464)
- Makefile: Add build and run targets [2405](https://github.com/keptn/keptn/issues/2405)

## Good to know / Known Limitations

