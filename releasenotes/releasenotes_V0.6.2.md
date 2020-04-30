# Release Notes 0.6.2

This release fixes several smaller issues and includes some quality-of-life improvements for API, Bridge and CLI.

## Keptn Specification

Implemented **Keptn spec** version: [0.1.3](https://github.com/keptn/spec/tree/0.1.3)

## New Features

* API
  * Introduce an API-gateway that proxies requests to configuration-service [#1510](https://github.com/keptn/keptn/issues/1510)
  * Query a list of projects [#1559](https://github.com/keptn/keptn/issues/1559)
  * Provide an endpoint for exposing Keptn's Bridge via Istio or nginx ingress [#1153](https://github.com/keptn/keptn/issues/1153)
* Bridge
  * Use icons for events [#1352](https://github.com/keptn/keptn/issues/1352)
  * Deep links [#1316](https://github.com/keptn/keptn/issues/1316)
  * Format SLOs as floats #1681
* CLI
  * Create a support-archive for troubleshooting [#1549](https://github.com/keptn/keptn/issues/1549)
  * Provide a CLI command for exposing Keptn's Bridge [#1560](https://github.com/keptn/keptn/issues/1560)
* Configuration-service
  * Improve troubleshooting for git related problems [#1637](https://github.com/keptn/keptn/issues/1637)
* Installer
  * Allow specifying a domain when installing Keptn (e.g., `keptn install --domain=127.0.0.1.nip.io`) [#1482](https://github.com/keptn/keptn/issues/1482)
  * Allow to re-use existing nginx-ingress installation [#1712](https://github.com/keptn/keptn/issues/1712)


## Fixed Issues

* API
  * Do not overwrite `source` property of cloud events [#1643](https://github.com/keptn/keptn/issues/1643)
* Configuration-service & Shipyard-service
  * Catch a not initialized Git repo by creating an initial commit [#1545](https://github.com/keptn/keptn/issues/1545)
  * Fixed error handling (issue with quality-gates multi-stage setups) [#1695](https://github.com/keptn/keptn/issues/1695)
* Installer
  * Check for ImagePullBackOff errors for the installer job [#1521](https://github.com/keptn/keptn/issues/1521)
  * Do not overwrite an existing Keptn installation [#1376](https://github.com/keptn/keptn/issues/1376)
* Bridge
  * Provide proper deep-link functionality for "Problem detected" events [#1557](https://github.com/keptn/keptn/issues/1557)
  * Bridge preselects wrong evaluation event in heatmap view [#1518](https://github.com/keptn/keptn/issues/1518) 
  * Heatmap shows undefined color for test results of type `fail` [#1580](https://github.com/keptn/keptn/issues/1580)
* Lighthouse
  * Evaluating "<=5%" was interpreted as "<=5" (missing percent sign) [#1498](https://github.com/keptn/keptn/issues/1498)
* Helm
  * Helm Service should not require an outbound Internet connection [#1532](https://github.com/keptn/keptn/issues/1532)

## Refactoring

* Refactor api-service and configuration-service [#1510](https://github.com/keptn/keptn/issues/1510)
* Refactor go-utils [#1492](https://github.com/keptn/keptn/issues/1492)
* Change APIVersion from apps/v1beta1 to apps/v1 [#1529](https://github.com/keptn/keptn/issues/1529)
* Added multiple unit tests to improve code coverage

## Development Workflow

* Improve travis-ci workflow
* Added GitHub actions for linting
* Updated contribution guide

## Good to know / known limitations
- For old limitations, please see [Release 0.6.1](https://github.com/keptn/keptn/releases/tag/0.6.1). 
