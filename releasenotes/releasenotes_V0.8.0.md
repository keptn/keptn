# Release Notes 0.8.0

Keptn 0.8 improves the core use cases of continuous delivery and automated operations by implementing the new Shipyard version [v0.2](https://github.com/keptn/spec/tree/0.2.0). This new Shipard version has been proposed and refined in [KEP 06](https://github.com/keptn/enhancement-proposals/pull/6).

---

**Key announcements:**

:rocket: *Separated and explicit Shipyard-defined processes (aka. task sequences) for continuous delivery / automated remediation*: With this Keptn release, it is possible to have multiple processes (called *task sequences*) within one stage. These task sequences are separated from each other and list the tasks that are triggered sequentially.

:tada: *Support individual tasks in sequences*: It is now possible to add custom tasks to a *task sequence* to address needs of delivery/remediation use-cases that go beyond the opinionated approach Keptn is offering.

:star: *Trigger of a sequence can be configured - allowing multiple parallel stages*: The new Shipyard supports the definition of triggers that launch the execution of a *task sequences*. This helps to make it explicit when a sequence gets triggers. Besides, this linking mechanism allows connecting multiple sequences (of different stages) to listen to the same trigger. Consequently, it is possible to connect multiple stages, which are on the same level, to one preceding stage.

:star2: *New types of events*: In course of implement the new Shipyard version in Keptn, the Keptn Cloud-events were streamlined and follow now a common pattern. Basically, Keptn just sends out an event of type: `sh.keptn.event.{task.name}.triggered` and other services react on: 
  * `sh.keptn.event.{task.name}.triggered`      > *sent out by Keptn*
  * `sh.keptn.event.{task.name}.started`        > *sent out by a Keptn-service when starting the tasks*
  * `sh.keptn.event.{task.name}.status.changed` > *sent out by a Keptn-service to inform about a status update*
  * `sh.keptn.event.{task.name}.finished`       > *sent out by a Keptn-service when the task is finished*

:sparkles: *Multi-cluster support*: Based on the implementation of Shipyard v0.2, Keptn - as a control-plane for delivery and remediation - is now capable of serving multiple clusters. This is known as the split between *Control plane* and *Execution plane*. For this use-case, the Keptn project offers to run the helm-service (to deploy) and jmeter-service (to test) on the *Execution plane*. This *Execution plane* can be on a cluster other than the cluster where Keptn is installed. 

:dizzy: *Sequence screen in Keptn Bridge*: The new capabilities of Keptn for dealing with task sequences received a dedicated screen in the Keptn Bridge. This screen provides filtering capabilities and a stage-divided view on the performed delivery or remediation tasks. 

> *Screenshot here*

---

**Supporting features:**

:tada: *Query usage statistics of your Keptn*: With this release, it is possible to retrieve usage statistics of a Keptn by using the `/api/statistics/v1` endpoint. This returns the number of events processed in the specified time frame. 

:star2: *Keptn CLI supports multiple Keptn installations*: The new Keptn CLI easies working with multiple Keptns since it recognizes switches between Kubernetes clusters and then asks for switching the context Keptn context too. Consequently, your CLI will be automatically connected to the Keptn running on another K8s cluster.   

:star: *Deployment of custom Helm Charts*: An extension of the helm-service allows to deploy custom Helm Charts meaning that the Helm Chart can contain any custom resource and is not limited to a *Kubernetes service* and *deployment*. *Note:* When using this option, the automatic rollback capability of Keptn is not supported and the Helm Chart is not under control by Keptn. Consequently, this feature is currently marked as experimental.

:sparkles: *SLI breakdown displayed as a table in Keptn Bridge*: For the quality gates capabilities of Keptn, the SLI breakdown is now displayed as a table given a better overview of the individual results. 

---

**Noteworthy changes and improvements:**

- Removed WebSocket communication between CLI and API
- Performance improvements in MongoDB
- Update of LICENSE file [2725](https://github.com/keptn/keptn/issues/2725)

Last but not least, many thanks to the community for the rich discussions around Keptn 0.8, the submitted [Keptn Enhancement Proposals](https://github.com/keptn/enhancement-proposals), and the implementation work!
 

## Keptn Specification

Implemented **Keptn spec** version: [0.2.0](https://github.com/keptn/spec/tree/0.2.0)


## Breaking changes

### API

- Introduction of shipyard controller API via: `/api/controlPlane/v1`
- Adding and updating resources works on the endpoint: `/api/configuration-service/v1`
- **Project** endpoints have been moved to: `/api/controlPlane/v1/project`
- **Stage** endpoints have been moved to: `/api/controlPlane/v1/stage`
- **Service** endpoints have been moved to: `/api/controlPlane/v1/service`
- **Evaluation** endpoint for triggering an evaluaiton has been moved to: `/api/v1/project​/{project}​/stage​/{stage}​/service​/{service}​/evaluation`
- **Events** /GET endpoint has been moved to: `/api/mongodb-datastore/event`

### CLI

- `keptn send event start-evaluation` to trigger an evaluation has been marked as deprecated. Use `keptn trigger evaluation` instead:

  ```
  keptn trigger evaluation --project=my-sockshop --service=foobar --stage=hardening
  ```

- `keptn send event new-artifact` to send a configuration change that triggers a delivery of a new artifact has been marked as deprecated. Use `kept trigger delivery` instead: 

  ```
  keptn trigger delivery --project=sockshop --service=carts-db --image=docker.io/mongo --tag=4.2.2 --sequence=delivery-direct
  ```

### Bridge

- The **Service** screen does not show the Keptn CloudEvents anymore since this information has moved to the new **Sequence** screen. 

- The old deep links still work but are adapted to the new screens as described [here](https://keptn.sh/docs/0.8.x/reference/bridge/deep_linking/#links-to-project-and-events)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Lower Kubernetes resource limits for distributors [2649](https://github.com/keptn/keptn/issues/2649) 
- Upgrade of NGNIX unprivileged to latest version [2653](https://github.com/keptn/keptn/issues/2653) 
- Test Keptn Keptn Control-plane for Kubernetes 1.19 using K3s [2411](https://github.com/keptn/keptn/issues/2411) 
- *Fixed*: `keptn install` hangs in case of ImagePullBackOff [2988](https://github.com/keptn/keptn/issues/2988) 

</p>
</details>

<details><summary>API</summary>
<p>

- Hide implementation details in the API [3001](https://github.com/keptn/keptn/issues/3001)
- Streamline Keptn API [2772](https://github.com/keptn/keptn/issues/2772)
- Remove uploading an Helm Chart on PUSH `/service` endpoint [3195](https://github.com/keptn/keptn/issues/3195)
- List services in alphabetical order on GET `/service` endpoint [2754](https://github.com/keptn/keptn/issues/2754)
- Parse shipyard and returns version or whether it is valid/invalid on GET `/project` endpoint [2804](https://github.com/keptn/keptn/issues/2804)
- Remove WebSocket communication between CLI and API [2727](https://github.com/keptn/keptn/issues/2727)
- *Fixed*: GET `/api/v1/metadata` returns null during K8s api downtime [2870](https://github.com/keptn/keptn/issues/2870)
- *Fixed*: API allows creating projects with special characters [2914](https://github.com/keptn/keptn/issues/2914)

</p>
</details>

<details><summary>CLI</summary>
<p>

- `keptn --help` Continue working with current Keptn context and remove Keptn context switch from [2721](https://github.com/keptn/keptn/issues/2721)
- `keptn create service` | `onboard service` | `delete service` - adapt CLI commands to use endpoint of the shipyard-controller [2557](https://github.com/keptn/keptn/issues/2557) 
- `keptn create project` - support for creating a project using the new shipyard spec [2266](https://github.com/keptn/keptn/issues/2266) 
- `keptn get event` - allow polling Keptn Cloud-events (e.g., by cloud-event-id) [2572](https://github.com/keptn/keptn/issues/2572)
- `keptn get event` - ensure compatibility with new cloud-events (e.g., evaluation.finished instead of evaluation-done) [2873](https://github.com/keptn/keptn/issues/2873)
- `keptn get project` - display shipyard version [2908](https://github.com/keptn/keptn/issues/2908)
- `keptn generate cloud-events-spec` - new command for generating Keptn Cloud-events specification [2926](https://github.com/keptn/keptn/issues/2926)
- `keptn install --help` - improved install message [2584](https://github.com/keptn/keptn/issues/2584) 
- `keptn send event new-artifact` - adapt CLI command to CloudEvents spec of 0.8.0 [2558](https://github.com/keptn/keptn/issues/2558)
- `keptn upgrade` - better instructions on how to download new CLI version  [2560](https://github.com/keptn/keptn/issues/2560)
- `keptn upgrade` - avoid the version check via a flag [2689](https://github.com/keptn/keptn/issues/2689)
- `keptn upgrade project` - upgrader for migrating from Shipyard v0.1 to Shipyard v0.2 [2500](https://github.com/keptn/keptn/issues/2500)
- `keptn version` - re-add the version check into the root command [2571](https://github.com/keptn/keptn/issues/2571)
- Add labels parameter to all keptn send events [2126](https://github.com/keptn/keptn/issues/2126)
- Removed outdated xip.io resolver [3058](https://github.com/keptn/keptn/issues/3058)
- Shell completion for Keptn CLI using Cobra [2539](https://github.com/keptn/keptn/issues/2539)
- Support for installing Keptn CLI via Homebrew [2864](https://github.com/keptn/keptn/issues/2864)
- Improvement to write version mismatch to std::err [2761](https://github.com/keptn/keptn/issues/2761)
- Improved post-installation steps by including Keptn API endpoint [2444](https://github.com/keptn/keptn/issues/2444)
- Keptn CLI support for multiple plans [1863](https://github.com/keptn/keptn/issues/1863) 
- YAML input support for URIs [1648](https://github.com/keptn/keptn/issues/1648) 
- Improved error message when no connection to Keptn API could be established [1349](https://github.com/keptn/keptn/issues/1349) 
- *Fixed*: Keptn tabular CLI output breaks automation with too long project, stage, or service names [2899](https://github.com/keptn/keptn/issues/2899)
- *Fixed*: Keptn 0.8.0-alpha CLI crashes for auth after upgrade from 0.7.3 [2912](https://github.com/keptn/keptn/issues/2912)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *api-gateway-nginx:*
  - Always restart api-gateway-nginx deployment on changes [3320](https://github.com/keptn/keptn/issues/3320)

- *configuration-service:*
  - Keep track of last processed artifact in materialized view [2692](https://github.com/keptn/keptn/issues/2692)
  - HEAD branch of Git repository not properly set [2735](https://github.com/keptn/keptn/issues/2735)
  - Include Git commit ID in GET `\resource` responses [2307](https://github.com/keptn/keptn/issues/2307)
  - *Fixed*: Updating existing upstream not working [2708](https://github.com/keptn/keptn/issues/2708)
  - *Fixed*: Pushing to upstream URL currently not working [3227](https://github.com/keptn/keptn/issues/3237)

- *distributor*:
  - Simplified event filter for distributor [3262](https://github.com/keptn/keptn/issues/3262)
  - Handle empty values of environment variables more reliably [2646](https://github.com/keptn/keptn/issues/2646) 
  - Removed subscription topic as requirement for the distributor to work [2562](https://github.com/keptn/keptn/issues/2562)
  - Extend distributor to bridge traffic from Keptn-service to Keptn API [2220](https://github.com/keptn/keptn/issues/2220)
  - Sidecar for polling open `*.triggered` events [2166](https://github.com/keptn/keptn/issues/2166)

- *eventbroker*:
  - Removed eventbroker from Keptn core [2254](https://github.com/keptn/keptn/issues/2254)

- *gatekeeper-service* --> *approval-service*:
  - Move gatekeeper-service to Keptn core and rename it to approval-service [3252](https://github.com/keptn/keptn/issues/3252)
  - Renamed to approval-service for automatic approvals [2533](https://github.com/keptn/keptn/issues/2533)

- *helm-service*: 
  - Support for `deployment_strategy: user_managed` that allows to deploy custom Helm charts [2764](https://github.com/keptn/keptn/issues/2764)
  - Check length of release name [2948](https://github.com/keptn/keptn/issues/2948)
  - Support https and x-token based communication with configuration endpoint [2841](https://github.com/keptn/keptn/issues/2841)
  - Make public deployment URI configurable [2362](https://github.com/keptn/keptn/issues/2362)
  - Created a sequence diagram for helm-service [2592](https://github.com/keptn/keptn/issues/2592)
  - Return Git commit ID in finished events [2531](https://github.com/keptn/keptn/issues/2531)
  - Increased test coverage for helm-service [2530](https://github.com/keptn/keptn/issues/2530)
  - Reacts on `release.triggered` and sends `release.started/finished` event [2265](https://github.com/keptn/keptn/issues/2265)
  - Reacts on `deployment.triggered` and sends `deployment.started/finished` event [2262](https://github.com/keptn/keptn/issues/2262)
  - *Fixed*: Fixed hostname template processing [2932](https://github.com/keptn/keptn/issues/2932)
  - *Verification*: How does helm-service behave on a faulty, user_managed Helm Chart? [3258](https://github.com/keptn/keptn/issues/3258)

- *jmeter-service*:
  - Loads JMeter extensions such as Prometheus or Dynatrace backend listener [2552](https://github.com/keptn/keptn/issues/2552)
  - Reacts on `test.triggered` and sends `test.started/finished` event [2263](https://github.com/keptn/keptn/issues/2263)

- *lighthouse-service*:
  - Support quality gates use-case with updated services [2724](https://github.com/keptn/keptn/issues/2724)
  - Reacts on `evaluation.triggered` and sends `evaluation.started/finished` event [2264](https://github.com/keptn/keptn/issues/2264)
  - *Fixed:* Needs to send previous payloads (e.g., "deployment") in `get-sli.triggered` [3411](https://github.com/keptn/keptn/issues/3411)

- *mongodb-datastore*:
  - Adapt query for excluding `evaluation.invalidated` events [3270](https://github.com/keptn/keptn/issues/2949)
  - Support backwards compatibility for `evaluation-done` events used in Keptn < 0.8 [2949](https://github.com/keptn/keptn/issues/2949)
  - Improve MongoDB datastore performance [2925](https://github.com/keptn/keptn/issues/2925)
  - Improved quering (root) events from mongodb-datastore when there are many events in the DB [2759](https://github.com/keptn/keptn/issues/2759)
  - *Fixed*: mongodb-datastore does not contain `triggeredid` in input [2514](https://github.com/keptn/keptn/issues/2514)

- *remediation-service*
  - Moved the storage of open remediations from *configuration-service* to *remediation-service* [2998](https://github.com/keptn/keptn/issues/2998)
  - Include `triggerid` property in `remediation.status.changed/finished` events [1917](https://github.com/keptn/keptn/issues/1917)
  - Support remediation use-case with updated services [2663](https://github.com/keptn/keptn/issues/2663)

- *shipyard-controller*:
  - Add `triggeredid` to finished event for a sequence [3329](https://github.com/keptn/keptn/issues/3329)
  - API returns shipyard version 0.1.7, although 0.2.0 is used [3325](https://github.com/keptn/keptn/issues/3325)
  - Keptn supports default sequences for "delivery", "evaluation" [3007](https://github.com/keptn/keptn/issues/3007)
  - Add keptn/spec version to metadata of Keptn CloudEvents [2983](https://github.com/keptn/keptn/issues/2983)
  - Removed `data.message` property from previous `.finished` event before sending next `.triggered` event [3043](https://github.com/keptn/keptn/issues/3043)
  - Propagate configurationChange through all tasks of a sequence [3199](https://github.com/keptn/keptn/issues/3199)
  - Allow filtering sequence triggers based on match properties [3028](https://github.com/keptn/keptn/issues/3028)
  - Trigger next stage regardless of evaluation result [3008](https://github.com/keptn/keptn/issues/3008)
  - Stops sequence when a task returns `result=fail` [3027](https://github.com/keptn/keptn/issues/3027)
  - Moved GET endpoints for project, stage, and service details from *configuration-service* to *shipyard-controller* [2999](https://github.com/keptn/keptn/issues/2999)
  - Checks whether the shipyard file is valid and has right version [2803](https://github.com/keptn/keptn/issues/2803)
  - Subscribes to trigger-events defined in the shipyard.yaml and provides a built-in task sequence for evaluations [2529](https://github.com/keptn/keptn/issues/2529)
  - Integrated into Travis CI build for release branches [2273](https://github.com/keptn/keptn/issues/2273)
  - Controls the task sequences defined in the Shipyard [2193](https://github.com/keptn/keptn/issues/2193)
  - Manages open `*.started` events in a MongoDB collection per project [2159](https://github.com/keptn/keptn/issues/2159)
  - Manages open `*.triggered` events in a MongoDB collection per project [2158](https://github.com/keptn/keptn/issues/2158)
  - *Fixed*: Do not return Internal server error when no matching `.triggered` event is available for a `.started/.finished` event [2956](https://github.com/keptn/keptn/issues/2956)
  - *Fixed*: Shipyard-controller does not set result field of next `.triggered` event [2816](https://github.com/keptn/keptn/issues/2816)

- *statistics-service*:
  - Moving the *statistics-service* to Keptn API endpoint [2809](https://github.com/keptn/keptn/issues/2809)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- *new Sequence screen:* 
  - Create stage-timeline component [2907](https://github.com/keptn/keptn/issues/2907)
  - Highlight the selected stage in the timeline [3049](https://github.com/keptn/keptn/issues/3049)
  - Add filter component and apply filter on list of sequences [2626](https://github.com/keptn/keptn/issues/2626)
  - Create sequence screen and load all triggers [2625](https://github.com/keptn/keptn/issues/2625)
  - Show task details in sequence details [2938](https://github.com/keptn/keptn/issues/2938)
  - Refinement of the sequence tile [2628](https://github.com/keptn/keptn/issues/2628)
- Replace occurances of old "send event" with the new "trigger" functionality [3332](https://github.com/keptn/keptn/issues/3332)
- Link back to evaluation from Environment [2696](https://github.com/keptn/keptn/issues/2696)
- Support deep links in Bridge for 0.8.x [3207](https://github.com/keptn/keptn/issues/3207)
- Adapt invalidation of events [3290](https://github.com/keptn/keptn/issues/3290)
- SLI breakdown in table [2478](https://github.com/keptn/keptn/issues/2478)
- Service screen refinement [3206](https://github.com/keptn/keptn/issues/3206)
- Update doc references for 0.8.x [3205](https://github.com/keptn/keptn/issues/3205)
- Misleading information in event stream of sequence screen [3016](https://github.com/keptn/keptn/issues/3016)
- Show shipyard version in project tile [2909](https://github.com/keptn/keptn/issues/2909)
- Sort the SLOs names in the Keptn Bridge Evaluation done results [1499](https://github.com/keptn/keptn/issues/1499)
- Highlight stages more prominent in eventflow [2229](https://github.com/keptn/keptn/issues/2229)
- Shows configuration-change events and its content [2872](https://github.com/keptn/keptn/issues/2872)
- Shows evaluations for new `evaluation.finished` events (instead of `evaluation-done`) [2866](https://github.com/keptn/keptn/issues/2866) 
- Use an HTTP-interceptor to add default headers and implement generic error handling [1987](https://github.com/keptn/keptn/issues/1987) 
- Added COPY button for SLO content [1997](https://github.com/keptn/keptn/issues/1997)
- *Refactoring*: Split project-board into smaller components [1989](https://github.com/keptn/keptn/issues/1989)
- *Refactoring*: Create view-component for environments tab [2939](https://github.com/keptn/keptn/issues/2939)
- *Refactoring*: Create view-component for integration tab [2942](https://github.com/keptn/keptn/issues/2942)
- *Refactoring*: Create stage-overview component [2943](https://github.com/keptn/keptn/issues/2943)
- *Refactoring*: Create stage-details component [2944](https://github.com/keptn/keptn/issues/2944)
- *Refactoring*: Create view-component for sequences tab [2941](https://github.com/keptn/keptn/issues/2941)
- *Refactoring*: Create view-component for services tab [2940](https://github.com/keptn/keptn/issues/2940)
- *Fixed*: Duplicate tasks showing up in Bridge [3382](https://github.com/keptn/keptn/issues/3382)
- *Fixed*: Sequence loading icon [3410](https://github.com/keptn/keptn/issues/3410)
- *Fixed*: Wrong score in SLI breakdown table [3383](https://github.com/keptn/keptn/issues/3223)
- *Fixed*: Root events are limited to 20 [3223](https://github.com/keptn/keptn/issues/3223)
- *Fixed*: Keptn Bridge: Deployed services is displayed as "not deployed" [3224](https://github.com/keptn/keptn/issues/3224
- *Fixed*: Manual approval does not trigger next task in sequence [3013](https://github.com/keptn/keptn/issues/3013)
- *Fixed*: ERROR TypeError: this.data.configurationChange.values.image is undefined [3021](https://github.com/keptn/keptn/issues/3021)
- *Fixed*: Approval not possible in cases when having the manual deployment strategy [2901](https://github.com/keptn/keptn/issues/2901)
- *Fixed*: Keptn Bridge is not showing notification about the new Keptn version [2693](https://github.com/keptn/keptn/issues/2693)
- *Fixed*: Keptn Bridge ignores deployed service artifact [2543](https://github.com/keptn/keptn/issues/2543)

</p>
</details>

## Miscellaneous

- Fixed several spelling mistakes [2849](https://github.com/keptn/keptn/issues/2849)
- Format the Go imports [3150](https://github.com/keptn/keptn/issues/3150)
- Test the linking of stages based on task sequence events: `sh.keptn.event.[stage].[sequence].finished` [2534](https://github.com/keptn/keptn/issues/2534)

<details><summary>Update of third-party dependencies to their latest version, most notable are:</summary>
<p>
 
* *Go* (Microservices)
  - google/uuid to 1.2.0
  - go.mongodb.org/monto-driver to 1.4.6
  - cloudevents/sdk-go (various versions needed)
  - nats-io/nats-server/v2 to 2.1.9
* *NodeJS* (Bridge)
  - marked to 2.0.0
  - higlights.js to 10.4.1

</p>
</details>

## Fixed Issues

- *Fixed*: Cannot run `helm-service` and `jmeter-service` on execution plane on a separate cluster/namespace [3418](https://github.com/keptn/keptn/issues/3418)
- *Fixed*: Upgrade from 0.7.3 to 0.8.0-rc1 failed (because of statistics-service) [3399](https://github.com/keptn/keptn/issues/3399)
- *Fixed*: Helm chart for continuous-delivery has dependencies to control-plane [2840](https://github.com/keptn/keptn/issues/2840)
- *Fixed*: Required flags are not validated before PreRunE is called [2729](https://github.com/keptn/keptn/issues/2729)
- *Fixed*: CLI does not work when using GPG pass [2638](https://github.com/keptn/keptn/issues/2638)
- *Fixed*: Storing the credentials does not work with Linux/OpenShift combination [2712](https://github.com/keptn/keptn/issues/2712)
- *Fixed*: Commands are taking too long to return when no connection to a cluster can be established [2505](https://github.com/keptn/keptn/issues/2505) 
- *Fixed*: Using `keptn create project` with `--shipyard` pointing to an URL does not properly work [2511](https://github.com/keptn/keptn/issues/2511) 

## Development Process / Testing

<details><summary>Moved CI builds and integration test from Travis-CI to GitHub Actions</summary>
<p>

- Travis-CI builds are disabled [2715](https://github.com/keptn/keptn/issues/2715)
- Migrate integration tests from Travis-CI to GitHub Actions [2811](https://github.com/keptn/keptn/issues/2811)
- Migrate go-utils and kubernetes-utils from Travis-CI to GitHub Actions [2796](https://github.com/keptn/keptn/issues/2796)
- Migrate CI from travis-ci.org to travis-ci.com (by Dec. 2020) [2356](https://github.com/keptn/keptn/issues/2356)
- Move Docker builds from Travis-CI to GitHub Actions [2752](https://github.com/keptn/keptn/issues/2752)
- Move check of deprecated K8s versions from Travis-CI to GitHub Actions [2717](https://github.com/keptn/keptn/issues/2717)
- Move unit test execution from TravisCI to GitHub Actions [2716](https://github.com/keptn/keptn/issues/2716)
- Remove hard-dependency of MacOS builds in Travis-CI [2719](https://github.com/keptn/keptn/issues/2719)
- Auto-updating go-utils and kubernetes-utils in keptn/keptn needs to be a signed commit (and moved to GitHub Actions) [2750](https://github.com/keptn/keptn/issues/2750)

</p>
</details>

<details><summary>Miscellaneous CI tasks (for build, test, quality checks)</summary>
<p>

- Added PAT to the create release branch workflow [3393](https://github.com/keptn/keptn/issues/3393)
- Multi-architecture build support for CLI (32 bit, ARM, ...) [2997](https://github.com/keptn/keptn/issues/2997)
- Add dependabot to keep dependencies up2date [2648](https://github.com/keptn/keptn/issues/2648)
- Switch from CLA Bot to DCO [2690](https://github.com/keptn/keptn/issues/2690)
- Solved problems with test coverage reporting [2929](https://github.com/keptn/keptn/issues/2929)
- Various improvements for GH actions [2824](https://github.com/keptn/keptn/issues/2824)
- Include pluto to automatically check for deprecated K8s apiVersions [2382](https://github.com/keptn/keptn/issues/2383)
- Integration tests: enable shielded GKE nodes for integration tests and in docs [2973](https://github.com/keptn/keptn/issues/2973)
- Integration tests: use newer Istio version [2976](https://github.com/keptn/keptn/issues/2976)
- Integration tests (GKE for 1.16): Legacy monitoring is not supported in this version [2789](https://github.com/keptn/keptn/issues/2789)
- After a feature/bug/patch/hotfix has been merged, the respective (temporary) images are deleted [1037](https://github.com/keptn/keptn/issues/1037)
- DockerHub: Stale images are going to be deleted soon [2710](https://github.com/keptn/keptn/issues/2710)
- Move tests for delivery assistant and self-healing to K3s [2771](https://github.com/keptn/keptn/issues/2771)
- Reduce the number of platform/integration tests on Travis-CI [2718](https://github.com/keptn/keptn/issues/2718)
- Makefile: *Fixed* - Build-CLI works, but the resulting binary is not [2504](https://github.com/keptn/keptn/issues/2504)
- Makefile: Add a way to build all Dockerfile [2464](https://github.com/keptn/keptn/issues/2464)
- Makefile: Add build and run targets [2405](https://github.com/keptn/keptn/issues/2405)

</p>
</details>

<details><summary>Fixed CI issues</summary>
<p>

- *Fixed*: Integration tests are failing (Minishift, self-healing) [3325](https://github.com/keptn/keptn/issues/3325)
- *Fixed*: Integration tests: GKE clusters are not deleted afterwards in some cases [3243](https://github.com/keptn/keptn/issues/3243)
- *Fixed*: Flaky integration tests: Integration tests fail (in unpredictable situations) [2149](https://github.com/keptn/keptn/issues/2149)
- *Fixed*: Integration test stalls at the Keptn auth command [2704](https://github.com/keptn/keptn/issues/2704)
- *Fixed*: Integration test: Setup of Keptn fails due to server version check [2701](https://github.com/keptn/keptn/issues/2701)
- *Fixed*: Unable to do remote debugging of mongodb-datastore due to liveness-probe [2536](https://github.com/keptn/keptn/issues/2536)
- *Fixed*: GitHub Action Reviewdog Fails: The `add-path` command is disabled [2694](https://github.com/keptn/keptn/issues/2694)

</p>
</details>

## Good to know / Known Limitations

This section lists bugs and limitations that are known but not fixed in this release. They will get addressed in one of the next releases.

- Keptn CLI can not be used for automation due to Kube context check [3208](https://github.com/keptn/keptn/issues/3208)
  - The workaround is explained [here](https://github.com/keptn/keptn/issues/3208#issuecomment-781982765)
- Creating a project fails on OpenShift due to missing write permissions [2453](https://github.com/keptn/keptn/issues/2453)
- Hovering over the score in an `approval.triggered` events in the Bridge leads to a scroll-up / jump-up in Firefox [#2369](https://github.com/keptn/keptn/issues/2369)
- Remove the functionality to listen to `sh.keptn.event.service.delete.finished` event from helm-service [2989](https://github.com/keptn/keptn/issues/2989)
  - The helm-service does not support listening on `sh.keptn.event.service.delete.finished` events when running on the execution plane. This leads in the limitation that deleting deployed services on the execution plan becomes a manual task. To delete a deployed service, execute: 
    ```
    helm ls -n <NAMESPACE>
    helm delete <HELM_RELEASE> -n <NAMESPACE>
    ```

## Upgrade to 0.8.0

- The upgrade from Keptn 0.7.3 to 0.8.0 is supported. Please find the documentation here: [Upgrade from Keptn 0.7.3 to 0.8.0](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-7-to-0-8)