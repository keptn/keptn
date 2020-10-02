# Release Notes 0.7.2

Keptn 0.7.2 improves the capabilities of the 0.7.0 release by adding ... 

**The key announcements of Keptn 0.7.2**:

:rocket: *Improved UX for Quality Gates*: This release introduces an `/evaluation` endpoint that brings the same user-experience to the API as used from the Keptn CLI. Hence, easily trigger a quality gate evaluation based on a small set of parameters, like: `start`, `end`, and `timeframe`. 

:sparkles: *Focused UI enhancements for Quality Gates*: The evaluation tile that displays the quality gate evaluation result has been improved in various ways. First, the Heatmap in this tile now highlights the currently selected evaluation result. Second, it also displays the evaluation results that were used for comparisons giving the user more insights into the quality gate. 

:tada: *Usage statistics for Keptn installation*: If you want to get more insights into your Keptn installation, feel free to deploy the [statistics-service](https://github.com/keptn-sandbox/statistics-service) from the Keptn sandbox. This service provides usage statistics based on events and service executions.

## Keptn Specification

Implemented **Keptn spec** version: [0.1.6](https://github.com/keptn/spec/tree/0.1.6)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- K8s 1.19 support [#2411](https://github.com/keptn/keptn/issues/2411)

</p>
</details>

<details><summary>API</summary>
<p>

- Introduced API endpoint `\evaluation` for triggering evaluations [#2387](https://github.com/keptn/keptn/issues/2387)
- Swagger automatically determines the scheme (https or http) [#2325](https://github.com/keptn/keptn/issues/2325)

</p>
</details>

<details><summary>CLI</summary>
<p>

- Improve robustness and UX of `keptn generate support-archive` [#2340](https://github.com/keptn/keptn/issues/2340)
- Point user to upgrade docs, fixed bug in `keptn update project` [#2293](https://github.com/keptn/keptn/issues/2293)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- Increased password and token entropy [#2295](https://github.com/keptn/keptn/issues/2295)

- *configuration-service*: 
  - Allow non-HTTPS connections to Git-upstream [#2336](https://github.com/keptn/keptn/issues/2336)

- *lighthouse-service*:
  - *Behavior change*: `include_result_with_score` just works on SLO-level and `sh.keptn.events.evaluation-done` returns compared evaluation results [#2388](https://github.com/keptn/keptn/issues/2388)
  - Send `sh.keptn.events.evaluation-done` events with error information when service/stage/project not found [#2365](https://github.com/keptn/keptn/issues/2365)
  - Use the ConfigMap `lighthouse-config` which refers to a default SLI provider [#2317](https://github.com/keptn/keptn/issues/2317)
  - Trigger SLI retrieval even though the SLO is empty or not available [#2318](https://github.com/keptn/keptn/issues/2318)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Evaluation tile rework [#2305](https://github.com/keptn/keptn/issues/2305)
- Heatmap highlights evaluation results that are used for comparison [#2389](https://github.com/keptn/keptn/issues/2389)
- Show the actual state of the approval in approval finished events [#2371](https://github.com/keptn/keptn/issues/2371)
- Enable highlighting of the currently selected evaluation result in Heatmap [#1640](https://github.com/keptn/keptn/issues/1640)
- Enable caching for static files in express [#2408](https://github.com/keptn/keptn/issues/2408)
- Provide access to up-stream configuration repository per project [#1335](https://github.com/keptn/keptn/issues/1335)
- Hide the API token and `keptn auth` command per default [#2257](https://github.com/keptn/keptn/issues/2257)
- Feature toggle for version check and API token info [#2320](https://github.com/keptn/keptn/issues/2320)
- Show a download link for Keptn CLI [#2319](https://github.com/keptn/keptn/issues/2319)

</p>
</details>

## Fixed Issues

- Fixed broken link to Keptn API in Bridge [#2430](https://github.com/keptn/keptn/issues/2430)
- Fixed infinite loop while confirming cluster information [#2425](https://github.com/keptn/keptn/issues/2425)
- Fixed wrong version number for API endpoints [#2315](https://github.com/keptn/keptn/issues/2315)
- Fixed bug: Bridge UI breaks on first open approval event on stage [#2354](https://github.com/keptn/keptn/issues/2354)

## Development Process / Testing

- Fixed Travis-CI integration tests [#2335](https://github.com/keptn/keptn/issues/2335)

## Good to know / Known Limitations

- The upgrade from 0.7.x to 0.7.2 is supported by the `keptn upgrade` command.
