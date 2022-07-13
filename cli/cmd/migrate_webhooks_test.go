package cmd

import (
	"testing"

	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/stretchr/testify/require"
)

func TestMigrateAlphaRequest(t *testing.T) {
	tests := []struct {
		input  string
		output *lib.Request
	}{
		{
			input: "curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" -H 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY -X POST",
			output: &lib.Request{
				URL:    "https://httpbin.org/post",
				Method: "POST",
				Headers: []lib.Header{
					{
						Key:   "Accept-Charset",
						Value: " utf-8",
					},
					{
						Key:   "Content-Type",
						Value: " application/json",
					},
				},
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "--some-random-options -YYY",
			},
		},
		{
			input: "curl -d '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" --header 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY --request POST",
			output: &lib.Request{
				URL:    "https://httpbin.org/post",
				Method: "POST",
				Headers: []lib.Header{
					{
						Key:   "Accept-Charset",
						Value: " utf-8",
					},
					{
						Key:   "Content-Type",
						Value: " application/json",
					},
				},
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "--some-random-options -YYY",
			},
		},
		{
			input: "curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" -H 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY",
			output: &lib.Request{
				URL:    "https://httpbin.org/post",
				Method: "GET",
				Headers: []lib.Header{
					{
						Key:   "Accept-Charset",
						Value: " utf-8",
					},
					{
						Key:   "Content-Type",
						Value: " application/json",
					},
				},
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "--some-random-options -YYY",
			},
		},
		{
			input: "curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' https://httpbin.org/post --some-random-options -YYY",
			output: &lib.Request{
				URL:     "https://httpbin.org/post",
				Method:  "GET",
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "--some-random-options -YYY",
			},
		},
		{
			input: "curl https://httpbin.org/post --some-random-options -YYY",
			output: &lib.Request{
				URL:     "https://httpbin.org/post",
				Method:  "GET",
				Options: "--some-random-options -YYY",
			},
		},
		{
			input: "curl https://httpbin.org/post",
			output: &lib.Request{
				URL:    "https://httpbin.org/post",
				Method: "GET",
			},
		},
		{
			input: "curl -d '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" --header 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY --request",
			output: &lib.Request{
				URL:    "https://httpbin.org/post",
				Method: "GET",
				Headers: []lib.Header{
					{
						Key:   "Accept-Charset",
						Value: " utf-8",
					},
					{
						Key:   "Content-Type",
						Value: " application/json",
					},
				},
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "--some-random-options -YYY --request",
			},
		},
		{
			input: "curl -d '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" --header 'Content-Type: application/json' httttps://httpbin.org/post --some-random-options -YYY --request",
			output: &lib.Request{
				URL:    "",
				Method: "GET",
				Headers: []lib.Header{
					{
						Key:   "Accept-Charset",
						Value: " utf-8",
					},
					{
						Key:   "Content-Type",
						Value: " application/json",
					},
				},
				Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
				Options: "httttps://httpbin.org/post --some-random-options -YYY --request",
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, _ := migrateAlphaRequest(tt.input)
			require.Equal(t, tt.output, res)

		})
	}

}

func TestMigrateAlphaWebhook(t *testing.T) {
	tests := []struct {
		input  *lib.WebHookConfig
		output *lib.WebHookConfig
	}{
		{
			input: &lib.WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1alpha1",
				Kind:       "WebhookConfig",
				Metadata: lib.Metadata{
					Name: "webhook-configuration",
				},
				Spec: lib.WebHookConfigSpec{
					Webhooks: []lib.Webhook{
						{
							Type:           "sh.keptn.event.webhook.triggered",
							SubscriptionID: "my-subscription-id",
							EnvFrom: []lib.EnvFrom{
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
			output: &lib.WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1beta1",
				Kind:       "WebhookConfig",
				Metadata: lib.Metadata{
					Name: "webhook-configuration",
				},
				Spec: lib.WebHookConfigSpec{
					Webhooks: []lib.Webhook{
						{
							Type:           "sh.keptn.event.webhook.triggered",
							SubscriptionID: "my-subscription-id",
							EnvFrom: []lib.EnvFrom{
								{
									Name: "mysecret",
								},
							},
							Requests: []interface{}{
								lib.Request{
									Method:  "GET",
									URL:     "http://localhost:8080",
									Options: "{{.data.project}} {{.env.mysecret}}",
								},
							},
						},
					},
				},
			},
		},
		{
			input: &lib.WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1alpha1",
				Kind:       "WebhookConfig",
				Metadata: lib.Metadata{
					Name: "webhook-configuration",
				},
				Spec: lib.WebHookConfigSpec{
					Webhooks: []lib.Webhook{
						{
							Type:           "sh.keptn.event.webhook.triggered",
							SubscriptionID: "my-subscription-id",
							EnvFrom: []lib.EnvFrom{
								{
									Name: "mysecret",
								},
							},
							Requests: []interface{}{
								"curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" -H 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY -X POST",
							},
						},
					},
				},
			},
			output: &lib.WebHookConfig{
				ApiVersion: "webhookconfig.keptn.sh/v1beta1",
				Kind:       "WebhookConfig",
				Metadata: lib.Metadata{
					Name: "webhook-configuration",
				},
				Spec: lib.WebHookConfigSpec{
					Webhooks: []lib.Webhook{
						{
							Type:           "sh.keptn.event.webhook.triggered",
							SubscriptionID: "my-subscription-id",
							EnvFrom: []lib.EnvFrom{
								{
									Name: "mysecret",
								},
							},
							Requests: []interface{}{
								lib.Request{
									URL:    "https://httpbin.org/post",
									Method: "POST",
									Headers: []lib.Header{
										{
											Key:   "Accept-Charset",
											Value: " utf-8",
										},
										{
											Key:   "Content-Type",
											Value: " application/json",
										},
									},
									Payload: "{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}",
									Options: "--some-random-options -YYY",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			res, _ := migrateAlphaWebhook(tt.input)
			require.Equal(t, tt.output, res)

		})
	}
}
