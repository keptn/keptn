module github.com/keptn/keptn/webhook-service

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1 // indirect
	github.com/keptn/go-utils v0.8.5-0.20210813063645-1bde8ec36705
	github.com/keptn/keptn/go-sdk v0.0.0-20210830082026-b7006e7b45bd
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
)

//replace github.com/keptn/keptn/go-sdk => ../go-sdk
