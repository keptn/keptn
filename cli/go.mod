module github.com/keptn/keptn/cli

go 1.13

require (
	github.com/alecthomas/jsonschema v0.0.0-20210214200137-e6fc2822d59d
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/danieljoos/wincred v1.1.0 // indirect
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/go-openapi/validate v0.19.5 // indirect
	github.com/go-test/deep v1.0.7
	github.com/google/uuid v1.2.0
	github.com/hashicorp/go-version v1.2.0
	github.com/keptn/go-utils v0.8.2-0.20210325092328-c8d257ddf2fc
	github.com/keptn/kubernetes-utils v0.8.1
	github.com/mattn/go-shellwords v1.0.11
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.5.1
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/cli-runtime v0.20.4
	k8s.io/client-go v0.20.4
	k8s.io/kubectl v0.20.4
	rsc.io/letsencrypt v0.0.3 // indirect
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
//replace (
//	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
//	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
//)
