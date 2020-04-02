module keptn/shipyard-service

go 1.12

require (
	github.com/Azure/go-autorest/autorest v0.2.0 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.6.0 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/keptn/go-utils v0.6.1-0.20200402095926-a4ac0d7044a0
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/time v0.0.0-20191023065245-6d3f0bb11be5 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1 // indirect
	gopkg.in/yaml.v2 v2.2.4
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab // indirect
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/helm v2.14.3+incompatible // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

// replace cloudevents/sdk-go with version 0.7.0
replace github.com/cloudevents/sdk-go => github.com/cloudevents/sdk-go v0.0.0-20190509003705-56931988abe3
