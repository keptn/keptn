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

// TestOnboardServiceWrongHelmChartPath tests the onboard service command.
func TestOnboardServiceWrongHelmChartPath(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	cmd := fmt.Sprintf("onboard service carts --project=sockshop --chart=cartsX")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected error event, but no one received.")
	}

	expected := "Provided Helm chart does not exist"
	if err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
	}
}

// TestOnboardServiceDeploymentStrategy tests the onboard service command.
func TestOnboardServiceDeploymentStrategy(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	cmd := fmt.Sprintf("onboard service carts --project=sockshop --deployment-strategy=directX")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected error event, but no one received.")
	}

	expected := "The provided deployment strategy is not supported. Select: [direct|blue_green_service]"
	if err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
	}
}
