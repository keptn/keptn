# Test your own Keptn build artefacts and/or custom helm chart configurations

## Use Cases
* As a developer I want to test my keptn fork build artefacts in a local kubernetes cluster environment.
* As a developer I want to build and deploy my customized keptn fork at my kubernetes cluster environment.
* As a keptn contributor I want to fork keptn, build, deploy and test my enhancements locally in my kubernetes cluster before I open a Pull Request.
* As a mac book pro m1 (arm64) user, I want to build keptn locally and deploy it to my kubernetes cluster.
* As a developer I want to set the repository and image tag at one point for all keptn artefacts.
* ...

## Customise Charts/Docker images and install with own values-local.yaml for local testing

1. Prerequisites
   - Install Cluster
2. Create `installer/manifests/keptn/values-local.yaml` file for your local values to be stored. The file should look like this:
   - global keptn configuration: Set `global.keptn.registry` and `global.keptn.tag` if you did a local full build of keptn artefacts
   - service configuration: If you only want to install one single artefact from your local build (e.g. apiService)

```yaml
#  # set global keptn registry and tag for completely overriding the keptn default config
global:
  keptn:
    registry: "testregistry/keptn"      # keptn registry/image name
    tag: "0.0.1"                        # keptn version/tag

#control-plane:
#  # service config: only set individual values if global.keptn.registry/tag are not satisfiying
#  apiService:
#    image:
#      registry: "my-local-api-reg"                             # Container Registry
#      tag: "my-local-api-tag"                                  # Container Tag
#
#  # only change if version at ./charts/control-plane/values.yaml --> apiGatewayNginx.registry/tag is not satisfying
#  apiGatewayNginx:
#    registry: this.is.a.test           # nginx registry/image name
#    tag: 10.0.0                        # ngnix version/tag
#
```
7. Create Namespace `kubectl create ns keptn`
8. Download Helm Dependencies `helm dependency update`
   - `installer/manifests/keptn/control-plane`
   - `helm-service/charts`
   - `jmeter-service/charts`
9. Install keptn in local cluster 
   - `installer/manifests/keptn` --> `helm upgrade --install -f values-local.yaml keptn . -n keptn`
   
## How to test helm charts locally
For local templating of helm charts to take a look about the changes use:
`helm template . -f values-local.yaml --name-template test-control-plane --output-dir ../../temp`
