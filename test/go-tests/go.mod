module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/google/uuid v1.2.0
	github.com/imroc/req v0.3.0
	github.com/keptn/go-utils v0.8.6-0.20210623084135-feb84d6fd126
	github.com/keptn/keptn/shipyard-controller v0.0.0-20210503133401-8c1194432b46
	github.com/keptn/kubernetes-utils v0.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/apimachinery v0.21.1
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
