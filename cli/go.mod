module github.com/keptn/keptn/cli

go 1.13

require (
	github.com/Masterminds/sprig/v3 v3.1.0 // indirect
	github.com/alecthomas/jsonschema v0.0.0-20201129101101-7b852d451add
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/danieljoos/wincred v1.1.0 // indirect
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/elazarl/goproxy v0.0.0-20180725130230-947c36da3153 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-test/deep v1.0.7
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/keptn/go-utils v0.8.0-alpha.0.20210203161317-67ac0f2ba06d
	github.com/keptn/kubernetes-utils v0.8.0-alpha
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mattn/go-shellwords v1.0.11
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/cli-runtime v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c // indirect
	k8s.io/kubectl v0.17.2
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
