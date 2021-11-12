package event_handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	event_handler_mock "github.com/keptn/keptn/lighthouse-service/event_handler/fake"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"

	keptnapi "github.com/keptn/go-utils/pkg/api/models"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
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
	wg := &sync.WaitGroup{}
	ctx := cloudevents.WithEncodingStructured(context.WithValue(context.Background(), GracefulShutdownKey, wg))
	var returnSlo bool
	var sloFileContent string
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
					encodedSLOContent := base64.StdEncoding.EncodeToString([]byte(sloFileContent))
					resourceURI := "slo.yaml"
					sloResource := &keptnapi.Resource{ResourceURI: &resourceURI, ResourceContent: encodedSLOContent}
					marshal, _ := json.Marshal(sloResource)
					w.WriteHeader(http.StatusOK)
					w.Write(marshal)
				} else if returnServiceNotFound {
					errObj := &keptnapi.Error{Code: 404, Message: stringp("Service not found")}
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
		Logger           *keptncommon.Logger
		Event            cloudevents.Event
		SLOFileRetriever SLOFileRetriever
	}
	tests := []struct {
		name                string
		fields              fields
		sloAvailable        bool
		sloFileContent      string
		serviceNotAvailable bool
		wantEventType       []string
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
			name: "No SLO file available -  send get-sli event",
			fields: fields{
				Event: getStartEvaluationEvent(),
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: &event_handler_mock.ResourceHandlerMock{GetServiceResourceFunc: func(project string, stage string, service string, resourceURI string) (*keptnapi.Resource, error) {
						return nil, nil
					}},
				},
			},
			sloAvailable:  false,
			wantEventType: []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName)},
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
			name: "Service not available - return evaluation.finished event",
			fields: fields{
				Event: getStartEvaluationEvent(),
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: api.NewResourceHandler(os.Getenv("CONFIGURATION_SERVICE")),
					ServiceHandler:  api.NewServiceHandler(os.Getenv("CONFIGURATION_SERVICE")),
				},
			},
			sloAvailable:        false,
			serviceNotAvailable: true,
			wantEventType:       []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)},
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
				Event: getStartEvaluationEvent(),
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: api.NewResourceHandler(os.Getenv("CONFIGURATION_SERVICE")),
					ServiceHandler:  api.NewServiceHandler(os.Getenv("CONFIGURATION_SERVICE")),
				},
			},
			sloAvailable:  false,
			wantEventType: []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName)},
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
		{
			name: "Error during SLO file parsing - send finished event with error",
			fields: fields{
				Event: getStartEvaluationEvent(),
				SLOFileRetriever: SLOFileRetriever{
					ResourceHandler: api.NewResourceHandler(os.Getenv("CONFIGURATION_SERVICE")),
					ServiceHandler:  api.NewServiceHandler(os.Getenv("CONFIGURATION_SERVICE")),
				},
			},
			sloAvailable:   true,
			sloFileContent: "invalid",
			wantEventType:  []string{keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)},
			wantErr:        false,
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
			sloFileContent = tt.sloFileContent
			returnServiceNotFound = tt.serviceNotAvailable
			eh := &StartEvaluationHandler{
				Event:        tt.fields.Event,
				KeptnHandler: keptnHandler,
				SLIProviderConfig: &MockSLIProviderConfig{
					ProjectSLIProvider: tt.ProjectSLIProvider,
					DefaultSLIProvider: tt.DefaultSLIProvider,
				},
				SLOFileRetriever: tt.fields.SLOFileRetriever,
			}
			if err := eh.HandleEvent(ctx); (err != nil) != tt.wantErr {
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
			Extensions: map[string]interface{}{
				"shkeptncontext": "my-context",
			},
		},
		DataEncoded: []byte(`{
    "project": "sockshop",
    "stage": "staging",
    "service": "carts",
    "testStrategy": "",
    "deploymentStrategy": "direct",
	"evaluation": {
		"timeframe": "5m"
    },
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

func Test_getEvaluationTimestamps(t *testing.T) {
	type args struct {
		e *keptnv2.EvaluationTriggeredEventData
	}
	tests := []struct {
		name      string
		args      args
		wantStart string
		wantEnd   string
		wantErr   bool
	}{
		{
			name: "get evaluation timestamps",
			args: args{
				e: &keptnv2.EvaluationTriggeredEventData{
					Test: keptnv2.Test{
						Start: "2006-01-02T15:04:05.000Z",
						End:   "2006-01-02T15:05:05.000Z",
					},
					Evaluation: keptnv2.Evaluation{
						Start: "2006-01-02T15:14:04.000Z",
						End:   "2006-01-02T15:15:05.000Z",
					},
				},
			},
			wantStart: "2006-01-02T15:14:04.000Z",
			wantEnd:   "2006-01-02T15:15:05.000Z",
			wantErr:   false,
		},
		{
			name: "get evaluation timestamp from timeframe",
			args: args{
				e: &keptnv2.EvaluationTriggeredEventData{
					Evaluation: keptnv2.Evaluation{
						Timeframe: "10m",
					},
				},
			},
			wantStart: timeutils.GetKeptnTimeStamp(time.Now().UTC().Add(-10 * time.Minute)),
			wantEnd:   timeutils.GetKeptnTimeStamp(time.Now().UTC()),
			wantErr:   false,
		},
		{
			name: "fallback to test timestamps",
			args: args{
				e: &keptnv2.EvaluationTriggeredEventData{
					Test: keptnv2.Test{
						Start: "2006-01-02T15:04:05.000Z",
						End:   "2006-01-02T15:05:05.000Z",
					},
				},
			},
			wantStart: "2006-01-02T15:04:05.000Z",
			wantEnd:   "2006-01-02T15:05:05.000Z",
			wantErr:   false,
		},
		{
			name: "no timestamps provided",
			args: args{
				e: &keptnv2.EvaluationTriggeredEventData{},
			},
			wantStart: "",
			wantEnd:   "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := getEvaluationTimestamps(tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEvaluationTimestamps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				require.NotNil(t, err)
				require.Empty(t, gotStart)
				require.Empty(t, gotEnd)
				return
			}

			timeFormat := timeutils.KeptnTimeFormatISO8601
			gotParsedStart, err := time.Parse(timeFormat, gotStart)
			require.Nil(t, err)

			gotParsedEnd, err := time.Parse(timeFormat, gotEnd)
			require.Nil(t, err)

			wantParsedStart, _ := time.Parse(timeFormat, tt.wantStart)
			wantParsedEnd, _ := time.Parse(timeFormat, tt.wantEnd)

			require.WithinDuration(t, wantParsedStart, gotParsedStart, time.Minute)
			require.WithinDuration(t, wantParsedEnd, gotParsedEnd, time.Minute)
		})
	}
}
