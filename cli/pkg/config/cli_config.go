package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/keptn/go-utils/pkg/common/fileutils"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const keptnFolderName = ".keptn"

// CLIConfig holds infos of the CLI config
type CLIConfig struct {
	AutomaticVersionCheck bool       `json:"automatic_version_check" mapstructure:"automatic_version_check"`
	KubeContextCheck      bool       `json:"kube_context_check" mapstructure:"kube_context_check"`
	LastVersionCheck      *time.Time `json:"last_version_check" mapstructure:"last_version_check"`
	CurrentContext        string     `json:"current-context" mapstructure:"current-context"`
}

// CLIConfigManager manages the path of the CLI config
type CLIConfigManager struct {
	CLIConfigPath string
	viper         *viper.Viper
}

var mgr *CLIConfigManager

// NewCLIConfigManager creates a new CLIConfigManager
func NewCLIConfigManager(cfgFile string) *CLIConfigManager {
	if mgr != nil {
		return mgr
	}
	if cfgFile == "" {
		var err error
		cfgFile, err = GetKeptnDefaultConfigPath()
		if err != nil {
			log.Fatal("could not get default config path", err)
		}

		logging.PrintLog(fmt.Sprintf("no value for --config-file was specified\n defaulting to: %s\n", cfgFile), logging.VerboseLevel)

	}

	mgr = &CLIConfigManager{}
	mgr.viper = viper.GetViper()
	mgr.viper.SetConfigFile(cfgFile)
	mgr.viper.SetConfigType("json")

	mgr.CLIConfigPath = cfgFile
	return mgr
}

// LoadCLIConfig loads the configuration from file
func (c *CLIConfigManager) LoadCLIConfig() (CLIConfig, error) {
	cliConfig := CLIConfig{AutomaticVersionCheck: true, KubeContextCheck: true}
	if !fileutils.FileExists(c.CLIConfigPath) {
		return cliConfig, nil
	}

	// Allow setting file path after CLIConfigManager is initialized
	c.viper.SetConfigFile(c.CLIConfigPath)
	if err := c.viper.ReadInConfig(); err != nil {
		return cliConfig, fmt.Errorf("error when reading config file: %w", err)
	}

	if err := c.viper.Unmarshal(&cliConfig, func(dConfig *mapstructure.DecoderConfig) {
		dConfig.DecodeHook = mapstructure.StringToTimeHookFunc(time.RFC3339)
	}); err != nil {
		return cliConfig, fmt.Errorf("error when unmarshalling config file: %w", err)
	}

	return cliConfig, nil
}

// GetCLIConfig gets the already loaded configuration
func (c *CLIConfigManager) GetCLIConfig() (CLIConfig, error) {
	cliConfig := CLIConfig{}

	if err := c.viper.Unmarshal(&cliConfig, func(dConfig *mapstructure.DecoderConfig) {
		dConfig.DecodeHook = mapstructure.StringToTimeHookFunc(time.RFC3339)
	}); err != nil {
		return cliConfig, fmt.Errorf("error when unmarshalling config file: %w", err)
	}

	return cliConfig, nil
}

// StoreCLIConfig stores the configuration into the file
func (c *CLIConfigManager) StoreCLIConfig(config CLIConfig) error {

	// Note 1: json.Marshal works better than converting the CLIConfig
	// to map[string]interface{} (used in Viper.MergeConfigMap).
	// json.Marshal does not need mapstructure DecodeHook
	// e.g., if json.Marshal is not used, we have to convert the
	// the CLIConfig struct to map[string]interface{}
	// this can be done with mapstructure.Decode(<config-struct-as-i/p>, <map-variable-as-o/p>)
	// but it does not convert time.Time to string (viper stores time as string internally)
	// This mismatch causes Viper.MergeConfigMap() to throw an error
	// Note 2: We are using json.Marshal because we have set viper config type to json
	// Please change this if we are using a different config type
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("error when marshalling config file: %w", err)
	}

	if err := c.viper.MergeConfig(bytes.NewReader(data)); err != nil {
		return fmt.Errorf("error when merging new config with the existing one: %w", err)
	}

	// Writes in os.FileMode(0644) by default
	// https://github.com/spf13/viper/blob/b1fdc47b0d05b6af898a3d50aefd62c5825a17fe/viper.go#L270
	if err := c.viper.WriteConfigAs(c.CLIConfigPath); err != nil {
		return fmt.Errorf("error when writing config file: %w", err)
	}

	return nil
}

// GetKeptnDefaultConfigPath returns default Keptn Config file path
func GetKeptnDefaultConfigPath() (string, error) {
	dir, err := GetKeptnDirectory()
	if err != nil {
		return "", err
	}
	return dir + "config", nil
}

// GetKeptnDirectory returns a path, which is used to store logs and possibly creds
func GetKeptnDirectory() (string, error) {

	keptnDir := fileutils.UserHomeDir() + string(os.PathSeparator) + keptnFolderName + string(os.PathSeparator)

	if _, err := os.Stat(keptnDir); os.IsNotExist(err) {
		err := os.MkdirAll(keptnDir, os.ModePerm)
		fmt.Println("keptn creates the folder " + keptnDir + " to store logs and possibly creds.")
		if err != nil {
			return "", err
		}
	}

	return keptnDir, nil
}
