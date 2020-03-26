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

const ce = `{
  "contenttype": "application/json",
  "data": {
    "deploymentURILocal": "http://carts-canary.sockshop-production",
    "deploymentstrategy": "blue_green_service",
    "eventContext": null,
    "project": "sockshop",
    "service": "carts",
    "stage": "production",
    "teststrategy": "performance"
  },
  "id": "fb02d165-53ce-4a28-80b0-62d7fc0c5087",
  "source": "https://github.com/keptn/keptn/api",
  "specversion": "0.2",
  "time": "2020-03-06T10:48:51.254Z",
  "type": "sh.keptn.events.deployment-finished",
  "shkeptncontext": "5403dc38-dc42-4218-a587-1b5973ac32fc"
}`

// TestOnboardServiceWrongHelmChartPath tests the onboard service command.
func TestSendEvent(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	resourceFileName := "ce.json"
	defer testResource(t, resourceFileName, ce)()

	cmd := fmt.Sprintf("send event --file=%s --mock", resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}
