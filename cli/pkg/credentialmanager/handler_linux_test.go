package credentialmanager

import (
	"io/ioutil"
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
		logging.Info.Printf("Expected endoint is %v but was %v", testEndPoint, endPoint)
		logging.Info.Printf("Expected secret is %v but was %v", testAPIToken, apiToken)
		t.Fatal("Readed creds do not match")
	}
}

func TestGetCredsFromFile(t *testing.T) {

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Note that this is no real domain nor token testAPIToken is only used for testing purpose
	content := testEndPoint.String() + "\n" + testAPIToken
	_, err = file.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}

	cm := NewCredentialManager()
	cm.apiTokenFile = file.Name()

	url, token, err := cm.GetCreds("")
	if err != nil {
		t.Fatal(err)
	}
	if url != testEndPoint {
		t.Fatal("URLs do not match")
	}
	if testAPIToken != token {
		t.Fatal("API tokens do not match")
	}
}
