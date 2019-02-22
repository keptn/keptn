package credentialmanager

import (
	"testing"
)

const testCred = "super-secret"

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
