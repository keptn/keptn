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
	github.com/google/uuid v1.2.0
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.8.1-0.20210324101115-21d3e25da004
	github.com/keptn/kubernetes-utils v0.8.1-0.20210308145316-785b480a4e52
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
	vbom.ml/util v0.0.0-20160121211510-db5cfe13f5cc // indirect
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
