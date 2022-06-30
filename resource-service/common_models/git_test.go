package common_models

import (
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
)

func TestGitCredentials_Validate(t *testing.T) {
	tests := []struct {
		name           string
		gitCredentials GitCredentials
		wantErr        bool
	}{
		{
			name: "valid credentials",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "my-token",
				},
			},
			wantErr: false,
		},
		{
			name: "empty token",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid URI",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my:repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "my-token",
				},
			},
			wantErr: true,
		},
		{
			name: "empty PrivateKey",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "ssh://my-repo",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey: "",
				},
			},
			wantErr: true,
		},
		{
			name: "valid PrivateKey",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "ssh://my-repo",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey: "privatekey",
				},
			},
			wantErr: false,
		},
		{
			name: "PrivateKey with https",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey: "",
				},
			},
			wantErr: true,
		},
		{
			name: "token with ssh",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "ssh://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "",
				},
			},
			wantErr: true,
		},
		{
			name: "http proxy",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "1.1.1.1:12",
						Scheme: "http",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "https proxy",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "1.1.1.1:12",
						Scheme: "https",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "proxy invalid scheme",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "1.1.1.1:12",
						Scheme: "fddd",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "proxy URL without port",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my-repo",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "token",
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "1.1.1.1",
						Scheme: "https",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid credentials",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "httpg://my-repo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.gitCredentials.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
