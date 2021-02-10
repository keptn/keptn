module github.com/keptn/keptn/api

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/spec v0.19.8
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/go-test/deep v1.0.7
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.8.0-alpha.0.20210203161317-67ac0f2ba06d
	github.com/keptn/kubernetes-utils v0.8.0-alpha.0.20210208085038-093c00d82da4
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
