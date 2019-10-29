module keptn/platform-support/openshift-route-service

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.3.0
	gopkg.in/yaml.v2 v2.2.4
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
