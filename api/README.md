# Keptn API Component

The api component is a Keptn core component and allows the communication with Keptn. Therefore, it provides a defined interface as shown in the `./swagger.json`. Besides, it maintains a websocket server to forward Keptn messages to the Keptn CLI, used by the end-user.

## Installation

The api component is installed as a part of [Keptn](https://keptn.sh).

## Deploy in your Kubernetes cluster

To deploy the current version of the api component in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed api component, use the file `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```