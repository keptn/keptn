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

// Oauthinfo is a wrapper for oauth related information that is
// managed by the OauthStore
type OauthInfo struct {
	// DiscoveryInfo holds information about Oauth location and capabilites
	DiscoveryInfo *OauthDiscoveryResult
	// ClientValues holds information about information needed from the user/client
	ClientValues *OauthClientValues
	// Token holds information about the fetched tokens
	Token *oauth2.Token
}

// TokenStore is used to store and read oauth related information
type TokenStore interface {
	StoreOauthInfo(*OauthInfo) error
	GetOauthInfo() (*OauthInfo, error)
	StoreTokenInfo(*oauth2.Token) error
	GetTokenInfo() (*oauth2.Token, error)
	StoreDiscoveryInfo(*OauthDiscoveryResult) error
	GetDiscoveryInfo() (*OauthDiscoveryResult, error)
	StoreClientInfo(*OauthClientValues) error
	GetClientInfo() (*OauthClientValues, error)
	Wipe() error
}

// TokenFileName is the name of the file containing the oauth token data
const TokenFileName = "tokens.json"
const DiscoveryResultFileName = "discovery.json"
const ClientValuesFileName = "client.json"

// NewLocalFileTokenStore creates a new LocalFileTokenStore
// The token store is persisted in the local keptn configuration directory ( ~/.keptn)
func NewLocalFileTokenStore() *LocalFileTokenStore {
	return &LocalFileTokenStore{
		location:             getDefaultTokenLocation(),
		discoveryLocation:    getDefaultDiscoveryResultLocation(),
		clientValuesLocation: getDefaultClientValuesLocation(),
	}
}

// LocalFileTokenStore is a simple token store implementation which stores its data as a
// JSON formatted file(s) on the local file system
type LocalFileTokenStore struct {
	location             string
	discoveryLocation    string
	clientValuesLocation string
}

// StoreOauthInfo persists all oauth related information in the store
func (t LocalFileTokenStore) StoreOauthInfo(i *OauthInfo) error {
	if err := t.StoreClientInfo(i.ClientValues); err != nil {
		return err
	}
	if err := t.StoreDiscoveryInfo(i.DiscoveryInfo); err != nil {
		return err
	}
	return t.StoreTokenInfo(i.Token)
}

// GetOauthInfo retrieves all oauth related information from the store
func (t LocalFileTokenStore) GetOauthInfo() (*OauthInfo, error) {
	clientValues, err := t.GetClientInfo()
	if err != nil {
		return nil, err
	}
	token, err := t.GetTokenInfo()
	if err != nil {
		return nil, err
	}
	discoveryInfo, err := t.GetDiscoveryInfo()
	if err != nil {
		return nil, err
	}

	return &OauthInfo{
		DiscoveryInfo: discoveryInfo,
		ClientValues:  clientValues,
		Token:         token,
	}, nil
}

// GetToken gets the oauth token from the token store
func (t LocalFileTokenStore) GetTokenInfo() (*oauth2.Token, error) {
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
func (t LocalFileTokenStore) StoreTokenInfo(token *oauth2.Token) error {
	tokenAsJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("unable to store token in local token store: %w", err)
	}
	if err := ioutil.WriteFile(t.location, tokenAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store token in local token store: %w", err)
	}
	return nil
}

// StoreDiscovery stores (overwrites) the token discovery information
func (t LocalFileTokenStore) StoreDiscoveryInfo(discoveryResult *OauthDiscoveryResult) error {
	discoveryResultAsJSON, err := json.Marshal(discoveryResult)
	if err != nil {
		return fmt.Errorf("unable to store discovery information in local token store: %w", err)
	}
	if err := ioutil.WriteFile(t.discoveryLocation, discoveryResultAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store discovery information in local token store: %w", err)
	}
	return nil
}

// getDiscovery gets the oauth token discovery information from the token store
func (t LocalFileTokenStore) GetDiscoveryInfo() (*OauthDiscoveryResult, error) {
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

func (t LocalFileTokenStore) StoreClientInfo(clientValues *OauthClientValues) error {
	clientValuesAsJSON, err := json.Marshal(clientValues)
	if err != nil {
		return fmt.Errorf("unable to store client values in local token store: %w", err)
	}
	if err := ioutil.WriteFile(t.clientValuesLocation, clientValuesAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store client values in local token store: %w", err)
	}
	return nil
}

func (t LocalFileTokenStore) GetClientInfo() (*OauthClientValues, error) {
	if !fileutils.FileExists(t.clientValuesLocation) {
		return nil, fmt.Errorf("unable to find local client values")
	}
	clientValuesFile, err := ioutil.ReadFile(t.clientValuesLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to read local client values: %w", err)
	}
	v := &OauthClientValues{}
	err = json.Unmarshal(clientValuesFile, v)
	if err != nil {
		return nil, fmt.Errorf("unable to read local client values: %w", err)
	}
	return v, nil
}

// Wipe wipes all persistent OAuth information from the store
func (t LocalFileTokenStore) Wipe() error {
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
	if fileutils.FileExists(t.clientValuesLocation) {
		if err := os.Remove(t.clientValuesLocation); err != nil {
			return fmt.Errorf("unable to delete client values from local token store: %w", err)
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

func getDefaultClientValuesLocation() string {
	configPath, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		return ClientValuesFileName
	}
	return configPath + ClientValuesFileName
}

// TokenStoreMock is an implementation of TokenStore usable as a mock in tests
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
