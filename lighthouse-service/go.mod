module github.com/keptn/keptn/lighthouse-service

go 1.13

require (
	github.com/Azure/go-autorest/autorest v0.9.2 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.0 // indirect
	github.com/cloudevents/sdk-go v0.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.4.1 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.1-a
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/gnatsd v1.4.1
	github.com/nats-io/nats-server v1.4.1
	github.com/stretchr/testify v1.4.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200327001022-6496210b90e8 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace github.com/keptn/go-utils => github.com/keptn/go-utils v0.6.1-0.20200331064125-beb163c41650
