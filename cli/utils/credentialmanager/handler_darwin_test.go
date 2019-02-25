package credentialmanager

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
)

const testCred = "super-secret"

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestSetAndGetCreds(t *testing.T) {

	if err := SetCreds(testCred); err != nil {
		t.Fatal(err)
	}

	secret, err := GetCreds()
	if err != nil {
		t.Fatal(err)
	}
	if testCred != secret {
		t.Fatal("Readed creds do not match")
	}
}

func TestOverwriteCreds(t *testing.T) {

	if err := SetCreds("old-secret"); err != nil {
		t.Fatal(err)
	}

	if err := SetCreds(testCred); err != nil {
		t.Fatal(err)
	}

	secret, err := GetCreds()
	if err != nil {
		t.Fatal(err)
	}
	if testCred != secret {
		t.Fatal("Readed creds do not match")
	}
}
