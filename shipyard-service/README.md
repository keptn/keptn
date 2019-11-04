# Shipyard Service

The *shipyard-service* is a Keptn core component. It is responsible for creating a project and processing a shipyard file that defines the stages each deployment has to go through until it is released to end-users. The definition of a shipyard file is provided [here](https://github.com/keptn/keptn/blob/develop/specification/shipyard.md).

The *shipyard-service* listens to Keptn events of type:
- [`sh.keptn.internal.events.project.create`](https://github.com/keptn/keptn/blob/develop/specification/cloudevents.md#create-project)

When receiving such an event, the *shipyard-service* processes the payload in the data block of the event. Thereby, it uses the API of the configuration-service to create the specified entities (i.e., project and stages) and to finally store the payload as shipyad.yaml.

## Installation

The *shipyard-service* is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the *shipyard-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *shipyard-service*, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```