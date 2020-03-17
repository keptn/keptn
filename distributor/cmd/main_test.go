package main

import (
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"net/url"
	"reflect"
	"testing"
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
