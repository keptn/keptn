# Resource Service :: The New Configuration Service

The *resource-service* is a Keptn core component used to manage resources for Keptn project-related entities,
i.e., project, stage, and service. The entity model is shown below. To store the resources with version control, a git
repository is used that is mounted as emptyDir volume.  Besides, this service has functionality to upload the git repository
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
When the resource-service is enabled by default, the Kubernetes `Service` will change name to `resource-service`.

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
In the backup data, you will find a folder for each Keptn project. For each project, open a shell in that directory and make sure the `Git` CLI is available.
Attach your upstream to the Keptn project via the Git CLI with `git remote add origin <remoteURL>`, where `<remoteURL>` is your Git upstream.
Afterward, run `git push --all` to synchronize your backup with your Git repository.
Finally, you can navigate to your Bridge installation and configure the upstream to your project.
