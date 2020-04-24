package event_handler

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
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
		Logger *keptnutils.Logger
		Event  cloudevents.Event
	}
	tests := []struct {
		name          string
		fields        fields
		sloAvailable  bool
		wantEventType string
		wantErr       bool
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
			name: "No SLO file available",
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
				},
			},
			sloAvailable:  false,
			wantEventType: keptnevents.EvaluationDoneEventType,
			wantErr:       false,
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
			eh := &StartEvaluationHandler{
				Logger:       tt.fields.Logger,
				Event:        tt.fields.Event,
				KeptnHandler: keptnHandler,
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

func stringp(s string) *string {
	return &s
}
