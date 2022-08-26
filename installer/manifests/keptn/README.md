## Keptn Installer
Cloud-native application life-cycle orchestration. Keptn automates your SLO-driven multi-stage delivery and operations & remediation of your applications.

## Parameters

### Global values

| Name                    | Description                                         | Value             |
| ----------------------- | --------------------------------------------------- | ----------------- |
| `global.keptn.registry` | Global Docker image registry                        | `docker.io/keptn` |
| `global.keptn.tag`      | The tag of Keptn that should be used for all images | `""`              |


### MongoDB

| Name                                                | Description                                                         | Value                 |
| --------------------------------------------------- | ------------------------------------------------------------------- | --------------------- |
| `mongo.enabled`                                     |                                                                     | `true`                |
| `mongo.host`                                        |                                                                     | `mongodb:27017`       |
| `mongo.architecture`                                |                                                                     | `standalone`          |
| `mongo.updateStrategy.type`                         | Set the update strategy for MongoDB                                 | `Recreate`            |
| `mongo.service.nameOverride`                        |                                                                     | `mongo`               |
| `mongo.service.ports.mongodb`                       | Port for MongoDB to listen at                                       | `27017`               |
| `mongo.auth.enabled`                                |                                                                     | `true`                |
| `mongo.auth.databases`                              |                                                                     | `["keptn"]`           |
| `mongo.auth.existingSecret`                         |                                                                     | `mongodb-credentials` |
| `mongo.auth.usernames`                              |                                                                     | `["keptn"]`           |
| `mongo.auth.password`                               |                                                                     | `nil`                 |
| `mongo.auth.rootUser`                               |                                                                     | `admin`               |
| `mongo.auth.rootPassword`                           |                                                                     | `nil`                 |
| `mongo.auth.bridgeAuthDatabase`                     |                                                                     | `keptn`               |
| `mongo.external.connectionString`                   |                                                                     | `nil`                 |
| `mongo.containerSecurityContext`                    | Container Security Context that should be used for all MongoDB pods |                       |
| `mongo.serviceAccount.automountServiceAccountToken` |                                                                     | `false`               |
| `mongo.resources`                                   | Define resources for MongoDB                                        |                       |


### Keptn Features

| Name                                        | Description                             | Value    |
| ------------------------------------------- | --------------------------------------- | -------- |
| `features.debugUI.enabled`                  |                                         | `false`  |
| `features.automaticProvisioning.serviceURL` | Service for provisioning remote git URL | `""`     |
| `features.automaticProvisioning.message`    | Message for provisioning remote git URL | `""`     |
| `features.automaticProvisioning.hideURL`    | Hide automatically provisioned URL      | `true`   |
| `features.swagger.hideDeprecated`           |                                         | `false`  |
| `features.oauth.enabled`                    | Enable OAuth for Keptn                  | `false`  |
| `features.oauth.prefix`                     |                                         | `keptn:` |
| `features.git.remoteURLDenyList`            |                                         | `""`     |


### NATS

| Name                                               | Description                                                                | Value        |
| -------------------------------------------------- | -------------------------------------------------------------------------- | ------------ |
| `nats.nameOverride`                                |                                                                            | `keptn-nats` |
| `nats.fullnameOverride`                            |                                                                            | `keptn-nats` |
| `nats.cluster.enabled`                             | Enable NATS clustering                                                     | `false`      |
| `nats.cluster.replicas`                            | Define the NATS cluster size                                               | `3`          |
| `nats.cluster.name`                                | Define the NATS cluster name                                               | `nats`       |
| `nats.securityContext`                             | Define security context settings for NATS                                  |              |
| `nats.nats.automountServiceAccountToken`           |                                                                            | `false`      |
| `nats.nats.resources`                              | Define resources for NATS                                                  |              |
| `nats.nats.healthcheck.startup.enabled`            | Enable NATS startup probe                                                  | `false`      |
| `nats.nats.jetstream.enabled`                      |                                                                            | `true`       |
| `nats.nats.jetstream.memStorage.enabled`           | Enable memory storage for NATS Jetstream                                   | `true`       |
| `nats.nats.jetstream.memStorage.size`              | Define the memory storage size for NATS Jetstream                          | `500Mi`      |
| `nats.nats.jetstream.fileStorage.enabled`          |                                                                            | `true`       |
| `nats.nats.jetstream.fileStorage.size`             |                                                                            | `5Gi`        |
| `nats.nats.jetstream.fileStorage.storageDirectory` |                                                                            | `/data/`     |
| `nats.nats.jetstream.fileStorage.storageClassName` |                                                                            | `""`         |
| `nats.nats.securityContext`                        | Define the container security context for NATS                             |              |
| `nats.natsbox.enabled`                             | Enable NATS Box utility container                                          | `false`      |
| `nats.reloader.enabled`                            | Enable NATS Config Reloader sidecar to reload configuration during runtime | `false`      |
| `nats.exporter.enabled`                            | Enable NATS Prometheus Exporter sidecar to emit prometheus metrics         | `false`      |


### API Gateway Nginx

| Name                                                       | Description                                             | Value                |
| ---------------------------------------------------------- | ------------------------------------------------------- | -------------------- |
| `apiGatewayNginx.type`                                     |                                                         | `ClusterIP`          |
| `apiGatewayNginx.port`                                     |                                                         | `80`                 |
| `apiGatewayNginx.targetPort`                               |                                                         | `8080`               |
| `apiGatewayNginx.nodePort`                                 |                                                         | `31090`              |
| `apiGatewayNginx.podSecurityContext.enabled`               | Enable the pod security context for the API Gateway     | `true`               |
| `apiGatewayNginx.podSecurityContext.defaultSeccompProfile` | Use the default seccomp profile for the API Gateway     | `true`               |
| `apiGatewayNginx.podSecurityContext.fsGroup`               | Filesystem group to be used by the API Gateway          | `101`                |
| `apiGatewayNginx.containerSecurityContext`                 | Define a container security context for the API Gateway |                      |
| `apiGatewayNginx.image.registry`                           | API Gateway image registry                              | `docker.io/nginxinc` |
| `apiGatewayNginx.image.repository`                         | API Gateway image repository                            | `nginx-unprivileged` |
| `apiGatewayNginx.image.tag`                                | API Gateway image tag                                   | `1.22.0-alpine`      |
| `apiGatewayNginx.nodeSelector`                             | API Gateway node labels for pod assignment              | `{}`                 |
| `apiGatewayNginx.gracePeriod`                              | API Gateway termination grace period                    | `60`                 |
| `apiGatewayNginx.preStopHookTime`                          | API Gateway pre stop timeout                            | `20`                 |
| `apiGatewayNginx.clientMaxBodySize`                        |                                                         | `5m`                 |
| `apiGatewayNginx.sidecars`                                 | Add additional sidecar containers to the API Gateway    | `[]`                 |
| `apiGatewayNginx.extraVolumeMounts`                        | Add additional volume mounts to the API Gateway         | `[]`                 |
| `apiGatewayNginx.extraVolumes`                             | Add additional volumes to the API Gateway               | `[]`                 |
| `apiGatewayNginx.resources`                                | Define resources for the API Gateway                    |                      |


### Remediation Service

| Name                                   | Description                                                  | Value                 |
| -------------------------------------- | ------------------------------------------------------------ | --------------------- |
| `remediationService.enabled`           | Enable Remediation Service                                   | `true`                |
| `remediationService.image.registry`    | Remediation Service image registry                           | `""`                  |
| `remediationService.image.repository`  | Remediation Service image repository                         | `remediation-service` |
| `remediationService.image.tag`         | Remediation Service image tag                                | `""`                  |
| `remediationService.nodeSelector`      | Remediation Service node labels for pod assignment           | `{}`                  |
| `remediationService.gracePeriod`       | Remediation Service termination grace period                 | `60`                  |
| `remediationService.preStopHookTime`   | Remediation Service pre stop timeout                         | `5`                   |
| `remediationService.sidecars`          | Add additional sidecar containers to the Remediation Service | `[]`                  |
| `remediationService.extraVolumeMounts` | Add additional volume mounts to the Remediation Service      | `[]`                  |
| `remediationService.extraVolumes`      | Add additional volumes to the Remediation Service            | `[]`                  |
| `remediationService.resources`         | Define resources for the Remediation Service                 |                       |


### API Service

| Name                                        | Description                                                                                                                                  | Value  |
| ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| `apiService.tokenSecretName`                | K8s secret to be used as API token in the API Service                                                                                        | `nil`  |
| `apiService.image.registry`                 | API Service image registry                                                                                                                   | `""`   |
| `apiService.image.repository`               | API Service image repository                                                                                                                 | `api`  |
| `apiService.image.tag`                      | API Service image tag                                                                                                                        | `""`   |
| `apiService.maxAuth.enabled`                | Enable API authentication rate limiting                                                                                                      | `true` |
| `apiService.maxAuth.requestsPerSecond`      | API authentication rate limiting requests per second                                                                                         | `1.0`  |
| `apiService.maxAuth.requestBurst`           | API authentication rate limiting requests burst                                                                                              | `2`    |
| `apiService.eventValidation.enabled`        | Enable stricter validation of inbound events via public the event endpoint                                                                   | `true` |
| `apiService.eventValidation.maxEventSizeKB` | specifies the max. size (in KB) of inbound event accepted by the public event endpoint. This check can be disabled by providing a value <= 0 | `64`   |
| `apiService.nodeSelector`                   | API Service node labels for pod assignment                                                                                                   | `{}`   |
| `apiService.gracePeriod`                    | API Service termination grace period                                                                                                         | `60`   |
| `apiService.preStopHookTime`                | API Service pre stop timeout                                                                                                                 | `5`    |
| `apiService.sidecars`                       | Add additional sidecar containers to the API Service                                                                                         | `[]`   |
| `apiService.extraVolumeMounts`              | Add additional volume mounts to the API Service                                                                                              | `[]`   |
| `apiService.extraVolumes`                   | Add additional volumes to the API Service                                                                                                    | `[]`   |
| `apiService.resources`                      | Define resources for the API Service                                                                                                         |        |


### Bridge

| Name                              | Description                                                                                                                                                                                                                                                                    | Value     |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------- |
| `bridge.image.registry`           | Bridge image registry                                                                                                                                                                                                                                                          | `""`      |
| `bridge.image.repository`         | Bridge image repository                                                                                                                                                                                                                                                        | `bridge2` |
| `bridge.image.tag`                | Bridge image tag                                                                                                                                                                                                                                                               | `""`      |
| `bridge.cliDownloadLink`          | Define an alternative download URL for the Keptn CLI                                                                                                                                                                                                                           | `nil`     |
| `bridge.secret.enabled`           | Enable bridge credentials for HTTP Basic Auth                                                                                                                                                                                                                                  | `true`    |
| `bridge.versionCheck.enabled`     | Enable check for updated versions of Keptn                                                                                                                                                                                                                                     | `true`    |
| `bridge.showApiToken.enabled`     | If disabled, the API token will not be shown in the Bridge info                                                                                                                                                                                                                | `true`    |
| `bridge.installationType`         | Can take the values: `QUALITY_GATES`, `CONTINUOUS_OPERATIONS`, `CONTINUOUS_DELIVERY` and determines the mode in which the Bridge will be started. If only `QUALITY_GATES` is set, only functionalities and data specific for the Quality Gates Only use case will be displayed | `nil`     |
| `bridge.lookAndFeelUrl`           | Define a different styling for the Bridge by providing a URL to a ZIP archive containing style files. This archive will be downloaded and used upon Bridge startup                                                                                                             | `nil`     |
| `bridge.podSecurityContext`       | Define a pod security context for the Bridge                                                                                                                                                                                                                                   |           |
| `bridge.containerSecurityContext` | Define a container security context for the Bridge                                                                                                                                                                                                                             |           |
| `bridge.oauth`                    | Configure OAuth settings for the Bridge                                                                                                                                                                                                                                        |           |
| `bridge.authMsg`                  |                                                                                                                                                                                                                                                                                | `""`      |
| `bridge.d3.enabled`               | Enable D3 basic heatmaps in the Bridge                                                                                                                                                                                                                                         | `true`    |
| `bridge.nodeSelector`             | Bridge node labels for pod assignment                                                                                                                                                                                                                                          | `{}`      |
| `bridge.sidecars`                 | Add additional sidecar containers to the Bridge                                                                                                                                                                                                                                | `[]`      |
| `bridge.extraVolumeMounts`        | Add additional volume mounts to the Bridge                                                                                                                                                                                                                                     | `[]`      |
| `bridge.extraVolumes`             | Add additional volumes to the Bridge                                                                                                                                                                                                                                           | `[]`      |
| `bridge.resources`                | Define resources for the Bridge                                                                                                                                                                                                                                                |           |


### Distributor

| Name                                         | Description                          | Value         |
| -------------------------------------------- | ------------------------------------ | ------------- |
| `distributor.metadata.hostname`              |                                      | `nil`         |
| `distributor.metadata.namespace`             |                                      | `nil`         |
| `distributor.image.registry`                 | Distributor image registry           | `""`          |
| `distributor.image.repository`               | Distributor image repository         | `distributor` |
| `distributor.image.tag`                      | Distributor image tag                | `""`          |
| `distributor.config.proxy.httpTimeout`       |                                      | `30`          |
| `distributor.config.proxy.maxPayloadBytesKB` |                                      | `64`          |
| `distributor.config.queueGroup.enabled`      | Enable queue groups for distributor  | `true`        |
| `distributor.config.oauth.clientID`          |                                      | `""`          |
| `distributor.config.oauth.clientSecret`      |                                      | `""`          |
| `distributor.config.oauth.discovery`         |                                      | `""`          |
| `distributor.config.oauth.tokenURL`          |                                      | `""`          |
| `distributor.config.oauth.scopes`            |                                      | `""`          |
| `distributor.resources`                      | Define resources for the Distributor |               |


### Shipyard Controller

| Name                                                      | Description                                                                      | Value                 |
| --------------------------------------------------------- | -------------------------------------------------------------------------------- | --------------------- |
| `shipyardController.image.registry`                       | Shipyard Controller image registry                                               | `""`                  |
| `shipyardController.image.repository`                     | Shipyard Controller image repository                                             | `shipyard-controller` |
| `shipyardController.image.tag`                            | Shipyard Controller image tag                                                    | `""`                  |
| `shipyardController.config.taskStartedWaitDuration`       |                                                                                  | `10m`                 |
| `shipyardController.config.uniformIntegrationTTL`         |                                                                                  | `48h`                 |
| `shipyardController.config.leaderElection.enabled`        | Enable leader election when multiple replicas of Shipyard Controller are running | `false`               |
| `shipyardController.config.replicas`                      | Number of replicas of Shipyard Controller                                        | `1`                   |
| `shipyardController.config.validation.projectNameMaxSize` | Maximum number of characters that a Keptn project name can have                  | `200`                 |
| `shipyardController.config.validation.serviceNameMaxSize` | Maximum number of characters that a service name can have                        | `43`                  |
| `shipyardController.nodeSelector`                         | Shipyard Controller node labels for pod assignment                               | `{}`                  |
| `shipyardController.gracePeriod`                          | Shipyard Controller termination grace period                                     | `60`                  |
| `shipyardController.preStopHookTime`                      | Shipyard Controller pre stop timeout                                             | `15`                  |
| `shipyardController.sidecars`                             | Add additional sidecar containers to Shipyard Controller                         | `[]`                  |
| `shipyardController.extraVolumeMounts`                    | Add additional volume mounts to Shipyard Controller                              | `[]`                  |
| `shipyardController.extraVolumes`                         | Add additional volumes to Shipyard Controller                                    | `[]`                  |
| `shipyardController.resources`                            | Define resources for Shipyard Controller                                         |                       |


### Secret Service

| Name                              | Description                                             | Value            |
| --------------------------------- | ------------------------------------------------------- | ---------------- |
| `secretService.image.registry`    | Secret Service image registry                           | `""`             |
| `secretService.image.repository`  | Secret Service image repository                         | `secret-service` |
| `secretService.image.tag`         | Secret Service image tag                                | `""`             |
| `secretService.nodeSelector`      | Secret Service node labels for pod assignment           | `{}`             |
| `secretService.gracePeriod`       | Secret Service termination grace period                 | `60`             |
| `secretService.preStopHookTime`   | Secret Service pre stop timeout                         | `20`             |
| `secretService.sidecars`          | Add additional sidecar containers to the Secret Service | `[]`             |
| `secretService.extraVolumeMounts` | Add additional volume mounts to the Secret Service      | `[]`             |
| `secretService.extraVolumes`      | Add additional volumes to the Secret Service            | `[]`             |
| `secretService.resources`         | Define resources for the Secret Service                 |                  |


### Resource Service

| Name                                            | Description                                                                | Value              |
| ----------------------------------------------- | -------------------------------------------------------------------------- | ------------------ |
| `resourceService.replicas`                      | Number of replicas of Resource Service                                     | `1`                |
| `resourceService.image.registry`                | Resource Service image registry                                            | `""`               |
| `resourceService.image.repository`              | Resource Service image repository                                          | `resource-service` |
| `resourceService.image.tag`                     | Resource Service image tag                                                 | `""`               |
| `resourceService.env.GIT_KEPTN_USER`            | Default git username for the Keptn configuration git repository            | `keptn`            |
| `resourceService.env.GIT_KEPTN_EMAIL`           | Default git email address for the Keptn configuration git repository       | `keptn@keptn.sh`   |
| `resourceService.env.DIRECTORY_STAGE_STRUCTURE` | Enable directory based structure in the Keptn configuration git repository | `false`            |
| `resourceService.nodeSelector`                  | Resource Service node labels for pod assignment                            | `{}`               |
| `resourceService.gracePeriod`                   | Resource Service termination grace period                                  | `60`               |
| `resourceService.fsGroup`                       | Configure file system group ID to be used in Resource Service              | `1001`             |
| `resourceService.preStopHookTime`               | Resource Service pre stop timeout                                          | `20`               |
| `resourceService.sidecars`                      | Add additional sidecar containers to the Resource Service                  | `[]`               |
| `resourceService.extraVolumeMounts`             | Add additional volume mounts to the Resource Service                       | `[]`               |
| `resourceService.extraVolumes`                  | Add additional volumes to the Resource Service                             | `[]`               |
| `resourceService.resources`                     | Define resources for the Resource Service                                  |                    |


### MongoDB Datastore

| Name                                 | Description                                                | Value               |
| ------------------------------------ | ---------------------------------------------------------- | ------------------- |
| `mongodbDatastore.image.registry`    | MongoDB Datastore image registry                           | `""`                |
| `mongodbDatastore.image.repository`  | MongoDB Datastore image repository                         | `mongodb-datastore` |
| `mongodbDatastore.image.tag`         | MongoDB Datastore image tag                                | `""`                |
| `mongodbDatastore.nodeSelector`      | MongoDB Datastore node labels for pod assignment           | `{}`                |
| `mongodbDatastore.gracePeriod`       | MongoDB Datastore termination grace period                 | `60`                |
| `mongodbDatastore.preStopHookTime`   | MongoDB Datastore pre stop timeout                         | `20`                |
| `mongodbDatastore.sidecars`          | Add additional sidecar containers to the MongoDB Datastore | `[]`                |
| `mongodbDatastore.extraVolumeMounts` | Add additional volume mounts to the MongoDB Datastore      | `[]`                |
| `mongodbDatastore.extraVolumes`      | Add additional volumes to the MongoDB Datastore            | `[]`                |
| `mongodbDatastore.resources`         | Define resources for the MongoDB Datastore                 |                     |


### Lighthouse Service

| Name                                  | Description                                                 | Value                |
| ------------------------------------- | ----------------------------------------------------------- | -------------------- |
| `lighthouseService.enabled`           | Enable Lighthouse Service                                   | `true`               |
| `lighthouseService.image.registry`    | Lighthouse Service image registry                           | `""`                 |
| `lighthouseService.image.repository`  | Lighthouse Service image repository                         | `lighthouse-service` |
| `lighthouseService.image.tag`         | Lighthouse Service image tag                                | `""`                 |
| `lighthouseService.nodeSelector`      | Lighthouse Service node labels for pod assignment           | `{}`                 |
| `lighthouseService.gracePeriod`       | Lighthouse Service termination grace period                 | `60`                 |
| `lighthouseService.preStopHookTime`   | Lighthouse Service pre stop timeout                         | `20`                 |
| `lighthouseService.sidecars`          | Add additional sidecar containers to the Lighthouse Service | `[]`                 |
| `lighthouseService.extraVolumeMounts` | Add additional volume mounts to the Lighthouse Service      | `[]`                 |
| `lighthouseService.extraVolumes`      | Add additional volumes to the Lighthouse Service            | `[]`                 |
| `lighthouseService.resources`         | Define resources for the Lighthouse Service                 |                      |


### Statistics Service

| Name                                  | Description                                                 | Value                |
| ------------------------------------- | ----------------------------------------------------------- | -------------------- |
| `statisticsService.enabled`           | Enable Statistics Service                                   | `true`               |
| `statisticsService.image.registry`    | Statistics Service image registry                           | `""`                 |
| `statisticsService.image.repository`  | Statistics Service image repository                         | `statistics-service` |
| `statisticsService.image.tag`         | Statistics Service image tag                                | `""`                 |
| `statisticsService.nodeSelector`      | Statistics Service node labels for pod assignment           | `{}`                 |
| `statisticsService.gracePeriod`       | Statistics Service termination grace period                 | `60`                 |
| `statisticsService.preStopHookTime`   | Statistics Service pre stop timeout                         | `20`                 |
| `statisticsService.sidecars`          | Add additional sidecar containers to the Statistics Service | `[]`                 |
| `statisticsService.extraVolumeMounts` | Add additional volume mounts to the Statistics Service      | `[]`                 |
| `statisticsService.extraVolumes`      | Add additional volumes to the Statistics Service            | `[]`                 |
| `statisticsService.resources`         | Define resources for the Statistics Service                 |                      |


### Approval Service

| Name                                | Description                                               | Value              |
| ----------------------------------- | --------------------------------------------------------- | ------------------ |
| `approvalService.enabled`           | Enable Approval Service                                   | `true`             |
| `approvalService.image.registry`    | Approval Service image registry                           | `""`               |
| `approvalService.image.repository`  | Approval Service image repository                         | `approval-service` |
| `approvalService.image.tag`         | Approval Service image tag                                | `""`               |
| `approvalService.nodeSelector`      | Approval Service node labels for pod assignment           | `{}`               |
| `approvalService.gracePeriod`       | Approval Service termination grace period                 | `60`               |
| `approvalService.preStopHookTime`   | Approval Service pre stop timeout                         | `5`                |
| `approvalService.sidecars`          | Add additional sidecar containers to the Approval Service | `[]`               |
| `approvalService.extraVolumeMounts` | Add additional volume mounts to the Approval Service      | `[]`               |
| `approvalService.extraVolumes`      | Add additional volumes to the Approval Service            | `[]`               |
| `approvalService.resources`         | Define resources for the Approval Service                 |                    |


### Webhook Service

| Name                               | Description                                              | Value             |
| ---------------------------------- | -------------------------------------------------------- | ----------------- |
| `webhookService.enabled`           | Enable Webhook Service                                   | `true`            |
| `webhookService.image.registry`    | Webhook Service image registry                           | `""`              |
| `webhookService.image.repository`  | Webhook Service image repository                         | `webhook-service` |
| `webhookService.image.tag`         | Webhook Service image tag                                | `""`              |
| `webhookService.nodeSelector`      | Webhook Service node labels for pod assignment           | `{}`              |
| `webhookService.gracePeriod`       | Webhook Service termination grace period                 | `60`              |
| `webhookService.preStopHookTime`   | Webhook Service pre stop timeout                         | `20`              |
| `webhookService.sidecars`          | Add additional sidecar containers to the Webhook Service | `[]`              |
| `webhookService.extraVolumeMounts` | Add additional volume mounts to the Webhook Service      | `[]`              |
| `webhookService.extraVolumes`      | Add additional volumes to the Webhook Service            | `[]`              |
| `webhookService.resources`         | Define resources for the Webhook Service                 |                   |


### Ingress

| Name                  | Description                            | Value    |
| --------------------- | -------------------------------------- | -------- |
| `ingress.enabled`     | Enable ingress configuration for Keptn | `false`  |
| `ingress.annotations` | Keptn Ingress annotations              | `{}`     |
| `ingress.host`        | Keptn Ingress host URL                 | `{}`     |
| `ingress.path`        | Keptn Ingress host path                | `/`      |
| `ingress.pathType`    | Keptn Ingress path type                | `Prefix` |
| `ingress.className`   | Keptn Ingress class name               | `""`     |
| `ingress.tls`         | Keptn Ingress TLS configuration        | `[]`     |


### Common settings

| Name                                    | Description                                                            | Value           |
| --------------------------------------- | ---------------------------------------------------------------------- | --------------- |
| `logLevel`                              | Global log level for all Keptn services                                | `info`          |
| `prefixPath`                            | URL prefix for all Keptn URLs                                          | `""`            |
| `keptnSpecVersion`                      | Version of the Keptn Spec definitions to be used                       | `latest`        |
| `strategy.type`                         | Strategy to use to replace existing Keptn pods                         | `RollingUpdate` |
| `strategy.rollingUpdate.maxSurge`       | Maximum number of additional pods to be spun up during rolling updates | `1`             |
| `strategy.rollingUpdate.maxUnavailable` | Maximum number of unavailable pods during rolling updates              | `0`             |
| `podSecurityContext`                    | Set the default pod security context for all pods                      |                 |
| `podSecurityContext.enabled`            | Enable the default pod security context for all pods                   | `true`          |
| `containerSecurityContext`              | Set the default container security context for all containers          |                 |
| `containerSecurityContext.enabled`      | Enable the default container security context for all containers       | `true`          |
| `nodeSelector`                          | Default node labels for pod assignment                                 | `{}`            |

