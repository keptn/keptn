# mongodb Datastore

The *mongodb-datastore* provides means to store and read data from a mongodb deployed in your Keptn cluster. In its current implementation, the service provides two endpoints:
- /events
- /logs

The endpoints are implemented in a REST-api manner. More information can be found by taking a look at the [generated swagger docs](#view-swagger-docs).

## Local development

### Generate source from Swagger

If the `swagger.json` is updated with new endpoints or models, generate the new source by executing:
```console
swagger generate server -A mongodb-datastore -f ./swagger.yaml
```

### View swagger docs

The swagger docs are exposed on http://localhost:8080/swagger-ui

### launch.json for VS Code

If you are using VS Code for your development, you can use the following launch configuration for your local deployment to start the process on the local port 8080.
```
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/mongo-db-datastore-server/main.go",
            "env": {
                "MONGO_DB_CONNECTION_STRING":"mongodb://user:password@localhost:27017/keptn",
                "MONGODB_DATABASE":"keptn"
            },
            "args": ["--port=8080"]
        }
    ]
}
```
