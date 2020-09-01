# Release Notes 0.7.1

Keptn 0.7.1 improves the capabilities of the 0.7 release by adding more flexibility to the tasks of installing/operating Keptn, introducing two new CLI commands and improving the UX in the Keptn Bridge. Especially the new feature of installing Keptn in different namespaces is a door opener for various use-cases since multiple Keptn deployments, e.g., one for each team, can be operated on one cluster. 

**The three key announcements of Keptn 0.7.1**:

:rocket: *More flexibility in installing/operating Keptn*: 
- `keptn install --namespace`: Allows you to specify the namespace to install Keptn in. 
- `keptn upgrade`: Instead of manually executing a K8s job as done by the previous upgrade processes, this CLI command performs the upgrade. 
- Set `root context`: As part of the installation Helm Chart of Keptn, a root context can be specified that allows you to customize the API endpoint of Keptn. 

:star2: *New CLI commands*:
- `keptn delete service`: This command will delete a service from a project and *undeploy* the service from the cluster. 
- `keptn get events`: This command is a generic implementation to retrieve Keptn events of any event type.  

:sparkles: *UX improvements in environment screen*: Based on feedback on the features of the Keptn delivery assistant, improvements regarding the user experience of the Keptn Bridge has been implemented.

*Additional note:* Added documentation of [GOVERNANCE](https://github.com/keptn/keptn/blob/0.7.1/GOVERNANCE.md) & [SECURITY](https://github.com/keptn/keptn/blob/0.7.1/GOVERNANCE.md) process

## Keptn Specification

Implemented **Keptn spec** version: [0.1.4](https://github.com/keptn/spec/tree/0.1.4)

## New Features

<details><summary>Platform Support / Installer</summary>
<p>

- Install Keptn in a namespace other than keptn using: `keptn install --namespace=` [#2195](https://github.com/keptn/keptn/issues/2195)
- Upgrade Keptn from 0.7 to 0.7.1 using: `keptn upgrade` [#2234](https://github.com/keptn/keptn/issues/2234)
- Make StorageClass and PersistentVolumeClaim configurable in Keptn installer [#2190](https://github.com/keptn/keptn/issues/2190)
- Allow to install Keptn with prefix in path (aka. context root) [#2124](https://github.com/keptn/keptn/issues/2124)

</p>
</details>

<details><summary>API</summary>
<p>

- Return 404 from `/event` endpoint when no event is found [#1655](https://github.com/keptn/keptn/issues/1655)

</p>
</details>

<details><summary>CLI</summary>
<p>

- `keptn delete` to delete a service from a Keptn project [#2199](https://github.com/keptn/keptn/issues/2199)
- Properly format the output of `keptn get event` command [#2207](https://github.com/keptn/keptn/issues/2207)
- `keptn get event` to get an event of any event type [#2171](https://github.com/keptn/keptn/issues/2171)
- `keptn add-resource` checks the number of arguments before executing the command [#1735](https://github.com/keptn/keptn/issues/1735)
- Immediately return an error if kube server version check error [#1944](https://github.com/keptn/keptn/issues/1944)
- Review of the description of all Keptn CLI commands [#1718](https://github.com/keptn/keptn/issues/1718)

</p>
</details>

<details><summary>Keptn Core</summary>
<p>

- *shipyard-controller*: (**not released**)
  - Manage open *.triggered events in a mongoDB collection per project [#2158](https://github.com/keptn/keptn/issues/2158)
  - Manage open *.started events in a mongoDB collection per project [#2159](https://github.com/keptn/keptn/issues/2159)
  - Control task sequences defined in the Shipyard [#2193](https://github.com/keptn/keptn/issues/2193)

- *distributor*:
  - Sidecar for polling open *.triggered events [#2166](https://github.com/keptn/keptn/issues/2166)

- *helm-service*:
  - Delete a service from the cluster when deleting it from a project [#2201](https://github.com/keptn/keptn/issues/2201)

- *lighthouse-service*: 
  - Mark info SLI correctly when empty pass/warning array is provided [#2231](https://github.com/keptn/keptn/issues/2231)
  - Change the comparison strategy to match the full quality gate result [#2224](https://github.com/keptn/keptn/issues/2224)

</p>
</details>

<details><summary>Bridge</summary>
<p>

- Icons in stage tile and labels work as filter [#2087](https://github.com/keptn/keptn/issues/2087)
- Provide API token and `keptn auth` command in user menu [#2197](https://github.com/keptn/keptn/issues/2197)
- Show remediation workflow in environment screen [#2085](https://github.com/keptn/keptn/issues/2085)
- Show failed quality gates in environment screen [#2086](https://github.com/keptn/keptn/issues/2086)
- Fixed misleading message in bridge if no deployment was done but evaluations happened [#2112](https://github.com/keptn/keptn/issues/2112)
- Improved color-coding in Keptn Bridge for `problem.resolved` event [#2139](https://github.com/keptn/keptn/issues/2139)
- Provider better indication and workflow for artifacts waiting for approval [#2142](https://github.com/keptn/keptn/issues/2142)
- Fix wrong version in environments overview when evaluation failed [#2133](https://github.com/keptn/keptn/issues/2133)
- Set height for evaluation chart and maxHeight for legend [#2150](https://github.com/keptn/keptn/issues/2150)
- Expand service tile automatically if there is an open approval [#2151](https://github.com/keptn/keptn/issues/2151)
- Show labels in event payload [#2138](https://github.com/keptn/keptn/issues/2138)
- Bridge code refactoring tasks: [#2000](https://github.com/keptn/keptn/issues/2000) [#2011](https://github.com/keptn/keptn/issues/2011) [#2012](https://github.com/keptn/keptn/issues/2012)

</p>
</details>

## Fixed Issues

- CLI: kubectl version check incorrectly reports if no connection to cluster could be made [#1944](https://github.com/keptn/keptn/issues/1944)

## Development Process / Testing

- Fix problems in TravisCI for building/testing [#2149](https://github.com/keptn/keptn/issues/2149)
- Fixed integration test function and replaced invalid error codes [#2162](https://github.com/keptn/keptn/issues/2162)
- Check if replicas of deployment are running [#2160](https://github.com/keptn/keptn/issues/2160)

## Good to know / Known Limitations

- The upgrade from 0.7 o 0.7.1 is supported by the `keptn upgrade` command.
