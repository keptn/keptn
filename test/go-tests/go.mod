module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/gin-gonic/gin v1.7.4 // indirect
	github.com/go-test/deep v1.0.5 // indirect
	github.com/google/uuid v1.3.0
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imroc/req v0.3.0
	github.com/jeremywohl/flatten v1.0.1 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/keptn/go-utils v0.9.0
	github.com/keptn/keptn/shipyard-controller v0.0.0-20210901055628-de03c96ac1a9
	github.com/keptn/kubernetes-utils v0.8.3
	github.com/ory/dockertest/v3 v3.6.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.0 // indirect
	go.mongodb.org/mongo-driver v1.7.1 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.2 // indirect
	k8s.io/klog/v2 v2.9.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
