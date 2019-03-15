package credentialmanager

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
)

func init() {
	utils.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestSetAndGetCreds(t *testing.T) {

	if err := SetCreds(testEndPoint, testAPIToken); err != nil {
		t.Fatal(err)
	}

	endPoint, apiToken, err := GetCreds()
	if err != nil {
		t.Fatal(err)
	}
	if testEndPoint != endPoint || testAPIToken != apiToken {
		utils.Info.Printf("Expected endoint is %v but was %v", testEndPoint, endPoint)
		utils.Info.Printf("Expected secret is %v but was %v", testAPIToken, apiToken)
		t.Fatal("Readed creds do not match")
	}
}
