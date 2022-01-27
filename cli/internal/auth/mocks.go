package auth

import (
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
	storeOauthInfoFn      func(*OauthInfo) error
	getOauthInfoFn        func() (*OauthInfo, error)
	getTokenFn            func() (*oauth2.Token, error)
	storeTokenFn          func(*oauth2.Token) error
	deleteTokenFn         func() error
	getTokenDiscoveryFn   func() (*OauthDiscoveryResult, error)
	storeTokenDiscoveryFn func(*OauthDiscoveryResult) error
	storeClientValuesFn   func(values *OauthClientValues) error
	getClientValuesFn     func() (*OauthClientValues, error)
}

func (t *TokenStoreMock) StoreOauthInfo(i *OauthInfo) error {
	t.Lock()
	defer t.Unlock()
	if t.storeOauthInfoFn != nil {
		return t.storeOauthInfoFn(i)
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
	if t.getOauthInfoFn != nil {
		return t.getOauthInfoFn()
	}
	return t.storedOauthInfo, nil
}

func (t *TokenStoreMock) GetTokenInfo() (*oauth2.Token, error) {
	t.Lock()
	defer t.Unlock()
	if t.getTokenFn != nil {
		return t.getTokenFn()
	}
	return t.storedToken, nil
}

func (t *TokenStoreMock) StoreTokenInfo(token *oauth2.Token) error {
	t.Lock()
	defer t.Unlock()
	if t.storeTokenFn != nil {
		return t.storeTokenFn(token)
	}
	t.storedToken = token
	return nil
}
func (t *TokenStoreMock) Wipe() error {
	t.Lock()
	defer t.Unlock()
	if t.Wipe() != nil {
		return t.deleteTokenFn()
	}
	return nil
}

func (t *TokenStoreMock) GetDiscoveryInfo() (*OauthDiscoveryResult, error) {
	t.Lock()
	defer t.Unlock()
	if t.getTokenDiscoveryFn != nil {
		return t.getTokenDiscoveryFn()
	}
	return t.storedTokenDiscovery, nil
}
func (t *TokenStoreMock) StoreDiscoveryInfo(discoveryResult *OauthDiscoveryResult) error {
	t.Lock()
	defer t.Unlock()
	if t.storeTokenDiscoveryFn != nil {
		return t.storeTokenDiscoveryFn(discoveryResult)
	}
	t.storedTokenDiscovery = discoveryResult
	return nil
}

func (t *TokenStoreMock) StoreClientInfo(values *OauthClientValues) error {
	t.Lock()
	defer t.Unlock()
	if t.storeClientValuesFn != nil {
		return t.storeClientValuesFn(values)
	}
	t.storedClientValues = values
	return nil
}

func (t *TokenStoreMock) GetClientInfo() (*OauthClientValues, error) {
	t.Lock()
	defer t.Unlock()
	if t.getClientValuesFn != nil {
		return t.getClientValuesFn()
	}
	return t.storedClientValues, nil
}
