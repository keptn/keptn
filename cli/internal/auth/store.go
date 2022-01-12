package auth

import (
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"sync"
)

// TokenStore is used to store and read an oauth token
type TokenStore interface {
	GetToken() (*oauth2.Token, error)
	StoreToken(token *oauth2.Token) error
	Reset() error
	GetTokenDiscovery() (*OauthDiscoveryResult, error)
	StoreTokenDiscovery(discoveryResult *OauthDiscoveryResult) error
}

// TokenFileName is the name of the file containing the oauth token data
const TokenFileName = "tokens.json"
const DiscoveryResultFileName = "discovery.json"

// NewLocalFileTokenStore creates a new LocalFileTokenStore
// The token store is persisted in the local keptn configuration directory ( ~/.keptn)
func NewLocalFileTokenStore() *LocalFileTokenStore {
	return &LocalFileTokenStore{
		location:          getDefaultTokenLocation(),
		discoveryLocation: getDefaultDiscoveryResultLocation(),
	}
}

// LocalFileTokenStore is a simple token store implementation which stores its data as a
// JSON formatted file(s) on the local file system
type LocalFileTokenStore struct {
	location          string
	discoveryLocation string
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
	tokenAsJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("unable to store token in local token store: %w", err)
	}
	if err := ioutil.WriteFile(t.location, tokenAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store token in local token store: %w", err)
	}
	return nil
}

// StoreTokenDiscovery stores (overwrites) the token discovery information
func (t LocalFileTokenStore) StoreTokenDiscovery(discoveryResult *OauthDiscoveryResult) error {
	discoveryResultAsJSON, err := json.Marshal(discoveryResult)
	if err != nil {
		return fmt.Errorf("unable to store discovery information in local token store: %w", err)
	}
	if err := ioutil.WriteFile(t.discoveryLocation, discoveryResultAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store discovery information in local token store: %w", err)
	}
	return nil
}

// GetTokenDiscovery gets the oauth token discovery information from the token store
func (t LocalFileTokenStore) GetTokenDiscovery() (*OauthDiscoveryResult, error) {
	if !fileutils.FileExists(t.discoveryLocation) {
		return nil, nil
	}
	discoveryFile, err := ioutil.ReadFile(t.discoveryLocation)
	if err != nil {
		return nil, err
	}
	d := &OauthDiscoveryResult{}
	err = json.Unmarshal(discoveryFile, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Reset wipes all persistent oauth information from the store
func (t LocalFileTokenStore) Reset() error {
	if fileutils.FileExists(t.location) {
		if err := os.Remove(t.location); err != nil {
			return fmt.Errorf("unable to delete token from local token store: %w", err)
		}
	}
	if fileutils.FileExists(t.discoveryLocation) {
		if err := os.Remove(t.discoveryLocation); err != nil {
			return fmt.Errorf("unable to delete token discovery from local token store: %w", err)
		}
	}
	return nil
}

// Location checks whether a oauth token is available and eventually returns its location on disk
func (t *LocalFileTokenStore) Location() (bool, string) {
	return fileutils.FileExists(t.location), t.location
}

func getDefaultTokenLocation() string {
	configPath, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return TokenFileName
	}
	return configPath + TokenFileName
}

func getDefaultDiscoveryResultLocation() string {
	configPath, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return DiscoveryResultFileName
	}
	return configPath + DiscoveryResultFileName
}

// TokenStoreMock is an implementation of TokenStore usable as a mock in tests
type TokenStoreMock struct {
	sync.Mutex
	storedToken           *oauth2.Token
	storedTokenDiscovery  *OauthDiscoveryResult
	getTokenFn            func() (*oauth2.Token, error)
	storeTokenFn          func(*oauth2.Token) error
	deleteTokenFn         func() error
	getTokenDiscoveryFn   func() (*OauthDiscoveryResult, error)
	storeTokenDiscoveryFn func(*OauthDiscoveryResult) error
}

func (t *TokenStoreMock) GetToken() (*oauth2.Token, error) {
	t.Lock()
	defer t.Unlock()
	if t.getTokenFn != nil {
		return t.getTokenFn()
	}
	return t.storedToken, nil
}

func (t *TokenStoreMock) StoreToken(token *oauth2.Token) error {
	t.Lock()
	defer t.Unlock()
	if t.storeTokenFn != nil {
		return t.storeTokenFn(token)
	}
	t.storedToken = token
	return nil
}
func (t *TokenStoreMock) Reset() error {
	t.Lock()
	defer t.Unlock()
	if t.Reset() != nil {
		return t.deleteTokenFn()
	}
	return nil
}

func (t *TokenStoreMock) GetTokenDiscovery() (*OauthDiscoveryResult, error) {
	t.Lock()
	defer t.Unlock()
	if t.getTokenDiscoveryFn != nil {
		return t.getTokenDiscoveryFn()
	}
	return t.storedTokenDiscovery, nil
}
func (t *TokenStoreMock) StoreTokenDiscovery(discoveryResult *OauthDiscoveryResult) error {
	t.Lock()
	defer t.Unlock()
	if t.storeTokenDiscoveryFn != nil {
		return t.storeTokenDiscoveryFn(discoveryResult)
	}
	t.storedTokenDiscovery = discoveryResult
	return nil
}
