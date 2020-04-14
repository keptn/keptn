# Keptn's Bridge

Keptn's bridge allows to browse the Keptn events.

## Installation

The Keptn's bridge is installed as a part of [Keptn](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the bridge in your Keptn Kubernetes cluster, use the file `deploy/bridge.yaml` from this repository and apply it:

```console
kubectl apply -f deploy/bridge.yaml
```

### Setting up Basic Authentication

Keptn's Bridge comes with a very simple basic authentication feature, which can be controlled by setting the following two environment variables:

* `BASIC_AUTH_USERNAME` - username
* `BASIC_AUTH_PASSWORD` - password

To enable it within your Kubernetes cluster, we recommend first creating a secret which holds the two variables, and then apply this secret within the Kubernetes deployment for Keptn's Bridge.

1. Create the secret using

    ```console
    kubectl -n keptn create secret generic bridge-credentials --from-literal="BASIC_AUTH_USERNAME=<USERNAME>" --from-literal="BASIC_AUTH_PASSWORD=<PASSWORD>"
    ```
    *Note: Replace `<USERNAME>` and `<PASSWORD>` with the desired credentials.*

2. In case you are using Keptn 0.6.1 or older, edit the deployment using

    ```console
    kubectl -n keptn edit deployment bridge
    ```
   
    and add the secret to the `bridge` container, such that the container spec looks like this:

    ```yaml
        ...
        spec:
          containers:
          - name: bridge
            image: keptn/bridge2:0.6.1
            imagePullPolicy: Always
            # EDIT STARTS HERE
            envFrom:
              - secretRef:
                  name: bridge-credentials
                  optional: true
            # EDIT ENDS HERE
            ports:
            - containerPort: 3000
            ...
    ```
   
**Note 1**: To disable authentication, just delete the secret using ``kubectl -n keptn delete secret bridge-credentials``.

**Note 2**: If you delete or edit the secret, you need to restart the respective pod by executing

```console
kubectl -n keptn scale deployment bridge --replicas=0
kubectl -n keptn scale deployment bridge --replicas=1
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
