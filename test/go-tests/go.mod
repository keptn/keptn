module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.6.1
	github.com/google/uuid v1.3.0
	github.com/imroc/req v0.3.1
	github.com/keptn/go-utils v0.10.1-0.20211103140516-b6f48452063e
	github.com/keptn/keptn/shipyard-controller v0.0.0-20210503133401-8c1194432b46
	github.com/keptn/kubernetes-utils v0.10.1-0.20211102080304-e59377afdc8b
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.3
	k8s.io/client-go v0.22.2
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
