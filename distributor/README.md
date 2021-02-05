# Distributor

A distributor queries event messages from NATS and sends the events to services that have a subscription to the event topic.
Thus, each service has its own distributor that is configured by the two environment variables:

- `KEPTN_API_ENDPOINT` - Keptn API Endpoint - needed when the distributor runs outside of the Keptn cluster. default = `"""`
- `KEPTN_API_TOKEN` - Keptn API Token - needed when the distributor runs outside of the Keptn cluster. default = `"""`
- `API_PROXY_PORT` - Port on which the distributor will listen for incoming Keptn API requests by its execution plane service. default = `8081`.
- `API_PROXY_PATH` - Path on which the distributor will listen for incoming Keptn API requests by its execution plane service. default = `/`.
- `HTTP_POLLING_INTERVAL` - Interval (in seconds) in which the distributor will check for new triggered events on the Keptn API. default = `10`
- `EVENT_FORWARDING_PATH` - Path on which the distributor will listen for incoming events from its execution plane service. default = `/event`
- `HTTP_SSL_VERIFY` - Determines whether the distributor should check the validity of SSL certificates when sending requests to a Keptn API endpoint via HTTPS. default = `true`
- `PUBSUB_URL` - The URL of the nats cluster the distributor should connect to when the distributor is running within the Keptn cluster. default = `nats://keptn-nats-cluster`
- `PUBSUB_TOPIC` - Comma separated list of topics (i.e. event types) the distributor should listen to (see https://github.com/keptn/keptn/blob/master/specification/cloudevents.md for details). When running within the Keptn cluster, it is possible to use NATS [Subject hierarchies](https://nats-io.github.io/docs/developer/concepts/subjects.html#matching-a-single-token). When running outside of the cluster (polling events via HTTP), wildcards can not be used. In this case, each specific topic has to be included in the list.
- `PUBSUB_RECIPIENT` - Hostname of the execution plane service the distributor should forward incoming CloudEvents to. default = `http://127.0.0.1`
- `PUBSUB_RECIPIENT_PORT` - Port of the execution plane service the distributor should forward incoming CloudEvents to. default = `8080`
- `PUBSUB_RECIPIENT_PATH` - Path of the execution plane service the distributor should forward incoming CloudEvents to. default = `/`

All cloud events specified in `PUBSUB_TOPIC` are forwarded to `http://{PUBSUB_RECIPIENT}:{PUBSUB_RECIPIENT_PORT}{PUBSUB_RECIPIENT_PATH}`, e.g.: `http://helm-service:8080`.

### Configuration examples

The above list of environment variables is pretty long, but in most scenarios only a few of them have to be set. The following examples show how to set the environment variables properly, depending on where the distributor and it's accompanying execution plane service should run:

**Configuring the distributor when running within the Keptn cluster**

In this case, usually only the `PUBSUB_TOPIC` has to be defined, e.g.:

```
PUBSUB_TOPIC: "sh.keptn.event.approval.triggered"
```

However, this is not necessary if the distributor is only used as a proxy for the Keptn API, and not needed for subscribing to any topic.

This will forward all incoming events of that topic to `http://127.0.0.1:8080` - which is the URL of the execution plane service running in the same pod as the distributor. If the execution plane service has a different hostname (e.g., when not running in the same pod), a different port, or listens for events on a different path, the env vars `PUBSUB_RECIPIENT`, `PUBSUB_RECIPIENT_PORT` and `PUBSUB_RECIPIENT_PATH` can be set to change this default URL, e.g.:

```
PUBSUB_RECIPIENT: "http://my-service
PUBSUB_RECIPIENT_PORT: "9000"
PUBSUB_RECIPIENT_PATH: "/event-path
```

This will cause the distributor to forward all incoming events for its subscribed topic to `http://my-service:9000/event-path`.

The execution plane service will then be able to access the distributor's Keptn API proxy at `http://localhost:8081/`, and can forward events by sending them to `http://localhost:8081/event`.
The Keptn API services will then be reachable for the execution plane service via the following URLs:


- Mongodb-datastore:
    - `http://localhost:8081/mongodb-datastore`
    - `http://localhost:8081/datastore`
    - `http://localhost:8081/event-store`

- Configuration-service:
    - `http://localhost:8081/configuration-service`
    - `http://localhost:8081/configuration`
    - `http://localhost:8081/config`

- Shipyard-controller:
    - `http://localhost:8081/shipyard-controller`
    - `http://localhost:8081/shipyard`

If the distributor should listen on a port other than `8081` (e.g. when that port is needed by the execution plane service), a different port can be set using the `API_PROXY_PORT` environment variable

**Configuring the distributor when running outside of the Keptn cluster**

In this case, the Keptn API URL and the API token, as well as a topic have to be defined:

```
KEPTN_API_ENDPOINT: "https://my-keptn-api:8080/api"
KEPTN_API_TOKEN: "my-keptn-api-token"
PUBSUB_TOPIC: "sh.keptn.event.approval.triggered" # can also be left empty in this case, if the distributor is only used as a proxy to interact with the Keptn API
```

If the endpoint specified by `KEPTN_API_ENDPOINT` does not provide a valid SSL certificate, the distributor will, per default, deny any requests to that endpoint. This behavior can be changed by setting the variable `HTTP_SSL_VERIFY` to `false`.

The remaining parameters, such as `PUBSUB_RECIPIENT`, `PUBSUB_RECIPIENT_PORT` and `PUBSUB_RECIPIENT_PATH`, as well as the `API_PROXY_PORT` can be configured as described above.

## Installation

Distributors are installed automatically as a part of [Keptn](https://keptn.sh). See
[core-distributors.yaml](/installer/manifests/keptn/core-distributors.yaml) for details.

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

## Create your own distributor

You can create your own distributor by writing a dedicated distributor deployment yaml:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: some-service-monitoring-configure-distributor
  namespace: keptn
spec:
  selector:
    matchLabels:
      run: distributor
  replicas: 1
  template:
    metadata:
      labels:
        run: distributor
    spec:
      containers:
        - name: distributor
          image: keptn/distributor:latest
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "32Mi"
              cpu: "50m"
            limits:
              memory: "128Mi"
              cpu: "500m"
          env:
            - name: PUBSUB_URL
              value: 'nats://keptn-nats-cluster'
            - name: PUBSUB_TOPIC
              value: 'sh.keptn.internal.event.some-event'
            - name: PUBSUB_RECIPIENT
              value: 'your-service'
```