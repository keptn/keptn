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

// TestConfigureMonitoringUnknownCommand
func TestConfigureMonitoringUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("configure monitoring prometheus someUnknownCommand --project=sockshop --service=helloservice")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "Requires a monitoring provider as argument"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestConfigureMonitoringUnknownParameter
func TestConfigureMonitoringUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("configure monitoring prometheus --projectt=sockshop --service=helloservice")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown flag: --projectt"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
