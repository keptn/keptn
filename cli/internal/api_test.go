package internal

import (
	"context"
	"errors"
	"github.com/keptn/keptn/cli/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HTTPClientIsUsedDirectly(t *testing.T) {

	t.Run("GetAPISet_WithClient", func(t *testing.T) {
		authorizationCodeFlow := &auth.OAuthAuthenticatorMock{}
		apiSet, err := getAPISetWithOauthGetter("", "", authorizationCodeFlow, &http.Client{})
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.False(t, authorizationCodeFlow.TokenStoreCalled)
		require.False(t, authorizationCodeFlow.AuthCalled)
		require.False(t, authorizationCodeFlow.GetAuthClientCalled)
		require.Empty(t, apiSet.Token())
	})

	t.Run("GetAPISet_WithClient_AndXToken", func(t *testing.T) {
		authorizationCodeFlow := &auth.OAuthAuthenticatorMock{}
		apiSet, err := getAPISetWithOauthGetter("", "XToken", authorizationCodeFlow, &http.Client{})
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.False(t, authorizationCodeFlow.TokenStoreCalled)
		require.False(t, authorizationCodeFlow.AuthCalled)
		require.False(t, authorizationCodeFlow.GetAuthClientCalled)
		require.Equal(t, "XToken", apiSet.Token())
	})

}

func Test_WhenHeadRequestFails_OAuthFlowIsTriggeredAgain(t *testing.T) {
	t.Run("GetAPISet_OAuth_RefreshTokenInvalid_StartAuth", func(t *testing.T) {
		authorizationCodeFlowStarted := false
		tokenStore := auth.TokenStoreMock{
			CreatedFn: func() bool { return true },
			GetOauthInfoFn: func() (*auth.OauthInfo, error) {
				return &auth.OauthInfo{
					DiscoveryInfo: &auth.OauthDiscoveryResult{},
					ClientValues:  &auth.OauthClientValues{},
				}, nil
			},
		}
		authenticator := &auth.OAuthAuthenticatorMock{
			TokenStoreFn:     func() auth.OauthStore { return &tokenStore },
			GetOauthClientFn: func(ctx context.Context) (*http.Client, error) { return &http.Client{}, nil },
			AuthFn: func(clientValues auth.OauthClientValues) error {
				authorizationCodeFlowStarted = true
				return nil
			},
		}

		apiSet, err := getAPISetWithOauthGetter("", "", authenticator)
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.True(t, authorizationCodeFlowStarted)
		require.Empty(t, apiSet.Token())

		apiSet, err = getAPISetWithOauthGetter("", "XToken", authenticator)
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.True(t, authorizationCodeFlowStarted)
		require.Equal(t, "XToken", apiSet.Token())
	})
	t.Run("GetAPISet_OAuth_RefreshTokenStillValid_DoesNotStartAuth", func(t *testing.T) {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) { w.WriteHeader(200) }))
		authorizationCodeFlowStarted := false
		tokenStore := auth.TokenStoreMock{
			CreatedFn: func() bool { return true },
			GetOauthInfoFn: func() (*auth.OauthInfo, error) {
				return &auth.OauthInfo{DiscoveryInfo: &auth.OauthDiscoveryResult{TokenEndpoint: ts.URL}, ClientValues: &auth.OauthClientValues{}}, nil
			},
		}
		authenticator := &auth.OAuthAuthenticatorMock{
			TokenStoreFn:     func() auth.OauthStore { return &tokenStore },
			GetOauthClientFn: func(ctx context.Context) (*http.Client, error) { return &http.Client{}, nil },
			AuthFn: func(clientValues auth.OauthClientValues) error {
				authorizationCodeFlowStarted = true
				return nil
			},
		}

		apiSet, err := getAPISetWithOauthGetter("", "", authenticator)
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.False(t, authorizationCodeFlowStarted)
		require.Empty(t, apiSet.Token())

		apiSet, err = getAPISetWithOauthGetter("", "XToken", authenticator)
		require.Nil(t, err)
		require.NotNil(t, apiSet)
		require.False(t, authorizationCodeFlowStarted)
		require.Equal(t, "XToken", apiSet.Token())
	})
}

func TestOnAPIError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{name: "test 401",
			err:     errors.New("error with status code 401"),
			wantErr: errors.New(ErrNotAuthenticated),
		},
		{name: "test 403",
			err:     errors.New("error with status code 403"),
			wantErr: errors.New(ErrForbidden),
		},
		{name: "test 500",
			err:     errors.New("error with status code 500"),
			wantErr: errors.New(ErrInternalServerError),
		},
		{name: "test not covered code 418",
			err:     errors.New("error with status code 418"),
			wantErr: errors.New("error with status code 418"),
		},
		{name: "test does not overwrite a covered code error with more stuff in it",
			err:     errors.New("This is an error coming from git and it tells you for project x you got error with status code 401"),
			wantErr: errors.New("This is an error coming from git and it tells you for project x you got error with status code 401"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, OnAPIError(tt.err), tt.wantErr)
		})
	}
}
