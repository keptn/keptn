# Release Notes 0.7.3

Keptn 0.7.3 provides a new feature that allows invalidating a quality gate evaluation result. Besides, this release addresses minor bugs as well as UI/UX improvements. 

**The key announcement of Keptn 0.7.3**:

:rocket: *Invalidate a quality gate evaluation*: This feature allows a user to manually invalidate an evaluation result meaning that it will be ignored for future comparison. The original evaluation result (i.e., the `sh.keptn.events.evaluation-done` event) will not be changed or deleted and can still be fetched from the Keptn API.

## Keptn Specification

Implemented **Keptn spec** version: [0.1.7](https://github.com/keptn/spec/tree/0.1.7)

- Added new event of type: `sh.keptn.events.evaluation.invalidated` - this event marks a quality gate evaluation result as invalid [#47](https://github.com/keptn/spec/issues/47) 

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Added installation type as environment variable to Bridge [#2606](https://github.com/keptn/keptn/issues/2606)
- Update NGINX version to 1.19.4-alpine [#2651](https://github.com/keptn/keptn/issues/2651)

</p>
</details>

<details><summary>CLI</summary>
<p>

- Show OS-specific quick access instructions [#2576](https://github.com/keptn/keptn/issues/2576)
- Added timeout of 5 seconds to server version check [#2589](https://github.com/keptn/keptn/issues/2589)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *configuration-service*: 
  - Update Git upstream information  in materialized view [#2577](https://github.com/keptn/keptn/issues/2577)
  - Fixed create project with upstream credentials [#2639](https://github.com/keptn/keptn/issues/2639)

- *helm-service*: 
  - Added timeout when waiting for deployment to be rolled out [#2578](https://github.com/keptn/keptn/issues/2578)  

- *lighthouse-service*: 
  -  Support the invalidation of evaluation results [#2449](https://github.com/keptn/keptn/issues/2449)

- *mongodb-datastore*: 
  - Removed cloudevents+json from list of produced content types [#2582](https://github.com/keptn/keptn/issues/2582)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Fixed issue of showing no information at the "Compared with" label [#2545](https://github.com/keptn/keptn/issues/2545)
- Show the correct evaluation timeframe when switching evaluations in the HeatMap [#2585](https://github.com/keptn/keptn/issues/2585)
- Show all core use-cases on integrations page depending on installation type[#2565](https://github.com/keptn/keptn/issues/2565)
- Allow invalidating evaluation results from Bridge [#2548](https://github.com/keptn/keptn/issues/2548)
- Fix base href in Bridge [#2564](https://github.com/keptn/keptn/issues/2564)

</p>
</details>

## Fixed Issues

- Allow loading more than 20 services/projects in Bridge [#2631](https://github.com/keptn/keptn/issues/2631)

## Development Process / Testing

- Migrated TravisCI from travis-ci.org to travis-ci.com [#2356](https://github.com/keptn/keptn/issues/2356)

## Good to know / Known Limitations

This section lists bugs and limitations that are known but not fixed in this release. They will get addressed in one of the next releases.

- Bridge ignores deployed service artifact [#2543](https://github.com/keptn/keptn/issues/2543)
  - The Bridge loads the last 20 triggers (aka. root events) and if this list does not contain a `sh.keptn.event.configuration.change` or `sh.keptn.event.start-evaluation` event, the label below the service name shows: *Service not deployed*.
- Hovering over the score in an `approval.triggered` events in the Bridge leads to a scroll-up / jump-up in Firefox [#2369](https://github.com/keptn/keptn/issues/2369)


## Upgrade to 0.7.3

- The upgrade from 0.7.x to 0.7.3 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.7.2 to 0.7.3](https://keptn.sh/docs/0.7.x/operate/upgrade/#upgrade-from-keptn-0-7-2-to-0-7-3)
