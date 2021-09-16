module github.com/keptn/keptn/remediation-service

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/ghodss/yaml v1.0.0
	github.com/keptn/go-utils v0.9.0
	github.com/keptn/keptn/go-sdk v0.0.0-20210830093722-f988e684471f
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

//replace github.com/keptn/keptn/go-sdk => ../go-sdk
