module github.com/keptn/keptn/lighthouse-service

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/go-test/deep v1.0.7
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.9.0
	github.com/nats-io/nats-server/v2 v2.4.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
)

replace github.com/keptn/go-utils => github.com/dynatrace-oss-contrib/go-utils v0.9.1-0.20210906144951-cd85cfcb4eb2
