module github.com/keptn/keptn/configuration-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.9.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.24
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/swag v0.19.14
	github.com/go-openapi/validate v0.20.2
	github.com/go-test/deep v1.0.5 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/martian v2.1.0+incompatible
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/keptn/go-utils v0.8.0-alpha.0.20210215140024-05aa809e2b21
	github.com/keptn/kubernetes-utils v0.8.0-alpha.0.20210210151638-4b5b09c2e79c
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/otiai10/copy v1.4.2
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/stretchr/testify v1.7.0
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
