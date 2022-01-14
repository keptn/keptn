module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/gin-gonic/gin v1.7.7 // indirect
	github.com/go-test/deep v1.0.8 // indirect
	github.com/google/uuid v1.3.0
	github.com/imroc/req v0.3.2
	github.com/jeremywohl/flatten v1.0.1 // indirect
	github.com/keptn/go-utils v0.11.1-0.20220110112638-4e1eff727e81
	github.com/keptn/kubernetes-utils v0.10.1-0.20211102080304-e59377afdc8b
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.8 // indirect
	github.com/tryvium-travels/memongo v0.3.2 // indirect
	go.mongodb.org/mongo-driver v1.7.5 // indirect
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
