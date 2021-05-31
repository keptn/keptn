# Release Notes 0.8.3

Keptn 0.8.3 implements the Keptn Enhancement Proposal [#37](https://github.com/keptn/enhancement-proposals/pull/37) for allowing a Keptn user to customize the remediation sequence in the Shipyard. Besides, improvements for the user experience in the Bridge are implemented like setting the Git upstream repository or linking various screens. 

---

**Key announcements:**

:tada: *Customization of auto-remediation sequences*: With this release, it is possible to customize the remediation sequences, which take care of resolving an open problem for a service. Therefore, the remediation sequence can be modeled in the Shipyard for a specific stage. Besides, it is possible to let the *action-providers* run on an execution plane. 

  - :warning: As part of the upgrade process to Keptn 0.8.3 and for utilizing the auto-remediation feature, please manually add the following sequence to the stage that should have auto-remediation enabled and replace the [STAGE-NAME] by the name of the stage you added it to. Without that sequence, no remediation will be triggered for an open problem! Please find here more information on how to upgrade the remediation use-case here: [Update your Shipyard for the Remediation Use-Case](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-2-to-0-8-3) 

  ```
  - name: remediation
    triggeredOn: 
    - event: [STAGE-NAME].remediation.finished
      selector:
        match:
          evaluation.result: fail
    tasks:
    - name: get-action 
    - name: action
    - name: evaluation
      triggeredAfter: "15m"
      properties:
        timeframe: "15m"
  ```

:star: *Length of service name increased to 43 characters*: The limitation of the service name length has been loosened to allow a length of 43 characters.

:rocket: *Support custom deployment URLs*: When deploying custom Helm Charts by using the `user_managed` deployment strategy of the *helm-service*, it is now possible to define a public and/or local deployment URL. Therefore, the file `endpoints.yaml` must be uploaded to the helm folder in the configuration repository. This file has to contain the `deploymentURIsLocal` and/or `deploymentURIsPublic`. For more details, please see the documentation [here](https://keptn.sh/docs/0.8.x/continuous_delivery/deployment_helm/#user-managed-deployments-experimental).

---

Many thanks to the community for the enhancements on this release! 
 
## Keptn Specification

Implemented **Keptn spec** version: [0.2.2](https://github.com/keptn/spec/tree/0.2.2)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- n/a

</p>
</details>

<details><summary>API</summary>
<p>

- Implementation of `/v1/sequence` endpoint [3796](https://github.com/keptn/keptn/issues/3796)
- Sequence endpoint: Allow to filter sequence states by name and status [3991](https://github.com/keptn/keptn/issues/3991)
- *Fixed*: `configure monitoring` not functioning according to spec [3638](https://github.com/keptn/keptn/issues/3638)

</p>
</details>

<details><summary>CLI</summary>
<p>

- Show bridge URL when executing keptn configure bridge --output [3688](https://github.com/keptn/keptn/issues/3688)
- Disable Kube context check [3666](https://github.com/keptn/keptn/issues/3666)
- Remove RemediationTriggered/Started/Finished, add GetActionTriggered/Started/Finished [4084](https://github.com/keptn/keptn/issues/4084)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *all*: 
  - Upgrade CLI and Keptn-services to latest Go release (e.g., go 1.16) [2936](https://github.com/keptn/keptn/issues/2936)
  - Length of service names is too restrictive [3585](https://github.com/keptn/keptn/issues/3585)

- *distributor*:
  - Distributor for remote execution plane needs to handle the case of slowly responding services [3893](https://github.com/keptn/keptn/issues/3893)
  - *Fixed*: Duplicated Helm Deployment.Started/Finished CloudEvents when using helm-service on a remote execution plane [3888](https://github.com/keptn/keptn/issues/3888)
  - *Fixed*: Distributor of shipyard-controller OOM crash after flooding it with events [3969](https://github.com/keptn/keptn/issues/3969)

- *helm-service/jmeter-service*:
  - Helm/JMeter charts do not honour 'remoteControlPlane.api.apiValidateTls: false' in template [3865](https://github.com/keptn/keptn/issues/3865)
  - Support custom deployment URLs for user-managed deployments [3757](https://github.com/keptn/keptn/issues/3757)

- *helm-service*:
  - *Fixed*: Helm service does not listen on sh.keptn.event.rollback.triggered events [4125](https://github.com/keptn/keptn/issues/4125)

- *jmeter-service*:
  - *Fixed*: JMeter service doesn't work for regular http/https URL as it infer the default http/https port from the URL [3916](https://github.com/keptn/keptn/issues/3916)

- *lighthouse-service*:
  - Allow timeframe to be passed via CloudEvent [4079](https://github.com/keptn/keptn/issues/4079)

- *remediation-service*:
  - Use rootCause field instead problemTitle [3755](https://github.com/keptn/keptn/issues/3755)
  - Clean-up remediation-service to not control the remediation sequence [3682](https://github.com/keptn/keptn/issues/3682)

- *shipyard-controller*:
  - Implement looping-mechanism via shipyard-controller [3683](https://github.com/keptn/keptn/issues/3683)
  - *Fixed*: Timestamps of delayed events are not set properly [4096](https://github.com/keptn/keptn/issues/4096)
  - *Fixed*: TriggeredID of `<stage>.<sequence>.finished` events not set properly [4091](https://github.com/keptn/keptn/issues/4091)
  - *Fixed*: Response time degradation at `/v1/event` endpoint [3962](https://github.com/keptn/keptn/issues/3962)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - No empty panel when no evaluation happened for a service in a stage [3941](https://github.com/keptn/keptn/issues/3941)
  - Switching projects without loosing context [3944](https://github.com/keptn/keptn/issues/3944)
  - Keptn bridge support for error pages [3925](https://github.com/keptn/keptn/issues/3925)
  - Bridge index.html should not be delivered with a 7 day cache header [3876](https://github.com/keptn/keptn/issues/3876)
  - Save service filter in environment screen project specific [3994](https://github.com/keptn/keptn/issues/3994)
  - Provide better navigation from full screen evaluation screen [3538](https://github.com/keptn/keptn/issues/3538)
  - Set Git upstream URL via Settings page [3417](https://github.com/keptn/keptn/issues/3417)
  - Service screen: Show stages the deployment went through [3713](https://github.com/keptn/keptn/issues/3713)
  - Environment screen: Support of more advanced staging environments [3647](https://github.com/keptn/keptn/issues/3647)
  - Environment screen: Click on sequence opens Sequence screen [3887](https://github.com/keptn/keptn/issues/3887)
  - Environment screen: Click on deployment opens Service screen [3760](https://github.com/keptn/keptn/issues/3760)
  - Environment screen: Filter for Service(s) [3759](https://github.com/keptn/keptn/issues/3759)
  - Bridge Menu: Use icons instead of text labels [3643](https://github.com/keptn/keptn/issues/3643)
  - OAuth/OpenID Connect based login for Keptn bridge [3448](https://github.com/keptn/keptn/issues/3448)

- *Fixes:*
  - *Fixed*: Bridge shows ${this.currentTime} instead of current time [3961](https://github.com/keptn/keptn/issues/3961)
  - *Fixed*: Service-filter and stage-details do not reset on project-change [3993](https://github.com/keptn/keptn/issues/3993)
  - *Fixed*: Incorrect heatmap width when switching chart type [3851](https://github.com/keptn/keptn/issues/3851)
  - *Fixed*: The notification messages in Bridge duplicate when the version check toggle is updated [3896](https://github.com/keptn/keptn/issues/3896)
  - *Fixed*: Show currently selected project [3912](https://github.com/keptn/keptn/issues/3912)
  - *Fixed*: Bridge in Quality gates only use-case breaks on the same sequence and task name "evaluation" [3927](https://github.com/keptn/keptn/issues/3927)
  - *Fixed*: Last Evaluation label is not visible in case of too many evaluations in the chart [3811](https://github.com/keptn/keptn/issues/3811)
  - *Fixed*: The chosen project is not selected and disappears after refresh (F5) [3853](https://github.com/keptn/keptn/issues/3853)
  - *Fixed*: Switch between the tabs Environment/Services, the expand/collapse icon is changed but Evaluation items remain expanded [3814](https://github.com/keptn/keptn/issues/3814)
  - *Fixed*: Bridge shows "started" wording on status.changed [3583](https://github.com/keptn/keptn/issues/3585)


</p>
</details>

## Miscellaneous


## Development Process / Testing

- Reviewdog fails to find .go files (after Go 1.16 upgrade) [4000](https://github.com/keptn/keptn/issues/4000)
- Introduce golangci-lint into build chain [3019](https://github.com/keptn/keptn/issues/3019)
- Deliver Keptn with build on master [3845](https://github.com/keptn/keptn/issues/3845)
- Integration tests: Download-artifact no longer works for PR builds [3897](https://github.com/keptn/keptn/issues/3897)
- Integration tests: Tests for continuous-delivery should fail when a delivery fails [3843](https://github.com/keptn/keptn/issues/3843)
- CodeCov Security Issue: Verify if we were affected and take remediation actions [3820](https://github.com/keptn/keptn/issues/3820)
- *Fixed*: Keptn install fails on master [3884](https://github.com/keptn/keptn/issues/3884)
- *Fixed*: Homebrew installed CLI fails install with 'Malformed constraint: ""' [3805](https://github.com/keptn/keptn/issues/3805)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - *Response time degradation in configuration-service* when using a Git Upstream (e.g., GitHub) [4066](https://github.com/keptn/keptn/issues/4066)
  - *Response time degradation in lighthouse-service* when spamming get-sli-events [4065](https://github.com/keptn/keptn/issues/4065)
  - Shipyard-controller keeps sending events for tasks with the same name indefinitely [4039](https://github.com/keptn/keptn/issues/4039)
  - Selection change in heatmap does not always update SLO table - needs second click [4007](https://github.com/keptn/keptn/issues/4007)
  - Mongodb OOM crash after flooding it with events [3968](https://github.com/keptn/keptn/issues/3968)
  - `keptn upgrade` getLatestKeptnRelease returns the wrong version [3841](https://github.com/keptn/keptn/issues/3841)
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
 

</p>
</details>

## Upgrade to 0.8.3

- The upgrade from 0.8.x to 0.8.3 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.3](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-2-to-0-8-3)

  - :warning: Please consider adding the *remediation sequence* to a stage for enabling the auto-remediation capabilities of Keptn. The instructions you will find in the upgrade guide: *Upgrade from Keptn 0.8.x to 0.8.3*
