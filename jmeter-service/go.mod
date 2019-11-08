module github.com/keptn/keptn/jmeter-service

go 1.12

require (
	cloud.google.com/go v0.38.0 // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.3.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	google.golang.org/appengine v1.5.0 // indirect
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
)

// replace cloudevents/sdk-go latest version with 0.7.0
replace (
	github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
	github.com/keptn/go-utils => ../../go-utils
)
