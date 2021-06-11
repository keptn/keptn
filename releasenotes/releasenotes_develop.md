# Release Notes develop

Keptn 0.8.4
---

**Key announcements:**

:tada: 

:star: 

:rocket: 

---

## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 46](https://github.com/keptn/enhancement-proposals/pull/46) & [KEP 45](https://github.com/keptn/enhancement-proposals/pull/45) 


## Keptn Specification

Implemented **Keptn spec** version: [0.2.3](https://github.com/keptn/spec/tree/0.2.3)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Add readinessProbe to Helm Chart of: keptn, jmeter-service, and helm-service [3648](https://github.com/keptn/keptn/issues/3648)

</p>
</details>

<details><summary>API</summary>
<p>

- List all secrets created by secret-service [4061](https://github.com/keptn/keptn/issues/4061)
- Register/Unregister endpoint for registering a Keptn-service that connects to Keptn control-plane [4041](https://github.com/keptn/keptn/issues/4041)

</p>
</details>

<details><summary>CLI</summary>
<p>

- `keptn upgrade`: Improve help messages [3479](https://github.com/keptn/keptn/issues/3479)
- Replace `exechelper.ExecuteCommand` with `keptnutils.ExecuteCommand` [4068](https://github.com/keptn/keptn/issues/4068)
- *Fixed*: Keptn configure bridge output shows error after disabling basic auth [4154](https://github.com/keptn/keptn/issues/4154)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *helm-service*: 
  - *Fixed*: Helm-service lost its resource requests/limits [4250](https://github.com/keptn/keptn/issues/4250)

- *shipyard-controller*: 
  - Define Uniform backend data model [4033](https://github.com/keptn/keptn/issues/4033)
  - *Fixed*: Keptn 0.8.3 shows that it uses specversion 0.2.1 instead of 0.2.2 [4192](https://github.com/keptn/keptn/issues/4192)
  - *Fixed*: Shipyard-controller keeps sending events for tasks with the same name indefinitely [4039](https://github.com/keptn/keptn/issues/4039)

- *lighthouse-service*:
  - *Fixed*: "Response time degradation in lighthouse-service" when spamming get-sil-events [4065](https://github.com/keptn/keptn/issues/4065)


</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - Environment layout improvement for service versions [4006](https://github.com/keptn/keptn/issues/4006)
  - Show uniform screen with data fetched from Uniform Backend [4034](https://github.com/keptn/keptn/issues/4034)
  - Improve status information in Bridge Service View for failed deployments [4002](https://github.com/keptn/keptn/issues/4002)
  - Show instructions or link for triggering evaluations in stage [4055](https://github.com/keptn/keptn/issues/4055)
  - Mark currently selected stage using a color [3948](https://github.com/keptn/keptn/issues/3948)
  - Update Service screen on a regular basis [4049](https://github.com/keptn/keptn/issues/4049)
  - Display running remediations in the service screen [3761](https://github.com/keptn/keptn/issues/3761)

- *Fixes:*
  - Bridge shows Configure monitoring succeeded, although dynatrace-service responded with result fail [4073](https://github.com/keptn/keptn/issues/4073)
  - Bridge breaks on "sh.keptn.event.evaluation.triggered" root event [4155](https://github.com/keptn/keptn/issues/4155)
  - Timelines show the wrong selection color for a running stage [4262](https://github.com/keptn/keptn/issues/4262)
  - Bridge runs version check although ENABLE_VERSION_CHECK env is set to "false" [4165](https://github.com/keptn/keptn/issues/4165)
  - Incorrect sequence filter if project is changed or the page is reloaded [4151](https://github.com/keptn/keptn/issues/4151)
  - Evaluation result can be viewed from Sequence but not from Service screen [4056](https://github.com/keptn/keptn/issues/4056)
  - Unexpected behavior of scrollbars in environment screen [4149](https://github.com/keptn/keptn/issues/4149)
  - Selection change in heatmap does not always update SLO table - needs second click [4007](https://github.com/keptn/keptn/issues/4007)
  - Environment panels are not updated on approval / finish [4048](https://github.com/keptn/keptn/issues/4048)
  - Sequence is only updated when detail is opened [4130](https://github.com/keptn/keptn/issues/4130)
  - Service tile breaks based on image:tag > `carts:353ff51.1` [4130](https://github.com/keptn/keptn/issues/4130)

</p>
</details>

## Miscellaneous

- Dependency incompatibility in services using helm library [4063](https://github.com/keptn/keptn/issues/4063)
- Add bridge and bridge server to dependabot [4077](https://github.com/keptn/keptn/issues/4077)

## Development Process / Testing

- Reduce dependabot to only post PRs once a week [4076](https://github.com/keptn/keptn/issues/4076)
- Selenium E2E tests for Bridge [4142](https://github.com/keptn/keptn/issues/4142)
- Introduce uitestid-s in Bridge [4038](https://github.com/keptn/keptn/issues/4038)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - Remediation-service lost fallback to `problem type: default` [4254](https://github.com/keptn/keptn/issues/4254)
  - Installing/Upgrading Keptn in an air-gapped environment does not work for `configuration-service` and `nats` [4183](https://github.com/keptn/keptn/issues/4183)
  - Selected service is not reset on project change [4166](https://github.com/keptn/keptn/issues/4166)
  - *Response time degradation in configuration-service* when using a Git Upstream (e.g., GitHub) [4066](https://github.com/keptn/keptn/issues/4066)
  - *Response time degradation in lighthouse-service* when spamming get-sli-events [4065](https://github.com/keptn/keptn/issues/4065)
  - Mongodb OOM crash after flooding it with events [3968](https://github.com/keptn/keptn/issues/3968)
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
 
</p>
</details>

## Upgrade to 0.8.4

- The upgrade from 0.8.x to 0.8.4 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.4](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-3-to-0-8-4)
