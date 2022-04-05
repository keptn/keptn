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
							Requests: []interface{}{
								"curl http://localhost:8080 {{.data.project}} {{.env.mysecret}}",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Beta1 version input",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080
          method: POST`),
			},
			want: &WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1beta1",
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
							Requests: []interface{}{
								Request{
									Method: "POST",
									URL:    "http://localhost:8080",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid Beta1 version input - full",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080
          method: POST
          payload: "some payload"
          options: "some options"
          headers:
            - value: value
              key: key`),
			},
			want: &WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1beta1",
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
							Requests: []interface{}{
								Request{
									Headers: []Header{
										{
											Key:   "key",
											Value: "value",
										},
									},
									Method:  "POST",
									Options: "some options",
									Payload: "some payload",
									URL:     "http://localhost:8080",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Beta1 version input - missing method",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Beta1 version input - invalid method",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080
          method: DELETE`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Beta1 version input - missing url",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - method: POST`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Beta1 version input - alpha version requests",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
			name: "Beta1 version input - missing headers value",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080
          method: POST
		  headers:
            - key: key`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Beta1 version input - missing headers key",
			args: args{
				webhookConfigYaml: []byte(`apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://localhost:8080
          method: POST
		  headers:
            - value: value`),
			},
			want:    nil,
			wantErr: true,
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
		Requests       []interface{}
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
		Requests       []interface{}
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
