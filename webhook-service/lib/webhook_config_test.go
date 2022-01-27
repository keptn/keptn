package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeWebHookConfigYAML(t *testing.T) {
	type args struct {
		webhookConfigYaml []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *WebHookConfig
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionID: "my-subscription-id"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`),
			},
			want: &WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1alpha1",
				Kind:       "WebhookConfig",
				Metadata: Metadata{
					Name: "webhook-configuration",
				},
				Spec: WebHookConfigSpec{
					Webhooks: []Webhook{
						{
							Type:           "sh.keptn.event.webhook.triggered",
							SubscriptionID: "my-subscription-id",
							EnvFrom: []EnvFrom{
								{
									Name: "mysecret",
								},
							},
							Requests: []string{
								"curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid input",
			args: args{
				webhookConfigYaml: []byte("hulumulu"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "bad padding invalid input",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1alpha1
						kind: WebhookConfig
						metadata:
						name: webhook-configuration
						spec:
						webhooks:
							- type: "sh.keptn.event.webhook.triggered"
							subscriptionID: "my-subscription-id"
							envFrom:
								- secretRef:
								name: mysecret
							requests:
								- "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "misspeling keyworkds invalid input",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionIDs: "my-subscription-id"
      envFrom:
        - secretRef:
          name: mysecret
      requests:
        - "curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}"`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing requests invalid input",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.webhook.triggered"
      subscriptionIDs: "my-subscription-id"
      envFrom:
        - secretRef:
          name: mysecret
      requests:`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeWebHookConfigYAML(tt.args.webhookConfigYaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeWebHookConfigYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestWebhook_ShouldSendStartedEvent(t *testing.T) {
	doSendStarted := true
	doNotSendStarted := false
	type fields struct {
		Type           string
		SubscriptionID string
		SendFinished   bool
		SendStarted    *bool
		EnvFrom        []EnvFrom
		Requests       []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "return default value",
			fields: fields{
				SendStarted: nil,
			},
			want: true,
		},
		{
			name: "return set value",
			fields: fields{
				SendStarted: &doSendStarted,
			},
			want: true,
		},
		{
			name: "return set value",
			fields: fields{
				SendStarted: &doNotSendStarted,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := Webhook{
				Type:           tt.fields.Type,
				SubscriptionID: tt.fields.SubscriptionID,
				SendFinished:   tt.fields.SendFinished,
				SendStarted:    tt.fields.SendStarted,
				EnvFrom:        tt.fields.EnvFrom,
				Requests:       tt.fields.Requests,
			}
			if got := wh.ShouldSendStartedEvent(); got != tt.want {
				t.Errorf("ShouldSendStartedEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhook_ShouldSendFinishedEvent(t *testing.T) {
	type fields struct {
		Type           string
		SubscriptionID string
		SendFinished   bool
		SendStarted    *bool
		EnvFrom        []EnvFrom
		Requests       []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "false",
			fields: fields{},
			want:   false,
		},
		{
			name: "true",
			fields: fields{
				SendFinished: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wh := Webhook{
				Type:           tt.fields.Type,
				SubscriptionID: tt.fields.SubscriptionID,
				SendFinished:   tt.fields.SendFinished,
				SendStarted:    tt.fields.SendStarted,
				EnvFrom:        tt.fields.EnvFrom,
				Requests:       tt.fields.Requests,
			}
			if got := wh.ShouldSendFinishedEvent(); got != tt.want {
				t.Errorf("ShouldSendFinishedEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
