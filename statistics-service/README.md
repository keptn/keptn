# Statistics Service

This service provides usage statistics about a Keptn installation.

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

You can access the service via the Keptn API under the `statistics` path, e.g.:

```
http://keptn-api-url.com/api/statistics
``` 

Or, if you would like to view the swagger UI of the service, you can use the following URL: 

```
http://keptn-api-url.com/api/swagger-ui/?urls.primaryName=statistics
```

You can then browse the API docs at by opening the Swagger docs in your [browser](http://localhost:8080/swagger-ui/index.html).

To retrieve usage statistics for a certain time frame, you need to provide the [Unix timestamps](https://www.epochconverter.com/) for the start and end of the time frame.
E.g.:

```
http://keptn-api-url.com/api/statistics/v1/statistics?from=1600656105&to=1600696105
```

cURL Example:

```
curl -X GET "http://keptn-api-url.com/api/statistics/v1/statistics?from=1600656105&to=1600696105" -H "accept: application/json" -H "x-token: <keptn-api-token>"
```

*Note*: You can generate timestamps using [epochconverter.com](https://www.epochconverter.com/).

### Configuring the service

By default, the service aggregates data with a granularity of 30 minutes. Whenever this period has passed, the service will create
a new entry in the Keptn-MongoDB within the Keptn cluster. If you would like to change how often statistics are stored, you can set the 
variable `AGGREGATION_INTERVAL_SECONDS` to your desired value.


