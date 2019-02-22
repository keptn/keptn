package credentialmanager

import (
	"fmt"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestZalandoKeyring(t *testing.T) {

	fmt.Println("Start...")
	service := "my-app"
	user := "anon"
	password := "secret"

	// set password
	err := keyring.Set(service, user, password)
	if err != nil {
		t.Fatal(err)
	}

	// get password
	secret, err := keyring.Get(service, user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(secret)
}
