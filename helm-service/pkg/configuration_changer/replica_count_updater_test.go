package configuration_changer

import (
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/pkg/helmtest"
	"helm.sh/helm/v3/pkg/chart"
	"reflect"
	"testing"
)

func TestIncreaseReplicaCount(t *testing.T) {

	const expectedPrimaryDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: carts-primary
spec:
  replicas: 3
  selector:
    matchLabels:
      app: carts-primary
  strategy:
    rollingUpdate:
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: carts-primary
    spec:
      containers:
      - image: docker.io/keptnexamples/carts:0.8.1
        imagePullPolicy: IfNotPresent
        name: carts
        resources: {}
status: {}
`
	expectedChart := chart.Chart{
		Raw: nil,
		Metadata: &chart.Metadata{
			Name:       "carts-generated",
			Version:    "0.1.0",
			Keywords:   []string{"deployment_strategy=" + keptnevents.Duplicate.String()},
			APIVersion: "v2",
		},
		Lock: nil,
		Templates: []*chart.File{
			{
				Name: "carts-canary-istio-destinationrule.yaml",
				Data: []byte(helm.GeneratedCanaryDestinationRule),
			},
			{
				Name: "carts-canary-service.yaml",
				Data: []byte(helm.GeneratedCanaryService),
			},
			{
				Name: "carts-istio-virtualservice.yaml",
				Data: []byte(helm.GeneratedVirtualService),
			},
			{
				Name: "carts-primary-deployment.yaml",
				Data: []byte(expectedPrimaryDeployment),
			},
			{
				Name: "carts-primary-istio-destinationrule.yaml",
				Data: []byte(helm.GeneratedPrimaryDestinationRule),
			},
			{
				Name: "carts-primary-service.yaml",
				Data: []byte(helm.GeneratedPrimaryService),
			},
		},
	}

	inputChart := helmtest.GetGeneratedChart()
	updater := NewReplicaCountUpdater(2)
	updater.Update(&inputChart)

	if !reflect.DeepEqual(inputChart, expectedChart) {
		t.Error("inputChart does not match expected chart")
	}
}
