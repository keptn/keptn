# Install master branch

**Attention:** This installs a potentially unstable and unreliable version of Keptn and is not meant to preview features, but for Keptn core/contrib/sandbox developers.

**DO NOT PERFORM THIS ON A PRODUCTION ENVIRONMENT**

1. Create a new cluster (e.g., using k3s)
1. Download latest CLI based on your platform: [Linux](https://storage.cloud.google.com/keptn-cli/latest/keptn-linux.zip) [Mac OS](https://storage.cloud.google.com/keptn-cli/latest/keptn-macOS.zip) [Windows](https://storage.cloud.google.com/keptn-cli/latest/keptn-windows.zip)
1. Unpack the binary and move it to a directory of your choice (e.g., `/usr/local/bin/` on Linux and MacOS); Alternatively, use `./keptn` or `./keptn.exe` where appropriate.
1. Verify that the installation has worked and that the version is correct by running `keptn version` 
1. Install keptn using `keptn install --chart-repo=https://storage.googleapis.com/keptn-installer/latest/keptn-0.1.0.tgz [--use-case=continuous-delivery]`
1. Authenticate to the Keptn Installation
1. Verify the Keptn version you installed by using `keptn version`

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
