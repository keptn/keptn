package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/strfmt"
	"github.com/kelseyhightower/envconfig"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_getPubSubRecipientURL(t *testing.T) {
	type args struct {
		recipientService string
		port             string
		path             string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple service name",
			args: args{
				recipientService: "lighthouse-service",
				port:             "",
				path:             "",
			},
			want:    "http://lighthouse-service:8080",
			wantErr: false,
		},
		{
			name: "simple service name with path (prepending slash)",
			args: args{
				recipientService: "lighthouse-service",
				port:             "",
				path:             "/event",
			},
			want:    "http://lighthouse-service:8080/event",
			wantErr: false,
		},
		{
			name: "simple service name with path (without prepending slash)",
			args: args{
				recipientService: "lighthouse-service",
				port:             "",
				path:             "event",
			},
			want:    "http://lighthouse-service:8080/event",
			wantErr: false,
		},
		{
			name: "simple service name with port",
			args: args{
				recipientService: "lighthouse-service",
				port:             "666",
				path:             "",
			},
			want:    "http://lighthouse-service:666",
			wantErr: false,
		},
		{
			name: "empty recipient name",
			args: args{
				recipientService: "",
				port:             "666",
				path:             "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "HTTPS recipient",
			args: args{
				recipientService: "https://lighthouse-service",
				port:             "",
				path:             "",
			},
			want:    "https://lighthouse-service:8080",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPubSubRecipientURL(tt.args.recipientService, tt.args.port, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPubSubRecipientURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getPubSubRecipientURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeCloudEvent(t *testing.T) {
	type args struct {
		data []byte
	}
	var tests = []struct {
		name    string
		args    args
		want    *cloudevents.Event
		wantErr bool
	}{
		{
			name: "Get V0.3 CloudEvent",
			args: args{
				data: []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "0.3",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`),
			},
			want: &cloudevents.Event{
				Context: &cloudevents.EventContextV02{
					SpecVersion: "0.3",
					Type:        "sh.keptn.events.deployment-finished",
					Source: types.URLRef{
						URL: url.URL{
							Scheme:     "",
							Opaque:     "helm-service",
							User:       nil,
							Host:       "",
							Path:       "",
							RawPath:    "",
							ForceQuery: false,
							RawQuery:   "",
							Fragment:   "",
						},
					},
					ID:        "6de83495-4f83-481c-8dbe-fcceb2e0243b",
					Time:      nil,
					SchemaURL: nil,
					Extensions: map[string]interface{}{
						"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb",
					},
				},
				Data:        []byte(`""`),
				DataEncoded: false,
				DataBinary:  false,
				FieldErrors: nil,
			},
			wantErr: false,
		},
		{
			name: "Get V0.3 CloudEvent",
			args: args{
				data: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeCloudEvent(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeCloudEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Context.GetSpecVersion(), tt.want.Context.GetSpecVersion()) {
					t.Errorf("decodeCloudEvent() specVersion: got = %v, want %v", got.Context.GetSpecVersion(), tt.want.Context.GetSpecVersion())
				}
				if !reflect.DeepEqual(got.Context.GetType(), tt.want.Context.GetType()) {
					t.Errorf("decodeCloudEvent() type: got = %v, want %v", got.Context.GetType(), tt.want.Context.GetType())
				}
			}
		})
	}
}

func Test_cleanSentEventList(t *testing.T) {
	type args struct {
		sentEvents []string
		events     []*keptnmodels.KeptnContextExtendedCE
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "remove no element from list",
			args: args{
				sentEvents: []string{"id-1"},
				events: []*keptnmodels.KeptnContextExtendedCE{
					{
						ID: "id-1",
					},
					{
						ID: "id-2",
					},
				},
			},
			want: []string{"id-1"},
		},
		{
			name: "remove element from list",
			args: args{
				sentEvents: []string{"id-3"},
				events: []*keptnmodels.KeptnContextExtendedCE{
					{
						ID: "id-1",
					},
					{
						ID: "id-2",
					},
				},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanSentEventList(tt.args.sentEvents, tt.args.events); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanSentEventList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasEventBeenSent(t *testing.T) {
	type args struct {
		sentEvents []string
		eventID    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "want true",
			args: args{
				sentEvents: []string{"sent-1", "sent-2"},
				eventID:    "sent-1",
			},
			want: true,
		},
		{
			name: "want false",
			args: args{
				sentEvents: []string{"sent-1", "sent-2"},
				eventID:    "sent-X",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasEventBeenSent(tt.args.sentEvents, tt.args.eventID); got != tt.want {
				t.Errorf("hasEventBeenSent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEventsFromEndpoint(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {}))

	type args struct {
		endpoint string
		token    string
		topic    string
	}
	tests := []struct {
		name              string
		args              args
		serverHandlerFunc http.HandlerFunc
		want              []*keptnmodels.KeptnContextExtendedCE
		wantErr           bool
	}{
		{
			name: "get all events",
			args: args{
				endpoint: ts.URL,
				token:    "",
				topic:    "my-topic",
			},
			serverHandlerFunc: func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json")
				events := keptnmodels.Events{
					Events: []*keptnmodels.KeptnContextExtendedCE{
						{
							ID: "id-1",
						},
						{
							ID: "id-2",
						},
					},
					NextPageKey: "",
					PageSize:    2,
					TotalCount:  2,
				}

				marshal, _ := json.Marshal(events)
				w.Write(marshal)
			},
			want: []*keptnmodels.KeptnContextExtendedCE{
				{
					ID: "id-1",
				},
				{
					ID: "id-2",
				},
			},
			wantErr: false,
		},
		{
			name: "get all events from paginated source",
			args: args{
				endpoint: ts.URL,
				token:    "",
				topic:    "my-topic",
			},
			serverHandlerFunc: func(w http.ResponseWriter, request *http.Request) {
				w.Header().Add("Content-Type", "application/json")

				var events keptnmodels.Events
				if request.FormValue("nextPageKey") == "" {
					events = keptnmodels.Events{
						Events: []*keptnmodels.KeptnContextExtendedCE{
							{
								ID: "id-1",
							},
							{
								ID: "id-2",
							},
						},
						NextPageKey: "2",
						PageSize:    2,
						TotalCount:  4,
					}
				} else if request.FormValue("nextPageKey") == "2" {
					events = keptnmodels.Events{
						Events: []*keptnmodels.KeptnContextExtendedCE{
							{
								ID: "id-3",
							},
							{
								ID: "id-4",
							},
						},
						NextPageKey: "",
						PageSize:    2,
						TotalCount:  4,
					}
				}

				marshal, _ := json.Marshal(events)
				w.Write(marshal)
			},
			want: []*keptnmodels.KeptnContextExtendedCE{
				{
					ID: "id-1",
				},
				{
					ID: "id-2",
				},
				{
					ID: "id-3",
				},
				{
					ID: "id-4",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts.Config.Handler = tt.serverHandlerFunc
			got, err := getEventsFromEndpoint(tt.args.endpoint, tt.args.token, tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEventsFromEndpoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEventsFromEndpoint() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pollEventsForTopic(t *testing.T) {

	var eventSourceReturnedPayload keptnmodels.Events
	eventSourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		marshal, _ := json.Marshal(eventSourceReturnedPayload)
		w.Write(marshal)
	}))

	recipientServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	}))

	parsedURL, _ := url.Parse(recipientServer.URL)
	split := strings.Split(parsedURL.Host, ":")
	os.Setenv("PUBSUB_RECIPIENT", split[0])
	os.Setenv("PUBSUB_RECIPIENT_PORT", split[1])

	type args struct {
		endpoint string
		token    string
		topic    string
	}
	tests := []struct {
		name                       string
		args                       args
		eventSourceReturnedPayload keptnmodels.Events
	}{
		{
			name: "",
			args: args{
				endpoint: eventSourceServer.URL,
				token:    "",
				topic:    "my-topic",
			},
			eventSourceReturnedPayload: keptnmodels.Events{
				Events: []*keptnmodels.KeptnContextExtendedCE{
					{
						Contenttype:    "application/json",
						Data:           "",
						Extensions:     nil,
						ID:             "1234",
						Shkeptncontext: "1234",
						Source:         stringp("my-source"),
						Specversion:    "0.3",
						Time:           strfmt.DateTime{},
						Triggeredid:    "1234",
						Type:           stringp("my-topic"),
					},
				},
				NextPageKey: "",
				PageSize:    1,
				TotalCount:  1,
			},
		},
	}
	for _, tt := range tests {
		eventSourceReturnedPayload = tt.eventSourceReturnedPayload
		t.Run(tt.name, func(t *testing.T) {
			client := createRecipientConnection()
			pollEventsForTopic(tt.args.endpoint, tt.args.token, tt.args.topic, client)
		})
	}
}

const TEST_PORT = 8370
const TEST_TOPIC = "test-topic"

func RunServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *server.Options) *server.Server {
	return natsserver.RunServer(opts)
}
func Test__main(t *testing.T) {
	messageReceived := make(chan bool)
	// Mock http server
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			messageReceived <- true
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	natsServer := RunServerOnPort(TEST_PORT)
	defer natsServer.Shutdown()
	natsURL := fmt.Sprintf("nats://127.0.0.1:%d", TEST_PORT)

	hostAndPort := strings.Split(ts.URL, ":")
	os.Setenv("PUBSUB_RECIPIENT", strings.TrimPrefix(hostAndPort[1], "//"))
	os.Setenv("PUBSUB_RECIPIENT_PORT", hostAndPort[2])
	os.Setenv("PUBSUB_TOPIC", "test-topic")
	os.Setenv("PUBSUB_URL", natsURL)

	natsPublisher, _ := nats.Connect(natsURL)

	go _main(nil, envConfig{})

	<-time.After(2 * time.Second)

	_ = natsPublisher.Publish(TEST_TOPIC, []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "0.3",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`))

	select {
	case <-messageReceived:
		return
	case <-time.After(5 * time.Second):
		t.Error("SubscribeToTopics(): timed out waiting for messages")
	}

	close <- true

	receivedMessage := make(chan bool)
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		t.Errorf("Failed to process env var: %s", err)
	}

	_ = os.Setenv("PUBSUB_URL", natsURL)

	natsClient, _ := nats.Connect(natsURL)
	defer natsClient.Close()

	_, _ = natsClient.Subscribe("sh.keptn.events.deployment-finished", func(m *nats.Msg) {
		receivedMessage <- true
	})

	go main()

	<-time.After(2 * time.Second)
	_, err := http.Post("http://127.0.0.1:"+strconv.Itoa(env.Port), "application/cloudevents+json", bytes.NewBuffer([]byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "0.3",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`)))

	if err != nil {
		t.Errorf("Could not send event")
	}
	select {
	case <-receivedMessage:
		t.Logf("Received event!")
	case <-time.After(5 * time.Second):
		t.Errorf("Message did not make it to the receiver")
	}
}
