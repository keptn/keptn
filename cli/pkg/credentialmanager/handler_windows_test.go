package credentialmanager

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestSetAndGetCreds(t *testing.T) {

	cm := NewCredentialManager()
	if err := cm.SetCreds(testEndPoint, testAPIToken, ""); err != nil {
		t.Fatal(err)
	}

	endPoint, apiToken, err := cm.GetCreds("")
	if err != nil {
		t.Fatal(err)
	}
	if testEndPoint != endPoint || testAPIToken != apiToken {
		t.Fatal("Readed creds do not match")
	}
}
