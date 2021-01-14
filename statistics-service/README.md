# Statistics Service

This service provides usage statistics about a Keptn installation.

## Compatibilty Matrix

| Keptn Version    | [Statistics Service](https://hub.docker.com/r/keptnsandbox/statistics-service/tags?page=1&ordering=last_updated) | Kubernetes Versions                      |
|:----------------:|:----------------------------------------:|:----------------------------------------:|
|       0.7.1      | keptnsandbox/statistics-service:0.1.0    | 1.14 - 1.19                              |
|       0.7.2      | keptnsandbox/statistics-service:0.1.1    | 1.14 - 1.19                              |
|       0.7.3      | keptnsandbox/statistics-service:0.2.0    | 1.14 - 1.19                              |


## Deploy in your Kubernetes cluster

Please note that the installation of the **statistics-service** differs slightly, depending on your installed Keptn version. Depending on your installed Keptn version, please follow the instructions below. 

### For Keptn versions < 0.8.0

To deploy the current version of the *statistics-service* in your Keptn Kubernetes cluster, use the file `deploy/service.yaml` from this repository and apply it.

```console
kubectl apply -f deploy/service.yaml -n keptn
```

### For Keptn versions >= 0.8.0

To deploy the current version of the *statistics-service* in your Keptn Kubernetes cluster, use the file `deploy/service_keptn_080.yaml` from this repository and apply it.

```console
kubectl apply -f deploy/service_keptn_080.yaml -n keptn
```

## Delete in your Kubernetes cluster

To delete a deployed *statistics-service*, use `deploy/service.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml -n keptn
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

## How to use the service

Once the service is deployed in your cluster, you can access it using `port-forward`:

```
kubectl port-forward -n keptn svc/statistics-service 8080
``` 

You can then browse the API docs at by opening the Swagger docs in your [browser](http://localhost:8080/swagger-ui/index.html).

To retrieve usage statistics for a certain time frame, you need to provide the [Unix timestamps](https://www.epochconverter.com/) for the start and end of the time frame.
E.g.:

```
http://localhost:8080/v1/statistics?from=1600656105&to=1600696105
```

cURL Example:

```
curl -X GET "http://localhost:8080/v1/statistics?from=1600656105&to=1600696105" -H "accept: application/json"
```

*Note*: You can generate timestamps using [epochconverter.com](https://www.epochconverter.com/).

### Configuring the service

By default, the service aggregates data with a granularity of 30 minutes. Whenever this period has passed, the service will create
a new entry in the Keptn-MongoDB within the Keptn cluster. If you would like to change how often statistics are stored, you can set the 
variable `AGGREGATION_INTERVAL_SECONDS` to your desired value.


