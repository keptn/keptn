# Release Notes 0.10.0

Keptn 0.10.0 gives you more control over sequence executions and allows creating/deleting a Keptn project in the Bridge.

---

**Key announcements:**



---

## Keptn Enhancement Proposals

This release implements the KEPs: [KEP 61](https://github.com/keptn/enhancement-proposals/pull/61) and parts of [KEP 48](https://github.com/keptn/enhancement-proposals/pull/48).

## New Features


<details><summary>mongodb-datastore</summary>
<p>
 - added dedicated get endpoint for readiness probe [#5499](https://github.com/keptn/keptn/issues/5499)
 - Fix mongodb-datastore resource requests and limits for skaffold setup [#5202](https://github.com/keptn/keptn/issues/5202)
</p>
</details>
<details><summary>general</summary>
<p>
 - Additional curl command validation to increase security [#5500](https://github.com/keptn/keptn/issues/5500)
 - Remove now unneeded spell checking files
 - Provide option to connect to external mongodb (#5369) [#5385](https://github.com/keptn/keptn/issues/5385)
 - fixed paths in commit messages [#5451](https://github.com/keptn/keptn/issues/5451)
 - fixed broken gosum in go-sdk module [#5463](https://github.com/keptn/keptn/issues/5463)
 - allow to disable sending the finished event in the webhook service (#5368) [#5418](https://github.com/keptn/keptn/issues/5418)
 - Option to disable automatic event response in SDK (#5368) [#5453](https://github.com/keptn/keptn/issues/5453)
 - handle error and use dedicated http err code when failing to update project due to wrong token [#5438](https://github.com/keptn/keptn/issues/5438)
 - Fixed leaking go routines in forwarder.go [#5404](https://github.com/keptn/keptn/issues/5404)
 - Filter Webhooks based on received subscription ID (#5264) [#5392](https://github.com/keptn/keptn/issues/5392)
 - fix test [#5390](https://github.com/keptn/keptn/issues/5390)
 - Fix bug where openshift route service go-utils were not upgraded during auto upgrade
 - Update Maintainers file [#5314](https://github.com/keptn/keptn/issues/5314)
 - Revert setup work for automated execution of tutorials [#3770](https://github.com/keptn/keptn/issues/3770)
 - Introduced WebHook Service (#4736) [#4938](https://github.com/keptn/keptn/issues/4938)
 - Added retry mechanism for creating projects in integration tests (#5241) [#5253](https://github.com/keptn/keptn/issues/5253)
 - fix issue of having no initial pubsub topic defined [#5230](https://github.com/keptn/keptn/issues/5230)
 - Added triscon to Adopters
 - Fall back to previous git credentials when updating upstream fails (#5064) [#5171](https://github.com/keptn/keptn/issues/5171)
 - Updated go-dependencies in integration tests (#5200) [#5205](https://github.com/keptn/keptn/issues/5205)
 - Use correct artifact name for mongodb-datastore in semantic PR setup
 - Increase memory limits for mongodb-datastore and mongodb (#5196) [#5197](https://github.com/keptn/keptn/issues/5197)
 - Add Brad McCoy to CONTRIBUTORS.md [#5067](https://github.com/keptn/keptn/issues/5067)
 - add disclamer to avoid security vulnerabilities to be reported reported as bugs [#5169](https://github.com/keptn/keptn/issues/5169)
 - Setup work for automated execution of tutorials [#3770](https://github.com/keptn/keptn/issues/3770)
 - add giovanni.liva as a security contact [#5165](https://github.com/keptn/keptn/issues/5165)
 - fixed dependency incompatibilities (#5078) [#5127](https://github.com/keptn/keptn/issues/5127)
 - customize helm chart image pull registry & pull secrets [#4984](https://github.com/keptn/keptn/issues/4984)
 - Updated dependencies in go-tests
 - Correct log level for storing root events [#5075](https://github.com/keptn/keptn/issues/5075)
</p>
</details>
<details><summary>bridge</summary>
<p>
 - initial integration tests [#5360](https://github.com/keptn/keptn/issues/5360)
 - Make session cookie timeout configurable and set default value to 60 minutes [#5455](https://github.com/keptn/keptn/issues/5455)
 - IDE ESLint setup [#4648](https://github.com/keptn/keptn/issues/4648)
 - Evaluation board only updates if there are new evaluations [#5396](https://github.com/keptn/keptn/issues/5396)
 - code style fixes [#4648](https://github.com/keptn/keptn/issues/4648)
 - Migrate to ESLint [#4648](https://github.com/keptn/keptn/issues/4648)
 - Allow to create secret with scope selection #5269 [#5388](https://github.com/keptn/keptn/issues/5388)
 - set latest sequence depending on the latest event [#5148](https://github.com/keptn/keptn/issues/5148)
 - Fix project settings page styles(#5382) [#5444](https://github.com/keptn/keptn/issues/5444)
 - include time zone for 'trigger evaluation' command [#5398](https://github.com/keptn/keptn/issues/5398)
 - removed deprecated links [#4612](https://github.com/keptn/keptn/issues/4612)
 - fixed task retrieval if shipyard does not contain any sequences [#5409](https://github.com/keptn/keptn/issues/5409)
 - Align the way how sequence states are displayed #5150 [#5376](https://github.com/keptn/keptn/issues/5376)
 - 'Show SLO' button was removed after loading evaluation results [#5393](https://github.com/keptn/keptn/issues/5393)
 - fixed shipyard file selection, if the same file was chosen again [#5380](https://github.com/keptn/keptn/issues/5380)
 - handle incorrect remediation sequences [#5383](https://github.com/keptn/keptn/issues/5383)
 - remove heatmap selection if deployment-sequence does not have an evaluation [#4636](https://github.com/keptn/keptn/issues/4636)
 - adapt retry-mechanism [#4867](https://github.com/keptn/keptn/issues/4867)
 - fixed redirect to login page if OAuth is configured [#5370](https://github.com/keptn/keptn/issues/5370)
 - Show a gray thick border when a running sequence is selected [#5141](https://github.com/keptn/keptn/issues/5141)
 - Configure webhook service in Bridge [#4750](https://github.com/keptn/keptn/issues/4750)
 - load sequence with more than 100 events correctly #5056 [#5308](https://github.com/keptn/keptn/issues/5308)
 - Show proper error messages if not OAUTH is configured and prevent login loop [#5086](https://github.com/keptn/keptn/issues/5086)
 - grouping sequence after pause #5154 [#5275](https://github.com/keptn/keptn/issues/5275)
 - fixed missing update on sequence screen [#5085](https://github.com/keptn/keptn/issues/5085)
 - Add cypress setup [#5190](https://github.com/keptn/keptn/issues/5190)
 - Show list of files and link to git repo per stage for a service (#4506) [#5193](https://github.com/keptn/keptn/issues/5193)
 - Set empty array when open remediations are not a sequence [#5217](https://github.com/keptn/keptn/issues/5217)
 - fixed error if sequence was not found [#5172](https://github.com/keptn/keptn/issues/5172)
 - Delete a service [#4380](https://github.com/keptn/keptn/issues/4380)
 - Create a service [#4500](https://github.com/keptn/keptn/issues/4500)
 - project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
 - polling of a project did not stop [#5094](https://github.com/keptn/keptn/issues/5094)
 - faded-out integrations were not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)
 - Redirect to service or sequence did not work on dashboard [#5126](https://github.com/keptn/keptn/issues/5126)
 - Redirect to service or sequence did not work on dashboard [#5126](https://github.com/keptn/keptn/issues/5126)
 - project delete dialog was not closed [#5091](https://github.com/keptn/keptn/issues/5091)
 - faded-out integrations where not excluded from unread-error-event check [#5118](https://github.com/keptn/keptn/issues/5118)
</p>
</details>
<details><summary>deps</summary>
<p>
 - Auto-update kubernetes-utils to latest version
 - Auto-update go-utils to latest version
 - bump actions/github-script from 4.1 to 5
 - bump actions/setup-node from 2.4.0 to 2.4.1
 - bump JasonEtco/create-an-issue from 2.5.0 to 2.6
 - Update go-utils and kubernetes-utils after history change
 - Auto-update kubernetes-utils to latest version
 - Auto-update go-utils to latest version
 - bump github.com/mitchellh/mapstructure from 1.4.1 to 1.4.2 in /cli
 - bump github.com/go-openapi/strfmt in /mongodb-datastore [#5015](https://github.com/keptn/keptn/issues/5015)
 - bump github.com/go-openapi/errors in /api [#5022](https://github.com/keptn/keptn/issues/5022)
 - bump axios from 0.21.1 to 0.21.4 in /bridge/server [#5178](https://github.com/keptn/keptn/issues/5178)
 - bump github.com/nats-io/nats.go in /distributor [#5176](https://github.com/keptn/keptn/issues/5176)
 - bump go.mongodb.org/mongo-driver in /mongodb-datastore [#5177](https://github.com/keptn/keptn/issues/5177)
 - bump go.mongodb.org/mongo-driver in /statistics-service [#5179](https://github.com/keptn/keptn/issues/5179)
 - bump go.mongodb.org/mongo-driver in /shipyard-controller [#5180](https://github.com/keptn/keptn/issues/5180)
 - bump actions/github-script from 4.0.2 to 4.1
 - bump jwalton/gh-find-current-pr from 1.1.0 to 1.2.0
 - bump k8s.io/kubectl from 0.21.3 to 0.22.1 in /helm-service [#5031](https://github.com/keptn/keptn/issues/5031)
 - bump k8s.io/cli-runtime in /helm-service [#5032](https://github.com/keptn/keptn/issues/5032)
 - bump github.com/go-openapi/strfmt in /configuration-service [#5034](https://github.com/keptn/keptn/issues/5034)
 - bump k8s.io/client-go in /configuration-service [#5035](https://github.com/keptn/keptn/issues/5035)
 - bump github.com/go-openapi/runtime in /api [#5116](https://github.com/keptn/keptn/issues/5116)
 - bump github.com/keptn/kubernetes-utils in /test/go-tests [#5113](https://github.com/keptn/keptn/issues/5113)
 - bump github.com/go-openapi/runtime [#5117](https://github.com/keptn/keptn/issues/5117)
 - bump k8s.io/client-go from 0.21.3 to 0.22.1 in /api [#5027](https://github.com/keptn/keptn/issues/5027)
 - bump github.com/go-openapi/errors in /configuration-service [#5036](https://github.com/keptn/keptn/issues/5036)
 - bump marocchino/sticky-pull-request-comment [#5103](https://github.com/keptn/keptn/issues/5103)
 - bump github.com/nats-io/nats-server/v2 in /distributor [#5104](https://github.com/keptn/keptn/issues/5104)
 - bump github.com/go-openapi/runtime in /mongodb-datastore [#5105](https://github.com/keptn/keptn/issues/5105)
 - bump github.com/nats-io/nats-server/v2 [#5107](https://github.com/keptn/keptn/issues/5107)
 - bump k8s.io/apimachinery in /test/go-tests [#5112](https://github.com/keptn/keptn/issues/5112)
</p>
</details>
<details><summary>remediation-service</summary>
<p>
 - adapt to recent changes in go sdk [#5464](https://github.com/keptn/keptn/issues/5464)
</p>
</details>
<details><summary>shipyard-controller</summary>
<p>
 - Reduce log noise for sequence watcher component [#5458](https://github.com/keptn/keptn/issues/5458)
 - More robust handling of multiple .started/finished events for the same task at the same time [#5440](https://github.com/keptn/keptn/issues/5440)
 - Remove log noise in sequence migrator [#5096](https://github.com/keptn/keptn/issues/5096)
 - Adapted sequence state representation for case where sequence can not be started (#5137) [#5194](https://github.com/keptn/keptn/issues/5194)
 - Return proper error message in case project is not available (#4399) [#5231](https://github.com/keptn/keptn/issues/5231)
 - Return error if a sequence for an unavailable stage is triggered (#4791) [#5069](https://github.com/keptn/keptn/issues/5069)
 - Adapted log output when no queued sequence is found (#5138) [#5167](https://github.com/keptn/keptn/issues/5167)
 - Adapted HTTP status codes of GET /event endpoint (#5132) [#5134](https://github.com/keptn/keptn/issues/5134)
 - Avoid endless loop (#5096) [#5124](https://github.com/keptn/keptn/issues/5124)
</p>
</details>
<details><summary>configuration-service</summary>
<p>
 - deprecate get default resources endpoints [#5443](https://github.com/keptn/keptn/issues/5443)
 - Make sure upstream changes are pulled when updating upstream creds (#5149) [#5224](https://github.com/keptn/keptn/issues/5224)
 - Implemented endpoints for deleting service and stage resources (#5136) [#5145](https://github.com/keptn/keptn/issues/5145)
</p>
</details>
<details><summary>installer</summary>
<p>
 - temporarily revert customization of repository string in chart [#5414](https://github.com/keptn/keptn/issues/5414)
 - Add option for ingress to control-plane helm chart keptn installer [#5066](https://github.com/keptn/keptn/issues/5066)
</p>
</details>
<details><summary>distributor</summary>
<p>
 - Ensure that the subscriptionID is passed to the event (#5405) [#5412](https://github.com/keptn/keptn/issues/5412)
 - pass along subscription id to service implementation [#5374](https://github.com/keptn/keptn/issues/5374)
 - Exclusive message processing for multiple distributors (#4689) [#5249](https://github.com/keptn/keptn/issues/5249)
 - Only interpret events with status=errored as error logs (#5170) [#5186](https://github.com/keptn/keptn/issues/5186)
</p>
</details>
<details><summary>deps-dev</summary>
<p>
 - bump typescript from 4.3.5 to 4.4.3 in /bridge/server
</p>
</details>
<details><summary>secret-service</summary>
<p>
 - added creation of rolebinding based on scope name [#5300](https://github.com/keptn/keptn/issues/5300)
 - Add list of keys within secrets created by the secret-service (#4749) [#5139](https://github.com/keptn/keptn/issues/5139)
</p>
</details>
<details><summary>cli</summary>
<p>
 - added zones to times format according to (ISO8601) [#4788](https://github.com/keptn/keptn/issues/4788)
 - Check if kubectl context matches Keptn CLI context before applying upgrade (#4583) [#5250](https://github.com/keptn/keptn/issues/5250)
 - skip version check on install [#5046](https://github.com/keptn/keptn/issues/5046)
</p>
</details>
<details><summary>lighthouse-service</summary>
<p>
 - calcscore missing error msg (#5142) [#5252](https://github.com/keptn/keptn/issues/5252)
 - Added error logs for failing monitoring configuration (#5088) [#5220](https://github.com/keptn/keptn/issues/5220)
 - Add message to event in case SLO parsing failed (#5130) [#5135](https://github.com/keptn/keptn/issues/5135)
</p>
</details>
<details><summary>jmeter-service</summary>
<p>
 - prevent failure if deploymentURIs does not end with a '/' [#3612](https://github.com/keptn/keptn/issues/3612)
</p>
</details>
<details><summary>api</summary>
<p>
 - Try to use X-real-ip and X-forwarded-for headers before Rem… [#5080](https://github.com/keptn/keptn/issues/5080)
 - Try to use X-real-ip and X-forwarded-for headers [#5082](https://github.com/keptn/keptn/issues/5082)
</p>
</details>
