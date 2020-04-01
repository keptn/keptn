module github.com/keptn/keptn/cli

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.10.0 // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/danieljoos/wincred v1.0.2 // indirect
	github.com/docker/docker-credential-helpers v0.6.3
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/go-version v1.2.0
	github.com/keptn/go-utils v0.6.1-a
	github.com/keptn/kubernetes-utils v0.0.0-20200331061313-d431a8f57bef
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-shellwords v1.0.10
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.2.8
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/helm v2.14.3+incompatible
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
)

replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200401063654-dcd515a62214
