Helm-service
===========

Helm Chart for the keptn helm-service

## Configuration

The following table lists the configurable parameters of the Helm-service chart and their default values.

| Parameter                               | Description                                                       | Default                                         |
|-----------------------------------------|-------------------------------------------------------------------|-------------------------------------------------|
| `nameOverride`                          | Override the name label and distributor subscription name         | `""`                                            |
| `global.keptn.registry`                 | Container repository/registry name. Will be set for all services. | `"docker.io/keptn/"`                            |
| `global.keptn.tag`                      | Container tag. Will be set for all services.                      | `""`                                            |
| `helmservice.image.registry`            | Container image name                                              | `"${global.keptn.image.registry}/helm-service"` |
| `helmservice.image.pullPolicy`          | Kubernetes image pull policy                                      | `"IfNotPresent"`                                |
| `helmservice.image.tag`                 | Container tag                                                     | `global.keptn.image.tag`                        |
| `helmservice.service.enabled`           | Creates a kubernetes service for the helm-service                 | `true`                                          |
| `distributor.stageFilter`               | Sets the stage this helm service belongs to                       | `""`                                            |
| `distributor.serviceFilter`             | Sets the service this helm service belongs to                     | `""`                                            |
| `distributor.projectFilter`             | Sets the project this helm service belongs to                     | `""`                                            |
| `distributor.image.registry`            | Container image name                                              | `"${global.keptn.image.registry}/distributor"`  |
| `distributor.image.pullPolicy`          | Kubernetes image pull policy                                      | `"IfNotPresent"`                                |
| `distributor.image.tag`                 | Container tag                                                     | `global.keptn.image.tag`                        |
| `remoteControlPlane.enabled`            | Enables remote execution plane mode                               | `false`                                         |
| `remoteControlPlane.api.protocol`       | Used protocol (http, https                                        | `"https"`                                       |
| `remoteControlPlane.api.hostname`       | Hostname of the control plane cluster (and port)                  | `""`                                            |
| `remoteControlPlane.api.apiValidateTls` | Defines if the control plane certificate should be validated      | `true`                                          |
| `remoteControlPlane.api.token`          | Keptn api token                                                   | `""`                                            |
| `imagePullSecrets`                      | Secrets to use for container registry credentials                 | `[]`                                            |
| `serviceAccount.create`                 | Enables the service account creation                              | `true`                                          |
| `serviceAccount.annotations`            | Annotations to add to the service account                         | `{}`                                            |
| `serviceAccount.name`                   | The name of the service account to use.                           | `""`                                            |
| `podAnnotations`                        | Annotations to add to the created pods                            | `{}`                                            |
| `podSecurityContext`                    | Set the pod security context (e.g. fsgroups)                      | `{}`                                            |
| `securityContext`                       | Set the security context (e.g. runasuser)                         | `{}`                                            |
| `resources`                             | Resource limits and requests                                      | `{}`                                            |
| `nodeSelector`                          | Node selector configuration                                       | `{}`                                            |
| `tolerations`                           | Tolerations for the pods                                          | `[]`                                            |
| `affinity`                              | Affinity rules                                                    | `{}`                                            |





