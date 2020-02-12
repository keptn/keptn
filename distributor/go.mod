module keptn/distributor

go 1.12

require (
	github.com/cloudevents/sdk-go v1.0.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/nats-io/nats.go v1.9.1
	github.com/nats-io/stan.go v0.6.0
)

// replace cloudevents/sdk-go latest version with 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
