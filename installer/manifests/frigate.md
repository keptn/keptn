
Keptn
===========

Cloud-native application life-cycle orchestration


## Configuration

The following table lists the configurable parameters of the Keptn chart and their default values.

| Parameter                | Description             | Default        |
| ------------------------ | ----------------------- | -------------- |
| `global.keptn.registry` |  | `"docker.io/keptn"` |
| `global.keptn.tag` |  | `""` |
| `continuousDelivery.enabled` |  | `false` |
| `continuousDelivery.ingressConfig.ingress_hostname_suffix` |  | `"svc.cluster.local"` |
| `continuousDelivery.ingressConfig.ingress_protocol` |  | `"http"` |
| `continuousDelivery.ingressConfig.ingress_port` |  | `"80"` |
| `continuousDelivery.ingressConfig.istio_gateway` |  | `"public-gateway.istio-system"` |
| `mongo.enabled` |  | `true` |
| `mongo.host` |  | `"mongodb:27017"` |
| `mongo.architecture` |  | `"standalone"` |
| `mongo.service.nameOverride` |  | `"mongo"` |
| `mongo.service.port` |  | `27017` |
| `mongo.auth.enabled` |  | `true` |
| `mongo.auth.databases` |  | `["keptn"]` |
| `mongo.auth.existingSecret` | If the password and rootPassword values below are used, remove this field. | `"mongodb-credentials"` |
| `mongo.auth.usernames` |  | `["keptn"]` |
| `mongo.auth.password` |  | `null` |
| `mongo.auth.rootUser` |  | `"admin"` |
| `mongo.auth.rootPassword` |  | `null` |
| `mongo.auth.bridgeAuthDatabase` |  | `"keptn"` |
| `mongo.external.connectionString` |  | `null` |
| `mongo.containerSecurityContext.allowPrivilegeEscalation` |  | `false` |
| `mongo.containerSecurityContext.capabilities.drop` |  | `["ALL"]` |
| `mongo.serviceAccount.automountServiceAccountToken` |  | `false` |
| `mongo.resources.requests.cpu` |  | `"200m"` |
| `mongo.resources.requests.memory` |  | `"100Mi"` |
| `mongo.resources.limits.cpu` |  | `"1000m"` |
| `mongo.resources.limits.memory` |  | `"500Mi"` |
| `prefixPath` |  | `""` |
| `keptnSpecVersion` |  | `"latest"` |
| `features.automaticProvisioning.serviceURL` |  | `""` |
| `features.automaticProvisioning.message` |  | `""` |
| `features.swagger.hideDeprecated` |  | `false` |
| `features.oauth.enabled` |  | `false` |
| `features.oauth.prefix` |  | `"keptn:"` |
| `nats.nameOverride` |  | `"keptn-nats"` |
| `nats.fullnameOverride` |  | `"keptn-nats"` |
| `nats.cluster.replicas` |  | `3` |
| `nats.cluster.name` |  | `"nats"` |
| `nats.securityContext.runAsNonRoot` |  | `true` |
| `nats.securityContext.runAsUser` |  | `10001` |
| `nats.securityContext.fsGroup` |  | `10001` |
| `nats.nats.resources.requests.cpu` |  | `"200m"` |
| `nats.nats.resources.requests.memory` |  | `"500Mi"` |
| `nats.nats.resources.limits.cpu` |  | `"500m"` |
| `nats.nats.resources.limits.memory` |  | `"1Gi"` |
| `nats.nats.healthcheck.startup.enabled` |  | `false` |
| `nats.nats.jetstream.enabled` |  | `true` |
| `nats.nats.jetstream.memStorage.enabled` |  | `true` |
| `nats.nats.jetstream.memStorage.size` |  | `"500Mi"` |
| `nats.nats.jetstream.fileStorage.enabled` |  | `true` |
| `nats.nats.jetstream.fileStorage.size` |  | `"5Gi"` |
| `nats.nats.jetstream.fileStorage.storageDirectory` |  | `"/data/"` |
| `nats.nats.jetstream.fileStorage.storageClassName` |  | `""` |
| `nats.nats.securityContext.readOnlyRootFilesystem` |  | `true` |
| `nats.nats.securityContext.allowPrivilegeEscalation` |  | `false` |
| `nats.nats.securityContext.runAsNonRoot` |  | `true` |
| `nats.nats.securityContext.capabilities.drop` |  | `["ALL"]` |
| `nats.natsbox.enabled` |  | `false` |
| `nats.reloader.enabled` |  | `false` |
| `nats.exporter.enabled` |  | `false` |
| `apiGatewayNginx.type` |  | `"ClusterIP"` |
| `apiGatewayNginx.port` |  | `80` |
| `apiGatewayNginx.targetPort` |  | `8080` |
| `apiGatewayNginx.nodePort` |  | `31090` |
| `apiGatewayNginx.podSecurityContext.enabled` |  | `true` |
| `apiGatewayNginx.podSecurityContext.defaultSeccompProfile` |  | `true` |
| `apiGatewayNginx.podSecurityContext.fsGroup` |  | `101` |
| `apiGatewayNginx.containerSecurityContext.enabled` |  | `true` |
| `apiGatewayNginx.containerSecurityContext.runAsNonRoot` |  | `true` |
| `apiGatewayNginx.containerSecurityContext.runAsUser` |  | `101` |
| `apiGatewayNginx.containerSecurityContext.readOnlyRootFilesystem` |  | `false` |
| `apiGatewayNginx.containerSecurityContext.allowPrivilegeEscalation` |  | `false` |
| `apiGatewayNginx.containerSecurityContext.privileged` |  | `false` |
| `apiGatewayNginx.containerSecurityContext.capabilities.drop` |  | `["ALL"]` |
| `apiGatewayNginx.image.registry` | Container Registry | `"docker.io/nginxinc"` |
| `apiGatewayNginx.image.repository` | Container Image Name | `"nginx-unprivileged"` |
| `apiGatewayNginx.image.tag` | Container Tag | `"1.22.0-alpine"` |
| `apiGatewayNginx.nodeSelector` |  | `{}` |
| `apiGatewayNginx.gracePeriod` |  | `60` |
| `apiGatewayNginx.preStopHookTime` |  | `20` |
| `apiGatewayNginx.clientMaxBodySize` |  | `"5m"` |
| `apiGatewayNginx.sidecars` |  | `[]` |
| `apiGatewayNginx.extraVolumeMounts` |  | `[]` |
| `apiGatewayNginx.extraVolumes` |  | `[]` |
| `apiGatewayNginx.resources.requests.memory` |  | `"64Mi"` |
| `apiGatewayNginx.resources.requests.cpu` |  | `"50m"` |
| `apiGatewayNginx.resources.limits.memory` |  | `"128Mi"` |
| `apiGatewayNginx.resources.limits.cpu` |  | `"100m"` |
| `remediationService.enabled` |  | `true` |
| `remediationService.image.registry` | Container Registry | `""` |
| `remediationService.image.repository` | Container Image Name | `"remediation-service"` |
| `remediationService.image.tag` | Container Tag | `""` |
| `remediationService.nodeSelector` |  | `{}` |
| `remediationService.gracePeriod` |  | `60` |
| `remediationService.preStopHookTime` |  | `5` |
| `remediationService.sidecars` |  | `[]` |
| `remediationService.extraVolumeMounts` |  | `[]` |
| `remediationService.extraVolumes` |  | `[]` |
| `remediationService.resources.requests.memory` |  | `"64Mi"` |
| `remediationService.resources.requests.cpu` |  | `"50m"` |
| `remediationService.resources.limits.memory` |  | `"1Gi"` |
| `remediationService.resources.limits.cpu` |  | `"200m"` |
| `apiService.tokenSecretName` |  | `null` |
| `apiService.image.registry` | Container Registry | `""` |
| `apiService.image.repository` | Container Image Name | `"api"` |
| `apiService.image.tag` | Container Tag | `""` |
| `apiService.maxAuth.enabled` |  | `true` |
| `apiService.maxAuth.requestsPerSecond` |  | `"1.0"` |
| `apiService.maxAuth.requestBurst` |  | `"2"` |
| `apiService.nodeSelector` |  | `{}` |
| `apiService.gracePeriod` |  | `60` |
| `apiService.preStopHookTime` |  | `5` |
| `apiService.sidecars` |  | `[]` |
| `apiService.extraVolumeMounts` |  | `[]` |
| `apiService.extraVolumes` |  | `[]` |
| `apiService.resources.requests.memory` |  | `"32Mi"` |
| `apiService.resources.requests.cpu` |  | `"50m"` |
| `apiService.resources.limits.memory` |  | `"64Mi"` |
| `apiService.resources.limits.cpu` |  | `"100m"` |
| `bridge.image.registry` | Container Registry | `""` |
| `bridge.image.repository` | Container Image Name | `"bridge2"` |
| `bridge.image.tag` | Container Tag | `""` |
| `bridge.cliDownloadLink` |  | `null` |
| `bridge.secret.enabled` |  | `true` |
| `bridge.versionCheck.enabled` |  | `true` |
| `bridge.showApiToken.enabled` |  | `true` |
| `bridge.installationType` |  | `null` |
| `bridge.lookAndFeelUrl` |  | `null` |
| `bridge.podSecurityContext.enabled` |  | `true` |
| `bridge.podSecurityContext.defaultSeccompProfile` |  | `true` |
| `bridge.podSecurityContext.fsGroup` |  | `65532` |
| `bridge.containerSecurityContext.enabled` |  | `true` |
| `bridge.containerSecurityContext.runAsNonRoot` |  | `true` |
| `bridge.containerSecurityContext.runAsUser` |  | `65532` |
| `bridge.containerSecurityContext.readOnlyRootFilesystem` |  | `true` |
| `bridge.containerSecurityContext.allowPrivilegeEscalation` |  | `false` |
| `bridge.containerSecurityContext.privileged` |  | `false` |
| `bridge.containerSecurityContext.capabilities.drop` |  | `["ALL"]` |
| `bridge.oauth.discovery` |  | `""` |
| `bridge.oauth.secureCookie` |  | `false` |
| `bridge.oauth.trustProxy` |  | `""` |
| `bridge.oauth.sessionTimeoutMin` |  | `""` |
| `bridge.oauth.sessionValidatingTimeoutMin` |  | `""` |
| `bridge.oauth.baseUrl` |  | `""` |
| `bridge.oauth.clientID` |  | `""` |
| `bridge.oauth.clientSecret` |  | `""` |
| `bridge.oauth.IDTokenAlg` |  | `""` |
| `bridge.oauth.scope` |  | `""` |
| `bridge.oauth.userIdentifier` |  | `""` |
| `bridge.oauth.mongoConnectionString` |  | `""` |
| `bridge.authMsg` |  | `""` |
| `bridge.d3heatmap.enabled` |  | `false` |
| `bridge.nodeSelector` |  | `{}` |
| `bridge.sidecars` |  | `[]` |
| `bridge.extraVolumeMounts` |  | `[]` |
| `bridge.extraVolumes` |  | `[]` |
| `bridge.resources.requests.memory` |  | `"64Mi"` |
| `bridge.resources.requests.cpu` |  | `"25m"` |
| `bridge.resources.limits.memory` |  | `"256Mi"` |
| `bridge.resources.limits.cpu` |  | `"200m"` |
| `distributor.metadata.hostname` |  | `null` |
| `distributor.metadata.namespace` |  | `null` |
| `distributor.image.registry` | Container Registry | `""` |
| `distributor.image.repository` | Container Image Name | `"distributor"` |
| `distributor.image.tag` | Container Tag | `""` |
| `distributor.config.proxy.httpTimeout` |  | `"30"` |
| `distributor.config.proxy.maxPayloadBytesKB` |  | `"64"` |
| `distributor.config.queueGroup.enabled` |  | `true` |
| `distributor.config.oauth.clientID` |  | `""` |
| `distributor.config.oauth.clientSecret` |  | `""` |
| `distributor.config.oauth.discovery` |  | `""` |
| `distributor.config.oauth.tokenURL` |  | `""` |
| `distributor.config.oauth.scopes` |  | `""` |
| `distributor.resources.requests.memory` |  | `"16Mi"` |
| `distributor.resources.requests.cpu` |  | `"25m"` |
| `distributor.resources.limits.memory` |  | `"32Mi"` |
| `distributor.resources.limits.cpu` |  | `"100m"` |
| `shipyardController.image.registry` | Container Registry | `""` |
| `shipyardController.image.repository` | Container Image Name | `"shipyard-controller"` |
| `shipyardController.image.tag` | Container Tag | `""` |
| `shipyardController.config.taskStartedWaitDuration` |  | `"10m"` |
| `shipyardController.config.uniformIntegrationTTL` |  | `"48h"` |
| `shipyardController.config.disableLeaderElection` |  | `true` |
| `shipyardController.config.replicas` |  | `1` |
| `shipyardController.config.validation.projectNameMaxSize` |  | `200` |
| `shipyardController.config.validation.serviceNameMaxSize` |  | `43` |
| `shipyardController.nodeSelector` |  | `{}` |
| `shipyardController.gracePeriod` |  | `60` |
| `shipyardController.preStopHookTime` |  | `15` |
| `shipyardController.sidecars` |  | `[]` |
| `shipyardController.extraVolumeMounts` |  | `[]` |
| `shipyardController.extraVolumes` |  | `[]` |
| `shipyardController.resources.requests.memory` |  | `"32Mi"` |
| `shipyardController.resources.requests.cpu` |  | `"50m"` |
| `shipyardController.resources.limits.memory` |  | `"128Mi"` |
| `shipyardController.resources.limits.cpu` |  | `"100m"` |
| `secretService.image.registry` | Container Registry | `""` |
| `secretService.image.repository` | Container Image Name | `"secret-service"` |
| `secretService.image.tag` | Container Tag | `""` |
| `secretService.nodeSelector` |  | `{}` |
| `secretService.gracePeriod` |  | `60` |
| `secretService.preStopHookTime` |  | `20` |
| `secretService.sidecars` |  | `[]` |
| `secretService.extraVolumeMounts` |  | `[]` |
| `secretService.extraVolumes` |  | `[]` |
| `secretService.resources.requests.memory` |  | `"32Mi"` |
| `secretService.resources.requests.cpu` |  | `"25m"` |
| `secretService.resources.limits.memory` |  | `"64Mi"` |
| `secretService.resources.limits.cpu` |  | `"200m"` |
| `configurationService.image.registry` | Container Registry | `""` |
| `configurationService.image.repository` | Container Image Name | `"configuration-service"` |
| `configurationService.image.tag` | Container Tag | `""` |
| `configurationService.storage` |  | `"100Mi"` |
| `configurationService.storageClass` |  | `null` |
| `configurationService.fsGroup` |  | `1001` |
| `configurationService.initContainer` |  | `true` |
| `configurationService.env.GIT_KEPTN_USER` |  | `"keptn"` |
| `configurationService.env.GIT_KEPTN_EMAIL` |  | `"keptn@keptn.sh"` |
| `configurationService.nodeSelector` |  | `{}` |
| `configurationService.gracePeriod` |  | `60` |
| `configurationService.preStopHookTime` |  | `20` |
| `configurationService.sidecars` |  | `[]` |
| `configurationService.extraVolumeMounts` |  | `[]` |
| `configurationService.extraVolumes` |  | `[]` |
| `configurationService.resources.requests.memory` |  | `"32Mi"` |
| `configurationService.resources.requests.cpu` |  | `"25m"` |
| `configurationService.resources.limits.memory` |  | `"256Mi"` |
| `configurationService.resources.limits.cpu` |  | `"100m"` |
| `resourceService.replicas` |  | `1` |
| `resourceService.image.registry` | Container Registry | `""` |
| `resourceService.image.repository` | Container Image Name | `"resource-service"` |
| `resourceService.image.tag` | Container Tag | `""` |
| `resourceService.env.GIT_KEPTN_USER` |  | `"keptn"` |
| `resourceService.env.GIT_KEPTN_EMAIL` |  | `"keptn@keptn.sh"` |
| `resourceService.env.DIRECTORY_STAGE_STRUCTURE` |  | `"false"` |
| `resourceService.nodeSelector` |  | `{}` |
| `resourceService.gracePeriod` |  | `60` |
| `resourceService.preStopHookTime` |  | `20` |
| `resourceService.sidecars` |  | `[]` |
| `resourceService.extraVolumeMounts` |  | `[]` |
| `resourceService.extraVolumes` |  | `[]` |
| `resourceService.resources.requests.memory` |  | `"32Mi"` |
| `resourceService.resources.requests.cpu` |  | `"25m"` |
| `resourceService.resources.limits.memory` |  | `"64Mi"` |
| `resourceService.resources.limits.cpu` |  | `"100m"` |
| `mongodbDatastore.image.registry` | Container Registry | `""` |
| `mongodbDatastore.image.repository` | Container Image Name | `"mongodb-datastore"` |
| `mongodbDatastore.image.tag` | Container Tag | `""` |
| `mongodbDatastore.nodeSelector` |  | `{}` |
| `mongodbDatastore.gracePeriod` |  | `60` |
| `mongodbDatastore.preStopHookTime` |  | `20` |
| `mongodbDatastore.sidecars` |  | `[]` |
| `mongodbDatastore.extraVolumeMounts` |  | `[]` |
| `mongodbDatastore.extraVolumes` |  | `[]` |
| `mongodbDatastore.resources.requests.memory` |  | `"32Mi"` |
| `mongodbDatastore.resources.requests.cpu` |  | `"50m"` |
| `mongodbDatastore.resources.limits.memory` |  | `"512Mi"` |
| `mongodbDatastore.resources.limits.cpu` |  | `"300m"` |
| `lighthouseService.enabled` |  | `true` |
| `lighthouseService.image.registry` | Container Registry | `""` |
| `lighthouseService.image.repository` | Container Image Name | `"lighthouse-service"` |
| `lighthouseService.image.tag` | Container Tag | `""` |
| `lighthouseService.nodeSelector` |  | `{}` |
| `lighthouseService.gracePeriod` |  | `60` |
| `lighthouseService.preStopHookTime` |  | `20` |
| `lighthouseService.sidecars` |  | `[]` |
| `lighthouseService.extraVolumeMounts` |  | `[]` |
| `lighthouseService.extraVolumes` |  | `[]` |
| `lighthouseService.resources.requests.memory` |  | `"128Mi"` |
| `lighthouseService.resources.requests.cpu` |  | `"50m"` |
| `lighthouseService.resources.limits.memory` |  | `"1Gi"` |
| `lighthouseService.resources.limits.cpu` |  | `"200m"` |
| `statisticsService.enabled` |  | `true` |
| `statisticsService.image.registry` | Container Registry | `""` |
| `statisticsService.image.repository` | Container Image Name | `"statistics-service"` |
| `statisticsService.image.tag` | Container Tag | `""` |
| `statisticsService.nodeSelector` |  | `{}` |
| `statisticsService.gracePeriod` |  | `60` |
| `statisticsService.preStopHookTime` |  | `20` |
| `statisticsService.sidecars` |  | `[]` |
| `statisticsService.extraVolumeMounts` |  | `[]` |
| `statisticsService.extraVolumes` |  | `[]` |
| `statisticsService.resources.requests.memory` |  | `"32Mi"` |
| `statisticsService.resources.requests.cpu` |  | `"25m"` |
| `statisticsService.resources.limits.memory` |  | `"64Mi"` |
| `statisticsService.resources.limits.cpu` |  | `"100m"` |
| `approvalService.enabled` |  | `true` |
| `approvalService.image.registry` | Container Registry | `""` |
| `approvalService.image.repository` | Container Image Name | `"approval-service"` |
| `approvalService.image.tag` | Container Tag | `""` |
| `approvalService.nodeSelector` |  | `{}` |
| `approvalService.gracePeriod` |  | `60` |
| `approvalService.preStopHookTime` |  | `5` |
| `approvalService.sidecars` |  | `[]` |
| `approvalService.extraVolumeMounts` |  | `[]` |
| `approvalService.extraVolumes` |  | `[]` |
| `approvalService.resources.requests.memory` |  | `"32Mi"` |
| `approvalService.resources.requests.cpu` |  | `"25m"` |
| `approvalService.resources.limits.memory` |  | `"128Mi"` |
| `approvalService.resources.limits.cpu` |  | `"100m"` |
| `webhookService.enabled` |  | `true` |
| `webhookService.image.registry` | Container Registry | `""` |
| `webhookService.image.repository` | Container Image Name | `"webhook-service"` |
| `webhookService.image.tag` | Container Tag | `""` |
| `webhookService.nodeSelector` |  | `{}` |
| `webhookService.gracePeriod` |  | `60` |
| `webhookService.preStopHookTime` |  | `20` |
| `webhookService.sidecars` |  | `[]` |
| `webhookService.extraVolumeMounts` |  | `[]` |
| `webhookService.extraVolumes` |  | `[]` |
| `webhookService.resources.requests.memory` |  | `"32Mi"` |
| `webhookService.resources.requests.cpu` |  | `"25m"` |
| `webhookService.resources.limits.memory` |  | `"64Mi"` |
| `webhookService.resources.limits.cpu` |  | `"100m"` |
| `ingress.enabled` |  | `false` |
| `ingress.annotations` |  | `{}` |
| `ingress.host` |  | `{}` |
| `ingress.path` |  | `"/"` |
| `ingress.pathType` |  | `"Prefix"` |
| `ingress.className` |  | `""` |
| `ingress.tls` |  | `[]` |
| `logLevel` |  | `"info"` |
| `strategy.type` |  | `"RollingUpdate"` |
| `strategy.rollingUpdate.maxSurge` |  | `1` |
| `strategy.rollingUpdate.maxUnavailable` |  | `0` |
| `podSecurityContext.enabled` |  | `true` |
| `podSecurityContext.defaultSeccompProfile` |  | `true` |
| `podSecurityContext.fsGroup` |  | `65532` |
| `containerSecurityContext.enabled` |  | `true` |
| `containerSecurityContext.runAsNonRoot` |  | `true` |
| `containerSecurityContext.runAsUser` |  | `65532` |
| `containerSecurityContext.readOnlyRootFilesystem` |  | `true` |
| `containerSecurityContext.allowPrivilegeEscalation` |  | `false` |
| `containerSecurityContext.privileged` |  | `false` |
| `containerSecurityContext.capabilities.drop` |  | `["ALL"]` |
| `nodeSelector` |  | `{}` |



---
_Documentation generated by [Frigate](https://frigate.readthedocs.io)._

