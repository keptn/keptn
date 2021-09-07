module github.com/keptn/keptn/distributor

go 1.16

require (
	github.com/cloudevents/sdk-go/protocol/nats/v2 v2.5.0
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.9.0
	github.com/nats-io/nats-server/v2 v2.4.0
	github.com/nats-io/nats.go v1.12.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
)

replace github.com/keptn/go-utils => github.com/dynatrace-oss-contrib/go-utils v0.9.1-0.20210907155952-fc60962cc8d1
