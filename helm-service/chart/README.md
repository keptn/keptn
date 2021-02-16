
Helm-service
===========

Helm Chart for the keptn helm-service


## Configuration

The following table lists the configurable parameters of the Helm-service chart and their default values.

| Parameter                | Description             | Default        |
| ------------------------ | ----------------------- | -------------- |
| `helmservice.image.repository` | Container image name | `"docker.io/keptn/helm-service"` |
| `helmservice.image.pullPolicy` | Kubernetes image pull policy | `"IfNotPresent"` |
| `helmservice.image.tag` | Container tag | `""` |
| `helmservice.service.enabled` | Creates a kubernetes service for the helm-service | `true` |
| `distributor.image.repository` | Container image name | `"docker.io/keptn/distributor"` |
| `distributor.image.pullPolicy` | Kubernetes image pull policy | `"IfNotPresent"` |
| `distributor.image.tag` | Container tag | `""` |
| `remoteControlPlane.enabled` | Enables remote execution plane mode | `false` |
| `remoteControlPlane.api.protocol` | Used protocol (http, https | `"http"` |
| `remoteControlPlane.api.hostname` | Hostname of the control plane cluster (and port) | `""` |
| `remoteControlPlane.api.apiValidateTls` | Defines if the control plane certificate should be validated | `true` |
| `remoteControlPlane.api.stageFilter` | Sets the stage this helm service belongs to | `""` |
| `remoteControlPlane.api.token` | Keptn api token | `""` |
| `imagePullSecrets` | Secrets to use for container registry credentials | `[]` |
| `serviceAccount.create` | Enables the service account creation | `true` |
| `serviceAccount.annotations` | Annotations to add to the service account | `{}` |
| `serviceAccount.name` | The name of the service account to use. | `""` |
| `podAnnotations` | Annotations to add to the created pods | `{}` |
| `podSecurityContext` | Set the pod security context (e.g. fsgroups) | `{}` |
| `securityContext.readOnlyRootFilesystem` |  | `true` |
| `securityContext.runAsNonRoot` |  | `true` |
| `securityContext.runAsUser` |  | `1000` |
| `resources` | Resource limits and requests | `{}` |
| `nodeSelector` | Node selector configuration | `{}` |
| `tolerations` | Tolerations for the pods | `[]` |
| `affinity` | Affinity rules | `{}` |





