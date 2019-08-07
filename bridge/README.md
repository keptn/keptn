# keptn's bridge

Right now, keptn's bridge lets you browse the keptn's log.

In the future it should provide realtime information and metrics about what's going on in your deployment.

## Set up

keptn's bridge is deployed as a part of [keptn](https://keptn.sh)

### Deploy in your k8s cluster

To deploy the current version of keptn's bridge in your keptn kubernetes cluster use the file `bridge.yaml` from this repository and apply it.
```
kubectl apply -f bridge.yaml
```

### Local development

1. Run `kubectl proxy` to create a proxy connection to your keptn cluster.
2. Edit `server/config/index.js` to define the elasticsearch API endpoint.
3. Run `npm install`.
4. Run `npm start` to start the express server that provides the API endpoints.
5. Run `npm run vue-dev` to start the development server.
6. Access the web through the url shown on the console.

### Production deployment

1. Run `npm install`
2. Run `npm run build`
3. Run `npm start`

By default, the process will listen on port 3000.
