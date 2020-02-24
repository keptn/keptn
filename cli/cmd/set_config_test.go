package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/keptn/keptn/cli/utils/config"

	"github.com/keptn/keptn/cli/pkg/logging"
	"gotest.tools/assert"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestSetCmdRunError(t *testing.T) {

	configMng = config.NewCLIConfigManager()
	err := setConfigCmd.RunE(nil, []string{"automaticversionchek", "true"})
	assert.Equal(t, err.Error(), "Unsupported key automaticversionchek", "Wrong error")
}

const testConfig = `{"automatic_version_check":true,"last_version_check":"2020-02-20T00:00:00Z"}`

func TestSetAutomaticVersionCheckFalse(t *testing.T) {

	configMng = config.NewCLIConfigManager()
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	defer os.RemoveAll(tmpDir)

	configMng.CLIConfigPath = filepath.Join(tmpDir, "config")
	ioutil.WriteFile(configMng.CLIConfigPath, []byte(testConfig), 0644)

	err = setConfigCmd.RunE(nil, []string{"automaticversioncheck", "false"})
	assert.Equal(t, err, nil, "Wrong error")

	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if cliConfig.AutomaticVersionCheck {
		t.Errorf("Automatic version check has to be false")
	}
}
