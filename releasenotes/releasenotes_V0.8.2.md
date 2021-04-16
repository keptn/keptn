# Release Notes 0.8.2

Keptn 0.8.2 

---

**Key announcements:**

---

Many thanks to the community for the enhancements on this release! 
 
## Keptn Specification

Implemented **Keptn spec** version: [0.2.1](https://github.com/keptn/spec/tree/0.2.1)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Run Keptn core containers as non-root [3764](https://github.com/keptn/keptn/issues/3764)
- Helm Chart (for control-plane) needs tag properties for deployments in values.yaml [3328](https://github.com/keptn/keptn/issues/3328)
- Split K8s role `keptn-configure-bridge` by secret and pod management [3767](https://github.com/keptn/keptn/issues/3767)
- Service account `keptn-configuration-service` does not need full permissions on secret management [3781](https://github.com/keptn/keptn/issues/3781)
- Dockerfile for Keptn Bridge in package.json usage needs improvement [3641](https://github.com/keptn/keptn/issues/3641)
- Improve handling of X-Forwarded-Proto header for Bridge [3672](https://github.com/keptn/keptn/issues/3672)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *helm-service*:
  - Smart Helm Chart values merger [3341](https://github.com/keptn/keptn/issues/3341)
  - *Fixed*: Not working parallel when deployed in the execution-plane [3427](https://github.com/keptn/keptn/issues/3427)
  - *Fixed*: Delivery failed with "Error when installing/upgrading chart" ... "has no deployed releases" [3407](https://github.com/keptn/keptn/issues/3407)

- *jmeter-service*:
  - Need better JMeter result other than just fail [3559](https://github.com/keptn/keptn/issues/3559)

- *lighthouse-service*:
  - *Fixed*: Properly set result, status, and message [3412](https://github.com/keptn/keptn/issues/3412)

- *shipyard-controller*:
  - *Fixed*: Only last `.finished` event for a task determines further sequence execution [3493](https://github.com/keptn/keptn/issues/3493)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Derive the list of deployments that are currently running for a service [3629](https://github.com/keptn/keptn/issues/3629)
- Uniform names of mock files [3714](https://github.com/keptn/keptn/issues/3714)
- Add "load older Sequences" button in Sequence screen [2280](https://github.com/keptn/keptn/issues/2280)
- Sequence icon colours represent status [3591](https://github.com/keptn/keptn/issues/3591)
- Show target values for criteria when hovering over values [2757](https://github.com/keptn/keptn/issues/2757)
- *Fixed*: Quality gate icon in the environment screen does not turn red [3592](https://github.com/keptn/keptn/issues/3592)
- *Fixed*: Some deep-links are broken [3631](https://github.com/keptn/keptn/issues/3631)
- *Fixed*: Problem filter in environment screen does not work [3652](https://github.com/keptn/keptn/issues/3652)

</p>
</details>

## Miscellaneous

- Delete helm-service and jmeter-service from continuous-delivery Helm Chart and adapt CLI accordingly [3350](https://github.com/keptn/keptn/issues/3350)
- Decouple unit tests from "get.keptn.sh/version.json" [3476](https://github.com/keptn/keptn/issues/3476)

## Development Process / Testing

- Create a GitHub Action to send CloudEvents to a Keptn installation [2797](https://github.com/keptn/keptn/issues/2797)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  - Auto-remediation does not work with remote execution plane [3498](https://github.com/keptn/keptn/issues/3498)
  - Bridge shows "started" wording on status.changed [3583](https://github.com/keptn/keptn/issues/3583)
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
  - API for Configure Monitoring not functioning as expected [3638](https://github.com/keptn/keptn/issues/3638)
  - Keptn CLI: Disable Kube context check [3666](https://github.com/keptn/keptn/issues/3666)

</p>
</details>

## Upgrade to 0.8.2

- The upgrade from 0.8.x to 0.8.2 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.2](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-x-to-0-8-2)