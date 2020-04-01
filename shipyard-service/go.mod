module keptn/shipyard-service

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-a
	github.com/keptn/kubernetes-utils v0.0.0-20200331061313-d431a8f57bef
	github.com/magiconair/properties v1.8.1
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200401063654-dcd515a62214
