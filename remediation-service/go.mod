module github.com/keptn/keptn/remediation-service

go 1.13

require (
	cloud.google.com/go v0.40.0 // indirect
	github.com/Azure/go-autorest/autorest v0.9.0 // indirect
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/strfmt v0.20.0
	github.com/go-openapi/validate v0.20.2 // indirect
	github.com/go-test/deep v1.0.7
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/gophercloud/gophercloud v0.9.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.8.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.17.0 // indirect
	k8s.io/apimachinery v0.17.0 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
)
