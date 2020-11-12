package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/keptn/keptn/cli/pkg/file"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"
)

// CLIConfig holds infos of the CLI config
type CLIConfig struct {
	AutomaticVersionCheck bool       `json:"automatic_version_check"`
	LastVersionCheck      *time.Time `json:"last_version_check"`
	CurrentContext        string     `json:"current-context"`
}

// CLIConfigManager manages the path of the CLI config
type CLIConfigManager struct {
	CLIConfigPath string
}

// NewCLIConfigManager creates a new CLIConfigManager
func NewCLIConfigManager() *CLIConfigManager {
	cliConfigManager := CLIConfigManager{}

	dir, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}
	cliConfigManager.CLIConfigPath = dir + "config"
	return &cliConfigManager
}

// LoadCLIConfig loads the configuration from file
func (c *CLIConfigManager) LoadCLIConfig() (CLIConfig, error) {

	cliConfig := CLIConfig{AutomaticVersionCheck: true}
	if !file.FileExists(c.CLIConfigPath) {
		return cliConfig, nil
	}

	data, err := file.ReadFile(c.CLIConfigPath)
	if err != nil {
		return cliConfig, fmt.Errorf("error when reading config file: %v", err)
	}
	if err := json.Unmarshal([]byte(data), &cliConfig); err != nil {
		return cliConfig, fmt.Errorf("error when unmarshalling config file: %v", err)
	}

	return cliConfig, nil
}

// StoreCLIConfig stores the configuration into the file
func (c *CLIConfigManager) StoreCLIConfig(config CLIConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error when marshalling config file: %v", err)
	}
	if err := ioutil.WriteFile(c.CLIConfigPath, []byte(data), 0644); err != nil {
		return fmt.Errorf("error when writing config file: %v", err)
	}
	return nil
}
