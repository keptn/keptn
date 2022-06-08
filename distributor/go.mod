module github.com/keptn/keptn/distributor

go 1.18

require (
	github.com/cloudevents/sdk-go/protocol/nats/v2 v2.10.0
	github.com/cloudevents/sdk-go/v2 v2.10.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.16.1
	github.com/keptn/keptn/cp-common v0.0.0-20220523070440-9a1e3bffda90
	github.com/nats-io/nats-server/v2 v2.8.4
	github.com/nats-io/nats.go v1.15.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	golang.org/x/oauth2 v0.0.0-20220524215830-622c5d57e401
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/cloudevents/sdk-go/observability/opentelemetry/v2 v2.0.0-20211001212819-74757a691209 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/klauspost/compress v1.14.4 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/nats-io/jwt/v2 v2.2.1-0.20220330180145-442af02fd36a // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.27.0 // indirect
	go.opentelemetry.io/otel v1.2.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/metric v0.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.2.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
)
