package db

import (
	"testing"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

func TestProjectCredentialsRepo_Transform(t *testing.T) {
	tests := []struct {
		name string
		in   *models.ExpandedProjectOld
		out  *apimodels.ExpandedProject
	}{
		{
			name: "project without credentials",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
			},
			out: nil,
		},
		{
			name: "project with empty credentials",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "",
				GitUser:          "",
			},
			out: nil,
		},
		{
			name: "project with ssh credentials - user",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "ssh://some-url",
				GitUser:          "user",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "ssh://some-url",
					User:      "user",
				},
			},
		},
		{
			name: "project with ssh credentials - no user",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "ssh://some-url",
				GitUser:          "",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "ssh://some-url",
				},
			},
		},
		{
			name: "project with ssh credentials - nil user",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "ssh://some-url",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "ssh://some-url",
				},
			},
		},
		{
			name: "project with http credentials - tls false",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitUser:          "user",
				InsecureSkipTLS:  false,
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					User:      "user",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
					},
				},
			},
		},
		{
			name: "project with http credentials - tls true",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitUser:          "user",
				InsecureSkipTLS:  true,
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					User:      "user",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: true,
					},
				},
			},
		},
		{
			name: "project with http credentials - tls undefined",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitUser:          "user",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					User:      "user",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
					},
				},
			},
		},
		{
			name: "project with http credentials - nil user",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
					},
				},
			},
		},
		{
			name: "project with http credentials - proxy",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitProxyURL:      "proxy-url",
				GitProxyScheme:   "http",
				GitProxyUser:     "user",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
						Proxy: &apimodels.ProxyGitAuthSecure{
							URL:    "proxy-url",
							Scheme: "http",
							User:   "user",
						},
					},
				},
			},
		},
		{
			name: "project with http credentials - proxy no user",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitProxyURL:      "proxy-url",
				GitProxyScheme:   "http",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
						Proxy: &apimodels.ProxyGitAuthSecure{
							URL:    "proxy-url",
							Scheme: "http",
						},
					},
				},
			},
		},
		{
			name: "project with http credentials - proxy empty",
			in: &models.ExpandedProjectOld{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitRemoteURI:     "http://some-url",
				GitProxyURL:      "",
				GitProxyScheme:   "",
			},
			out: &apimodels.ExpandedProject{
				CreationDate:     "date",
				LastEventContext: nil,
				ProjectName:      "project",
				Shipyard:         "shippy",
				ShipyardVersion:  "ship",
				Stages:           nil,
				GitCredentials: &apimodels.GitAuthCredentialsSecure{
					RemoteURL: "http://some-url",
					HttpsAuth: &apimodels.HttpsGitAuthSecure{
						InsecureSkipTLS: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := TransformGitCredentials(tt.in)
			require.Equal(t, tt.out, out)
		})
	}
}
