# Test your own Keptn build artefacts and/or custom helm chart configurations

## Use Cases
* As a developer I want to test my keptn fork build artefacts in a local kubernetes cluster environment.
* As a developer I want to build and deploy my customized keptn fork at my kubernetes cluster environment.
* As a keptn contributor I want to fork keptn, build, deploy and test my enhancements locally in my kubernetes cluster before I open a Pull Request.
* As a Macbook Pro M1 (arm64) user, I want to build keptn locally and deploy it to my kubernetes cluster.
* As a developer I want to set the repository and image tag globally for all keptn artefacts.
* ...

## Customise Charts/Docker images and install with own values-local.yaml for local testing

1. Prerequisites
   - [Install Cluster](local_cluster_setup.md)
2. Create `installer/manifests/keptn/values-local.yaml` file for your local values to be stored. The file should look like this:
   - global keptn configuration: Set `global.keptn.registry` and `global.keptn.tag` if you did a local full build of keptn artefacts
   - service configuration: If you only want to install one single artefact from your local build (e.g. apiService)
```yaml
# set global keptn registry and tag for completely overriding the keptn default config
global:
  keptn:
    registry: "k3d-container-registry.127.0.0.1.nip.io:12345/keptn"      # keptn registry/image name
    tag: "local-snapshot"                                                # keptn version/tag
```
3. Test helm charts locally 
   - For local templating of helm charts to take a look about the changes use:
   ```shell
   helm template . -f values-local.yaml --name-template test-control-plane --output-dir ../../temp+
   ```
4. Create Namespace
   ```shell
   kubectl create ns keptn
   ```
5. Download Helm Dependencies `helm dependency update`
   - `installer/manifests/keptn`
   - `helm-service/chart`
   - `jmeter-service/chart`
6. Install keptn in local cluster
   Go to `installer/manifests/keptn`
   ```shell
   helm upgrade --install -f values-local.yaml keptn . -n keptn
   ```
7. Open a new terminal and type:
   ```shell
   kubectl -n keptn port-forward service/api-gateway-nginx 8080:80
   ```
8. Authenticate Keptn:
   ```shell
   keptn auth --endpoint=http://127.0.0.1:8080/api --api-token=$(kubectl get secret keptn-api-token -n keptn -ojsonpath={.data.keptn-api-token} | base64 --decode)
   ```
9. Verify the installation has worked
   ```shell
   keptn status
   ```
10. Verify which images have been deployed
    ```shell
    kubectl -n keptn get deployments
    ```
11. Run tests (e.g., UniformRegistration):
   ```shell
   cd test/go-tests && KEPTN_ENDPOINT="http://127.0.0.1:8080/api" go test ./...
   ```
   **Note**: If you want to run a single test, (e.g. BackupTestore_Test), please add `_test` suffix to the test file name, so it becomes executable. Otherwise, you cano run only `testsuite_*_test.go` files. For running a single test use:
   ```shell
   cd test/go-tests && KEPTN_ENDPOINT="http://127.0.0.1:8080/api" go test ./... -v -run <NameOfTheTest>
   ```

   

