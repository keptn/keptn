module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/google/uuid v1.3.0
	github.com/imroc/req v0.3.0
	github.com/keptn/go-utils v0.9.0
	github.com/keptn/keptn/shipyard-controller v0.0.0-20210503133401-8c1194432b46
	github.com/keptn/kubernetes-utils v0.8.3
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.22.1
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
