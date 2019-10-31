# Distributor

A distributor queries event messages from NATS and sends the events to services that have a subscription to the event topic. 
Thus, each service has its own distributor that is configured by the two environment variables:

- `PUBSUB_TOPIC` -  e.g., `sh.keptn.events.configuration.change` (see https://github.com/keptn/keptn/blob/master/specification/cloudevents.md for details)
  - [Subject hierarchies](https://nats-io.github.io/docs/developer/concepts/subjects.html#matching-a-single-token) are also possible, e.g.: `sh.keptn.internal.event.project.>`
- `PUBSUB_RECIPIENT` -  e.g., `helm-service`
- `PUBSUB_RECIPIENT_PORT` (optional; default: `8080`)
- `PUBSUB_RECIPIENT_PATH` (optional; default: `"""` empty string)

In addition, the following environment variable defines the URL of the NATS server:
- `PUBSUB_URL` - e.g., `nats://keptn-nats-cluster`

All cloud events specified in `PUBSUB_TOPIC` are forwarded to `http://{PUBSUB_RECIPIENT}:{PUBSUB_RECIPIENT_PORT}{PUBSUB_RECIPIENT_PATH}`, e.g.: `http://helm-service:8080`.

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