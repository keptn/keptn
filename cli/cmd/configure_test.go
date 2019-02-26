package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestConfigureCmd(t *testing.T) {

	credentialmanager.MockCreds = true

	org, usr, tok, err := readGitTokenForTest()
	if err != nil {
		t.Error(err)
		return
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"configure",
		fmt.Sprintf("--org=%s", org),
		fmt.Sprintf("--user=%s", usr),
		fmt.Sprintf("--token=%s", tok),
	}
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}

func readGitTokenForTest() (string, string, string, error) {
	const gitTokenFile = "gitToken.txt"
	if _, err := os.Stat(gitTokenFile); os.IsNotExist(err) {
		return "", "", "", err
	}

	data, err := ioutil.ReadFile(gitTokenFile)
	if err != nil {
		return "", "", "", err
	}
	creds := strings.Split(string(data), "\n")
	if len(creds) != 3 {
		return "", "", "", errors.New("Format of gitToken.txt file is invalid")
	}
	return creds[0], creds[1], creds[2], nil
}
