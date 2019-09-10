# Distributor

A distributor queries event messages from NATS and sends the events to services that have a subscription to the event topic. Thus, each service has its own distributor that is configured by the two environment variables:
- `PUBSUB_TOPIC:` e.g., `sh.keptn.events.new-artifact`
- `PUBSUB_RECIPIENT:` e.g., `helm-service`

## Installation

Distributors are installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of a *distributor* in your Keptn Kubernetes cluster, use the file `deploy/distributor.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *distributor*, use the file `deploy/distributor.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```