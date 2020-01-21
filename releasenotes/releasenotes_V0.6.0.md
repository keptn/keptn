# Release Notes 0.6.0

## Keptn Specification

Implemented [Keptn spec](https://github.com/keptn/spec/tree/0.1.2) version: 0.1.2

## New Features
- Self-healing use-case with Dynatrace [#1185](https://github.com/keptn/keptn/issues/1185)
- Self-healing use-case based on feature toggle remediation with Unleash [#1104](https://github.com/keptn/keptn/issues/1104)
- Service Level Indicators (SLI) are stored in Keptn's git repository [#1192](https://github.com/keptn/keptn/issues/1192)
- Single Istio gateway accepting HTTPS traffic [#1231](https://github.com/keptn/keptn/issues/1231)
- Check if Istio is already installed before installing it [#1208](https://github.com/keptn/keptn/issues/1208)
- Avoid installing incompatible Keptn version [#1162](https://github.com/keptn/keptn/issues/1162)
- Servicenow-service and prometheus-service have been removed from Keptn's uniform [#1302](https://github.com/keptn/keptn/issues/1302)
- *jmeter-service*: Sends a test-finished with status `fail` if test execution failed [#542](https://github.com/keptn/keptn/issues/542)

## CLI Enhancements
- Setup Dynatrace monitoring using the Keptn CLI `configure monitoring` command [#443](https://github.com/keptn/keptn/issues/443)
- Configure domain assumes dedicated Keptn gateway [#1295](https://github.com/keptn/keptn/issues/1295)

## Fixed Issues
- Retrieving hostname/ip of ingress improved [#1199](https://github.com/keptn/keptn/issues/1199)
- Keptn API: Swagger UI refers to https instead of http [#1012](https://github.com/keptn/keptn/issues/1012)
- Allow project, stage, and service names with hyphens [#1166](https://github.com/keptn/keptn/issues/1166)
- Check whether Keptn-service name matches Kubernetes-service name [#1261](https://github.com/keptn/keptn/issues/1261)
- api.keptn.DOMAIN/v1/events endpoints requiers query parameters [#1210](https://github.com/keptn/keptn/issues/1210)
- *configuration-service*: API improvements [#1275](https://github.com/keptn/keptn/issues/1275)
- *lighthouse-service*: Prevent previous failed SLI results to be used for comparison-based evaluation [#1263](https://github.com/keptn/keptn/issues/1263)
- *lighthouse-service*: can deal with an SLO file that has no SLI criteria defined [#1213](https://github.com/keptn/keptn/issues/1213)
- *mongodb-datastore*: eventContext removed from events that are read from mongodb-datastore [#1300](https://github.com/keptn/keptn/issues/1300)

## Known Limitations
