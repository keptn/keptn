package auth

import (
	"context"
	"golang.org/x/oauth2"
	"sync"
)

type TokenNotifyFunc func(*oauth2.Token) error

// NotifyRefreshTokenSource is an implementation of TokenSource and is used to refresh the access token
// using a refresh token from the local token store
type NotifyRefreshTokenSource struct {
	tokenStore TokenStore
	mu         sync.Mutex
	config     *oauth2.Config
}

// Token tries to get an already existing token from the local token store.
// If it is still valid, it is returned. Otherwise, a new token will eventually retrieved and stored
// in the local token store
func (s *NotifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.tokenStore.GetTokenInfo()
	if err != nil {
		return nil, err
	}
	if t.Valid() {
		return t, nil
	}
	tokenSource := s.config.TokenSource(context.TODO(), t)
	t, err = tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return t, s.tokenStore.StoreTokenInfo(t)
}
