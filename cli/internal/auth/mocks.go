package auth

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
	"sync"
)

// BrwoserMock is an implementation of URLOpener
// usable mocking in tests
type BrowserMock struct {
	openFn func(string) error
}

// Open calls the mocked function of the BrwoserMock
func (b *BrowserMock) Open(url string) error {
	if b != nil && b.openFn != nil {
		return b.openFn(url)
	}
	return nil
}

// HTTPClientMock is mock implementation of http client usable for testing
type HTTPClientMock struct {
	DoFunc func(r *http.Request) (*http.Response, error)
}

func (h HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return h.DoFunc(r)
}

// RedirectHandlerMock is a mocked implementation of TokenGetter
// usable in tests
type RedirectHandlerMock struct {
	handleFn func([]byte, *oauth2.Config, string) (*oauth2.Token, error)
}

func (t RedirectHandlerMock) Handle(codeVerifier []byte, oauthConfig *oauth2.Config, state string) (*oauth2.Token, error) {
	if t.handleFn != nil {
		return t.handleFn(codeVerifier, oauthConfig, state)
	}
	panic("handleFn of RedirectHandlerMock not set")
}

// TokenStoreMock is an implementation of OauthStore usable as a mock in tests
type TokenStoreMock struct {
	sync.Mutex
	storedOauthInfo       *OauthInfo
	storedToken           *oauth2.Token
	storedTokenDiscovery  *OauthDiscoveryResult
	storedClientValues    *OauthClientValues
	created               bool
	CreatedFn             func() bool
	StoreOauthInfoFn      func(*OauthInfo) error
	GetOauthInfoFn        func() (*OauthInfo, error)
	GetTokenFn            func() (*oauth2.Token, error)
	StoreTokenFn          func(*oauth2.Token) error
	DeleteTokenFn         func() error
	GetTokenDiscoveryFn   func() (*OauthDiscoveryResult, error)
	StoreTokenDiscoveryFn func(*OauthDiscoveryResult) error
	StoreClientValuesFn   func(values *OauthClientValues) error
	GetClientValuesFn     func() (*OauthClientValues, error)
}

func (t *TokenStoreMock) StoreOauthInfo(i *OauthInfo) error {
	t.Lock()
	defer t.Unlock()
	if t.StoreOauthInfoFn != nil {
		return t.StoreOauthInfoFn(i)
	}
	t.storedOauthInfo = i
	t.storedToken = i.Token
	t.storedTokenDiscovery = i.DiscoveryInfo
	t.storedClientValues = i.ClientValues

	return nil
}

func (t *TokenStoreMock) GetOauthInfo() (*OauthInfo, error) {
	t.Lock()
	defer t.Unlock()
	if t.GetOauthInfoFn != nil {
		return t.GetOauthInfoFn()
	}
	return t.storedOauthInfo, nil
}

func (t *TokenStoreMock) GetTokenInfo() (*oauth2.Token, error) {
	t.Lock()
	defer t.Unlock()
	if t.GetTokenFn != nil {
		return t.GetTokenFn()
	}
	return t.storedToken, nil
}

func (t *TokenStoreMock) StoreTokenInfo(token *oauth2.Token) error {
	t.Lock()
	defer t.Unlock()
	if t.StoreTokenFn != nil {
		return t.StoreTokenFn(token)
	}
	t.storedToken = token
	return nil
}
func (t *TokenStoreMock) Wipe() error {
	t.Lock()
	defer t.Unlock()
	if t.Wipe() != nil {
		return t.DeleteTokenFn()
	}
	return nil
}

func (t *TokenStoreMock) Created() bool {
	t.Lock()
	defer t.Unlock()
	if t.CreatedFn != nil {
		return t.CreatedFn()
	}
	return t.created
}

func (t *TokenStoreMock) GetDiscoveryInfo() (*OauthDiscoveryResult, error) {
	t.Lock()
	defer t.Unlock()
	if t.GetTokenDiscoveryFn != nil {
		return t.GetTokenDiscoveryFn()
	}
	return t.storedTokenDiscovery, nil
}
func (t *TokenStoreMock) StoreDiscoveryInfo(discoveryResult *OauthDiscoveryResult) error {
	t.Lock()
	defer t.Unlock()
	if t.StoreTokenDiscoveryFn != nil {
		return t.StoreTokenDiscoveryFn(discoveryResult)
	}
	t.storedTokenDiscovery = discoveryResult
	return nil
}

func (t *TokenStoreMock) StoreClientInfo(values *OauthClientValues) error {
	t.Lock()
	defer t.Unlock()
	if t.StoreClientValuesFn != nil {
		return t.StoreClientValuesFn(values)
	}
	t.storedClientValues = values
	return nil
}

func (t *TokenStoreMock) GetClientInfo() (*OauthClientValues, error) {
	t.Lock()
	defer t.Unlock()
	if t.GetClientValuesFn != nil {
		return t.GetClientValuesFn()
	}
	return t.storedClientValues, nil
}

type OAuthAuthenticatorMock struct {
	AuthCalled, GetAuthClientCalled, TokenStoreCalled bool
	AuthFn                                            func(clientValues OauthClientValues) error
	GetOauthClientFn                                  func(ctx context.Context) (*http.Client, error)
	TokenStoreFn                                      func() OauthStore
}

func (a OAuthAuthenticatorMock) Auth(clientValues OauthClientValues) error {
	a.AuthCalled = true
	if a.AuthFn != nil {
		return a.AuthFn(clientValues)
	}
	panic("AuthFn called on mock but not set")
}

func (a OAuthAuthenticatorMock) GetOauthClient(ctx context.Context) (*http.Client, error) {
	a.GetAuthClientCalled = true
	if a.GetOauthClientFn != nil {
		return a.GetOauthClientFn(ctx)
	}
	panic("GetOauthClientFn called on mock but not set")
}

func (a OAuthAuthenticatorMock) TokenStore() OauthStore {
	a.TokenStoreCalled = true
	if a.TokenStoreFn != nil {
		return a.TokenStoreFn()
	}
	panic("TokenStoreFn called on mock but not set")
}
