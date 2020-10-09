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
	"github.com/keptn/keptn/helm-service/pkg/helm"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

func mockChartResourceEndpoints() *httptest.Server {

	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet &&
				strings.Contains(r.RequestURI, "/v1/project/sockshop/stage/production/service/carts/resource/helm%2Fcarts-generated.tgz") {
				defer r.Body.Close()

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)

				ch := helm.GetTestGeneratedChart()
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

	inData := keptnv2.EventData{
		Project: "sockshop",
		Stage:   "production",
		Service: "carts",
		Labels:  nil,
		Status:  "",
		Result:  "",
		Message: "",
	}

	tests := []struct {
		name                 string
		actionTriggeredEvent keptnv2.ActionTriggeredEventData
		wanted               keptnv2.ActionFinishedEventData
	}{
		{
			name: "validAction",
			actionTriggeredEvent: keptnv2.ActionTriggeredEventData{
				EventData: inData,
				Action: keptnv2.ActionInfo{
					Name:        "my-scaling-action",
					Action:      "scaling",
					Description: "this is a unit test",
					Value:       "1",
				},
				Problem: keptnv2.ProblemDetails{},
			},
			wanted: keptnv2.ActionFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Stage:   "production",
					Service: "carts",
					Labels:  nil,
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultPass,
					Message: "Successfully executed scaling action",
				},
				Action: keptnv2.ActionData{
					GitCommit: "123-456",
				},
			},
		},
		{
			name: "invalidAction",
			actionTriggeredEvent: keptnv2.ActionTriggeredEventData{
				EventData: inData,
				Action: keptnv2.ActionInfo{
					Name:        "my-scaling-action",
					Action:      "scaling",
					Description: "this is a unit test",
					Value:       "byOne",
				},
				Problem: keptnv2.ProblemDetails{},
			},
			wanted: keptnv2.ActionFinishedEventData{
				EventData: keptnv2.EventData{
					Project: "sockshop",
					Stage:   "production",
					Service: "carts",
					Labels:  nil,
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultFailed,
					Message: "could not parse action.value to int",
				},
				Action: keptnv2.ActionData{
					GitCommit: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ce := cloudevents.NewEvent()
			ce.SetData(cloudevents.ApplicationJSON, tt.actionTriggeredEvent)

			keptnHandler, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

			mockHandler := &HandlerBase{
				keptnHandler:     keptnHandler,
				helmExecutor:     helm.NewHelmMockExecutor(),
				configServiceURL: ts.URL,
			}

			a := ActionTriggeredHandler{
				Handler:mockHandler,
			}

			resp := a.handleScaling(tt.actionTriggeredEvent)
			if !reflect.DeepEqual(resp, tt.wanted) {
				t.Error("unexpected action.finished response")
			}
		})
	}
}
