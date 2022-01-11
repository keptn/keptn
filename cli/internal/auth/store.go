package auth

import (
	"encoding/json"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

// TokenStore is used to store and read an oauth token
type TokenStore interface {
	// GetToken gets a token from the token store
	GetToken() (*oauth2.Token, error)
	// StoreToken stores (or overwrites) the token in the token store
	StoreToken(token *oauth2.Token) error
	// DeleteToken deletes the token from the token store
	DeleteToken() error
}

// TokenFileName is the name of the file containing the oauth token data
const TokenFileName = "tokens.json"

// NewLocalFileTokenStore creates a new LocalFileTokenStore
// The token store is persisted in the local keptn configuration directory ( ~/.keptn)
func NewLocalFileTokenStore() *LocalFileTokenStore {
	location := getDefaultLocation()
	return &LocalFileTokenStore{location: location}
}

// LocalFileTokenStore is a simple token store implementation which stores its data as a
// JSON formated file on the local file system
type LocalFileTokenStore struct {
	location string
}

// GetToken gets the oauth token from the token store
func (t LocalFileTokenStore) GetToken() (*oauth2.Token, error) {
	tokenFile, err := ioutil.ReadFile(t.location)
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

// StoreToken stores (overwrites) an oauth token in the token store
func (t LocalFileTokenStore) StoreToken(token *oauth2.Token) error {
	tokenMarshalled, err := json.Marshal(token)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(t.location, tokenMarshalled, 0600)
	if err != nil {
		return err
	}
	return nil
}

// DeleteToken deletes the oauth token from the token store
func (t LocalFileTokenStore) DeleteToken() error {
	if fileutils.FileExists(t.location) {
		if err := os.Remove(t.location); err != nil {
			return err
		}
	}
	return nil
}

// Location checks whether a oauth token is available and eventually returns its location on disk
func (t *LocalFileTokenStore) Location() (bool, string) {
	return fileutils.FileExists(t.location), t.location
}

func getDefaultLocation() string {
	configPath, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return TokenFileName
	}
	return configPath + TokenFileName
}

// TokenStoreMock is an implementation of TokenStore usable as a mock in tests
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
func (t *TokenStoreMock) DeleteToken() error {
	return nil
}
