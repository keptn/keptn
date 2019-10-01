# Release Notes 0.5.0

## New Features
- Upstream mechanism to any Git repo [#671](https://github.com/keptn/keptn/issues/671)
- Verification of remediation action by pitometer service [#850](https://github.com/keptn/keptn/issues/850)
- Untar Helm chart in configuration-service [#870](https://github.com/keptn/keptn/issues/870)
- Added additional features to Keptn's Bridge [#875](https://github.com/keptn/keptn/issues/875)
    - Keptn entry points indicate whether it is a configuration-change or remediation-action
    - Include complete JSON payload of events
    - Show evaluation results of pitometer
    - Context-sensitive information for events (e.g. canary action, promotion to next stage, test duration)
- Use applied manifests for generating the generated charts [#900](https://github.com/keptn/keptn/issues/900)

## CLI Enhancements
- Improve checks for Helm charts [#672](https://github.com/keptn/keptn/issues/672)
- Allow to onboard Helm folder [#868](https://github.com/keptn/keptn/issues/868)
- Delete installer job in uninstall commando [#873](https://github.com/keptn/keptn/issues/873)
- Provide delete project command [#887](https://github.com/keptn/keptn/issues/887)
- Rename flag for deployment strategy in onboard service [#889](https://github.com/keptn/keptn/issues/889)
- Check project and stage names [#745](https://github.com/keptn/keptn/issues/745)


## New Services
-

## Fixed Issues
- Retry of authentication in configure domain [#846](https://github.com/keptn/keptn/issues/846)
- Avoid duplicate entries in Keptn's Bridge [#851](https://github.com/keptn/keptn/issues/851)
- Do not uninstall Tiller and Istio [#857](https://github.com/keptn/keptn/issues/857)
- Log Helm output [#574](https://github.com/keptn/keptn/issues/574)
- Remove read limit in websocket communication [#847](https://github.com/keptn/keptn/issues/847)
- Update keptn specification to be 0.5.0 compatible [#844](https://github.com/keptn/keptn/issues/844)
- Make websocket communication optional [#728](https://github.com/keptn/keptn/issues/728)

## Known Limitations
