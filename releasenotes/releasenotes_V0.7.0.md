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

* **Upgrade from 0.6.2 to 0.7:** *Keptn 0.7 uses Helm 3 while previous Keptn releases rely on Helm 2*. To upgrade  your Helm releases from Helm 2 to 3, two options are provided: 
  1. *Job without Helm 3 Upgrade:* This option is needed when the cluster contains Helm releases not managed by Keptn. If this job is executed, it is necessary to manually converted the releases from Helm 2 to 3 as explained on [keptn.sh/docs](https://keptn.sh/docs/0.7.0/operate/upgrade/#job-without-helm-3-0-upgrade).
  1. *Job with Helm 3 Upgrade:* Full automation of Helm upgrade for installations were just Keptn is installed. If this job is executed, **all** Helm releases on the cluster are converted from Helm 2 to 3 and Tiller will be removed.
  