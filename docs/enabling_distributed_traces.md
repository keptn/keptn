# Enabling distributed traces for Keptn

The Keptn services are instrumented using OpenTelemetry. 

> If this is your first time hearing about OpenTelemetry, take a look at the official documentation: [What is OpenTelemetry](https://opentelemetry.io/docs/concepts/what-is-opentelemetry/).

The effort to add instrumentation to all services is a working in progress. As of now, the following services are *partially* instrumented:

- Distributor

- Shipyard controller

- Lighthouse service

In a nutshell, the services collect the spans and via the OTLP exporter they are sent to a [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/). The job of the collector (among many others) is to export the spans it receives to any back-end configured, for example, [Jaeger](https://www.jaegertracing.io/).

If you want to collect and see the traces produced by the Keptn services, follow the next steps where we will deploy a OpenTelemetry collector + a Jaeger instance in the cluster.

## Deploying and configuring a OpenTelemetry collector

For simplicity, we are going to deploy the collector on the same cluster where Keptn is running.

> The collector is extensible and offers a lot of options for configuration and modes of deployment. The intention of this guide is to be a quick starting point for Keptn users. You can learn more and adapt your collector deployment by looking at the [official documentation](https://opentelemetry.io/docs/collector/getting-started/#deployment).


### 1. Namespace

Let's create a namespace to better organize our observability services:

```shell
kubectl create namespace observability
```

### 2. Deploy the collector

Next, run the following `kubectl` command to create the deployment of our OpenTelemetry Collector: 

```shell
kubectl apply -n observability -f <<EOF
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: otel-collector-conf
  namespace: observability
  labels:
    app: opentelemetry
    component: otel-collector-conf
data:
  otel-collector-config: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:
    processors:
    exporters:
      logging:
      jaeger:
        endpoint: "simplest-collector-headless.observability:14250"
        insecure: "true"
    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: []
          exporters: [logging, jaeger]
---
apiVersion: v1
kind: Service
metadata:
  name: otel-collector
  namespace: observability
  labels:
    app: opentelemetry
    component: otel-collector
spec:
  ports:
  - name: otlpgrpc # Default endpoint for OpenTelemetry gRPC receiver.
    port: 4317
    protocol: TCP
    targetPort: 4317
  - name: otlphttp # Default endpoint for OpenTelemetry HTTP receiver.
    port: 4318
    protocol: TCP
    targetPort: 4318
  selector:
    component: otel-collector
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: otel-collector
  namespace: observability
  labels:
    app: opentelemetry
    component: otel-collector
spec:
  replicas: 1
  selector:
    matchLabels:
      app: opentelemetry
      component: otel-collector
  template:
    metadata:
      labels:
        app: opentelemetry
        component: otel-collector
    spec:
      containers:
      - name: otel-collector
        args:
        - --config=/conf/otel-collector-config.yaml
        image: otel/opentelemetry-collector:0.33.0
        ports:
          - containerPort: 4317 
          - containerPort: 4318
        volumeMounts:
        - name: otel-collector-config-vol
          mountPath: /conf
      volumes:
      - configMap:
          name: otel-collector-conf
          items:
            - key: otel-collector-config
              path: otel-collector-config.yaml
        name: otel-collector-config-vol
EOF
```

The command above created a `ConfigMap`, `Service` and `Deployment`. The `ConfigMap` is the most interesting part. There we configured the `receivers` and `exporters`. In this case, we are going to receive in [OTLP](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/otlp.md) and export it to the console and to Jaeger (we will deploy Jaeger next).

> Note: We configured Jaeger endpoint with `simplest-collector-headless.observability`, this will be the name we will give next when deploying it. If you already have one running and want to use it instead, make sure to change it in the ConfigMap before.

To make sure it worked, you can run `kubectl get pods -n observability` and check you have a `otel-collector-<hash>` running.

### 3. Deploy Jaeger

We are going to deploy the "all-in-one" Jaeger image. We'll do that by installing the [Jaeger Operator](https://www.jaegertracing.io/docs/1.26/operator/).

> All-in-one is an executable designed for quick local testing, launches the Jaeger UI, collector, query, and agent, with an in memory storage component. Check the documentation as well for production grade deployments.

1. Install the operator:

```shell
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml
```

2. Deploy the all-in-one instance

```shell
kubectl apply -n observability -f - <<EOF
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: simplest
EOF
```

3. Create an ingress to expose the Jaeger UI (using Istio)

```shell
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: istio
  name: jaeger-ingress
  namespace: observability
spec:
  rules:
  - host: jaeger.127.0.0.1.nip.io
    http:
      paths:
      - pathType: Prefix
        path: /
        backend:
          service:
            name: simplest-query
            port:
              number: 16686
EOF
```
Make sure Jaeger is running with `kubectl get jaegers`.

### 4. Configure Keptn services

Now that we have the OpenTelemetry collector running in our cluster, we have to tell the Keptn services where to find it. This is done by the environment variable `OTEL_COLLECTOR_ENDPOINT`. 

The way we can set this environment variable is by using `Helm upgrade`. This will restart the services and inject the environment variable into the pod's containers. 

If you followed along, the collector endpoint should be `otel-collector.observability:4317` (<servicename.namespace:port>)

```shell
helm upgrade keptn keptn -n keptn --version=0.9.0 --repo=https://storage.googleapis.com/keptn-installer --set global.observability.otelCollectorEndpoint=otel-collector.observability:4317
```

Wait for the pods to restart. To make sure the environment variable was correctly set, you can get a shell inside a container:

```shell
$ k exec -it shipyard-controller-5b7d8765f8-qz22b -c distributor -n keptn -- sh

$ printenv | grep OTEL
OTEL_COLLECTOR_ENDPOINT=otel-collector.observability:4317

```

### Checking the traces

Navigate to the exposed Jaeger UI (in this case `http://jaeger.127.0.0.1.nip.io:8082/`) and start checking your traces!

> The istio/ingress is based on on the getting started guide script from https://raw.githubusercontent.com/keptn/keptn.github.io/master/content/docs/quickstart/exposeKeptnConfigureIstio.sh. It uses port 8082 by default. If you have configured something different make sure to use that instead to reach the Jaeger UI.

