# Release Notes 0.8.1

Keptn 0.8.1

---

**Key announcements:**

:rocket: 

:tada: 

:star: 

:dizzy: 

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

- *Fixed*: GET services from a stage endpoint requires stage but contains service in the path [3456](https://github.com/keptn/keptn/issues/3456)
- *Fixed*: Endpoint is missing path parameter and mismatch between parameter name [3489](https://github.com/keptn/keptn/issues/3489)

</p>
</details>

<details><summary>CLI</summary>
<p>

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
  - *Fixed*: `deploymentURI` shows up twice in shipyard controller test.triggered event [3449](https://github.com/keptn/keptn/issues/3449)
  - *Fixed*: Fixed errors in swagger definition of shipyard-controller [3530](https://github.com/keptn/keptn/pull/3530)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Show SLI with display name or "smart SLI name" [3345](https://github.com/keptn/keptn/issues/3345)
- Stage tile supports many services [2289](https://github.com/keptn/keptn/issues/2289)
- Show evaluation result on Service tile (next to stage) [3425](https://github.com/keptn/keptn/issues/3425)
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
- Create a workflow step to automatically create GitHub issue with the label "type:critical" if CI or integration tests on master/release branch fail [3166](https://github.com/keptn/keptn/issues/3166)
- *Fixed*: Debugging helm-service is no longer possible [3421](https://github.com/keptn/keptn/issues/3421)

## Good to know / Known Limitations
