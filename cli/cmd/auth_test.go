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

// TestAuthCmd tests the auth command.
func TestAuthCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	endPoint, apiToken, err := credentialmanager.NewCredentialManager().GetCreds()
	if err != nil {
		t.Error(err)
		return
	}

	cmd := fmt.Sprintf("auth --endpoint=%s --api-token=%s --mock", endPoint.String(), apiToken)
	_, err = executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}
