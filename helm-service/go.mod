module github.com/keptn/keptn/helm-service

go 1.13

require (
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.4.4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.8.0-alpha.0.20210215163126-910e83f1bd1c
	github.com/keptn/kubernetes-utils v0.8.0-alpha.0.20210210151638-4b5b09c2e79c
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.7.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.20.3
	k8s.io/apimachinery v0.20.3
	k8s.io/cli-runtime v0.20.3
	k8s.io/client-go v0.20.3
	k8s.io/kubectl v0.20.3
	sigs.k8s.io/yaml v1.2.0
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
