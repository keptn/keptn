module keptn/gatekeeper-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/uuid v1.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.3-0.20201021140127-a974f80c5982
	github.com/keptn/kubernetes-utils v0.2.1-0.20201102114243-598bca4e91ad
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
