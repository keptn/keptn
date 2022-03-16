//go:build !nokubectl
// +build !nokubectl

package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestGenerateSupportArchiveDirectoryDoesNotExist(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("generate support-archive --dir=%s", "does/not/exist")
	_, err := executeActionCommandC(cmd)
	assert.Equal(t, err.Error(), "Error trying to access directory does/not/exist. Please make sure the directory exists.", "Received unexpected error")
}

// TestGenerateSupportArchive tests the generation of a support archive
func TestGenerateSupportArchive(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	// create tempo directory
	dname := t.TempDir()

	cmd := fmt.Sprintf("generate support-archive --dir=%s --mock", dname)
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestGenerateSupportArchiveUnknownCommand
func TestGenerateSupportArchiveUnknownCommand(t *testing.T) {
	testInvalidInputHelper("generate support-archive someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn generate support-archive\"", t)
}

// TestGenerateSupportArchiveUnknownParameter
func TestGenerateSupportArchiveUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("generate support-archive --project=sockshop", "unknown flag: --project", t)
}
