module github.com/keptn/keptn/webhook-service

go 1.16

require (
	github.com/keptn/go-utils v0.8.5-0.20210813063645-1bde8ec36705
	github.com/keptn/keptn/go-sdk v0.0.0-00010101000000-000000000000
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
)

replace github.com/keptn/keptn/go-sdk => ./../go-sdk
