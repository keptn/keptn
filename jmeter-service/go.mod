module github.com/keptn/keptn/jmeter-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-a
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
)

replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200331064125-beb163c41650
