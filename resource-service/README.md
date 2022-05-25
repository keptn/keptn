# Resource Service :: The New Configuration Service

The *resource-service* is a Keptn core component used to manage resources for Keptn project-related entities,
i.e., project, stage, and service. The entity model is shown below. To store the resources with version control, a Git
repository is used that is mounted as emptyDir volume.  Besides, this service has functionality to upload the Git repository
to any Git-based service such as GitLab, GitHub, Bitbucket, etc.

The *resource-service* has been designed from the ground up to work with a remote upstream.
Hence, Keptn projects must always have a Git repository configured. Furthermore, the *resource-service* does **not** have the requirement of using uninitialized repositories.
These changes allow the service implementation to be more flexible and faster in retrieving and storing Keptn data comparing it to the *configuration-service*.

## Entity model

```
------------          ------------          ------------
|          | 1        |          | 1        |          |
| Project  |----------|  Stage   |----------| Service  |
|          |        * |          |        * |          |
------------          ------------          ------------
  1 \                   1  \                   1  \
     \ *                    \ *                    \ *
   ------------           ------------           ------------
   |          |           |          |           |          |
   | Resource |           | Resource |           | Resource |
   |          |           |          |           |          |
   ------------           ------------           ------------
```

## Installation

The *resource-service* replaces the *configuration-service*, hence only one of the two can be run at the same time.
The *resource-service* can be enabled during the installation of Keptn setting the Helm value `control-plane.resourceService.enabled` to `true`.
This flag changes the *configuration-service* `Service` to point towards the *resource-service* `Pod`.
In the future, the *resource-service* will be enabled by default. With this, we will remove the `configuration-service` Kubernetes `Service` in favor of a `resource-service` Kubernetes `Service`.

### Deploy it directly into your Kubernetes cluster

To deploy the current version of the *resource-service* in your Keptn Kubernetes cluster,
use the file `deploy/service.yaml` from this repository and apply it.

```console
kubectl apply -f deploy/service.yaml
```

### Delete it from your Kubernetes cluster

To delete a deployed *resource-service*, use the file `deploy/service.yaml` from this repository
and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## Migration from the configuration-service

Before migrating from the *configuration-service* to the *resource-service* it is recommended to (i) attach an upstream to your Keptn projects and (ii) do a [backup](https://keptn.sh/docs/0.15.x/operate/backup_and_restore/#back-up-configuration-service). If you set an upstream for all your Keptn projects, no additional steps are required.

Suppose you need the additional features provided by the *resource-service*,  such as HTTPS/SSH or Proxy, to configure your Keptn project with an upstream. In that case,
you can also deploy the *resource-service* and configure the Git repositories later. For this, a backup is necessary.

1. Back up of the [configuration-service](https://keptn.sh/docs/0.15.x/operate/backup_and_restore/#back-up-configuration-service).
2. For each Keptn project in the backup data open a shell in that directory and make sure the `Git` CLI is available.
3. Attach your upstream to the Keptn project via the Git CLI with `git remote add origin <remoteURL>`, where `<remoteURL>` is your Git upstream.
4. Run `git push --all` to synchronize your backup with your Git repository.
5. Install Keptn with the *resource-service* enabled
6. Navigate to your Bridge installation and configure an upstream to the Keptn projects.

