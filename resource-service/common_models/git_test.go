package common_models

import "testing"

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
				HttpsAuth: &HttpsGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
					Token: "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid URI",
			gitCredentials: GitCredentials{
				User:      "my-user",
				RemoteURL: "https://my:repo",
				HttpsAuth: &HttpsGitAuth{
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
				SshAuth: &SshGitAuth{
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
				SshAuth: &SshGitAuth{
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
				SshAuth: &SshGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
					Token: "token",
					Proxy: &ProxyGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
					Token: "token",
					Proxy: &ProxyGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
					Token: "token",
					Proxy: &ProxyGitAuth{
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
				HttpsAuth: &HttpsGitAuth{
					Token: "token",
					Proxy: &ProxyGitAuth{
						URL:    "1.1.1.1",
						Scheme: "https",
					},
				},
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
