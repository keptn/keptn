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

func TestDeleteServiceCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("delete service %s --project=%s --mock", "carts", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}
