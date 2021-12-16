package auth

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type LocalFileTokenStore struct {
}

type TokenStore interface {
	GetToken() (*oauth2.Token, error)
	StoreToken(token *oauth2.Token) error
}

func (t LocalFileTokenStore) GetToken() (*oauth2.Token, error) {
	tokenFile, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.Unmarshal(tokenFile, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t LocalFileTokenStore) StoreToken(token *oauth2.Token) error {
	tokenMarshalled, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("tokens.json", tokenMarshalled, 0600)
	if err != nil {
		return err
	}
	// persist token
	return nil // or error
}

type TokenStoreMock struct {
	storedToken  *oauth2.Token
	getTokenFn   func() (*oauth2.Token, error)
	storeTokenFn func(*oauth2.Token) error
}

func (t *TokenStoreMock) GetToken() (*oauth2.Token, error) {
	if t != nil && t.getTokenFn != nil {
		return t.getTokenFn()
	}
	return t.storedToken, nil
}

func (t *TokenStoreMock) StoreToken(token *oauth2.Token) error {
	if t != nil && t.storeTokenFn != nil {
		return t.storeTokenFn(token)
	}
	t.storedToken = token
	return nil
}
