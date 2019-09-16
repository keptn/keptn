# Pitometer Service

The *pitometer-service* is a Keptn core component. It is responsible for validating a test execution based on monitoring data. 

The *pitometer-service* listens to Keptn events of type:
- `sh.keptn.events.tests-finished`

When it receives such an event, the *pitometer-service* queries data from a monitoring solution using the grader module. Based on the data, the service can make the decision whether the test passed or failed. This decision is then sent to the Keptn's eventbroker via a `sh.keptn.events.evaluation-done` event.

## Installation

The *pitometer-service* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *pitometer-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *pitometer-service*, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```