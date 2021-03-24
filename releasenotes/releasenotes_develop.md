# Release Notes 0.8.1

Keptn 0.8.1 improves the user experience of Keptn by allowing custom names for SLOs, showing a stage-wide overview of currently running sequences including quality gate evaluations, and offering an API/CLI support to create secrets (on top of *Kubernetes Secrets*). Besides, this release addresses encountered bugs and issues encountered.

---

**Key announcements:**

:tada: *API/CLI enhancements for creating secrets*:  This release introduces the new feature to create a secret on the Keptn control-plane, which is then stored as [Kubernetes Secret](https://kubernetes.io/docs/concepts/configuration/secret/). Therefore, the Keptn API and CLI provide the required functionality. Please see [keptn create secret](https://keptn.sh/docs/0.8.x/reference/cli/commands/keptn_create_secret/) to learn how to use this feature. 

:star: *Bridge improvements for SLO names and stage overview*: The SLO spec allows adding a `displayName` for an SLO. This name is optional but will be used in the Bridge when available; please see the snippet below. Additionally, the Bridge provides enhancements for the environment screen where an overview of the currently running sequences is given and the evaluation of a quality gate is displayed:

> *add Screenshot here*

<details><summary>Snippet: SLO.yaml file</summary>
<p>

```
...
objectives:
  - sli: "response_time_p95"
    displayName: "Response time P95"
    key_sli: false
    pass:             
      - criteria:
          - "<600"    
    warning:          
      - criteria:
          - "<=800"
    weight: 1
...
```
</p>
</details>

:dizzy: *(UI mockup for the Keptn Uniform) Bridge displays deployed Keptn-services (integrations) and their subscriptions*: In this release, a UI mockup is provided that should provide a look-and-feel on how to display Keptn-services that are connected to the control-plane. To take a look at this mockup, open the Bridge and navigate to: `your.keptn.endpoint/bridge/project/{project_name}/ff-uniform` 

---

Many thanks to the community for the enhancements on this release! 
 
## Keptn Specification

Implemented **Keptn spec** version: [0.2.1](https://github.com/keptn/spec/tree/0.2.1)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Support for Kubernetes 1.20 [3495](https://github.com/keptn/keptn/issues/3495)

</p>
</details>

<details><summary>API</summary>
<p>

- Create/Delete/Update secret using Keptn API/CLI [3465](https://github.com/keptn/keptn/pull/3465)
- *Fixed*: GET services from a stage endpoint requires stage but contains service in path [3456](https://github.com/keptn/keptn/issues/3456)
- *Fixed*: Endpoint is missing path parameter and mismatch between parameter name [3489](https://github.com/keptn/keptn/issues/3489)

</p>
</details>

<details><summary>CLI</summary>
<p>

- `keptn create secret`: Commands for managing secrets [3596](https://github.com/keptn/keptn/pull/3596)
- CLI & Bridge: Automatically determine doc version [2863](https://github.com/keptn/keptn/issues/2863)
- *Fixed*: CLI in alpine docker image not working [3475](https://github.com/keptn/keptn/issues/3475)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *configuration-service*:
  - *Fixed*: Cannot checkout main|master in AWS CodeCommit [3403](https://github.com/keptn/keptn/issues/3403)

- *distributor*:
  - Allow comma-separated list on event filters for distributors [3577](https://github.com/keptn/keptn/issues/3577)

- *helm-service & jmeter-service*: 
  - Add Helm schema validation support for a 'remoteControlPlane.api.hostname' port value [3450](https://github.com/keptn/keptn/issues/3450)
  - Allow helm-service to work without admin permissions [3511](https://github.com/keptn/keptn/issues/3511)

- *shipyard-controller*:
  - *Fixed*: Upgrade Shipyard: shipyardVersion in GET /project response not updated immediately [3384](https://github.com/keptn/keptn/issues/3384)
  - *Fixed*: `deploymentURI` shows up twice in shipyard-controller `test.triggered` event [3449](https://github.com/keptn/keptn/issues/3449)
  - *Fixed*: Fixed errors in swagger definition of shipyard-controller [3530](https://github.com/keptn/keptn/pull/3530)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Mockup to show installed Keptn-services (aka. Uniform) and latest version available [1280](https://github.com/keptn/keptn/issues/1280)
- Show SLI with display name or "smart SLI name" [3345](https://github.com/keptn/keptn/issues/3345)
- Stage tile supports many services [2289](https://github.com/keptn/keptn/issues/2289)
- Show evaluation result on Service tile (next to stage) [3425](https://github.com/keptn/keptn/issues/3425)
- *Fixed*: Navigation with smart linking [3578](https://github.com/keptn/keptn/issues/3578)
- *Fixed*: Approval events sent by bridge should only include approval-related properties [3557](https://github.com/keptn/keptn/issues/3557)
- *Fixed*: Bridge no longer shows link to deployment urls in environment screen [3535](https://github.com/keptn/keptn/issues/3535)
- *Fixed*: Evaluation component in *Service screen* does not show all labels as compared to full screen view [3537](https://github.com/keptn/keptn/issues/3537)
- *Fixed*: Bridge shows empty test events due to wrong order of events (test.started timestamp < test.triggered timestamp) [3435](https://github.com/keptn/keptn/issues/3435) 
- *Fixed*: Bridge does not list failed quality gate evaluations in *Environment screen* [3438](https://github.com/keptn/keptn/issues/3438)
- *Fixed*: Version check failed [3446](https://github.com/keptn/keptn/issues/3446)
- *Fixed*: Approvals are not working [3477](https://github.com/keptn/keptn/issues/3477)

</p>
</details>

## Miscellaneous

- Disable SpellCheck for *_test.go files [3543](https://github.com/keptn/keptn/issues/3543)
- Added spell checker - big thanks to @jsoref :tada: [3234](https://github.com/keptn/keptn/issues/3234)

## Development Process / Testing

- Integration Test: Run jmeter and helm in execution plane (independent namespace) [3214](https://github.com/keptn/keptn/issues/3214)
- Integration Test: Test the "multiple parallel stages" with cont.delivery [2970](https://github.com/keptn/keptn/issues/2970)
- Automatically create GitHub issue with the label "type:critical" if integration tests on master/release branch fail [3166](https://github.com/keptn/keptn/issues/3166)
- *Fixed*: Debugging helm-service is no longer possible [3421](https://github.com/keptn/keptn/issues/3421)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  - Lighthouse-service needs to properly set result, status and message [3412](https://github.com/keptn/keptn/issues/3412)
  - Helm-service is not working parallel when deployed in execution-plane [3427](https://github.com/keptn/keptn/issues/3427)
  - Shipyard-controller: Only last `.finished` event for a task determines further sequence execution [3493](https://github.com/keptn/keptn/issues/3493)
  - Auto-remediation does not work with remote execution plane [3498](https://github.com/keptn/keptn/issues/3498)
  - Quality gate icon in environment screen does not turn red [3592](https://github.com/keptn/keptn/issues/3592)

</p>
</details>

## Upgrade to 0.8.1

- The upgrade from 0.8.0 to 0.8.1 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.0 to 0.8.1](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-0-to-0-8-1)