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
