module github.com/keptn/keptn/helm-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-0.20200123122855-3dd223374d76
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.4.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.0.2
	k8s.io/api v0.0.0-20191016110408-35e52d86657a // kubernetes-1.16.2
	k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8 // kubernetes-1.16.2
	sigs.k8s.io/yaml v1.1.0
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3

replace github.com/docker/docker => github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
