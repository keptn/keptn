# Release Notes 0.10.0

Keptn 0.10.0 gives you the possibility to integrate Keptn with external services via WebHooks.

---

**Key announcements:**

Webhook
SDK (Experimental)

---

## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 61](https://github.com/keptn/enhancement-proposals/pull/61) and parts of [KEP 48](https://github.com/keptn/enhancement-proposals/pull/48), [KEP 53](https://github.com/keptn/enhancement-proposals/pull/53), and [KEP 54](https://github.com/keptn/enhancement-proposals/pull/54)

## New Features


<details><summary>Keptn Core</summary>
<p>

- *configuration-service*:
  - *Deprecated*: GET default resources endpoints [#5443](https://github.com/keptn/keptn/issues/5443)
  - Make sure upstream changes are pulled when updating upstream creds (#5149) [#5224](https://github.com/keptn/keptn/issues/5224)
  - Implemented endpoints for deleting service and stage resources (#5136) [#5145](https://github.com/keptn/keptn/issues/5145)
  - Handle error and use dedicated http err code when failing to update project due to wrong token [#5438](https://github.com/keptn/keptn/issues/5438)
  - Fall back to previous git credentials when updating upstream fails (#5064) [#5171](https://github.com/keptn/keptn/issues/5171)

- *distributor*:
  - Ensure that the subscriptionID is passed to the event (#5405) [#5412](https://github.com/keptn/keptn/issues/5412)
  - Pass along subscription id to service implementation [#5374](https://github.com/keptn/keptn/issues/5374)
  - Exclusive message processing for multiple distributors (#4689) [#5249](https://github.com/keptn/keptn/issues/5249)
  - Only interpret events with status=errored as error logs (#5170) [#5186](https://github.com/keptn/keptn/issues/5186)
  - Fixed leaking go routines in forwarder.go [#5404](https://github.com/keptn/keptn/issues/5404)
  - fix issue of having no initial pubsub topic defined [#5230](https://github.com/keptn/keptn/issues/5230)

- *helm-service*:
  - Customize helm chart image pull registry & pull secrets [#4984](https://github.com/keptn/keptn/issues/4984)

- *jmeter-service*:
  - Prevent failure if deploymentURIs does not end with a '/' [#3612](https://github.com/keptn/keptn/issues/3612)

- *lighthouse-service*:
  - Calcscore missing error msg (#5142) [#5252](https://github.com/keptn/keptn/issues/5252)
  - Added error logs for failing monitoring configuration (#5088) [#5220](https://github.com/keptn/keptn/issues/5220)
  - Add message to event in case SLO parsing failed (#5130) [#5135](https://github.com/keptn/keptn/issues/5135)

- *mongodb-datastore*:
  - Added dedicated get endpoint for readiness probe [#5499](https://github.com/keptn/keptn/issues/5499)
  - *Fixed:* mongodb-datastore resource requests and limits for skaffold setup [#5202](https://github.com/keptn/keptn/issues/5202)
  - Provide option to connect to external mongodb (#5369) [#5385](https://github.com/keptn/keptn/issues/5385)
  - Increase memory limits for mongodb-datastore and mongodb (#5196) [#5197](https://github.com/keptn/keptn/issues/5197)
  - Correct log level for storing root events [#5075](https://github.com/keptn/keptn/issues/5075)

- *remediation-service*:
  - Adapt to recent changes in go sdk [#5464](https://github.com/keptn/keptn/issues/5464)

- *shipyard-controller*:
  - Reduce log noise for sequence watcher component [#5458](https://github.com/keptn/keptn/issues/5458)
  - More robust handling of multiple .started/finished events for the same task at the same time [#5440](https://github.com/keptn/keptn/issues/5440)
  - Remove log noise in sequence migrator [#5096](https://github.com/keptn/keptn/issues/5096)
  - Adapted sequence state representation for case where sequence can not be started (#5137) [#5194](https://github.com/keptn/keptn/issues/5194)
  - Return proper error message in case project is not available (#4399) [#5231](https://github.com/keptn/keptn/issues/5231)
  - Return error if a sequence for an unavailable stage is triggered (#4791) [#5069](https://github.com/keptn/keptn/issues/5069)
  - Adapted log output when no queued sequence is found (#5138) [#5167](https://github.com/keptn/keptn/issues/5167)
  - Adapted HTTP status codes of GET /event endpoint (#5132) [#5134](https://github.com/keptn/keptn/issues/5134)
  - Avoid endless loop (#5096) [#5124](https://github.com/keptn/keptn/issues/5124)
  - Fixed dependency incompatibilities (#5078) [#5127](https://github.com/keptn/keptn/issues/5127)

- *secret-service*:
  - Added creation of rolebinding based on scope name [#5300](https://github.com/keptn/keptn/issues/5300)
  - Add list of keys within secrets created by the secret-service (#4749) [#5139](https://github.com/keptn/keptn/issues/5139)

- *webhook-service*:
  - Introduced WebHook Service (#4736) [#4938](https://github.com/keptn/keptn/issues/4938)
  - Additional curl command validation to increase security [#5500](https://github.com/keptn/keptn/issues/5500)
  - Allow to disable sending the finished event in the webhook service (#5368) [#5418](https://github.com/keptn/keptn/issues/5418)
  - Filter Webhooks based on received subscription ID (#5264) [#5392](https://github.com/keptn/keptn/issues/5392)


</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Enhancements:*
  - Initial integration tests [#5360](https://github.com/keptn/keptn/issues/5360)
  - Make session cookie timeout configurable and set default value to 60 minutes [#5455](https://github.com/keptn/keptn/issues/5455)
  - Align the way how sequence states are displayed #5150 [#5376](https://github.com/keptn/keptn/issues/5376)
  - Evaluation board only updates if there are new evaluations [#5396](https://github.com/keptn/keptn/issues/5396)
  - Allow to create secret with scope selection #5269 [#5388](https://github.com/keptn/keptn/issues/5388)
  - set latest sequence depending on the latest event [#5148](https://github.com/keptn/keptn/issues/5148)
  - Include time zone for 'trigger evaluation' command [#5398](https://github.com/keptn/keptn/issues/5398)
  - 'Show SLO' button was removed after loading evaluation results [#5393](https://github.com/keptn/keptn/issues/5393)
  - Handle incorrect remediation sequences [#5383](https://github.com/keptn/keptn/issues/5383)
  - Remove heatmap selection if deployment-sequence does not have an evaluation [#4636](https://github.com/keptn/keptn/issues/4636)
  - Show a gray thick border when a running sequence is selected [#5141](https://github.com/keptn/keptn/issues/5141)
  - Configure webhook service in Bridge [#4750](https://github.com/keptn/keptn/issues/4750)
  - Load sequence with more than 100 events correctly #5056 [#5308](https://github.com/keptn/keptn/issues/5308)
  - Show proper error messages if not OAUTH is configured and prevent login loop [#5086](https://github.com/keptn/keptn/issues/5086)
  - Grouping sequence after pause #5154 [#5275](https://github.com/keptn/keptn/issues/5275)
  - Show list of files and link to git repo per stage for a service (#4506) [#5193](https://github.com/keptn/keptn/issues/5193)
  - Set empty array when open remediations are not a sequence [#5217](https://github.com/keptn/keptn/issues/5217)
  - Delete a service [#4380](https://github.com/keptn/keptn/issues/4380)
  - Create a service [#4500](https://github.com/keptn/keptn/issues/4500)

- *Refactoring:*
  - Removed deprecated links [#4612](https://github.com/keptn/keptn/issues/4612)
  - code style fixes [#4648](https://github.com/keptn/keptn/issues/4648)
  - Migrate to ESLint [#4648](https://github.com/keptn/keptn/issues/4648)
  - IDE ESLint setup [#4648](https://github.com/keptn/keptn/issues/4648)
  - Adapt retry-mechanism [#4867](https://github.com/keptn/keptn/issues/4867)
  - Add cypress setup [#5190](https://github.com/keptn/keptn/issues/5190)


- *Fixes:*
  - Project settings page styles(#5382) [#5444](https://github.com/keptn/keptn/issues/5444)
  - Task retrieval if shipyard does not contain any sequences [#5409](https://github.com/keptn/keptn/issues/5409)
  - Shipyard file selection, if the same file was chosen again [#5380](https://github.com/keptn/keptn/issues/5380)
  - Redirect to login page if OAuth is configured [#5370](https://github.com/keptn/keptn/issues/5370)
  - Fixed missing update on sequence screen [#5085](https://github.com/keptn/keptn/issues/5085)
  - Fixed error if sequence was not found [#5172](https://github.com/keptn/keptn/issues/5172)
  - Project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
  - Polling of a project did not stop [#5094](https://github.com/keptn/keptn/issues/5094)
  - Faded-out integrations were not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)
  - Redirect to service or sequence did not work on dashboard [#5126](https://github.com/keptn/keptn/issues/5126)
  - Redirect to service or sequence did not work on dashboard [#5126](https://github.com/keptn/keptn/issues/5126)
  - Project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
  - Faded-out integrations where not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)

</p>
</details>


<details><summary>Platform Support / Installer</summary>
<p>
 - Temporarily revert customization of repository string in chart [#5414](https://github.com/keptn/keptn/issues/5414)
 - Add option for ingress to control-plane helm chart keptn installer [#5066](https://github.com/keptn/keptn/issues/5066)
 - Fix bug where openshift route service go-utils were not upgraded during auto upgrade
</p>
</details>


<details><summary>CLI</summary>
<p>
 - Added zones to times format according to (ISO8601) [#4788](https://github.com/keptn/keptn/issues/4788)
 - Check if kubectl context matches Keptn CLI context before applying upgrade (#4583) [#5250](https://github.com/keptn/keptn/issues/5250)
 - Skip version check on install [#5046](https://github.com/keptn/keptn/issues/5046)
</p>
</details>


<details><summary>API</summary>
<p>
 - Try to use X-real-ip and X-forwarded-for headers [#5082](https://github.com/keptn/keptn/issues/5082)
 - Fixed broken gosum in go-sdk module [#5463](https://github.com/keptn/keptn/issues/5463)
 - Option to disable automatic event response in SDK (#5368) [#5453](https://github.com/keptn/keptn/issues/5453)
</p>
</details>


## Development Process / Testing

- Fixed paths in commit messages [#5451](https://github.com/keptn/keptn/issues/5451)
- Fix integration tests [#5390](https://github.com/keptn/keptn/issues/5390)
- Added retry mechanism for creating projects in integration tests (#5241) [#5253](https://github.com/keptn/keptn/issues/5253)
- Updated go-dependencies in integration tests (#5200) [#5205](https://github.com/keptn/keptn/issues/5205)
- Add disclamer to avoid security vulnerabilities to be reported reported as bugs [#5169](https://github.com/keptn/keptn/issues/5169)
- Update Maintainers file [#5314](https://github.com/keptn/keptn/issues/5314)

## Good to know / Known Limitations

## Upgrade to 0.10.0

- The upgrade from 0.9.x to 0.10.0 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.9.x to 0.10.0](https://keptn.sh/docs/0.10.x/operate/upgrade/#upgrade-from-keptn-0-9-x-to-0-10-0)
