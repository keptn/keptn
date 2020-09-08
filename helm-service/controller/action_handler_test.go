package controller

import (
	"encoding/base64"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller/helm"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	"helm.sh/helm/v3/pkg/chart"
)

func getGeneratedChart() chart.Chart {
	return chart.Chart{
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
				Data: []byte(helm.GeneratedPrimaryDeployment),
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
}

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

	a := &ActionTriggeredHandler{
		helmExecutor:     helm.NewHelmMockExecutor(),
		configServiceURL: "",
	}

	inputChart := getGeneratedChart()
	a.increaseReplicaCount(&inputChart, 2)

	if !reflect.DeepEqual(inputChart, expectedChart) {
		t.Error("inputChart does not match expected chart")
	}
}

func mockChartResourceEndpoints() *httptest.Server {

	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet &&
				strings.Contains(r.RequestURI, "/v1/project/sockshop/stage/production/service/carts/resource/helm%2Fcarts-generated.tgz") {
				defer r.Body.Close()

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)

				ch := getGeneratedChart()
				chPackage, _ := keptnutils.PackageChart(&ch)

				resp := models.Resource{
					ResourceContent: base64.StdEncoding.EncodeToString(chPackage),
					ResourceURI:     stringp("helm/carts-generated.tgz"),
				}
				data, _ := json.Marshal(resp)
				w.Write(data)
			} else if r.Method == http.MethodPost &&
				strings.Contains(r.RequestURI, "v1/project/sockshop/stage/production/service/carts/resource") {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)

				resp := models.Version{
					Version: "123-456",
				}
				data, _ := json.Marshal(resp)
				w.Write(data)
			}
		}),
	)
}

func TestHandleScaling(t *testing.T) {
	ts := mockChartResourceEndpoints()
	defer ts.Close()

	tests := []struct {
		name                 string
		actionTriggeredEvent keptnevents.ActionTriggeredEventData
		wanted               keptnevents.ActionFinishedEventData
	}{
		{
			name: "validAction",
			actionTriggeredEvent: keptnevents.ActionTriggeredEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "production",
				Action: keptnevents.ActionInfo{
					Name:        "my-scaling-action",
					Action:      "scaling",
					Description: "this is a unit test",
					Value:       "1",
				},
				Problem: keptnevents.ProblemDetails{},
				Labels:  nil,
			},
			wanted: keptnevents.ActionFinishedEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "production",
				Action: keptnevents.ActionResult{
					Result: "pass",
					Status: keptnevents.ActionStatusSucceeded,
				},
				Labels: nil,
			},
		},
		{
			name: "invalidAction",
			actionTriggeredEvent: keptnevents.ActionTriggeredEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "production",
				Action: keptnevents.ActionInfo{
					Name:        "my-scaling-action",
					Action:      "scaling",
					Description: "this is a unit test",
					Value:       "byOne",
				},
				Problem: keptnevents.ProblemDetails{},
				Labels:  nil,
			},
			wanted: keptnevents.ActionFinishedEventData{
				Project: "sockshop",
				Service: "carts",
				Stage:   "production",
				Action: keptnevents.ActionResult{
					Result: "strconv.Atoi: parsing \"byOne\": invalid syntax",
					Status: keptnevents.ActionStatusErrored,
				},
				Labels: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ce := cloudevents.NewEvent()
			ce.SetData(cloudevents.ApplicationJSON, tt.actionTriggeredEvent)

			keptnHandler, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

			a := &ActionTriggeredHandler{
				helmExecutor:     helm.NewHelmMockExecutor(),
				keptnHandler:     keptnHandler,
				configServiceURL: ts.URL,
			}

			resp := a.handleScaling(tt.actionTriggeredEvent)
			if !reflect.DeepEqual(resp, tt.wanted) {
				t.Error("unexpected action.finished response")
			}
		})
	}
}
