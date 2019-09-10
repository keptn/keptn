# Eventbroker

The *eventbroker* is a Keptn core component that is responsible for receiving all events, transferring non-Keptn events into valid Keptn Cloud Events, and sending those into NATS. 

## Installation

The *eventbroker* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *eventbroker* in your Keptn Kubernetes cluster, use the file `deploy/eventbroker.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/eventbroker.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *eventbroker*, use the file `deploy/eventbroker.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/eventbroker.yaml
```