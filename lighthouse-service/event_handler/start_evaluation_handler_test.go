package event_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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
		Logger *keptncommon.Logger
		Event  cloudevents.Event
	}
	tests := []struct {
		name               string
		fields             fields
		sloAvailable       bool
		wantEventType      []string
		wantErr            bool
		ProjectSLIProvider struct {
			val string
			err error
		}
		DefaultSLIProvider struct {
			val string
			err error
		}
	}{
		{
			name: "No SLO file available -  send get-sli event",
			fields: fields{
				Logger: keptncommon.NewLogger("", "", ""),
				Event:  getStartEvaluationEvent(),
			},
			sloAvailable:  false,
			wantEventType: []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnevents.InternalGetSLIEventType},
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
			name: "No SLI provider configured for project - use default",
			fields: fields{
				Logger: keptncommon.NewLogger("", "", ""),
				Event:  getStartEvaluationEvent(),
			},
			sloAvailable:  false,
			wantEventType: []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnevents.InternalGetSLIEventType},
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

			keptnHandler, _ := keptnv2.NewKeptn(&tt.fields.Event, keptncommon.KeptnOpts{
				EventBrokerURL:          os.Getenv("EVENTBROKER"),
				ConfigurationServiceURL: os.Getenv("CONFIGURATION_SERVICE"),
			})
			returnSlo = tt.sloAvailable
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

			receivedEvents := []string{}
			receivedExpected := 0
			for {
				select {
				case msg := <-ch:
					t.Logf("Received event type: %v", msg)
					receivedEvents = append(receivedEvents, msg)

					// check if all expected events have been received
					for _, want := range tt.wantEventType {
						found := false
						for _, rec := range receivedEvents {
							if rec == want {
								found = true
								break
							}
						}
						if found {
							receivedExpected = receivedExpected + 1
							break
						}
					}
					if receivedExpected == len(tt.wantEventType) {
						// received all events
						return
					}

					// check if no unexpected event has been received
					for _, rec := range receivedEvents {
						found := false
						for _, want := range tt.wantEventType {
							if want == rec {
								found = true
							}
						}
						if !found {
							t.Errorf("HandleEvent() sent event type = %v, wantEventType %v", receivedEvents, tt.wantEventType)
						}
					}

				case <-time.After(5 * time.Second):
					t.Errorf("Expected messages did not make it to the receiver")
					t.Errorf("HandleEvent() sent event type = %v, wantEventType %v", receivedEvents, tt.wantEventType)
					return
				}
			}

		})
	}
}

func getStartEvaluationEvent() cloudevents.Event {
	return cloudevents.Event{
		Context: &cloudevents.EventContextV1{
			Type:            keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName),
			Source:          types.URIRef{},
			ID:              "",
			Time:            nil,
			DataContentType: stringp("application/json"),
			Extensions:      nil,
		},
		DataEncoded: []byte(`{
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
		DataBase64: false,
	}
}

func stringp(s string) *string {
	return &s
}
