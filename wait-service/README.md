# Wait Service
The *wait-service* is a Keptn core component that is responsible for delaying a workflow for a specified duration. The duration is set by the environment variable `WAIT_DURATION`. The value of this variable must follow the pattern: **[duration][unit]**, e.g., 1h, 5m, 30s.

The *wait-service* listens to Keptn events of type:
- `sh.keptn.events.deployment-finished`

When receiving such an event, the *wait-service* sleeps for the specified duration. After sleeping for this time, a `sh.keptn.events.tests-finished` event will be sent to Keptn's eventbroker.

## Installation

The *wait-service* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *wait-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *wait-service*, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```