# Release Notes 0.7.0

:rocket: *Delivery Assistant:* - [SPEC 26](https://github.com/keptn/spec/pull/26)

:sparkles: *Closed-loop Remediation with custom Integrations:* - [KEP 09](https://github.com/keptn/enhancement-proposals/pull/9) | [SPEC 31](https://github.com/keptn/spec/pull/31)

:rocket: *Improved automation support with API extensions:* - [KEP 10](https://github.com/keptn/enhancement-proposals/pull/10)

:star2: *Upgrade from Helm 2.0 to 3.0:*

:hammer: *Hardening Keptn:*

## Keptn Specification

Implemented **Keptn spec** version: [master](https://github.com/keptn/spec/tree/master)

## New Features


## New Services


## Fixed Issues


## Development Process


## Good to know / Known Limitations

* **Upgrade from 0.6.2 to 0.7:** *Keptn 0.7 uses Helm 3.0 while previous Keptn releases rely on Helm 2.0*. By using the provided upgrader, **all** Helm releases are upgraded from Helm 2.0 to 3.0. This also includes Helm releases that are not managed by Keptn. If you have Helm releases on your cluster that are on version 2.0 and you do not want to upgrade, don't use the upgrader. Please take into account that the end-of-life period of Helm 2.0 begins on [August 13th, 2020](https://helm.sh/blog/covid-19-extending-helm-v2-bug-fixes/).  