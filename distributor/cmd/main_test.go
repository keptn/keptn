package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-openapi/strfmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"

	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
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
			want:    "http://127.0.0.1:666",
			wantErr: true,
		},
		{
			name: "HTTPS recipient",
			args: args{
				recipientService: "https://lighthouse-service",
				port:             "",
				path:             "",
			},
			want: "https://lighthouse-service:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.recipientService != "" {
				os.Setenv("PUBSUB_RECIPIENT", tt.args.recipientService)
			} else {
				os.Unsetenv("PUBSUB_RECIPIENT")
			}
			if tt.args.port != "" {
				os.Setenv("PUBSUB_RECIPIENT_PORT", tt.args.port)
			} else {
				os.Unsetenv("PUBSUB_RECIPIENT_PORT")
			}
			if tt.args.path != "" {
				os.Setenv("PUBSUB_RECIPIENT_PATH", tt.args.path)
			} else {
				os.Unsetenv("PUBSUB_RECIPIENT_PATH")
			}

			env = envConfig{}
			_ = envconfig.Process("", &env)

			got := getPubSubRecipientURL()
			if got != tt.want {
				t.Errorf("getPubSubRecipientURL() got = %v, want1 %v", got, tt.want)
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
			name: "Get V1.0 CloudEvent",
			args: args{
				data: []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`),
			},
			want:    getExpectedCloudEvent(),
			wantErr: false,
		},
		{
			name: "Get V1.0 CloudEvent",
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
					t.Errorf("decodeCloudEvent() specVersion: got = %v, want1 %v", got.Context.GetSpecVersion(), tt.want.Context.GetSpecVersion())
				}
				if !reflect.DeepEqual(got.Context.GetType(), tt.want.Context.GetType()) {
					t.Errorf("decodeCloudEvent() type: got = %v, want1 %v", got.Context.GetType(), tt.want.Context.GetType())
				}
			}
		})
	}
}

func getExpectedCloudEvent() *cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("helm-service")
	event.SetType("sh.keptn.events.deployment-finished")
	event.SetID("6de83495-4f83-481c-8dbe-fcceb2e0243b")
	event.SetExtension("shkeptncontext", "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb")
	event.SetData(cloudevents.TextPlain, `""`)
	return &event
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
			name: "want1 true",
			args: args{
				sentEvents: []string{"sent-1", "sent-2"},
				eventID:    "sent-1",
			},
			want: true,
		},
		{
			name: "want1 false",
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
				t.Errorf("hasEventBeenSent() = %v, want1 %v", got, tt.want)
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
				t.Errorf("getEventsFromEndpoint() got = %v, want1 %v", got, tt.want)
			}
		})
	}
}

// Test_pollAndForwardEventsForTopic tests the polling and forwarding mechanism (in combination of ceCache)
func Test_pollAndForwardEventsForTopic(t *testing.T) {

	var eventSourceReturnedPayload keptnmodels.Events
	var recipientSleepTimeSeconds int

	// store number of received CloudEvents for the recipient server
	var recipientReceivedCloudEvents int

	// mock the server where we poll CloudEvents from
	eventSourceServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		marshal, _ := json.Marshal(eventSourceReturnedPayload)
		w.Write(marshal)
	}))

	// mock the recipient server where CloudEvents are sent to
	recipientServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Second * time.Duration(recipientSleepTimeSeconds))
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		recipientReceivedCloudEvents += 1
	}))

	parsedURL, _ := url.Parse(recipientServer.URL)
	split := strings.Split(parsedURL.Host, ":")
	os.Setenv("PUBSUB_RECIPIENT", split[0])
	os.Setenv("PUBSUB_RECIPIENT_PORT", split[1])

	env.PubSubRecipient = split[0]
	env.PubSubRecipientPort = split[1]

	// define CloudEvents that are provided by the polling mechanism
	cloudEventsToSend := []*keptnmodels.KeptnContextExtendedCE{
		{
			Contenttype:    "application/json",
			Data:           "",
			Extensions:     nil,
			ID:             "1234",
			Shkeptncontext: "1234",
			Source:         stringp("my-source"),
			Specversion:    "1.0",
			Time:           strfmt.DateTime{},
			Triggeredid:    "1234",
			Type:           stringp("my-topic"),
		},
		{
			Contenttype:    "application/json",
			Data:           "",
			Extensions:     nil,
			ID:             "3456",
			Shkeptncontext: "1234",
			Source:         stringp("my-source"),
			Specversion:    "1.0",
			Time:           strfmt.DateTime{},
			Triggeredid:    "1234",
			Type:           stringp("my-topic"),
		},
		{
			Contenttype:    "application/json",
			Data:           "",
			Extensions:     nil,
			ID:             "7890",
			Shkeptncontext: "1234",
			Source:         stringp("my-source"),
			Specversion:    "1.0",
			Time:           strfmt.DateTime{},
			Triggeredid:    "1234",
			Type:           stringp("my-topic"),
		},
	}

	type args struct {
		endpoint string
		token    string
		topic    string
	}
	tests := []struct {
		name                       string
		args                       args
		eventSourceReturnedPayload keptnmodels.Events
		recipientSleepTimeSeconds  int
	}{
		{
			name: "",
			args: args{
				endpoint: eventSourceServer.URL,
				token:    "",
				topic:    "my-topic",
			},
			eventSourceReturnedPayload: keptnmodels.Events{
				// incoming events (topic: my-topic)
				Events:      cloudEventsToSend,
				NextPageKey: "",
				PageSize:    3,
				TotalCount:  3,
			},
			recipientSleepTimeSeconds: 2,
		},
	}
	for _, tt := range tests {
		eventSourceReturnedPayload = tt.eventSourceReturnedPayload
		recipientSleepTimeSeconds = tt.recipientSleepTimeSeconds
		recipientReceivedCloudEvents = 0
		t.Run(tt.name, func(t *testing.T) {
			setupCEClient()
			// poll events
			pollEventsForTopic(tt.args.endpoint, tt.args.token, tt.args.topic)

			// assert that the events above are present in ceCache
			assert.True(t, ceCache.Contains("my-topic", "1234"), "Event with ID 1234 not in ceCache")
			assert.True(t, ceCache.Contains("my-topic", "3456"), "Event with ID 3456 not in ceCache")
			assert.True(t, ceCache.Contains("my-topic", "7890"), "Event with ID 7890 not in ceCache")

			// assert that the correct number of events is in ceCache
			assert.Equal(t, ceCache.Length("my-topic"), 3)

			// however, due to recipientSleepTimeSeconds no events should be received by the recipient yet
			assert.Equal(t, recipientReceivedCloudEvents, 0, "The recipient should not have received any CloudEvents")

			// poll again
			pollEventsForTopic(tt.args.endpoint, tt.args.token, tt.args.topic)

			// verify that there is still only 3 events in ceCache
			assert.Equal(t, ceCache.Length("my-topic"), 3)

			// and there still should be no events received by the recipient yet
			assert.Equal(t, recipientReceivedCloudEvents, 0, "The recipient should not have received any CloudEvents")

			// Okay, now we have to wait a little bit, until the recipient service has processed everything
			time.Sleep(time.Second * 1)

			// verify that recipientServer has processed 3 CloudEvents eventually
			assert.Eventually(t, func() bool {
				if recipientReceivedCloudEvents == 3 {
					return true
				}
				return false
			}, time.Second*time.Duration(tt.recipientSleepTimeSeconds), 100*time.Millisecond)

			// wait a little bit longer, and verify that it is still only 3 CloudEvents
			time.Sleep(time.Second * time.Duration(tt.recipientSleepTimeSeconds) * 2)
			assert.Equal(t, 3, recipientReceivedCloudEvents)

			// verify that there is still only 3 events in ceCache
			assert.Equal(t, ceCache.Length("my-topic"), 3)
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
	env = envConfig{}
	if err := envconfig.Process("", &env); err != nil {
		t.Errorf("Failed to process env var: %s", err)
	}
	env.APIProxyPort = TEST_PORT + 1
	go _main(env)

	<-time.After(2 * time.Second)

	_ = natsPublisher.Publish(TEST_TOPIC, []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`))

	select {
	case <-messageReceived:
		t.Logf("Received event!")
	case <-time.After(5 * time.Second):
		t.Error("SubscribeToTopics(): timed out waiting for messages")
	}

	receivedMessage := make(chan bool)

	_ = os.Setenv("PUBSUB_URL", natsURL)

	natsClient, err := nats.Connect(natsURL)
	if err != nil {
		t.Errorf("could not initialize nats client: %s", err.Error())
	}
	defer natsClient.Close()

	_, _ = natsClient.Subscribe("sh.keptn.events.deployment-finished", func(m *nats.Msg) {
		receivedMessage <- true
	})

	<-time.After(2 * time.Second)
	_, err = http.Post("http://127.0.0.1:"+strconv.Itoa(env.APIProxyPort)+"/event", "application/cloudevents+json", bytes.NewBuffer([]byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "1.0",
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

	_, err = http.Post("http://127.0.0.1:"+strconv.Itoa(env.APIProxyPort)+env.APIProxyPath+"/datastore?foo=bar", "application/json", bytes.NewBuffer([]byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "1.0",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`)))
	if err != nil {
		t.Errorf("Could not handle API request")
	}

	closeChan <- true
}

func Test_getProxyRequestURL(t *testing.T) {
	type args struct {
		endpoint string
		path     string
	}
	tests := []struct {
		name             string
		args             args
		wantScheme       string
		wantHost         string
		wantPath         string
		externalEndpoint string
	}{
		{
			name: "Get internal Datastore",
			args: args{
				endpoint: "",
				path:     "/mongodb-datastore/event/type/sh.keptn.event.evaluation.finished",
			},
			wantScheme: "http",
			wantHost:   "mongodb-datastore:8080",
			wantPath:   "event/type/sh.keptn.event.evaluation.finished",
		},
		{
			name: "Get internal configuration service",
			args: args{
				endpoint: "",
				path:     "/configuration-service",
			},
			wantScheme: "http",
			wantHost:   "configuration-service:8080",
		},
		{
			name: "Get configuration service",
			args: args{
				endpoint: "",
				path:     "/configuration-service",
			},
			wantScheme: "http",
			wantHost:   "configuration-service:8080",
		},
		{
			name: "Get configuration service via public API",
			args: args{
				endpoint: "",
				path:     "/configuration-service",
			},
			wantScheme:       "http",
			wantHost:         "external-api.com",
			wantPath:         "/api/configuration-service/",
			externalEndpoint: "http://external-api.com/api",
		},
		{
			name: "Get configuration service via public API with API prefix",
			args: args{
				endpoint: "",
				path:     "/configuration-service",
			},
			wantScheme:       "http",
			wantHost:         "external-api.com",
			wantPath:         "/my/path/prefix/api/configuration-service/",
			externalEndpoint: "http://external-api.com/my/path/prefix/api",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.KeptnAPIEndpoint = tt.externalEndpoint
			scheme, host, path := getProxyHost(tt.args.path)

			if scheme != tt.wantScheme {
				t.Errorf("getProxyHost(); host = %v, want %v", scheme, tt.wantScheme)
			}

			if host != tt.wantHost {
				t.Errorf("getProxyHost(); path = %v, want %v", host, tt.wantHost)
			}

			if path != tt.wantPath {
				t.Errorf("getProxyHost(); path = %v, want %v", path, tt.wantPath)
			}
		})
	}
}

func Test_getHTTPPollingEndpoint(t *testing.T) {
	tests := []struct {
		name              string
		apiEndpointEnvVar string
		want              string
	}{
		{
			name:              "get internal endpoint",
			apiEndpointEnvVar: "",
			want:              "http://shipyard-controller:8080/v1/event/triggered",
		},
		{
			name:              "get external endpoint",
			apiEndpointEnvVar: "https://my-keptn.com/api",
			want:              "https://my-keptn.com/api/controlPlane/v1/event/triggered",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.KeptnAPIEndpoint = tt.apiEndpointEnvVar
			if got := getHTTPPollingEndpoint(); got != tt.want {
				t.Errorf("getHTTPPollingEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getCloudEventWithEventData(eventData keptnv2.EventData) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("helm-service")
	event.SetType("sh.keptn.events.deployment-finished")
	event.SetID("6de83495-4f83-481c-8dbe-fcceb2e0243b")
	event.SetExtension("shkeptncontext", "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb")
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetData(cloudevents.ApplicationJSON, eventData)
	return event
}

func Test_matchesFilter(t *testing.T) {
	type args struct {
		e cloudevents.Event
	}
	tests := []struct {
		name          string
		args          args
		projectFilter string
		stageFilter   string
		serviceFilter string
		want          bool
	}{
		{
			name: "no filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "",
			serviceFilter: "",
			want:          true,
		},
		{
			name: "project filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-project",
			stageFilter:   "",
			serviceFilter: "",
			want:          true,
		},
		{
			name: "project filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-project,my-project-2,my-project-3",
			stageFilter:   "",
			serviceFilter: "",
			want:          true,
		},
		{
			name: "project filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-other-project",
			stageFilter:   "",
			serviceFilter: "",
			want:          false,
		},
		{
			name: "project filter (comma-separated list) - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-other-project,my-second-project",
			stageFilter:   "",
			serviceFilter: "",
			want:          false,
		},
		{
			name: "stage filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "my-stage",
			serviceFilter: "",
			want:          true,
		},
		{
			name: "stage filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "my-first-stage,my-stage",
			serviceFilter: "",
			want:          true,
		},
		{
			name: "stage filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "my-other-stage",
			serviceFilter: "",
			want:          false,
		},
		{
			name: "service filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "",
			serviceFilter: "my-service",
			want:          true,
		},
		{
			name: "service filter (comma-separated list) - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "",
			serviceFilter: "my-other-service,my-service",
			want:          true,
		},
		{
			name: "service filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "",
			stageFilter:   "",
			serviceFilter: "my-other-service",
			want:          false,
		},
		{
			name: "combined filter - should match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-project",
			stageFilter:   "my-stage",
			serviceFilter: "my-service",
			want:          true,
		},
		{
			name: "combined filter - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-other-project",
			stageFilter:   "my-stage",
			serviceFilter: "my-service",
			want:          false,
		},
		{
			name: "combined filter (comma-separated list) - should not match",
			args: args{
				getCloudEventWithEventData(keptnv2.EventData{
					Project: "project",
					Stage:   "my-stage",
					Service: "my-service",
				}),
			},
			projectFilter: "my-project,project-1",
			stageFilter:   "my-stage",
			serviceFilter: "my-service",
			want:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env.ProjectFilter = tt.projectFilter
			env.StageFilter = tt.stageFilter
			env.ServiceFilter = tt.serviceFilter
			if got := matchesFilter(tt.args.e); got != tt.want {
				t.Errorf("matchesFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}
