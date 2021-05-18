# Release Notes 0.8.3

Keptn 0.8.3 
---

**Key announcements:**


:tada: *Closed-loop remediation modelled via Shipyard v2*: 

:star: *Length of service name increased*: 

:rocket: *Support custom deployment URLs*: 

---

Many thanks to the community for the enhancements on this release! 
 
## Keptn Specification

Implemented **Keptn spec** version: [0.2.2](https://github.com/keptn/spec/tree/0.2.2)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>


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
  - Support custom deployment URLs for user manged deployments [3757](https://github.com/keptn/keptn/issues/3757)

- *jmeter-service*:
  - *Fixed*: JMeter service doesnt work for regular http/https URL as it infer the default http/https port from the URL [3916](https://github.com/keptn/keptn/issues/3916)

- *ligthouse-service*:
  - Allow timeframe to be passed via CloudEvent [4079](https://github.com/keptn/keptn/issues/4079)

- *remediation-service*:
  - Use rootCause field instead problemTitle [3755](https://github.com/keptn/keptn/issues/3755)
  - Clean-up remediation-service to not control the remediation sequence [3682](https://github.com/keptn/keptn/issues/3682)

- *shipyard-controller*:
  - Implement looping-mechanism via shipyard-controller [3683](https://github.com/keptn/keptn/issues/3683)
  - *Fixed*: Shipyard-controller `/v1/event` endpoint - Response time degredation [3962](https://github.com/keptn/keptn/issues/3962)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- No empty panel when no evaluation happened for a service in a stage [3941](https://github.com/keptn/keptn/issues/3941)
- Switching projects without loosing context [3944](https://github.com/keptn/keptn/issues/3944)
- Keptn bridge support for error pages [3925](https://github.com/keptn/keptn/issues/3925)
- Bridge index.html should not be delivered with a 7 day cache header [3876](https://github.com/keptn/keptn/issues/3876)
- Save service filter in environment screen project specific [3994](https://github.com/keptn/keptn/issues/3994)
- Provide better navigation from full screen evaluation screen [3538](https://github.com/keptn/keptn/issues/3538)
- Set Git upstream URL via Settings page [3417](https://github.com/keptn/keptn/issues/3417)
- Service screen: Show stages the deployment went through [3713](https://github.com/keptn/keptn/issues/3713)
- Environment screen: Click on sequence opens Sequence screen [3887](https://github.com/keptn/keptn/issues/3887)
- Environment screen: Click on deployment opens Service screen [3760](https://github.com/keptn/keptn/issues/3760)
- Environment screen: Filter for Service(s) [3759](https://github.com/keptn/keptn/issues/3759)
- Bridge Menu: Use icons instead of text labels [3643](https://github.com/keptn/keptn/issues/3643)
- OAuth/OpenID Connect based login for Keptn bridge [3448](https://github.com/keptn/keptn/issues/3448)
- *Fixed*: Bridge shows ${this.currentTime} instead of current time [3961](https://github.com/keptn/keptn/issues/3961)
- *Fixed*: Service-filter and stage-details do not reset on project-change [3993](https://github.com/keptn/keptn/issues/3993)
- *Fixed*: Incorrect heatmap width when switching chart type [3851](https://github.com/keptn/keptn/issues/3851)
- *Fixed*: The notification messages in Bridge duplicate when the version check toggle is updated [3896](https://github.com/keptn/keptn/issues/3896)
- *Fixed*: Show currently selected project [3912](https://github.com/keptn/keptn/issues/3912)
- *Fixed*: Bridge in Quality gates only use-case breaks on same sequence and task name "evaluation" [3927](https://github.com/keptn/keptn/issues/3927)
- *Fixed*: Last Evaluation label is not visible in case of too many evaluations in the chart [3811](https://github.com/keptn/keptn/issues/3811)
- *Fixed*: Chosen project is not selected and disappears after refresh (F5) [3853](https://github.com/keptn/keptn/issues/3853)
- *Fixed*: Switch between the tabs Environment-Services, the expand/collapse icon is changed but Evaluation items remain expanded [3814](https://github.com/keptn/keptn/issues/3814)
- *Fixed*: Bridge shows "started" wording on status.changed [3583](https://github.com/keptn/keptn/issues/3585)



</p>
</details>

## Miscellaneous


## Development Process / Testing

- Reviewdog fails to find .go files (after Go 1.16 upgrade) [4000](https://github.com/keptn/keptn/issues/4000)
- Introduce golangci-lint into build chain [3019](https://github.com/keptn/keptn/issues/3019)
- Deliver Keptn with build on master [3845](https://github.com/keptn/keptn/issues/3845)
- Integration tests: Download-artifact no longer works for PR builds [3897](https://github.com/keptn/keptn/issues/3897)
- Integration tests: Tests for continuous-delivery should fail when delivery fails [3843](https://github.com/keptn/keptn/issues/3843)
- CodeCov Security Issue: Verify if we were affected and take remediation actions [3820](https://github.com/keptn/keptn/issues/3820)
- *Fixed*: Keptn install fails on master [3884](https://github.com/keptn/keptn/issues/3884)
- *Fixed*: Homebrew installed CLI fails install with 'Malformed constraint: ""' [3805](https://github.com/keptn/keptn/issues/3805)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
  - Keptn CLI: Disable Kube context check [3666](https://github.com/keptn/keptn/issues/3666)

</p>
</details>

## Upgrade to 0.8.3

- The upgrade from 0.8.1 to 0.8.2 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.2](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-1-to-0-8-2)