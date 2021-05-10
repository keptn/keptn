module github.com/keptn/keptn/configuration-service

go 1.16

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.9.0 // indirect
	github.com/go-openapi/errors v0.19.9
	github.com/go-openapi/loads v0.20.2
	github.com/go-openapi/runtime v0.19.28
	github.com/go-openapi/spec v0.20.3
	github.com/go-openapi/strfmt v0.20.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.2
	github.com/google/martian v2.1.0+incompatible
	github.com/jessevdk/go-flags v1.4.0
	github.com/keptn/go-utils v0.8.4-0.20210506073402-95bb36d6f884
	github.com/keptn/kubernetes-utils v0.8.2-0.20210506073412-a06581ae8a26
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/otiai10/copy v1.6.0
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/stretchr/testify v1.7.0
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
