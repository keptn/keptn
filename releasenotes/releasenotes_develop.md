# Release Notes 0.8.2

Keptn 0.8.2 is a hardening release and incorporates changes that are required to deploy Keptn using Keptn. Besides, the *Sequence* and *Service* screens in the Bridge are improved to provide a holistic view of the services on the one hand, and a clear separation of sequentially executed sequences on the other hand.

---

**Key announcements:**

:cocktail: *Keptn drinks its own champaign*: Each merge on the *master* (aka. main) branch of the keptn/keptn repository triggers a Keptn to deploy this new version into a development environment. This helps the Keptn project to have the latest and greatest version of Keptn running and to do feature presentations using this deployment. 

> A GitHub action has been implemented that allows sending a Keptn event to a Keptn installation: [gh-action-send-event](https://github.com/keptn/gh-action-send-event). Do wait to integrate it into your CI on GitHub today!

:hammer: *Hardening of Keptn*: Containers of Keptn core services are not running as root users and a security context has been defined. *Please note*: Defining the security context for NATS is not possible yet, since we are waiting for a new release containing the PR: https://github.com/nats-io/k8s/pull/222.  

:star: *Service/Sequence screen enhancements*: Instead of listing all sequences that were executed on a service, the *Service* screen now displays the deployments of a service that are currently running in the different stages. Besides, the *Sequence* screen improved in displaying sequences that are executed sequentially in the same stage.

> *add Screenshot here*

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

<details><summary>API</summary>
<p>

- *Fixed*: GET `/project/{project}` returns 200 instead of 404 [3699](https://github.com/keptn/keptn/issues/3699)

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

- OAuth/OpenID connect based login for Keptn bridge [3448](https://github.com/keptn/keptn/issues/3448)
- Text in the notification is hard to read and the link is not working [3634](https://github.com/keptn/keptn/issues/3634)
- Chart tooltip should show all active metrics [2546](https://github.com/keptn/keptn/issues/2546)
- Better sequence visualization when having more than one sequence in a stage [3542](https://github.com/keptn/keptn/issues/3542)
- Load details of deployment in service screen [3706](https://github.com/keptn/keptn/issues/3706)
- Derive the list of deployments that are currently running for a service [3629](https://github.com/keptn/keptn/issues/3629)
- Uniform names of mock files [3714](https://github.com/keptn/keptn/issues/3714)
- Add "load older Sequences" button in Sequence screen [2280](https://github.com/keptn/keptn/issues/2280)
- Sequence icon colors represent status [3591](https://github.com/keptn/keptn/issues/3591)
- Show target values for criteria when hovering over values [2757](https://github.com/keptn/keptn/issues/2757)
- *Fixed*: Chart in Heatmap shows incorrect Y-axis and incorrect tooltip [3645](https://github.com/keptn/keptn/issues/3645)
- *Fixed*: Daily version check was disabled [3668](https://github.com/keptn/keptn/issues/3668)
- *Fixed*: Quality gate icon in the environment screen does not turn red [3592](https://github.com/keptn/keptn/issues/3592)
- *Fixed*: Some deep-links are broken [3631](https://github.com/keptn/keptn/issues/3631)
- *Fixed*: Problem filter in environment screen does not work [3652](https://github.com/keptn/keptn/issues/3652)

</p>
</details>

## Miscellaneous

- Delete helm-service and jmeter-service from continuous-delivery Helm Chart and adapt CLI accordingly [3350](https://github.com/keptn/keptn/issues/3350)
- Decouple unit tests from "get.keptn.sh/version.json" [3476](https://github.com/keptn/keptn/issues/3476)

## Development Process / Testing

- GitHub Action to setup a Keptn deployment [3554](https://github.com/keptn/keptn/issues/3554)
- GitHub Action to send CloudEvents to a Keptn installation [2797](https://github.com/keptn/keptn/issues/2797)
- Use github.com/stretchr/testify/require instead of general asserts [3149](https://github.com/keptn/keptn/issues/3149)

## Good to know / Known Limitations

- See the know limitations from [0.8.0](https://github.com/keptn/keptn/releases/tag/0.8.0)

<details><summary>Open issues that will be fixed in upcoming releases</summary>
<p>

  <!--TODO: final check-->
  - Auto-remediation does not work with remote execution plane [3498](https://github.com/keptn/keptn/issues/3498)
  - Bridge shows "started" wording on status.changed [3583](https://github.com/keptn/keptn/issues/3583)
  - Inconsistent usage of user-managed and user_managed causing issues [3624](https://github.com/keptn/keptn/issues/3624)
  - Keptn API: `configure monitoring` not functioning according to spec [3638](https://github.com/keptn/keptn/issues/3638)
  - Keptn CLI: Disable Kube context check [3666](https://github.com/keptn/keptn/issues/3666)
  - Homebrew installed CLI fails install with 'Malformed constraint: ""' [3805](https://github.com/keptn/keptn/issues/3805)

</p>
</details>

## Upgrade to 0.8.2

- The upgrade from 0.8.x to 0.8.2 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.2](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-1-to-0-8-2)