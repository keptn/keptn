module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/google/uuid v1.3.0
	github.com/imroc/req v0.3.2
	github.com/keptn/go-utils v0.11.1-0.20220112141105-c3a5144daf3a
	github.com/keptn/keptn/shipyard-controller v0.0.0-00010101000000-000000000000
	github.com/keptn/kubernetes-utils v0.10.1-0.20211102080304-e59377afdc8b
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
