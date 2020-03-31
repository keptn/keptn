module github.com/keptn/keptn/helm-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-0.20200331064125-beb163c41650
	github.com/keptn/kubernetes-utils v0.0.0-20200331061313-d431a8f57bef
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.4.0
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/helm v2.14.3+incompatible
	sigs.k8s.io/yaml v1.1.0
)

// TODO: remove this replace instruction when the go-utils have been updated officially
replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200331064125-beb163c41650
