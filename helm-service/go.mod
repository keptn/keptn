module github.com/keptn/keptn/helm-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.7.1
	github.com/keptn/kubernetes-utils v0.2.1-0.20201023093247-d74fe9c6b4b7
	github.com/kinbiko/jsonassert v1.0.1
	github.com/stretchr/testify v1.5.1
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/cli-runtime v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/kubectl v0.17.2
	sigs.k8s.io/yaml v1.2.0
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
