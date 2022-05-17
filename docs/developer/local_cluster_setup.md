# Set up a local cluster for testing

1. Set ENV Vars
   ````
   export AIRGAPPED_REGISTRY_URL=k3d-container-registry.127.0.0.1.nip.io:12345
   export KEPTN_NAMESPACE=keptn
   ````
2. Build keptn artefacts
    - Build all (call from keptn root folder)
      ```shell
      VERSION=local-snapshot ; for d in $(find . -name "Dockerfile" | sed -e "s/\.\/\(.*\)\/Dockerfile$/\1/g") ; do echo "building dir $d" ; cd $d ; docker build . -t "$AIRGAPPED_REGISTRY_URL/keptndev/$d:$VERSION" ; cd .. ; done
      ```
    - Build just one artefact/keptn-service (go to the service folder where the Dockerfile is located)
      ```shell
      docker build . -t $AIRGAPPED_REGISTRY_URL/keptndev/api:local-snapshot
      ```
3. Build Keptn CLI:
    - go to `/keptn/cli` and run `go build -o keptn main.go`
4. Create new Registry and push images
   ```shell
   k3d registry create container-registry.127.0.0.1.nip.io --port 12345`
   for d in $(docker images "keptn/*:local-snapshot" --format "{{.Repository}}:{{.Tag}}"); do docker push $AIRGAPPED_REGISTRY_URL/$d ; done
   ```
5. Create a new cluster (e.g., using k3d)
    ```shell
     k3d cluster create mykeptn -p "8082:80@loadbalancer" --k3s-arg "--no-deploy=traefik@server:*" --agents 1 --registry-use "$AIRGAPPED_REGISTRY_URL"
     kubectl config use-context k3d-mykeptn
    ```
    - Verify that everything has worked using `kubectl get nodes`
6. (optional) Load locally build images in k3d cluster (if no local registry is used)
    - If you want to load all keptn images `for i in $(docker images "keptn/*" --format "{{.Repository}}"); do k3d image load $i\:local-snapshot -c mykeptn ; done`
    - If you want to load only single images e.g. `k3d image load keptn/bridge-service:local-snapshot -c mykeptn`
    - For further information see [k3d image import](https://k3d.io/v5.2.0/usage/commands/k3d_image_import/)
7. Install Gitea
    - `.github/workflows/integration_tests.yaml#Install Gitea`
8. Install Mockserver
    - `.github/workflows/integration_tests.yaml#Install Mockserver`