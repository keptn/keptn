# JMeter Service

The *jmeter-service* is a Keptn core component and used for triggering JMeter tests.

The *jmeter-service* listens to Keptn events of type:
- `sh.keptn.events.deployment-finished`

In case the tests succeeed, this service sends a `sh.keptn.events.test-finished` event. In case the tests do not succeed (e.g., the error rate is too high), this service sends an `sh.keptn.events.evaluation-done` event with the data `evaluationpassed=false`.

## Installation

The *jmeter-service* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *jmeter-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *jmeter-service*, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```
