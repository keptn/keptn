# Release Notes 0.8.2-hotfix1

This is a hotfix release for Keptn 0.8.2 `helm-service` if deployed as an execution plane service.

---

## Fixes

- Fix duplicated Helm Deployment.Started/Finished CloudEvents when using helm-service as a remote execution plane [3888](https://github.com/keptn/keptn/issues/3888)

## Upgrade to 0.8.2-hotfix1

You only need to upgrade `helm-service` if [deployed in execution-plane (Multi-cluster setup)](https://keptn.sh/docs/0.8.x/operate/multi_cluster/) as follows:

```console
helm upgrade helm-service https://github.com/keptn/keptn/releases/download/0.8.2-hotfix1/helm-service-0.8.2-hotfix1.tgz -n keptn-exec
```

It is not required to upgrade the cli or any other services for this hotfix release.
