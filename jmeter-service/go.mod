module github.com/keptn/keptn/jmeter-service

go 1.18

require (
	github.com/cloudevents/sdk-go/v2 v2.10.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.17.1-0.20220718124953-6b126790c54c
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cloudevents/sdk-go/observability/opentelemetry/v2 v2.0.0-20211001212819-74757a691209 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.0 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace (
	github.com/emicklei/go-restful/v3 => github.com/emicklei/go-restful/v3 v3.8.0
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/net => golang.org/x/net v0.0.0-20220708220712-1185a9018129
	golang.org/x/text => golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 => gopkg.in/yaml.v3 v3.0.1
)
