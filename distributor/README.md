# Distributor

A distributor subscribes a Keptn service with the Keptn Control Plane.
Both local and remote subscriptions are supported:

- Local (Keptn service runs in the same local Kubernetes cluster
as the Keptn Control Plane) --
it queries event messages from NATS
and sends the events to services that have a subscription to the event topic.
- Remote (Keptn service runs in a remote "execution plane") --
subscriptions are implemented using the Keptn Subscription API.

Each service has its own distributor
that is configured by the two environment variables:

- `KEPTN_API_ENDPOINT` - Keptn API Endpoint - needed when the distributor runs outside of the Keptn cluster. default = `""`
- `KEPTN_API_TOKEN` - Keptn API Token - needed when the distributor runs outside of the Keptn cluster. default = `""`

Additional environment variables configure other information for the distributor:

- `API_PROXY_PORT` - Port on which the distributor listens for incoming Keptn API requests by its execution plane service. default = `8081`.
- `API_PROXY_PATH` - Path on which the distributor listens for incoming Keptn API requests by its execution plane service. default = `/`.
- `API_PROXY_HTTP_TIMEOUT` - Timeout value (in seconds) for the API Proxy's HTTP Client. default = `30`.
- `HTTP_POLLING_INTERVAL` - Interval (in seconds) in which the distributor checks for new triggered events on the Keptn API. default = `10`
- `EVENT_FORWARDING_PATH` - Path on which the distributor listens for incoming events from its execution plane service. default = `/event`
- `HTTP_SSL_VERIFY` - Determines whether the distributor should check the validity of SSL certificates when sending requests to a Keptn API endpoint via HTTPS. default = `true`
- `PUBSUB_URL` - The URL of the nats cluster the distributor should connect to when the distributor is running within the Keptn cluster. default = `nats://keptn-nats`
- `PUBSUB_TOPIC` - Comma separated list of topics (i.e. event types) the distributor should listen to (see https://github.com/keptn/keptn/blob/master/specification/cloudevents.md for details). When running within the Keptn cluster, it is possible to use NATS [Subject hierarchies](https://nats-io.github.io/docs/developer/concepts/subjects.html#matching-a-single-token). When running outside of the cluster (polling events via HTTP), wildcards can not be used. In this case, each specific topic has to be included in the list.
- `PUBSUB_RECIPIENT` - Hostname of the execution plane service the distributor should forward incoming CloudEvents to. default = `http://127.0.0.1`
- `PUBSUB_RECIPIENT_PORT` - Port of the execution plane service the distributor should forward incoming CloudEvents to. default = `8080`
- `PUBSUB_RECIPIENT_PATH` - Path of the execution plane service the distributor should forward incoming CloudEvents to. default = `/`
- `PUBSUB_GROUP` - Used to join a group for receiving messages from the message broker. Note, that only **one** instance of a distributor in a set of distributors having the same `PUBSUB_GROUP` can receive the event. default = `""`
- `PROJECT_FILTER` - Filter events for a specific project. default = `""` (all); supports a comma-separated list of projects.

- `STAGE_FILTER` - Filter events for a specific stage. default = `""` (all); supports a comma-separated list of stages.
- `SERVICE_FILTER` - Filter events for a specific service. default = `""` (all); supports a comma-separated list of services.
- `DISABLE_REGISTRATION` - Disables automatic registration of the Keptn integration to the control plane. default = `false`
- `REGISTRATION_INTERVAL` - Time duration between trying to re-register to the Keptn control plane. default =`10s`
- `LOCATION` - Location where the distributor is running, e.g. "executionPlane-A". default = `""`
- `DISTRIBUTOR_VERSION` - The software version of the distributor. default = `""`
- `VERSION` - The version of the Keptn integration. default = `""`
- `K8S_DEPLOYMENT_NAME` - Kubernetes deployment name of the Keptn integration. default = `""`
- `K8S_POD_NAME` -  Kubernetes deployment name of the Keptn integration. default = `""`
- `K8S_NAMESPACE` - Kubernetes namespace of the Keptn integration. default = `""`
- `K8S_NODE_NAME` - Kubernetes node name the Keptn integration is running on. default = `""`
- `MAX_HEARTBEAT_RETRIES` - Maximum number of times the distributor tries to do its heartbeat before it gives up. default=`10`
- `HEARTBEAT_INTERVAL` - TIme duration between each heartbeat.  default:`10s`
- `MAX_REGISTRATION_RETRIES` - Maximum number of times the distributor is trying to register itself to the control plane when started. default:`10`
- `REGISTRATION_INTERVAL` - Time duration between trying to re-register to the control plane. default =`10s`
- `OAUTH_CLIENT_ID` - OAuth client ID used when performing Oauth Client Credentials Flow. default = `""`
- `OAUTH_CLIENT_SECRET` - OAuth client ID used when performing Oauth Client Credentials Flow. default = `""`
- `OAUTH_DISCOVERY` - Discovery URL called by the distributor to obtain further information for the OAuth Client Credentials Flow, e.g. the token URL. default = `""`
- `OAUTH_TOKEN_URL` - Url to obtain the access token. If set, this overrides `OAUTH_DISCOVERY` meaning, that no discovery will happen. default = `""`
- `OAUTH_SCOPES` - Comma separated list of tokens to be used during the OAuth Client Credentials Flow. =`""`

All cloud events specified in `PUBSUB_TOPIC` and matching the filters are forwarded to `http://{PUBSUB_RECIPIENT}:{PUBSUB_RECIPIENT_PORT}{PUBSUB_RECIPIENT_PATH}`, e.g.: `http://helm-service:8080`.

### Configuration examples

The above list of environment variables is pretty long, but in most scenarios only a few of them have to be set. The following examples show how to set the environment variables properly, depending on where the distributor and it's accompanying execution plane service should run:

**Configuring the distributor when running within the Keptn cluster**

In this case, usually only the `PUBSUB_TOPIC` has to be defined, e.g.:

```
PUBSUB_TOPIC: "sh.keptn.event.approval.triggered"
```

However, this is not necessary if the distributor is only used as a proxy for the Keptn API, and not needed for subscribing to any topic.

This forwards all incoming events of that topic to `http://127.0.0.1:8080` - which is the URL of the execution plane service running in the same pod as the distributor. If the execution plane service has a different hostname (e.g., when not running in the same pod), a different port, or listens for events on a different path, the env vars `PUBSUB_RECIPIENT`, `PUBSUB_RECIPIENT_PORT` and `PUBSUB_RECIPIENT_PATH` can be set to change this default URL, e.g.:

```
PUBSUB_RECIPIENT: "http://my-service
PUBSUB_RECIPIENT_PORT: "9000"
PUBSUB_RECIPIENT_PATH: "/event-path
```

This causes the distributor to forward all incoming events for its subscribed topic to `http://my-service:9000/event-path`.

The execution plane service can then access the distributor's Keptn API proxy at `http://localhost:8081/`, and can forward events by sending them to `http://localhost:8081/event`.
The Keptn API services are then reachable for the execution plane service via the following URLs:


- Mongodb-datastore:
  - `http://localhost:8081/mongodb-datastore`

- Configuration-service:
  - `http://localhost:8081/configuration-service`

- Shipyard-controller:
  - `http://localhost:8081/controlPlane`

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

## Filtering for a set of stages, projects, or services

The STAGE_FILTER, PROJECT_FILTER, and SERVICE_FILTER environment variables
control the Keptn service's subscription to events with Keptn's Control Plane.
The values of these environment variables are set by fields in the values.yaml file for the service;
by default, all stages, projects, and services are subscribed.
Provide a comma-separated list of stages, projects, or services to the appropriate variable
to filter the set.
Define the value of these variables in the appropriate field of the *value.yaml* file for the service;
that populates the value of the environment variables that the Distributor uses.

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
          image: keptndev/distributor:latest
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
              value: 'nats://keptn-nats'
            - name: PUBSUB_TOPIC
              value: 'sh.keptn.internal.event.some-event'
            - name: PUBSUB_RECIPIENT
              value: 'your-service'
```
