module github.com/keptn/keptn/distributor

go 1.16

require (
	github.com/benbjohnson/clock v1.3.0
	github.com/cloudevents/sdk-go/protocol/nats/v2 v2.8.0
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.13.1-0.20220318125157-fe974e59cc65
	github.com/nats-io/nats-server/v2 v2.7.2
	github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.27.0
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
)
