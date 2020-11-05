package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/keptn/keptn/cli/pkg/file"
)

const testConfig = `{"automatic_version_check":true,"last_version_check":"2020-02-20T00:00:00Z","current-context":""}`

var testTime time.Time

func init() {
	testTime = time.Date(2020, time.February, 20, 0, 0, 0, 0, time.UTC)
}

func TestLoadNonExistingCLIConfig(t *testing.T) {

	mng := NewCLIConfigManager()
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	mng.CLIConfigPath = filepath.Join(tmpDir, "config")

	cliConfig, err := mng.LoadCLIConfig()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if !cliConfig.AutomaticVersionCheck {
		t.Errorf("Unexpected value of AutomaticVersionCheck")
	}
	if cliConfig.LastVersionCheck != nil {
		t.Errorf("Unexpected value of LastVersionCheck")
	}
}

func TestStoreCLIConfig(t *testing.T) {

	mng := NewCLIConfigManager()
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	mng.CLIConfigPath = filepath.Join(tmpDir, "config")

	cliConfig := CLIConfig{AutomaticVersionCheck: true, LastVersionCheck: &testTime}

	err = mng.StoreCLIConfig(cliConfig)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	data, err := file.ReadFile(mng.CLIConfigPath)
	if data != testConfig {
		t.Errorf("Different config stored")
	}
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
}

func TestLoadCLIConfig(t *testing.T) {

	mng := NewCLIConfigManager()
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	mng.CLIConfigPath = filepath.Join(tmpDir, "config")
	ioutil.WriteFile(mng.CLIConfigPath, []byte(testConfig), 0644)

	cliConfig, err := mng.LoadCLIConfig()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if !cliConfig.AutomaticVersionCheck {
		t.Errorf("Different config read")
	}
	if cliConfig.LastVersionCheck == nil || *cliConfig.LastVersionCheck != testTime {
		t.Errorf("Different config read")
	}
}
