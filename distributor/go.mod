module github.com/keptn/keptn/distributor

go 1.16

require (
	github.com/cloudevents/sdk-go/protocol/nats/v2 v2.6.1
	github.com/cloudevents/sdk-go/v2 v2.6.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.11.1-0.20211124091533-fe1fb0c473a4
	github.com/nats-io/nats-server/v2 v2.6.6
	github.com/nats-io/nats.go v1.13.1-0.20211122170419-d7c1d78a50fc
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
)

replace (
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59 => golang.org/x/crypto v0.0.0-20201203163018-be400aefbc4c
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
)
