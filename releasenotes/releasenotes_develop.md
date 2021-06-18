# Release Notes 0.8.4

*Gain visibility into integrations connected to Keptn* - Keptn 0.8.4 starts to implement the Uniform mockup that has been released/presented with Keptn 0.8.1. This new Bridge screen brings insights into the Keptn-services (aka Integrations) that are connected to a Keptn control-plane, allows troubleshooting by retrieving their error logs, and enables creating/deleting secrets for integrations.

---

**Key announcements:**

:tada: *Troubleshooting support for Integrations*: To support troubleshooting integrations without connecting to the environment that runs them, errors are sent to Keptn and displayed in the Uniform screen of a project.

:star: *Creating/Deleting Secrets for Integrations*: To not rely on the Keptn CLI to manage secrets for integrations, Bridge allows creating/deleting secrets. This is supported for integrations that are running on a Keptn control-plane since the public Keptn API does yet not allow querying secrets.  

:rocket: *Customization of Bridge*: With this release, Keptn Bridge can get a custom *look-and-feel* by providing a custom logo, title, and/or stylesheet. More details on this feature are available [here](https://github.com/keptn/keptn/tree/0.8.4/bridge#custom-look-and-feel).

---

*Note*: If you are a maintainer of an Integration that is hosted on [github.com/keptn-contrib](https://github.com/keptn-contrib) or [github.com/keptn-sandbox](https://github.com/keptn-sandbox), you will receive an issue explaining how to upgrade your integration; especially, the distributor. With this upgrade, your integration will then benefit from the new feature and will be displayed in Keptn Bridge.  

## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 45](https://github.com/keptn/enhancement-proposals/pull/45) & [KEP 46](https://github.com/keptn/enhancement-proposals/pull/46) 

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

- Log ingest endpoint for a Keptn-Service [4032](https://github.com/keptn/keptn/issues/4032)
- List all secrets created by secret-service [4061](https://github.com/keptn/keptn/issues/4061)
- Register/Unregister endpoint for registering a Keptn-service that connects to Keptn control-plane [4041](https://github.com/keptn/keptn/issues/4041)

</p>
</details>

<details><summary>CLI</summary>
<p>

- `keptn upgrade`: Improve help messages [3479](https://github.com/keptn/keptn/issues/3479)
- Replace `exechelper.ExecuteCommand` with `keptnutils.ExecuteCommand` [4068](https://github.com/keptn/keptn/issues/4068)
- *Fixed*: Keptn configure bridge output shows error after disabling basic auth [4154](https://github.com/keptn/keptn/issues/4154)
- *Fixed*: Trying to install a different keptn version on the cluster results in error [3959](https://github.com/keptn/keptn/issues/3959)
- *Fixed*: `keptn upgrade` getLatestKeptnRelease returns the wrong version [3841](https://github.com/keptn/keptn/issues/3841)
- *Fixed*: `keptn generate support-archive` not working on windows [4225](https://github.com/keptn/keptn/issues/4225)
- *Fixed*: `keptn uninstall` does not have any effect on cluster [3958](https://github.com/keptn/keptn/issues/3958) 

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *general*:
  - `shkeptnspecversion` missing in many Keptn CloudEvents [3408](https://github.com/keptn/keptn/issues/3408)

- *distributor*:
  - Forward log messages of execution plane Keptn-services to Keptn core [4030](https://github.com/keptn/keptn/issues/4030)
  - Send data of subscribed Keptn-services (via distributors) to uniform [4031](https://github.com/keptn/keptn/issues/4031)

- *helm-service*: 
  - *Fixed*: Helm-service lost its resource requests/limits [4250](https://github.com/keptn/keptn/issues/4250)

- *lighthouse-service*:
  - *Fixed*: "Response time degradation in lighthouse-service" when spamming get-sli-events [4065](https://github.com/keptn/keptn/issues/4065)

- *remediation-service*:
  - *Fixed*: Remediation-service lost fallback to `problem type: default` [4254](https://github.com/keptn/keptn/issues/4254)

- *shipyard-controller*: 
  - Define Uniform backend data model [4033](https://github.com/keptn/keptn/issues/4033)
  - *Fixed*: Keptn 0.8.3 shows that it uses specversion 0.2.1 instead of 0.2.2 [4192](https://github.com/keptn/keptn/issues/4192)
  - *Fixed*: Shipyard-controller keeps sending events for tasks with the same name indefinitely [4039](https://github.com/keptn/keptn/issues/4039)


</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - List, create and delete Secrets [4062](https://github.com/keptn/keptn/issues/4062)
  - Bridge downloads and uses customized look and feel on startup [4095](https://github.com/keptn/keptn/issues/4095)
  - Environment layout improvement for service versions [4006](https://github.com/keptn/keptn/issues/4006)
  - Show *Uniform screen* with data fetched from Uniform Backend [4034](https://github.com/keptn/keptn/issues/4034)
  - Improve status information in *Service screen* for failed deployments [4002](https://github.com/keptn/keptn/issues/4002)
  - Show instructions or link for triggering evaluations in stage [4055](https://github.com/keptn/keptn/issues/4055)
  - Mark currently selected stage using a color [3948](https://github.com/keptn/keptn/issues/3948)
  - Update *Service screen* on a regular basis [4049](https://github.com/keptn/keptn/issues/4049)
  - Display running remediations in the *Service screen* [3761](https://github.com/keptn/keptn/issues/3761)

- *Fixes:*
  - Bridge shows `Configure monitoring succeeded`, although dynatrace-service responded with result fail [4073](https://github.com/keptn/keptn/issues/4073)
  - Bridge breaks on "sh.keptn.event.evaluation.triggered" root event [4155](https://github.com/keptn/keptn/issues/4155)
  - Timelines show the wrong selection color for a running stage [4262](https://github.com/keptn/keptn/issues/4262)
  - Bridge runs version check although ENABLE_VERSION_CHECK env is set to "false" [4165](https://github.com/keptn/keptn/issues/4165)
  - Incorrect sequence filter if project is changed or the page is reloaded [4151](https://github.com/keptn/keptn/issues/4151)
  - Evaluation result can be viewed from Sequence but not from *Service screen* [4056](https://github.com/keptn/keptn/issues/4056)
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

- Integration Test for Uniform and Log Ingest Feature [4103](https://github.com/keptn/keptn/issues/4103)
- For non-release-builds, use timestamps for containers in Helm Charts [4191](https://github.com/keptn/keptn/issues/4191)
- Integration test: create an issue if integration tests on master branch are failing [3772](https://github.com/keptn/keptn/issues/3772)
- Reduce dependabot to only post PRs once a week [4076](https://github.com/keptn/keptn/issues/4076)
- Selenium E2E tests for Bridge [4142](https://github.com/keptn/keptn/issues/4142)
- Introduce uitestid-s in Bridge [4038](https://github.com/keptn/keptn/issues/4038)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - Shipyard-controller and Bridge run into errors when using an `image` object for a configuration change [4348](https://github.com/keptn/keptn/issues/4348)
  - Installing/Upgrading Keptn in an air-gapped environment does not work for `configuration-service` and `nats` [4183](https://github.com/keptn/keptn/issues/4183)
  - Selected service is not reset on project change [4166](https://github.com/keptn/keptn/issues/4166)
  - *Response time degradation in configuration-service* when using a Git Upstream (e.g., GitHub) [4066](https://github.com/keptn/keptn/issues/4066)
  - Mongodb OOM crash after flooding it with events [3968](https://github.com/keptn/keptn/issues/3968)
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
 
</p>
</details>

## Upgrade to 0.8.4

- The upgrade from 0.8.x to 0.8.4 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.4](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-3-to-0-8-4)
