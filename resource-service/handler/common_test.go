package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	gittransport "github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/resource-service/common_models"
)

func Test_getAuthMethod(t *testing.T) {
	tests := []struct {
		name              string
		gitCredentials    *common_models.GitCredentials
		wantErr           bool
		expectedGoGitAuth gittransport.AuthMethod
	}{
		{
			name:              "no credentials",
			gitCredentials:    &common_models.GitCredentials{},
			wantErr:           false,
			expectedGoGitAuth: nil,
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
			expectedGoGitAuth: &githttp.BasicAuth{
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
			expectedGoGitAuth: &githttp.BasicAuth{
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
			expectedGoGitAuth: &githttp.BasicAuth{
				Username: "keptnuser",
				Password: "some-token",
			},
		},
		{
			name: "credentials without token",
			gitCredentials: &common_models.GitCredentials{
				RemoteURL: "https://some.url",
				HttpsAuth: &apimodels.HttpsGitAuth{
					InsecureSkipTLS: false,
				},
				User: "user",
			},
			wantErr:           false,
			expectedGoGitAuth: &githttp.BasicAuth{Username: "user", Password: ""},
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
			wantErr:           true,
			expectedGoGitAuth: nil,
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
			wantErr:           true,
			expectedGoGitAuth: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := getAuthMethod(tt.gitCredentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAuthMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectedGoGitAuth == nil {
				require.Nil(t, auth.GoGitAuth)
				return
			}
			if err != nil && auth.GoGitAuth != tt.expectedGoGitAuth {
				t.Errorf("getAuthMethod() auth = %v, expectedGoGitAuth %v", err, tt.wantErr)
			}
			if auth != nil {
				require.NotNil(t, auth.Git2GoAuth.CredCallback)
				if tt.gitCredentials.SshAuth != nil {
					require.NotNil(t, auth.Git2GoAuth.CertCallback)
				}
				if tt.gitCredentials.HttpsAuth.Proxy != nil {
					require.NotNil(t, auth.Git2GoAuth.ProxyOptions)
				}
			}
		})
	}
}

func TestOnAPIError(t *testing.T) {

	ginContext := func(w *httptest.ResponseRecorder) *gin.Context {
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			Header: make(http.Header),
		}

		return ctx
	}

	tests := []struct {
		name           string
		recorder       *httptest.ResponseRecorder
		err            error
		wantStatusCode int
	}{
		{
			name:           "ErrProjectAlreadyExists -> 409 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrProjectAlreadyExists,
			wantStatusCode: http.StatusConflict,
		},
		{
			name:           "errors.ErrProjectRepositoryNotEmpty -> 409 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrProjectRepositoryNotEmpty,
			wantStatusCode: http.StatusConflict,
		},
		{
			name:           "errors.ErrInvalidGitToken -> 424 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrInvalidGitToken,
			wantStatusCode: http.StatusFailedDependency,
		},
		{
			name:           "errors.ErrAuthenticationRequired -> 424 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrAuthenticationRequired,
			wantStatusCode: http.StatusFailedDependency,
		},
		{
			name:           "errors.ErrAuthorizationFailed -> 424 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrAuthorizationFailed,
			wantStatusCode: http.StatusFailedDependency,
		},
		{
			name:           " errors.ErrCredentialsNotFound -> 404 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrCredentialsNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "errors.ErrMalformedCredentials -> 424 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrMalformedCredentials,
			wantStatusCode: http.StatusFailedDependency,
		},
		{
			name:           "errors.ErrCredentialsInvalidRemoteURL -> 400 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrCredentialsInvalidRemoteURL,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "errors.ErrCredentialsTokenMustNotBeEmpty -> 400 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrCredentialsTokenMustNotBeEmpty,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "errors.ErrRepositoryNotFound -> 404 ",
			recorder:       httptest.NewRecorder(),
			err:            errors.ErrRepositoryNotFound,
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "error -> 500 ",
			recorder:       httptest.NewRecorder(),
			err:            fmt.Errorf("some other err"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OnAPIError(ginContext(tt.recorder), tt.err)
			assert.Equal(t, tt.wantStatusCode, tt.recorder.Code)
		})
	}
}
