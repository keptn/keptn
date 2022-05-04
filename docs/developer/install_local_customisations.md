# Test your own Keptn build artefacts and/or custom helm chart configurations

## Use Cases
* As a developer I want to test my keptn fork build artefacts in a local kubernetes cluster environment.
* As a developer I want to build and deploy my customized keptn fork at my kubernetes cluster environment.
* As a keptn contributor I want to fork keptn, build, deploy and test my enhancements locally in my kubernetes cluster before I open a Pull Request.
* As a mac book pro m1 (arm64) user, I want to build keptn locally and deploy it to my kubernetes cluster.
* As a developer I want to set the repository and image tag at one point for all keptn artefacts.
* ...

## Customise Charts/Docker images and install with own values-local.yaml for local testing
1. Set ENV Vars
   ````
   export AIRGAPPED_REGISTRY_URL=k3d-container-registry.localhost:12345
   export KEPTN_NAMESPACE=keptn
   ````
2. Build keptn artefacts
   - Build all (call from keptn root folder)
      - `VERSION=local-snapshot ; for d in $(find . -name "Dockerfile" | sed -e "s/\.\/\(.*\)\/Dockerfile$/\1/g") ; do echo "building dir $d" ; cd $d ; docker build . -t "$AIRGAPPED_REGISTRY_URL/keptndev/$d:$VERSION" ; cd .. ; done`
   - Just only the one you want to test
      - `docker build . -t $AIRGAPPED_REGISTRY_URL/keptndev/api:local-snapshot`
3. Build Keptn CLI:
   - go to `/keptn/cli` and run `go build  -o keptn main.go`
4. Create a new cluster (e.g., using k3d)
   - ```shell
     k3d registry create container-registry.localhost --port 12345
     k3d cluster create mykeptn -p"8082:80@loadbalancer" --k3s-arg "--no-deploy=traefik@server:*" --agents 1 --registry-use "$AIRGAPPED_REGISTRY_URL"
     kubectl config use-context k3d-mykeptn
     ```
   - Verify that everything has worked using `kubectl get nodes`  
5. Load locally build images in k3d cluster
   - If you want to load all keptn images `for i in $(docker images "keptn/*" --format "{{.Repository}}"); do k3d image load $i\:local-snapshot -c mykeptn ; done`
   - If you want to load only single images e.g. `k3d image load keptn/bridge-service:local-snapshot -c mykeptn`
   - For further information see [k3d image import](https://k3d.io/v5.2.0/usage/commands/k3d_image_import/)
7. Create `installer/manifests/keptn/values-local.yaml` file for your local values to be stored. The file should look like this:
   - global keptn configuration: Set global.keptn.registry/tag if you did a local full build of keptn artefacts
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
10. Install Gitea
    - TODO: like at .github/workflows/integration_tests.yaml#Install Gitea
11. Install Mockserver
    - TODO: like at .github/workflows/integration_tests.yaml#Install Mockserver


## How to test helm charts locally
For local templating of helm charts to take a look about the changes use:
`helm template . -f values-local.yaml --name-template test-control-plane --output-dir ../../temp`
