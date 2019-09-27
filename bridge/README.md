# Keptn's Bridge

Right now, Keptn's bridge lets you browse the Keptn's log.

In the future it should provide realtime information and metrics about what is going on in your deployment.

## Installation

The Keptn's bridge is installed as a part of [Keptn](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the bridge in your Keptn Kubernetes cluster, use the file `deploy/bridge.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/bridge.yaml
```

### Delete in your Kubernetes cluster

To delete a deployed bridge, use the file `deploy/bridge.yaml` from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/bridge.yaml
```

## Local development

1. Run `kubectl proxy` to create a proxy connection to your Keptn cluster.
2. Edit `server/config/index.js` to define the keptn datastore API endpoint.
3. Run `npm install`.
4. Run `npm start` to start the express server that provides the API endpoints.
5. Run `npm run vue-dev` to start the development server.
6. Access the web through the url shown on the console.

## Production deployment

1. Run `npm install`
2. Run `npm run build`
3. Run `npm start`

By default, the process will listen on port 3000.
