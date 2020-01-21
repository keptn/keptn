module github.com/keptn/keptn/helm-service

go 1.12

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.0
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271 // indirect
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/helm v2.14.3+incompatible
	sigs.k8s.io/yaml v1.1.0
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
