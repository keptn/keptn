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

func TestConfigureMonitoringCmdForPrometheus(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring prometheus --project=%s --service=%s --mock", "sockshop", "carts")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForDatadog(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring datadog --project=%s --service=%s --mock", "sockshop", "carts")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForSumologic(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring sumologic --project=%s --service=%s --mock", "sockshop", "carts")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForPrometheusWithWrongArgs(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring prometheus --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Please specify a service" {
		t.Errorf(unexpectedErrMsg, err)
	}

	// Note: params are defined in configure_monitoring.go
	// params.Project and params.Service are set for all tests every time we run a test
	// which executes `configure monitoring` command.
	// We have to reset them every time before a test which runs `configure monitoring`
	*params.Project = ""
	*params.Service = ""
	cmd = fmt.Sprintf("configure monitoring prometheus --service=%s --mock", "carts")
	_, err = executeActionCommandC(cmd)
	if err.Error() != "Please specify a project" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForDatadogWithWrongArgs(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring datadog --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Please specify a service" {
		t.Errorf(unexpectedErrMsg, err)
	}

	*params.Project = ""
	*params.Service = ""
	cmd = fmt.Sprintf("configure monitoring datadog --service=%s --mock", "carts")
	_, err = executeActionCommandC(cmd)
	if err.Error() != "Please specify a project" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureMonitoringCmdForSumologicWithWrongArgs(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*params.Project = ""
	*params.Service = ""
	cmd := fmt.Sprintf("configure monitoring sumologic --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Please specify a service" {
		t.Errorf(unexpectedErrMsg, err)
	}

	*params.Project = ""
	*params.Service = ""
	cmd = fmt.Sprintf("configure monitoring sumologic --service=%s --mock", "carts")
	_, err = executeActionCommandC(cmd)
	if err.Error() != "Please specify a project" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestConfigureMonitoringUnknownCommand
func TestConfigureMonitoringUnknownCommand(t *testing.T) {
	testInvalidInputHelper("configure monitoring prometheus someUnknownCommand --project=sockshop --service=helloservice", "Requires a monitoring provider as argument", t)
}

// TestConfigureMonitoringUnknownParameter
func TestConfigureMonitoringUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("configure monitoring prometheus --projectt=sockshop --service=helloservice", "unknown flag: --projectt", t)
}
