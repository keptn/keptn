# Shipyard Controller

## Installation

The *shipyard-controller* is installed as a part of [keptn](https://keptn.sh)

## Deploy in your Kubernetes cluster

This service should be automatically deployed when executing `keptn install` or installing Keptn from a Helm chart. If
you still want to deploy it manually in your Keptn Kubernetes, you can either

* use the manifest `deploy/service.yaml` from this repository and apply it
  ```console
  kubectl apply -f deploy/service.yaml
  ```
* build locally (from source) and deploy it using
  ```console
  export SKAFFOLD_DEFAULT_REPO=containerregistry # e.g., docker.io/username
  skaffold run --tail  
  ```
  **Note**: When stopping skaffold (e.g., using CTRL-C) the deployment will automatically be cleaned up

## Delete in your Kubernetes cluster

To delete a deployed *shipyard-controller*, use the manifest `deploy/service.yaml` from this repository and delete the
Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## Generate  Swagger doc from source

1. Download and install Swag for Go by calling `go install github.com/swaggo/swag/cmd/swag` in fresh terminal.
2. `cd` to the Shipyard Controller's root folder and run `swag init --parseDependency`

## Technical Details

### Sequence handling
The shipyard controller orchestrates sequences by listening for certain events, such as `.triggered` events that should trigger a sequence,
or events that indicate the start/completion of a task execution by one of Keptns execution plane services. Further, it is responsible for the 
following tasks:

- Dispatching sequences while ensuring that there are no multiple sequences running in the same stage for the same service at any given point in time
- Sending out `.triggered` events that indicate that a task within a sequence should be executed
- Cancelling sequences when a timeout for a task has been detected

The following flow charts illustrate the workflow of how these responsibilities are handled by the shipyard controller:

**Reception of a .triggered event:**

![handleTriggeredEvent](assets/handleTriggeredEvent.png?raw=true "handleTriggeredEvent")

**Dispatching Sequences:**

![sequenceDispatcher](assets/sequenceDispatcher.png?raw=true "sequenceDispatcher")

**Watching for timed out tasks:**

![sequenceWatcher](assets/sequenceWatcher.png?raw=true "sequenceWatcher")

**Keep track of .started events:**

![handleStartedEvent](assets/handleStartedEvent.png?raw=true "handleStartedEvent")

**Keep track of .finished events:**

![handleFinishedEvent](assets/handleFinishedEvent.png?raw=true "handleFinishedEvent")
