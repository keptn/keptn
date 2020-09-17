package event_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/swag"
	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

const TEST_PORT = 8370

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}

type MockSLIProviderConfig struct {
	ProjectSLIProvider struct {
		val string
		err error
	}
	DefaultSLIProvider struct {
		val string
		err error
	}
}

func (m *MockSLIProviderConfig) GetDefaultSLIProvider() (string, error) {
	return m.DefaultSLIProvider.val, m.DefaultSLIProvider.err
}

func (m *MockSLIProviderConfig) GetSLIProvider(project string) (string, error) {
	return m.ProjectSLIProvider.val, m.DefaultSLIProvider.err
}

func TestStartEvaluationHandler_HandleEvent(t *testing.T) {

	type ceTypeEvent struct {
		Type string `json:"type"`
	}
	ch := make(chan string)

	var returnSlo bool
	var returnServiceNotFound bool
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodPost && strings.Contains(r.RequestURI, "/events") {
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {

				}
				fmt.Printf("Received request: %v\n", string(body))

				event := &ceTypeEvent{}
				_ = json.Unmarshal(body, &event)

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{}`))
				go func() { ch <- event.Type }()
			} else if strings.Contains(r.RequestURI, "/configuration") {
				if returnSlo {

				} else if returnServiceNotFound {
					errObj := &keptnapi.Error{Code: 404, Message: swag.String("Service not found")}
					marshal, _ := json.Marshal(errObj)
					w.WriteHeader(404)
					w.Write(marshal)
				} else {
					w.WriteHeader(404)
					w.Write([]byte(``))
				}
			}
		}),
	)
	defer ts.Close()

	os.Setenv("EVENTBROKER", ts.URL+"/events")
	os.Setenv("CONFIGURATION_SERVICE", ts.URL+"/configuration")

	////////// TEST DEFINITION ///////////
	type fields struct {
		Logger *keptnutils.Logger
		Event  cloudevents.Event
	}
	tests := []struct {
		name                string
		fields              fields
		sloAvailable        bool
		serviceNotAvailable bool
		wantEventType       string
		wantErr             bool
		ProjectSLIProvider  struct {
			val string
			err error
		}
		DefaultSLIProvider struct {
			val string
			err error
		}
	}{
		{
			name: "No test strategy set",
			fields: fields{
				Logger: keptnutils.NewLogger("", "", ""),
				Event: cloudevents.Event{
					Context: &cloudevents.EventContextV02{
						SpecVersion: "0.2",
						Type:        "sh.keptn.events.tests-finished",
						Source:      types.URLRef{},
						ID:          "",
						Time:        nil,
						SchemaURL:   nil,
						ContentType: stringp("application/json"),
						Extensions:  nil,
					},
					Data: []byte(`{
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "",
    "deploymentStrategy": "direct",
    "start": "2019-09-01 12:00:00",
    "end": "2019-09-01 12:05:00",
    "labels": {
      "testid": "12345",
      "buildnr": "build17",
      "runby": "JohnDoe"
    },
    "result": "pass"
  }`),
					DataEncoded: false,
				},
			},
			sloAvailable:  false,
			wantEventType: keptnevents.EvaluationDoneEventType,
			wantErr:       false,
		},
		{
			name: "No SLO file available -  send get-sli event",
			fields: fields{
				Logger: keptnutils.NewLogger("", "", ""),
				Event:  getStartEvaluationEvent(),
			},
			sloAvailable:  false,
			wantEventType: keptnevents.InternalGetSLIEventType,
			wantErr:       false,
			ProjectSLIProvider: struct {
				val string
				err error
			}{
				val: "my-sli-provider",
				err: nil,
			},
			DefaultSLIProvider: struct {
				val string
				err error
			}{},
		},
		{
			name: "Service not available - return evaluation.done event",
			fields: fields{
				Logger: keptnutils.NewLogger("", "", ""),
				Event:  getStartEvaluationEvent(),
			},
			sloAvailable:        false,
			serviceNotAvailable: true,
			wantEventType:       keptnevents.InternalGetSLIEventType,
			wantErr:             false,
			ProjectSLIProvider: struct {
				val string
				err error
			}{
				val: "my-sli-provider",
				err: nil,
			},
			DefaultSLIProvider: struct {
				val string
				err error
			}{},
		},
		{
			name: "No SLI provider configured for project - use default",
			fields: fields{
				Logger: keptnutils.NewLogger("", "", ""),
				Event:  getStartEvaluationEvent(),
			},
			sloAvailable:  false,
			wantEventType: keptnevents.InternalGetSLIEventType,
			wantErr:       false,
			ProjectSLIProvider: struct {
				val string
				err error
			}{
				val: "",
				err: errors.New(""),
			},
			DefaultSLIProvider: struct {
				val string
				err error
			}{
				val: "default-sli-provider",
				err: nil,
			},
		},
	}
	////////// TEST EXECUTION ///////////
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			keptnHandler, _ := keptnutils.NewKeptn(&tt.fields.Event, keptnutils.KeptnOpts{
				EventBrokerURL:          os.Getenv("EVENTBROKER"),
				ConfigurationServiceURL: os.Getenv("CONFIGURATION_SERVICE"),
			})
			returnSlo = tt.sloAvailable
			returnServiceNotFound = tt.serviceNotAvailable
			eh := &StartEvaluationHandler{
				Event:        tt.fields.Event,
				KeptnHandler: keptnHandler,
				SLIProviderConfig: &MockSLIProviderConfig{
					ProjectSLIProvider: tt.ProjectSLIProvider,
					DefaultSLIProvider: tt.DefaultSLIProvider,
				},
			}
			if err := eh.HandleEvent(); (err != nil) != tt.wantErr {
				t.Errorf("HandleEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			select {
			case msg := <-ch:
				t.Logf("Received event type: %v", msg)
				if msg != tt.wantEventType {
					t.Errorf("HandleEvent() sent event type = %v, wantEventType %v", msg, tt.wantEventType)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Message did not make it to the receiver")
			}

		})
	}
}

func getStartEvaluationEvent() cloudevents.Event {
	return cloudevents.Event{
		Context: &cloudevents.EventContextV02{
			SpecVersion: "0.2",
			Type:        "sh.keptn.events.tests-finished",
			Source:      types.URLRef{},
			ID:          "",
			Time:        nil,
			SchemaURL:   nil,
			ContentType: stringp("application/json"),
			Extensions:  nil,
		},
		Data: []byte(`{
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "performance",
    "deploymentStrategy": "direct",
    "start": "2019-09-01 12:00:00",
    "end": "2019-09-01 12:05:00",
    "labels": {
      "testid": "12345",
      "buildnr": "build17",
      "runby": "JohnDoe"
    },
    "result": "pass"
  }`),
		DataEncoded: false,
	}
}

func stringp(s string) *string {
	return &s
}
