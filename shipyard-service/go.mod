module keptn/shipyard-service

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.6.1-0.20200310164337-e6ef292da6a0
	github.com/magiconair/properties v1.8.1
	gopkg.in/yaml.v2 v2.2.4
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
