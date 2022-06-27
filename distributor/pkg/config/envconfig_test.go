package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

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
			env := EnvConfig{
				KeptnAPIEndpoint: tt.externalEndpoint,
			}
			scheme, host, path := env.ProxyHost(tt.args.path)

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
			env := EnvConfig{
				KeptnAPIEndpoint: tt.apiEndpointEnvVar,
			}
			if got := env.HTTPPollingEndpoint(); got != tt.want {
				t.Errorf("PollingEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

			env := EnvConfig{}
			_ = envconfig.Process("", &env)

			got := env.PubSubRecipientURL()
			if got != tt.want {
				t.Errorf("getPubSubRecipientURL() got = %v, want1 %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateKeptnAPIEndpointURL(t *testing.T) {
	// valid
	config := EnvConfig{KeptnAPIEndpoint: "http:1.2.3.4.nip.io/some-path"}
	assert.Nil(t, config.ValidateKeptnAPIEndpointURL())
	// not valid
	config = EnvConfig{KeptnAPIEndpoint: "d"}
	assert.NotNil(t, config.ValidateKeptnAPIEndpointURL())
	// not given
	config = EnvConfig{KeptnAPIEndpoint: ""}
	assert.Nil(t, config.ValidateKeptnAPIEndpointURL())
}

func Test_GetPubSubTopics(t *testing.T) {
	// multiple topics
	config := EnvConfig{PubSubTopic: "a,b,c"}
	assert.Equal(t, 3, len(config.PubSubTopics()))

	// zero topics
	config = EnvConfig{}
	assert.Equal(t, 0, len(config.PubSubTopics()))
}

func Test_OAuthEnabled(t *testing.T) {
	tests := []struct {
		input EnvConfig
		want  bool
	}{
		{
			input: EnvConfig{},
			want:  false,
		},
		{
			input: EnvConfig{
				OAuthClientID: "sso-client-id",
			},
			want: false,
		},
		{
			input: EnvConfig{
				OAuthClientID:     "sso-client-id",
				OAuthClientSecret: "sso-client-secret",
				OAuthScopes:       nil,
				OAuthDiscovery:    "",
				OauthTokenURL:     "",
			},
			want: false,
		},
		{
			input: EnvConfig{
				OAuthClientID:     "sso-client-id",
				OAuthClientSecret: "sso-client-secret",
				OAuthScopes:       []string{"scope"},
				OAuthDiscovery:    "",
				OauthTokenURL:     "",
			},
			want: false,
		},
		{
			input: EnvConfig{
				OAuthClientID:     "sso-client-id",
				OAuthClientSecret: "sso-client-secret",
				OAuthScopes:       []string{"scope"},
				OAuthDiscovery:    "http://some-url.com",
				OauthTokenURL:     "",
			},
			want: true,
		},
		{
			input: EnvConfig{
				OAuthClientID:     "sso-client-id",
				OAuthClientSecret: "sso-client-secret",
				OAuthScopes:       []string{"scope"},
				OAuthDiscovery:    "",
				OauthTokenURL:     "http://some-url.com",
			},
			want: true,
		},
		{
			input: EnvConfig{
				OAuthClientID:     "sso-client-id",
				OAuthClientSecret: "sso-client-secret",
				OAuthScopes:       []string{"scope"},
				OAuthDiscovery:    "http://some-url.com",
				OauthTokenURL:     "http://some-url.com",
			},
			want: true,
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.want, tc.input.OAuthEnabled())
	}
}

func TestEnvConfig_GetAPIProxyHTTPTimeout(t *testing.T) {
	type fields struct {
		APIProxyHTTPTimeout string
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name:   "Get default timeout",
			fields: fields{},
			want:   30 * time.Second,
		},
		{
			name: "Get configured timeout",
			fields: fields{
				APIProxyHTTPTimeout: "5",
			},
			want: 5 * time.Second,
		},
		{
			name: "Get default timeout if invalid value",
			fields: fields{
				APIProxyHTTPTimeout: "invalid",
			},
			want: 30 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := &EnvConfig{
				APIProxyHTTPTimeout: tt.fields.APIProxyHTTPTimeout,
			}
			if got := env.GetAPIProxyHTTPTimeout(); got != tt.want {
				t.Errorf("GetAPIProxyHTTPTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvConfig_GetAPIProxyMaxBytes(t *testing.T) {
	type fields struct {
		APIProxyMaxBytesKB int
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "64KB",
			fields: fields{
				APIProxyMaxBytesKB: 64,
			},
			want: 65536,
		},
		{
			name: "no limit",
			fields: fields{
				APIProxyMaxBytesKB: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := &EnvConfig{
				APIProxyMaxPayloadBytesKB: tt.fields.APIProxyMaxBytesKB,
			}
			assert.Equalf(t, tt.want, env.GetAPIProxyMaxBytes(), "GetAPIProxyMaxBytes()")
		})
	}
}
