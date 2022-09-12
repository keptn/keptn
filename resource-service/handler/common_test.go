package handler

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/resource-service/common_models"
)

func Test_getAuthMethod(t *testing.T) {
	tests := []struct {
		name           string
		gitCredentials *common_models.GitCredentials
		wantErr        bool
		expectedOutput transport.AuthMethod
	}{
		{
			name:           "no credentials",
			gitCredentials: &common_models.GitCredentials{},
			wantErr:        false,
			expectedOutput: nil,
		},
		{
			name: "valid credentials",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "https://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "some-token",
					InsecureSkipTLS: false,
				},
				User: "user",
			},
			wantErr: false,
			expectedOutput: &http.BasicAuth{
				Username: "user",
				Password: "some-token",
			},
		},
		{
			name: "valid credentials proxy",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "https://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "some-token",
					InsecureSkipTLS: false,
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "proxy-url",
						Scheme:   "http",
						User:     "proxy-user",
						Password: "proxy-password",
					},
				},
				User: "user",
			},
			wantErr: false,
			expectedOutput: &http.BasicAuth{
				Username: "user",
				Password: "some-token",
			},
		},
		{
			name: "valid credentials no user",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "https://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "some-token",
					InsecureSkipTLS: false,
				},
				User: "",
			},
			wantErr: false,
			expectedOutput: &http.BasicAuth{
				Username: "keptnuser",
				Password: "some-token",
			},
		},
		{
			name: "invalid credentials",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "https://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					InsecureSkipTLS: false,
				},
				User: "user",
			},
			wantErr:        false,
			expectedOutput: nil,
		},
		{
			name: "invalid ssh credentials",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "ssh://some.url",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "private-key",
					PrivateKeyPass: "password",
				},
				User: "user",
			},
			wantErr:        true,
			expectedOutput: nil,
		},
		{
			name: "dumb credentials",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "ssh://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "some",
					InsecureSkipTLS: false,
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "",
						Scheme:   "",
						User:     "hate",
						Password: "",
					},
				},
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "",
					PrivateKeyPass: "password",
				},
				User: "user",
			},
			wantErr:        true,
			expectedOutput: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := getAuthMethod(tt.gitCredentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && auth != tt.expectedOutput {
				t.Errorf("getAuthMethod() auth = %v, expectedOutput %v", err, tt.wantErr)
			}
		})
	}
}
