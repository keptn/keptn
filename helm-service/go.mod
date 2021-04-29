module github.com/keptn/keptn/helm-service

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/validate v0.19.5 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.4.4
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.8.4-0.20210429120423-086f16f75234
	github.com/keptn/kubernetes-utils v0.8.1
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.7.0
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/cli-runtime v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/kubectl v0.20.4
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
	sigs.k8s.io/yaml v1.2.0
	vbom.ml/util v0.0.0-20160121211510-db5cfe13f5cc // indirect
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
