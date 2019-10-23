module keptn/eventbroker

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/btree v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.2.3
)

// replace latest version with 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
