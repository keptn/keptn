module keptn/gatekeeper-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/uuid v1.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-a
	github.com/keptn/kubernetes-utils v0.0.0-20200401103501-ae44a5ee0656
	k8s.io/client-go v11.0.0+incompatible // indirect
)

replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200331064125-beb163c41650
