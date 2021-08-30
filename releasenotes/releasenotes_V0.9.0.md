# Release Notes 0.9.0

Keptn 0.9.0 gives you more control over sequence executions and allows creating/deleting a Keptn project in the Bridge.

---

**Key announcements:**

:tada: *Advanced sequences handling*: Keptn core provides new capabilities for handling a task sequence execution:

  * *Controls for sequence executions*: Controls for `pausing`, `resuming`, and `aborting` a task sequence have been implemented. These controls can be used either by the CLI or directly in the Bridge.

  * *Smart defaults for sequence executions*:  Keptn automatically terminates a sequence execution when now Keptn-service is acting upon a triggered task. Besides, it queues a delivery or remediation sequence if there is currently one running for a service in the same stage. 

:star: *Improved UX in Keptn Bridge*: With this release, creating or deleting a Keptn project is now possible via the Bridge. Additionally, cross-linking elements have been added to components for optimizing user flows.

:pick: *Hardening Keptn Core / Bridge*: Hardening of Keptn core components has been conducted concerning resource optimization and enhancing readiness probes. Besides, refactoring in the Bridge was performed in order to reliably deal with custom sequences and tasks.

:apple: *MacOS M1 Support*: For all Apple users, the Keptn CLI is now ready for your new MacOS M1.

:ship: *Docker Registries*: Release images are now also on quay.io and ghcr.io

---


## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 39](https://github.com/keptn/enhancement-proposals/pull/39) and parts of [KEP 53](https://github.com/keptn/enhancement-proposals/pull/53)

## Keptn Specification

Implemented **Keptn spec** version: [0.2.3](https://github.com/keptn/spec/tree/0.2.3)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Reduce K8s resource limits/requests of Keptn core services [3018](https://github.com/keptn/keptn/issues/3018)
- Host Keptn Release Docker Images on multiple container registries/repositories [3314](https://github.com/keptn/keptn/issues/3314)
- Enhance readiness probes of Keptn core services [4518](https://github.com/keptn/keptn/issues/4518)
- Create a list of dependencies of Keptn Core [4409](https://github.com/keptn/keptn/issues/4409)
- Migrate old sequences to materialized view [4140](https://github.com/keptn/keptn/issues/4140)
- *Fixed:* Installing/Upgrading Keptn in an air-gapped environment does not work for configuration-service and nats [4183](https://github.com/keptn/keptn/issues/4183)
- *Fixed:* Bridge `LOOK_AND_FEEL_URL` is missing in Keptn installer Helm Chart [4476](https://github.com/keptn/keptn/issues/4476)

</p>
</details>

<details><summary>API</summary>
<p>

- Provide description for `From` and `To` on GET `/api/statistics/v1` endpoint [3921](https://github.com/keptn/keptn/issues/3921)
- Add parameter `keptnContext` for `/sequence` endpoint [4433](https://github.com/keptn/keptn/issues/4433)
- Extend Uniform to support subscription management [4437](https://github.com/keptn/keptn/issues/4437)
- Remove the `/config/bridge/` endpoint [4589](https://github.com/keptn/keptn/issues/4589)
- Introduce rate limiation on `/auth` endpoint: 429 responses contain information on whether token was valid or not [4906](https://github.com/keptn/keptn/issues/4906)

</p>
</details>

<details><summary>CLI</summary>
<p>

- Support of MacOS M1/Apple Silicon Build [3987](https://github.com/keptn/keptn/issues/3987)
- Commands for pausing/resuming/aborting task sequences [3785](https://github.com/keptn/keptn/issues/3785)
- Adapt output for configure bridge command [4435](https://github.com/keptn/keptn/issues/4435)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *secret-service*:
  - Provide default scope when creating secrets [4281](https://github.com/keptn/keptn/issues/4281)

- *approval-service*:
  - Excluded open approvals from task timeout [4620](https://github.com/keptn/keptn/issues/4620)
  - *Fixed:* Approval-service does not automatically approve in case it is the first task in a sequence [4391](https://github.com/keptn/keptn/issues/4391)

- *distributor*:
  - Allow setting environment details sent by the distributor [4590](https://github.com/keptn/keptn/issues/4590)

- *helm-service & jmeter-service*:
  - Cleanup of README.md and Manifests for jmeter-service/helm-service [4503](https://github.com/keptn/keptn/issues/4503)
  - jmeter-service/helm-service are missing timestamp in tag [4403](https://github.com/keptn/keptn/issues/4403)
  - *Fixed:* Installing jmeter-service/helm-service from a registry with a non-default port does not work [4422](https://github.com/keptn/keptn/issues/4422)

- *lighthouse-service*:
  - Remove override of evaluation result using previous test result [4930](https://github.com/keptn/keptn/issues/4930)

- *remediation-service*:
  - Improve error messages for remediation-services [4412](https://github.com/keptn/keptn/issues/4412)

- *shipyard-controller*:
  - Handle sequences sequentially per stage [3776](https://github.com/keptn/keptn/issues/3776)
  - Termination of orphaned tasks [3778](https://github.com/keptn/keptn/issues/3778)
  - *Fixed:* Run into errors when using an image object for a configuration change [4384](https://github.com/keptn/keptn/issues/4384)
  - *Fixed:* Panics with out of range error [4772](https://github.com/keptn/keptn/issues/4772)
  - *Fixed:* Crashes when receiving event for non-existent project [4797](https://github.com/keptn/keptn/issues/4797)
  - *Fixed:* Race condition in sequence state [4969](https://github.com/keptn/keptn/issues/4969)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - Provide a warning if data will be lost in project creation [4677](https://github.com/keptn/keptn/issues/4677)
  - Add weight to the SLI breakdown [4758](https://github.com/keptn/keptn/issues/4758)
  - Prevent from expanding tile when there is no content [4057](https://github.com/keptn/keptn/issues/4057)
  - Text inside "View Evaluation" is cropped [4760](https://github.com/keptn/keptn/issues/4760)
  - Indicate errors happening in integrations [4381](https://github.com/keptn/keptn/issues/4381)
  - Show service name at sequence tile [4653](https://github.com/keptn/keptn/issues/4653)
  - Show action name and description for a remediation action [4410](https://github.com/keptn/keptn/issues/4410)
  - Rename `Error logs` to `Error events` [4426](https://github.com/keptn/keptn/issues/4426)
  - Delete project via Bridge [4379](https://github.com/keptn/keptn/issues/4379) 
  - Show recent task sequences on project level [2506](https://github.com/keptn/keptn/issues/2506)
  - Create project with shipyard [4493](https://github.com/keptn/keptn/issues/4493) 
  - Show waiting status of successive sequence executions [3777](https://github.com/keptn/keptn/issues/3777)
  - Improve layout of configuring Git upstream repository [4623](https://github.com/keptn/keptn/issues/4623)
  - Show alt text / tooltip for icon buttons [3803](https://github.com/keptn/keptn/issues/3803)
  - Display comparison value and absolute/relative delta of SLI [4305](https://github.com/keptn/keptn/issues/4305)
  - Environment screen always has scrollbars when having more than 2 stages [4146](https://github.com/keptn/keptn/issues/4146)
  - Collapsevaluation heatmap to top 10 [4255](https://github.com/keptn/keptn/issues/4255)
  - Show `keptn create service` when Bridge is used for quality gates only use case [4172](https://github.com/keptn/keptn/issues/4172)
  - Better UX to show which sequence is currently selected [3976](https://github.com/keptn/keptn/issues/3976)
  - Project does not reflect current status after creating a service [4170](https://github.com/keptn/keptn/issues/4170)
  - Add `X-Frame-Options` header to Bridge responses [4257](https://github.com/keptn/keptn/issues/4257)
  - Show subscriptions of integrations [4436](https://github.com/keptn/keptn/issues/4436)
  - Adding / Deleting / Updating subscription [4572](https://github.com/keptn/keptn/issues/4572)
  - Add service name for running sequences on the stage tile [4733](https://github.com/keptn/keptn/issues/4733)
  - Introduce settings navigation [4501](https://github.com/keptn/keptn/issues/4501)
  - Controls for pausing/resuming/aborting task sequences [3798](https://github.com/keptn/keptn/issues/3798)

- *Refactoring:*
  - Add null-check to tsconfig [4628](https://github.com/keptn/keptn/issues/4628)
  - Update Bridge server to TS and ESM [4443](https://github.com/keptn/keptn/issues/4443)
  - Refactor Angular router usage [4022](https://github.com/keptn/keptn/issues/4022)
  - Refactor observables inside of router parameter subscription [4188](https://github.com/keptn/keptn/issues/4188)
  - Migrate testing framework to Jest [4841](https://github.com/keptn/keptn/issues/4841)

- *Fixes:*
  - *OAuth:* Regenerating the session cookie after login [4947](https://github.com/keptn/keptn/issues/4947)
  - *Service Screen:* Keptn context in URI is not properly updated [4912](https://github.com/keptn/keptn/issues/4912)
  - *Sequence screen:* Is blank caused by JavaScript error [4442](https://github.com/keptn/keptn/issues/4442)
  - *Environment screen:* The sequences of services are not loaded [4667](https://github.com/keptn/keptn/issues/4667)
  - *Environment screen:* Is broken caused by JavaScript error [4446](https://github.com/keptn/keptn/issues/4446)
  - *Integration screen:* Update URL for API calls [4830](https://github.com/keptn/keptn/issues/4830)
  - Evaluation results chart is being hidden after page refresh [4927](https://github.com/keptn/keptn/issues/4927)
  - Update message should not print all possible upgradable versions [4831](https://github.com/keptn/keptn/issues/4831)
  - Settings screen is not updated when the project is changed [4781](https://github.com/keptn/keptn/issues/4781)
  - Redirect to dashboard if project is deleted [4765](https://github.com/keptn/keptn/issues/4765)
  - The environment does not always show the right information [4538](https://github.com/keptn/keptn/issues/4538)
  - Report the project on the evaluation page [4759](https://github.com/keptn/keptn/issues/4759)
  - Setting upstream failed, but the error is not shown  [4374](https://github.com/keptn/keptn/issues/4374)
  - Bridge maps deployment event to wrong stage in case of multiple parallel stages with approval [4392](https://github.com/keptn/keptn/issues/4392)
  - Wrong stage focused on deployment selection of a service [4438](https://github.com/keptn/keptn/issues/4438)
  - SLI results in the heatmap are missing if the chart is collapsed [4569](https://github.com/keptn/keptn/issues/4569)
  - Selected deployment gets lost if project is updated [4396](https://github.com/keptn/keptn/issues/4396)
  - Selected service is not reset on project change [4166](https://github.com/keptn/keptn/issues/4166)
  - Bridge throws JavaScript errors - not showing approval option until refresh [4521](https://github.com/keptn/keptn/issues/4521)
  - Cannot read property of 'score' undefined [Bridge] [4936](https://github.com/keptn/keptn/issues/4936)

</p>
</details>


## Development Process / Testing

- *Keptn on Keptn:* Mark test as failed if delivery fails [4186](https://github.com/keptn/keptn/issues/4186)
- *Integration test report:* Note down the actor or whether it was scheduled [3999](https://github.com/keptn/keptn/issues/3999)
- *go-utils/kubernetes-utils:* Add go get and go mod tidy to auto-upgrade PR [4341](https://github.com/keptn/keptn/issues/4341)
- Unify Go versions in CI workflow [4331](https://github.com/keptn/keptn/issues/4331)
- Fix random error messages in integration test workflow [4385](https://github.com/keptn/keptn/issues/4385)
- Create GH-action bot-user for Keptn organization [3825](https://github.com/keptn/keptn/issues/3825)
- Create Auto-PR for Keptn Spec and Keptn Docs [2927](https://github.com/keptn/keptn/issues/2927)
- Create draft release on triggered release workflow for keptn/go-utils [4762](https://github.com/keptn/keptn/issues/4762)
- Create a PR that updates the specification subfolder in keptn/keptn when a new Keptn spec is released [4366](https://github.com/keptn/keptn/issues/4366)
- Activate Snyk protection [417](https://github.com/keptn/keptn/issues/417)
- Unify automated PRs to have the pipeline logic in the target repo [4533](https://github.com/keptn/keptn/issues/4533)
 
## Good to know / Known Limitations

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - `keptn upgrade` does not respect cluster choice [4583](https://github.com/keptn/keptn/issues/4583)
  - Vague error message when setting Git upstream [4399](https://github.com/keptn/keptn/issues/4399)
  - Response time degradation in configuration-service when using a Git upstream (e.g., GitHub) [4066](https://github.com/keptn/keptn/issues/4066)
  - Prometheus self-healing example based on response time does not work [3439](https://github.com/keptn/keptn/issues/3439)
  - Registrations might lose their current subscription if they are scheduled in a different node [4437](https://github.com/keptn/keptn/issues/4437)
  - If a stage cannot be found, the sequence needs to be stopped manually [4791](https://github.com/keptn/keptn/issues/4791)
</p>
</details>

## Upgrade to 0.9.0

- The upgrade from 0.8.x to 0.9.0 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.9.0](https://keptn.sh/docs/0.9.x/operate/upgrade/#upgrade-from-keptn-0-8-x-to-0-9-0)
