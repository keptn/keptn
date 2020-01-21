module keptn/shipyard-service

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.5.1-0.20200114112413-96a2056eae33
	github.com/magiconair/properties v1.8.1
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/helm v2.14.3+incompatible // indirect
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3

replace github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
