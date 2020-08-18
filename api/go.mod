module github.com/keptn/keptn/api

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/gbrlsnchs/jwt/v2 v2.0.0
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/go-openapi/validate v0.19.5
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/keptn/go-utils v0.6.3-0.20200818082305-f51180cb3d15
	github.com/keptn/kubernetes-utils v0.2.0
	github.com/kinbiko/jsonassert v1.0.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magiconair/properties v1.8.1
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
