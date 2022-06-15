package db

import (
	"fmt"
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	common_mock "github.com/keptn/keptn/shipyard-controller/common/fake"
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

func TestMongoDBSecretCredentialsRepo_UpdateSecret(t *testing.T) {
	tests := []struct {
		name        string
		secretStore common_mock.SecretStoreMock
		in          *ExpandedProjectOld
		err         error
	}{
		{
			name: "valid old",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credentials": []byte(`{"user":"gitea_admin","token":"af05ab26ef932b7ec955da2b18faa7e6c28c538a","remoteURI":"http://gitea-http:3000/gitea_admin/keptn-test-2223-1-backup-restore"}`)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: nil,
		},
		{
			name: "valid new",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credentials": []byte(`{"user":"gitea_admin","remoteURI":"http://gitea-http:3000/gitea_admin/keptn-test-2223-1-backup-restore", "https":{"token":"af05ab26ef932b7ec955da2b18faa7e6c28c538a"}}`)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: nil,
		},
		{
			name: "non readable credentials",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credential": []byte(`{"user":"gitea_admin","remoteURI":"http://gitea-http:3000/gitea_admin/keptn-test-2223-1-backup-restore", "https":{"token":"af05ab26ef932b7ec955da2b18faa7e6c28c538a"}}`)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: nil,
		},
		{
			name: "empty credentials",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credentials": []byte(``)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: fmt.Errorf("failed to unmarshal git-credentials secret during migration for project project"),
		},
		{
			name: "invalid credentials",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credentials": []byte(`invalid`)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: fmt.Errorf("failed to unmarshal git-credentials secret during migration for project project"),
		},
		{
			name: "cannot get secret",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return nil, fmt.Errorf("some err")
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: fmt.Errorf("failed to get git-credentials secret during migration for project project"),
		},
		{
			name: "no secret",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return nil, nil
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: nil,
		},
		{
			name: "update err",
			secretStore: common_mock.SecretStoreMock{
				GetSecretFunc: func(name string) (map[string][]byte, error) {
					return map[string][]byte{"git-credentials": []byte(`{"user":"gitea_admin","token":"af05ab26ef932b7ec955da2b18faa7e6c28c538a","remoteURI":"http://gitea-http:3000/gitea_admin/keptn-test-2223-1-backup-restore"}`)}, nil
				},
				UpdateSecretFunc: func(name string, content map[string][]byte) error {
					return fmt.Errorf("some err")
				},
			},
			in: &ExpandedProjectOld{
				ProjectName: "project",
			},
			err: fmt.Errorf("could not store git credentials during migration for project project: some err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMongoDBSecretCredentialsRepo(&tt.secretStore)
			err := repo.UpdateSecret(tt.in)
			require.Equal(t, tt.err, err)
		})
	}
}
