package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestCreateProjectCmd tests the default use of the create project command
func TestCreateServiceCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("create service carts --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}
}
