module keptn/distributor

go 1.12

require (
	github.com/cloudevents/sdk-go v0.9.2
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.0.0
)

// replace cloudevents/sdk-go latest version with 0.7.0
replace (
	github.com/keptn/go-utils => ../../go-utils
  	github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
)	  

