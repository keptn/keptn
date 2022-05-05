# Install master branch

**Attention:** This installs a potentially unstable and unreliable version of Keptn and is not meant to preview features, but for Keptn core/contrib/sandbox developers.

**DO NOT PERFORM THIS ON A PRODUCTION ENVIRONMENT**

1. Create a new cluster (e.g., using k3s)
2. - Build Keptn CLI from master:
      - download source, 
      - go to `/keptn/cli` and run `go build  -o keptn main.go`
   - Or use an older CLI from a previous release, e.g., `curl -sL https://get.keptn.sh | KEPTN_VERSION=0.12.0 bash`
4. Add the Keptn developer helm repository:
```
helm repo add keptn-dev https://charts-dev.keptn.sh
```
5. Run `helm repo update`
6. Show all developer images tag with `helm search repo keptn-dev --devel -l `
7. For current master use `<current_version>-dev` images, for images built in a PR add `-PR-<PR_#>` to the tag,
8. To install version 0.13.0-dev-PR-1234 
``` 
helm upgrade --install keptn keptn-dev/keptn -n keptn --create-namespace --set=continuousDelivery.enabled=true --wait --version 0.13.0-dev-PR-1234
helm upgrade --install jmeter-service keptn-dev/jmeter-service -n keptn --wait --version 0.13.0-dev-PR-1234
helm upgrade --install helm-service keptn-dev/helm-service -n keptn --wait --version 0.13.0-dev-PR-1234

```
7. Since this is a dev version, check the deployments, if needed, kill pods that are stuck.
8. Authenticate to the Keptn Installation
9. Verify the Keptn version you installed by using `keptn version`


<details>
<summary>Table of Contents</summary>

<!-- toc -->

- [Known Issues](#known-issues)
  * [Docker Image Pull Rate](#docker-image-pull-rate)
  * [Old Images](#old-images)

<!-- tocstop -->

</details>

## Known Issues

### Docker Image Pull Rate
As of November 2nd, 2020, Docker introduces an image pull rate-throttling of 100 pulls per 6 hours per IP address for anonymous users.
This will result in `ImagePullBackOff` errors in your cluster. There is nothing we can do about this as of right now.

### Old Images
Depending on your cluster configuration and the used Keptn Version you might end up with a state that's not the current master branch.

We use `imagePullPolicy: IfNotPresent` and as we install images with the `:latest` tag, you might end up with and old version of Keptn on your Kubernetes cluster (in the Docker cache).

You can verify this by looking at the log output of one of the services, which should print a datetime header on when it was installed, e.g.:
```
$> kubectl -n keptn logs deployments/api-service api-service

##########
branch: release-0.7.2
repository: https://github.com/keptn/keptn
commitlink: https://github.com/keptn/keptn/commit/d469876ecfe95d467872921ae6e9ce877f2ccca6
repolink: https://github.com/keptn/keptn/tree/d469876ecfe95d467872921ae6e9ce877f2ccca6
travisbuild: https://travis-ci.org/keptn/keptn/jobs/735015708
timestamp: 20201012.1651
##########

```

If this is the case for you, you can manually edit the deployments and set image pull policy to always, or alternatively, create a new Kubernetes cluster (which should clear the Docker cache).
