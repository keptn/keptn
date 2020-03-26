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

func TestConfigureMonitoringCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring prometheus --project=%s --service=%s --mock", "sockshop", "carts")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForPrometheus(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring prometheus --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Please specify a service" {
		t.Errorf(unexpectedErrMsg, err)
	}
}
