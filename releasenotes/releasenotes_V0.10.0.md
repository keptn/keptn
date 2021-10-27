# Release Notes 0.10.0

Keptn 0.10.0 provides a native way for integrating your tooling by just calling their Webhooks. This is a great enabler for various delivery and operational use cases that can be implemented without writing custom code. Just two steps and your tool is integrated: (1) define the sequence task that works as a trigger (2) define the HTTP request endpoint and payload of the Webhook:

![webhook](https://user-images.githubusercontent.com/729071/136449846-756723c5-e42f-4699-8121-e3255754a117.png)

---

**Key announcements:**

:tada: *Seamless integration of DevOps landscape* using Webhooks: This release is a major step towards the seamless integration of DevOps tooling for your continuous delivery or operational use cases. Therefore, Keptn 0.10 ships a webhook-service with Keptn core that allows the call of external tools using HTTP. To customize this HTTP request, the Bridge provides the corresponding interface and the secret management has been extended in this regard.

:star: *Create/Delete a service* via Bridge: Next to the Webhook configuration, the Bridge allows creating and deleting a service.

:gift: Our [new integrations page](https://keptn.sh/docs/integrations/) got a facelift a let's you explore and search available Keptn integrations. All powered by the ArtifactHub.

:information_source: Keptn provides an internal Git repository for each Keptn project regardless of whether a Git upstream is configured. This internal Git repository will become deprecated in an upcoming Keptn release; more detail will follow. Consequently, it is recommended to set a Git upstream to your own, publically accessible Git repository today. Therefore, use the Keptn [CLI](https://keptn.sh/docs/0.10.x/reference/cli/commands/keptn_update_project/) or [Bridge](https://keptn.sh/docs/0.10.x/reference/bridge/manage_projects/). If there are specific requirements to connect to an own repository, please reach out on Slack: [keptn.slack.com](https://keptn.slack.com)

---

## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 61](https://github.com/keptn/enhancement-proposals/pull/61) and parts of [KEP 48](https://github.com/keptn/enhancement-proposals/pull/48), [KEP 53](https://github.com/keptn/enhancement-proposals/pull/53), and [KEP 54](https://github.com/keptn/enhancement-proposals/pull/54)

## Keptn Specification

Implemented **Keptn spec** version: [0.2.3](https://github.com/keptn/spec/tree/0.2.3)

## New Features

<details><summary>Keptn Core</summary>
<p>

- *configuration-service*:
  - *Deprecated*: GET default resources endpoints: `/project/{projectName}/service/{serviceName}/resource` [#5443](https://github.com/keptn/keptn/issues/5443)
  - Make sure upstream changes are pulled when updating upstream creds [#5224](https://github.com/keptn/keptn/issues/5224)
  - Implemented endpoints for deleting service and stage resources [#5145](https://github.com/keptn/keptn/issues/5145)
  - Handle error and use dedicated HTTP error code when failing to update project due to wrong token [#5438](https://github.com/keptn/keptn/issues/5438)
  - Fall back to previous git credentials when updating upstream fails [#5171](https://github.com/keptn/keptn/issues/5171)
  - *Fix* updating upstream to uninitialized repo [#5569](https://github.com/keptn/keptn/issues/5569)

- *distributor*:
  - Ensure that the subscriptionId is passed to the event [#5412](https://github.com/keptn/keptn/issues/5412)
  - Pass along subscriptionId to service implementation [#5374](https://github.com/keptn/keptn/issues/5374)
  - Exclusive message processing for multiple distributors [#5249](https://github.com/keptn/keptn/issues/5249)
  - Only interpret events with status=errored as error logs [#5186](https://github.com/keptn/keptn/issues/5186)
  - Hardening of ce cache [#5736](https://github.com/keptn/keptn/issues/5736)
  - *Fixed:* Leaking go routines in forwarder.go [#5404](https://github.com/keptn/keptn/issues/5404)
  - *Fixed:* Fails when having no initial PubSub topic defined [#5230](https://github.com/keptn/keptn/issues/5230)
  - *Fixed:* Potential timing issue in distributor unit tests [#5538](https://github.com/keptn/keptn/issues/5538)
  - *fixed:* Send event once for each matching subscription [#5681](https://github.com/keptn/keptn/issues/5681)

- *helm-service*:
  - Customize Helm Chart image pull registry & pull secrets [#4984](https://github.com/keptn/keptn/issues/4984)
  - Revert upgrade to helm v3.7.0 because of memory issues [#5588](https://github.com/keptn/keptn/issues/5588)
  - Increase resource limits to avoid OOM crashes [#5572](https://github.com/keptn/keptn/issues/5572)
  - *Fixed:* Use `user_managed` instead of `user-managed` [#3624](https://github.com/keptn/keptn/issues/3624)

- *jmeter-service*:
  - Prevent failure if deploymentURIs does not end with a '/' [#3612](https://github.com/keptn/keptn/issues/3612)
  - Implement a retry loop for `checkEndpointAvailability` [#5619](https://github.com/keptn/keptn/issues/5619)

- *lighthouse-service*:
  - Calcscore missing error msg [#5252](https://github.com/keptn/keptn/issues/5252)
  - Added error logs for failing monitoring configuration [#5220](https://github.com/keptn/keptn/issues/5220)
  - Add message to event in case SLO parsing failed [#5135](https://github.com/keptn/keptn/issues/5135)
  - *Fixed:* Check for `nil` entries in SLO objectives [#5522](https://github.com/keptn/keptn/issues/5522)
  - *Fixed:* Return the wrong error message if it fails to read slo.yaml [#5549](https://github.com/keptn/keptn/issues/5549)

- *mongodb-datastore*:
  - Added dedicated GET endpoint for readiness probe [#5499](https://github.com/keptn/keptn/issues/5499)
  - Provide option to connect to external MongoDB [#5385](https://github.com/keptn/keptn/issues/5385)
  - Increase memory limits for mongodb-datastore and mongodb [#5197](https://github.com/keptn/keptn/issues/5197)
  - Correct log level for storing root events [#5075](https://github.com/keptn/keptn/issues/5075)
  - *Fixed:* mongodb-datastore resource requests and limits for skaffold setup [#5202](https://github.com/keptn/keptn/issues/5202)

- *remediation-service*:
  - Adapt to recent changes in go SDK [#5464](https://github.com/keptn/keptn/issues/5464)

- *shipyard-controller*:
  - Allow to abort queued sequences [#5472](https://github.com/keptn/keptn/issues/5472)
  - Reduce log noise for sequence watcher component [#5458](https://github.com/keptn/keptn/issues/5458)
  - Remove log noise in sequence migrator [#5096](https://github.com/keptn/keptn/issues/5096)
  - More robust handling of multiple `.started`/`.finished` events for the same task at the same time [#5440](https://github.com/keptn/keptn/issues/5440)
  - Adapted sequence state representation when sequence can not be started [#5194](https://github.com/keptn/keptn/issues/5194)
  - Return proper error message in case project is not available [#5231](https://github.com/keptn/keptn/issues/5231)
  - Return error if a sequence for an unavailable stage is triggered [#5069](https://github.com/keptn/keptn/issues/5069)
  - Adapted log output when no queued sequence is found [#5167](https://github.com/keptn/keptn/issues/5167)
  - Adapted HTTP status codes of GET /event endpoint [#5134](https://github.com/keptn/keptn/issues/5134)
  - Avoid endless loop [#5124](https://github.com/keptn/keptn/issues/5124)
  - Clean up list of open `.triggered` events when completing a sequence [#5601](https://github.com/keptn/keptn/issues/5601)
  - Correctly handle time format in evaluation manager [#5633](https://github.com/keptn/keptn/issues/5633)
  - Ensure list of open `.triggered` events is cleaned up when deleting project [#5502](https://github.com/keptn/keptn/issues/5502)
  - Use timestamp of incoming events to queue sequences [#5620](https://github.com/keptn/keptn/issues/5620)
  - Check for existence of stages in shipyard.yaml when creating a project [#5774](https://github.com/keptn/keptn/issues/5774)
  - *Fixed:* Dependency incompatibilities [#5127](https://github.com/keptn/keptn/issues/5127)
  - *Fixed:* Evaluation score should be computed based only on lighthouse events [#5640](https://github.com/keptn/keptn/issues/5640)

- *secret-service*:
  - Creation of RoleBinding based on scope name [#5300](https://github.com/keptn/keptn/issues/5300)
  - Add list of keys within secrets created by the secret-service [#5139](https://github.com/keptn/keptn/issues/5139)
  - *Fixed:* Correct HTTP status code for invalid key or name [#5479](https://github.com/keptn/keptn/issues/5479)

- *webhook-service*:
  - Introduced webhook-service in Keptn core [#4938](https://github.com/keptn/keptn/issues/4938)
  - Additional curl command validation to increase security [#5500](https://github.com/keptn/keptn/issues/5500)
  - Allow to disable sending the finished event in the webhook-service [#5418](https://github.com/keptn/keptn/issues/5418)
  - Filter Webhooks based on received subscription ID [#5392](https://github.com/keptn/keptn/issues/5392)
  - Allow to control if webhook-service is installed [#5574](https://github.com/keptn/keptn/issues/5574)
  - Add required scope to secret created for webhook integration test [#5594](https://github.com/keptn/keptn/issues/5594)
  - Allow to control if the webhook-service is installed [#5556](https://github.com/keptn/keptn/issues/5556)
</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - Initial integration tests [#5360](https://github.com/keptn/keptn/issues/5360)
  - Make session cookie timeout configurable and set default value to 60 minutes [#5455](https://github.com/keptn/keptn/issues/5455)
  - Align the way how sequence states are displayed [#5376](https://github.com/keptn/keptn/issues/5376)
  - Evaluation board only updates if there are new evaluations [#5396](https://github.com/keptn/keptn/issues/5396)
  - Create secret with scope selection [#5388](https://github.com/keptn/keptn/issues/5388)
  - Set latest sequence depending on the latest event [#5148](https://github.com/keptn/keptn/issues/5148)
  - Include time zone for `trigger evaluation` command [#5398](https://github.com/keptn/keptn/issues/5398)
  - Handle incorrect remediation sequences [#5383](https://github.com/keptn/keptn/issues/5383)
  - Remove HeatMap selection if deployment-sequence does not have an evaluation [#4636](https://github.com/keptn/keptn/issues/4636)
  - Show a gray thick border when a running sequence is selected [#5141](https://github.com/keptn/keptn/issues/5141)
  - Configure webhook-service in Bridge [#4750](https://github.com/keptn/keptn/issues/4750)
  - Load sequence with more than 100 events correctly [#5308](https://github.com/keptn/keptn/issues/5308)
  - Show proper error messages if not OAuth is configured and prevent login loop [#5086](https://github.com/keptn/keptn/issues/5086)
  - Grouping sequence after pause [#5275](https://github.com/keptn/keptn/issues/5275)
  - Show list of files and link to git repo per stage for a service [#5193](https://github.com/keptn/keptn/issues/5193)
  - Set empty array when open remediations are not a sequence [#5217](https://github.com/keptn/keptn/issues/5217)
  - Delete a service [#4380](https://github.com/keptn/keptn/issues/4380)
  - Create a service [#4500](https://github.com/keptn/keptn/issues/4500)
  - Show loading bar only on initial data fetch #4910 [#5586](https://github.com/keptn/keptn/issues/5586)
  - Show loading indicator in environment screen until data is fetched [#5417](https://github.com/keptn/keptn/issues/5417)
  - Show payload of last event in subscription configuration [#5585](https://github.com/keptn/keptn/issues/5585)
  - Make all project tiles same height [#5577](https://github.com/keptn/keptn/issues/5577)
  - Tooltips for heatmap [#4523](https://github.com/keptn/keptn/issues/4523)
  - Dynamically set SLI button positions [#5416](https://github.com/keptn/keptn/issues/5416)
  - Support also clone urls for creating the git repo link [#5391](https://github.com/keptn/keptn/issues/5391)
  - Update webhook with right subscription property, fix stuck subscription update [#5582](https://github.com/keptn/keptn/issues/5582)
  - Heatmap did not correctly change on stage change [#5578](https://github.com/keptn/keptn/issues/5578)
  - Allow multiple webhooks with same subscription configuration [#5267](https://github.com/keptn/keptn/issues/5267)
  - Add secrets to webhook configuration [#4751](https://github.com/keptn/keptn/issues/4751)
  - Validate secret name length [#5478](https://github.com/keptn/keptn/issues/5478)
  - Add ability to configure feature flags [#5211](https://github.com/keptn/keptn/issues/5211)

- *Refactoring:*
  - Removed deprecated links [#4612](https://github.com/keptn/keptn/issues/4612)
  - Code style fixes [#4648](https://github.com/keptn/keptn/issues/4648)
  - Migration to ESLint [#4648](https://github.com/keptn/keptn/issues/4648)
  - IDE ESLint setup [#4648](https://github.com/keptn/keptn/issues/4648)
  - Adapt retry-mechanism [#4867](https://github.com/keptn/keptn/issues/4867)
  - Add cypress setup [#5190](https://github.com/keptn/keptn/issues/5190)


- *Fixes:*
  - 'Show SLO' button disappeared after loading evaluation results [#5393](https://github.com/keptn/keptn/issues/5393)
  - Project settings page styles[#5444](https://github.com/keptn/keptn/issues/5444)
  - Task retrieval if shipyard does not contain any sequences [#5409](https://github.com/keptn/keptn/issues/5409)
  - Shipyard file selection, if the same file was chosen again [#5380](https://github.com/keptn/keptn/issues/5380)
  - Redirect to login page if OAuth is configured [#5370](https://github.com/keptn/keptn/issues/5370)
  - Fixed missing update on sequence screen [#5085](https://github.com/keptn/keptn/issues/5085)
  - Fixed error if sequence was not found [#5172](https://github.com/keptn/keptn/issues/5172)
  - Project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
  - Polling of a project did not stop [#5094](https://github.com/keptn/keptn/issues/5094)
  - Faded-out integrations were not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)
  - Redirect to service or sequence did not work on dashboard [#5126](https://github.com/keptn/keptn/issues/5126)
  - Project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
  - Faded-out integrations where not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)
  - Fixed SLI compared value [#5460](https://github.com/keptn/keptn/issues/5460)
  - Fixed missing view updates when sending an approval [#5505](https://github.com/keptn/keptn/issues/5505)
  - Service incorrectly shows that there are open remediations [#5688](https://github.com/keptn/keptn/issues/5688)
  - Catch error only in interceptor and show toast [#5213](https://github.com/keptn/keptn/issues/5213)
</p>
</details>


<details><summary>Platform Support / Installer</summary>
<p>
 - Temporarily revert customization of repository string in chart [#5414](https://github.com/keptn/keptn/issues/5414)
 - Add option for Ingress to control-plane Helm Chart Keptn installer [#5066](https://github.com/keptn/keptn/issues/5066)
 - Use correct images in airgapped installation [#5532](https://github.com/keptn/keptn/issues/5532)
 - Bump nginx image version to 1.21.3-alpine [#5564](https://github.com/keptn/keptn/issues/5564)
 - *Fix* bug where OpenShift route service go-utils were not upgraded during auto upgrade
</p>
</details>


<details><summary>CLI</summary>
<p>
 - Added zones to times format according to (ISO8601) [#4788](https://github.com/keptn/keptn/issues/4788)
 - Check if kubectl context matches Keptn CLI context before applying upgrade [#5250](https://github.com/keptn/keptn/issues/5250)
 - Skip version check on install [#5046](https://github.com/keptn/keptn/issues/5046)
 - Remove the upgrade available message while upgrading Keptn [#5276](https://github.com/keptn/keptn/issues/5276)
 - Configure automatic version check based on config [#5290](https://github.com/keptn/keptn/issues/5290)
 - Option to continue install/upgrade if K8s version is higher than the supported one [#5698](https://github.com/keptn/keptn/issues/5698)
</p>
</details>


<details><summary>API</summary>
<p>
 - Try to use X-real-ip and X-forwarded-for headers [#5082](https://github.com/keptn/keptn/issues/5082)
 - *Fixed* broken go-sum in go-sdk module [#5463](https://github.com/keptn/keptn/issues/5463)
 - Option to disable automatic event response in SDK [#5453](https://github.com/keptn/keptn/issues/5453)
</p>
</details>


## Development Process / Testing

- *Fixed* paths in commit messages [#5451](https://github.com/keptn/keptn/issues/5451)
- *Fixed* integration tests [#5390](https://github.com/keptn/keptn/issues/5390)
- Added retry mechanism for creating projects in integration tests [#5253](https://github.com/keptn/keptn/issues/5253)
- Updated go-dependencies in integration tests [#5205](https://github.com/keptn/keptn/issues/5205)
- Add disclamer to avoid security vulnerabilities to be reported reported as bugs [#5169](https://github.com/keptn/keptn/issues/5169)
- Update Maintainers file [#5314](https://github.com/keptn/keptn/issues/5314)


## Good to know / Known Limitations

- Aborting a pending deployment sequence in helm-service leads to failure until the aborted sequence finally finishes [#5557](https://github.com/keptn/keptn/issues/5557)
- The following characters/strings are forbitten in the WebHook payload: `$`, `|`, `;`, `>`, `$(`, ` &`, `&&`, \`, `/var/run`


## Upgrade to 0.10.0

- The upgrade from 0.9.x to 0.10.0 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.9.x to 0.10.0](https://keptn.sh/docs/0.10.x/operate/upgrade/#upgrade-from-keptn-0-9-x-to-0-10-0)