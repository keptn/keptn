# Gatekeeper Service

The *gatekeeper-service* is a Keptn core component and implements the quality gate in each stage, i.e., depending on the evaluation result it either promotes an artifact to the next stage or not.

The *gatekeeper-service* listens to Keptn events of type:
- `sh.keptn.events.evaluation-done`

The `evaluation-done` contains the result of the evaluation. If the evaluation result is positive, this service sends a `new-artifact` event for the next stage. If the evaluation result is negative and the service is deployed with a blue/green strategy, this service changes the configuration back to the old version and sends a `configuration-changed` event.

## Installation

The *gatekeeper-service* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *gatekeeper-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *gatekeeper-service*, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```