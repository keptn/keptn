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

// TestCreateProjectCmd tests the default use of the update project command
func TestUpdateProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	cmd := fmt.Sprintf("update project sockshop -t token -u user -r https:// --mock")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestUpdateProjectIncorrectProjectNameCmd tests whether the update project command aborts
// due to a project name with upper case character
func TestUpdateProjectIncorrectProjectNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	cmd := fmt.Sprintf("update project Sockshop -t token -u user -r https://github.com/user/upstream.git --mock")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "contains upper case letter(s) or special character(s)") {
		t.Errorf("missing expected error, but got %v", err)
	}
}
