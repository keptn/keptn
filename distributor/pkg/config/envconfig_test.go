package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
			scheme, host, path := env.GetProxyHost(tt.args.path)

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
			if got := env.GetHTTPPollingEndpoint(); got != tt.want {
				t.Errorf("GetHTTPPollingEndpoint() = %v, want %v", got, tt.want)
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

			got := env.GetPubSubRecipientURL()
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
