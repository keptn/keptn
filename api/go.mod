module github.com/keptn/keptn/api

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.24
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.20.2
	github.com/go-test/deep v1.0.7
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.8.0-alpha.0.20210212091339-6378ba149260
	github.com/keptn/kubernetes-utils v0.8.0-alpha.0.20210210151638-4b5b09c2e79c
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.3
	k8s.io/apimachinery v0.20.3
	k8s.io/client-go v0.20.3
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
