# Shipyard Controller

## Installation

The *shipyard-controller* is installed as a part of [keptn](https://keptn.sh)

## Deploy in your Kubernetes cluster

To deploy the current version of the *shipyard-controller* in your Keptn Kubernetes cluster, use the files `deploy/pvc.yaml` and `deploy/service.yaml` from this repository and apply it.

```console
kubectl apply -f deploy/service.yaml
```

## Delete in your Kubernetes cluster

To delete a deployed *shipyard-controller*, use the files `deploy/pvc.yaml` and `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

### Generate  Swagger doc from source

First, the following go modules have to be installed:

```
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

If the `swagger.yaml` should be updated with new endpoints or models, generate the new source by executing:

```console
swag init
```
