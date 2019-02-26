package credentialmanager

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
)

const testEndPoint = "my-endpoint"
const testCred = "super-secret"

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestSetAndGetCreds(t *testing.T) {

	if err := SetCreds(testEndPoint, testCred); err != nil {
		t.Fatal(err)
	}

	endPoint, secret, err := GetCreds()s
	if err != nil {
		t.Fatal(err)
	}
	if testEndPoint != endPoint || testCred != secret {
		t.Fatal("Readed creds do not match")
	}
}
