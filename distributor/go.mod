module github.com/keptn/keptn/distributor

go 1.20

require (
	github.com/cloudevents/sdk-go/protocol/nats/v2 v2.12.0
	github.com/cloudevents/sdk-go/v2 v2.13.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.20.1
	github.com/nats-io/nats-server/v2 v2.9.15
	github.com/nats-io/nats.go v1.24.0
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.8.2
	golang.org/x/oauth2 v0.7.0
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cloudevents/sdk-go/observability/opentelemetry/v2 v2.13.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.0 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.3.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.4 // indirect
	go.opentelemetry.io/otel v1.11.1 // indirect
	go.opentelemetry.io/otel/metric v0.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	golang.org/x/crypto => golang.org/x/crypto v0.8.0
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)
