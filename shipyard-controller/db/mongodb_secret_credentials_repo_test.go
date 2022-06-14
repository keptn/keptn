package db

import (
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

func TestMongoDBSecretCredentialsRepo_Transform(t *testing.T) {
	tests := []struct {
		name string
		in   *GitOldCredentials
		out  *apimodels.GitAuthCredentials
	}{
		{
			name: "project with new format credentials",
			in: &GitOldCredentials{
				RemoteURI: "",
			},
			out: nil,
		},
		{
			name: "project with ssh credentials",
			in: &GitOldCredentials{
				RemoteURI:         "ssh://some-url",
				User:              "user",
				GitPrivateKey:     "key",
				GitPrivateKeyPass: "pass",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "ssh://some-url",
				User:      "user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "key",
					PrivateKeyPass: "pass",
				},
			},
		},
		{
			name: "project with ssh credentials",
			in: &GitOldCredentials{
				RemoteURI:         "ssh://some-url",
				GitPrivateKey:     "key",
				GitPrivateKeyPass: "pass",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "ssh://some-url",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey:     "key",
					PrivateKeyPass: "pass",
				},
			},
		},
		{
			name: "project with ssh credentials",
			in: &GitOldCredentials{
				RemoteURI:     "ssh://some-url",
				User:          "user",
				GitPrivateKey: "key",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "ssh://some-url",
				User:      "user",
				SshAuth: &apimodels.SshGitAuth{
					PrivateKey: "key",
				},
			},
		},
		{
			name: "project with https credentials",
			in: &GitOldCredentials{
				RemoteURI:       "https://some-url",
				User:            "user",
				Token:           "token",
				InsecureSkipTLS: false,
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "token",
					InsecureSkipTLS: false,
				},
			},
		},
		{
			name: "project with https credentials - certificate",
			in: &GitOldCredentials{
				RemoteURI:         "https://some-url",
				User:              "user",
				Token:             "token",
				GitPemCertificate: "cert",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:       "token",
					Certificate: "cert",
				},
			},
		},
		{
			name: "project with https credentials",
			in: &GitOldCredentials{
				RemoteURI: "https://some-url",
				User:      "user",
				Token:     "token",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token: "token",
				},
			},
		},
		{
			name: "project with https credentials - proxy",
			in: &GitOldCredentials{
				RemoteURI:        "https://some-url",
				User:             "user",
				Token:            "token",
				InsecureSkipTLS:  false,
				GitProxyURL:      "url",
				GitProxyScheme:   "http",
				GitProxyUser:     "user",
				GitProxyPassword: "pass",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "token",
					InsecureSkipTLS: false,
					Proxy: &apimodels.ProxyGitAuth{
						URL:      "url",
						Scheme:   "http",
						User:     "user",
						Password: "pass",
					},
				},
			},
		},
		{
			name: "project with https credentials - proxy",
			in: &GitOldCredentials{
				RemoteURI:       "https://some-url",
				User:            "user",
				Token:           "token",
				InsecureSkipTLS: false,
				GitProxyURL:     "url",
				GitProxyScheme:  "http",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "token",
					InsecureSkipTLS: false,
					Proxy: &apimodels.ProxyGitAuth{
						URL:    "url",
						Scheme: "http",
					},
				},
			},
		},
		{
			name: "project with https credentials - proxy",
			in: &GitOldCredentials{
				RemoteURI:       "https://some-url",
				User:            "user",
				Token:           "token",
				InsecureSkipTLS: false,
				GitProxyURL:     "url",
			},
			out: &apimodels.GitAuthCredentials{
				RemoteURL: "https://some-url",
				User:      "user",
				HttpsAuth: &apimodels.HttpsGitAuth{
					Token:           "token",
					InsecureSkipTLS: false,
					Proxy: &apimodels.ProxyGitAuth{
						URL: "url",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := transformSecret(tt.in)
			require.Equal(t, tt.out, out)
		})
	}
}
