# Release Notes 0.8.7

This is a bug fix release for Keptn 0.8.4 - 0.8.6, containing fixes for the Bridge and management of Keptn-service registrations.

---
## Improvements and Fixes 

<details><summary>Installer</summary>
<p>

- *Fixes:*
  - Upstream sent too big header while reading response header from upstream error while SSO logout [4662](https://github.com/keptn/keptn/issues/4662)

</p>
</details>

<details><summary>CLI</summary>
<p>

- *Fixes:*
  - Removed warning when `KUBECONFIG` file is missing [4553](https://github.com/keptn/keptn/issues/4553)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *Improvement:*
  - Throttled authenticating a user on Bridge or Keptn CLI on `/auth` endpoint  [4323](https://github.com/keptn/keptn/issues/4323)

- *Fixes:*
  - *shipyard-controller*: `nil` pointer dereferences when receiving events with missing fields [4652](https://github.com/keptn/keptn/issues/4652)
  - De-registration of uniform services fails on upgrade scenario [4615](https://github.com/keptn/keptn/issues/4615)
  - Support `https` for uniform registration outside the Keptn cluster [4516](https://github.com/keptn/keptn/issues/4516)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- *Fixes:*
  - The bridge breaks if the first key-value pair is deleted [4622](https://github.com/keptn/keptn/issues/4622)
  - Secrets list is not updated after deleting a secret [4633](https://github.com/keptn/keptn/issues/4633)
  - Secrets view is empty after deleting one secret [4660](https://github.com/keptn/keptn/issues/4660)
  - Bridge server returns wrong http-status-code [4658](https://github.com/keptn/keptn/issues/4658)
  - Service evaluation screen refreshes every time after polling data [4491](https://github.com/keptn/keptn/issues/4491)

</p>
</details>

## Upgrade to 0.8.7

- The upgrade from 0.8.x to 0.8.7 is supported by the `keptn upgrade` command. Find the documentation here: [Upgrade from Keptn 0.8.x to 0.8.7](https://keptn.sh/docs/0.8.x/operate/upgrade/#upgrade-from-keptn-0-8-6-to-0-8-7)