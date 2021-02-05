package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestGenerateDocsDirectoryDoesNotExist(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("generate docs --dir=%s --mock", "does/not/exist")
	_, err := executeActionCommandC(cmd)
	assert.Equal(t, err.Error(), "Error trying to access directory does/not/exist. Please make sure the directory exists.", "Received unexpected error")
}

// Tests generating docs in a temp directory
func TestGenerateDocs(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	// create tempo directory
	dname, err := ioutil.TempDir("", "docs_temp")
	defer os.RemoveAll(dname)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	cmd := fmt.Sprintf("generate docs --dir=%s --mock", dname)
	_, err = executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}
