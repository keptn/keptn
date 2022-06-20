package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/common/kubeutils"
	"golang.org/x/oauth2"
)

// OauthInfo is a wrapper for oauth related information that is
// managed by the OauthStore
type OauthInfo struct {
	// DiscoveryInfo holds information about Oauth location and capabilities
	DiscoveryInfo *OauthDiscoveryResult
	// ClientValues holds information about information needed from the user/client
	ClientValues *OauthClientValues
	// Token holds information about the fetched tokens
	Token *oauth2.Token
}

// OauthStore is used to store and read oauth related information
type OauthStore interface {
	StoreOauthInfo(*OauthInfo) error
	GetOauthInfo() (*OauthInfo, error)
	StoreTokenInfo(*oauth2.Token) error
	GetTokenInfo() (*oauth2.Token, error)
	StoreDiscoveryInfo(*OauthDiscoveryResult) error
	GetDiscoveryInfo() (*OauthDiscoveryResult, error)
	StoreClientInfo(*OauthClientValues) error
	GetClientInfo() (*OauthClientValues, error)
	Wipe() error
	Created() bool
}

// TokenFileName is the name of the file containing the oauth token data
const TokenFileName = "tokens.json"
const DiscoveryResultFileName = "discovery.json"
const ClientValuesFileName = "client.json"

// NewLocalFileOauthStore creates a new LocalFileOauthStore which sotres its data inside the local Keptn
// configuration directory (~/.keptn)
func NewLocalFileOauthStore() *LocalFileOauthStore {
	return &LocalFileOauthStore{
		location:             getDefaultTokenLocation(),
		discoveryLocation:    getDefaultDiscoveryResultLocation(),
		clientValuesLocation: getDefaultClientValuesLocation(),
	}
}

// LocalFileOauthStore is a local file based implementation of OauthStore which stores its data as a
// JSON formatted file(s) on the local file system
type LocalFileOauthStore struct {
	location             string
	discoveryLocation    string
	clientValuesLocation string
}

// StoreOauthInfo persists all oauth related information in the store
func (t LocalFileOauthStore) StoreOauthInfo(i *OauthInfo) error {
	if err := t.StoreClientInfo(i.ClientValues); err != nil {
		return err
	}
	if err := t.StoreDiscoveryInfo(i.DiscoveryInfo); err != nil {
		return err
	}
	return t.StoreTokenInfo(i.Token)
}

// GetOauthInfo retrieves all oauth related information from the store
func (t LocalFileOauthStore) GetOauthInfo() (*OauthInfo, error) {
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

// GetTokenInfo retrieves the oauth token
func (t LocalFileOauthStore) GetTokenInfo() (*oauth2.Token, error) {
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

// StoreTokenInfo stores (overwrites) an oauth token
func (t LocalFileOauthStore) StoreTokenInfo(token *oauth2.Token) error {
	tokenAsJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}
	if err := ioutil.WriteFile(t.location, tokenAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}
	return nil
}

// StoreDiscoveryInfo stores (overwrites) the token discovery information
func (t LocalFileOauthStore) StoreDiscoveryInfo(discoveryResult *OauthDiscoveryResult) error {
	discoveryResultAsJSON, err := json.Marshal(discoveryResult)
	if err != nil {
		return fmt.Errorf("unable to store discovery information: %w", err)
	}
	if err := ioutil.WriteFile(t.discoveryLocation, discoveryResultAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store discovery information: %w", err)
	}
	return nil
}

// GetDiscoveryInfo retrieves the oauth token discovery information
func (t LocalFileOauthStore) GetDiscoveryInfo() (*OauthDiscoveryResult, error) {
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

// StoreClientInfo stores client information
func (t LocalFileOauthStore) StoreClientInfo(clientValues *OauthClientValues) error {
	clientValuesAsJSON, err := json.Marshal(clientValues)
	if err != nil {
		return fmt.Errorf("unable to store client values: %w", err)
	}
	if err := ioutil.WriteFile(t.clientValuesLocation, clientValuesAsJSON, 0600); err != nil {
		return fmt.Errorf("unable to store client values: %w", err)
	}
	return nil
}

// GetClientInfo retrieves client information
func (t LocalFileOauthStore) GetClientInfo() (*OauthClientValues, error) {
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
func (t LocalFileOauthStore) Wipe() error {
	if fileutils.FileExists(t.location) {
		if err := os.Remove(t.location); err != nil {
			return fmt.Errorf("unable to delete token: %w", err)
		}
	}
	if fileutils.FileExists(t.discoveryLocation) {
		if err := os.Remove(t.discoveryLocation); err != nil {
			return fmt.Errorf("unable to delete token discovery: %w", err)
		}
	}
	if fileutils.FileExists(t.clientValuesLocation) {
		if err := os.Remove(t.clientValuesLocation); err != nil {
			return fmt.Errorf("unable to delete client values: %w", err)
		}
	}
	return nil
}

// Created checks whether a LocalFileOauthStore was created or not
func (t *LocalFileOauthStore) Created() bool {
	return fileutils.FileExists(t.location) && fileutils.FileExists(t.discoveryLocation) && fileutils.FileExists(t.clientValuesLocation)
}

func getDefaultTokenLocation() string {
	configPath, err := kubeutils.GetKeptnDirectory()
	if err != nil {
		return TokenFileName
	}
	return configPath + TokenFileName
}

func getDefaultDiscoveryResultLocation() string {
	configPath, err := kubeutils.GetKeptnDirectory()
	if err != nil {
		return DiscoveryResultFileName
	}
	return configPath + DiscoveryResultFileName
}

func getDefaultClientValuesLocation() string {
	configPath, err := kubeutils.GetKeptnDirectory()
	if err != nil {
		return ClientValuesFileName
	}
	return configPath + ClientValuesFileName
}
