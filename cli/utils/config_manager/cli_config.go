package config_manager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/keptn/keptn/cli/utils"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

// CLIConfig holds infos of the CLI config
type CLIConfig struct {
	AutomaticVersionCheck bool       `json:"automatic_version_check"`
	LastVersionCheck      *time.Time `json:"last_version_check"`
}

type CLIConfigManager struct {
	cliConfigPath string
}

func newCLIConfigManager() *CLIConfigManager {
	cliConfigManager := CLIConfigManager{}

	dir, err := keptnutils.GetKeptnDirectory()
	if err != nil {
		log.Fatal(err)
	}
	cliConfigManager.cliConfigPath = dir + "config"
	return &cliConfigManager
}

// LoadCLIConfig loads the configuration from file
func (c *CLIConfigManager) LoadCLIConfig() (CLIConfig, error) {

	cliConfig := CLIConfig{}
	if !utils.FileExists(c.cliConfigPath) {
		return cliConfig, nil
	}

	data, err := utils.ReadFile(c.cliConfigPath)
	if err != nil {
		return cliConfig, err
	}
	err = json.Unmarshal([]byte(data), &cliConfig)
	return cliConfig, err
}

// StoreCLIConfig stores the configuration into the file
func (c *CLIConfigManager) StoreCLIConfig(config CLIConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.cliConfigPath, []byte(data), 0644)
}
