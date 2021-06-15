module github.com/keptn/keptn/helm-service

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/godbus/dbus v0.0.0-20190422162347-ade71ed3457e // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.4.4
	github.com/golangplus/bytes v0.0.0-20160111154220-45c989fe5450 // indirect
	github.com/golangplus/fmt v0.0.0-20150411045040-2a5d6d7d2995 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.8.5-0.20210526102329-8b22839909f4
	github.com/keptn/kubernetes-utils v0.8.3-0.20210615103534-2f8cee055a20
	github.com/kinbiko/jsonassert v1.0.1
	github.com/opencontainers/runtime-tools v0.0.0-20181011054405-1d69bd0f9c39 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/gocapability v0.0.0-20170704070218-db04d3cc01c8 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.6.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/cli-runtime v0.21.0
	k8s.io/client-go v0.21.1
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kubectl v0.21.0
	k8s.io/kubernetes v1.13.0 // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	sigs.k8s.io/yaml v1.2.0
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
)
