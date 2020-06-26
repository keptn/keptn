package main

import (
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
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
			name: "Get V0.2 CloudEvent",
			args: args{
				data: []byte(`{
				"data": "",
				"id": "6de83495-4f83-481c-8dbe-fcceb2e0243b",
				"source": "helm-service",
				"specversion": "0.2",
				"type": "sh.keptn.events.deployment-finished",
				"shkeptncontext": "3c9ffbbb-6e1d-4789-9fee-6e63b4bcc1fb"
			}`),
			},
			want: &cloudevents.Event{
				Context: &cloudevents.EventContextV02{
					SpecVersion: "0.2",
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
			name: "Get V0.2 CloudEvent",
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
				"specversion": "0.2",
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
}
