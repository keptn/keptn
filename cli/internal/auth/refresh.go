package auth

import (
	"context"
	"golang.org/x/oauth2"
	"sync"
)

type TokenNotifyFunc func(*oauth2.Token) error

type NotifyRefreshTokenSource struct {
	tokenStore TokenStore
	mu         sync.Mutex
	config     *oauth2.Config
}

func (s *NotifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.tokenStore.GetToken()
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
	return t, s.tokenStore.StoreToken(t)
}
