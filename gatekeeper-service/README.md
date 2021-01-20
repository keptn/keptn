# Gatekeeper Service

The *gatekeeper-service* is a Keptn core component and implements the quality gate in each stage, i.e., depending on the evaluation result it either promotes an artifact to the next stage or not.

The *gatekeeper-service* listens to Keptn events of type:
- `sh.keptn.event.approval.triggered`

The `approval.triggered` contains the approval strategy, as well as the result of the previous service execution (e.g., an evaluation). If the result is positive (e.g., 
 `result = "pass" || result = "warning"`), and the approval strategyy is set to `automatic`, the service will automatically send out a `approval.finished` event to continue the task sequence for the associated Keptn context.
 If the strategy is set to `manual`, the service will not respond with any further events. In this case, the user is responsible for sending an `approval.finished` events to continue the task sequence for the associated Keptn context.
 If the rsult of the previous service execution is set to `fail`, the gatekeeper service will automatically send an `approval.finished` event with the result set to `fail`.

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