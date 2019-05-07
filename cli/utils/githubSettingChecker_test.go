package utils

import (
	"os"
	"testing"
)

func TestOrgExists(t *testing.T) {

	token := os.Getenv("GITHUB_TOKEN_NIGHTLY")
	if token == "" {
		t.Errorf("Test failed because it requires an environment variable `GITHUB_TOKEN_NIGHTLY` for the GitHub personal access token")
		t.FailNow()
	}

	const existingOrg = "keptn-nightly"

	res, err := OrgExists(token, existingOrg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !res {
		t.Errorf("Organization " + existingOrg + " not found")
	}

	const nonExistingOrg = "kpn-nightly"
	res, err = OrgExists(token, nonExistingOrg)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if res {
		t.Errorf("Organization " + nonExistingOrg + " should not exist")
	}

	const existingOrgWithoutRights = "keptn-tiger"
	res, err = OrgExists(token, existingOrgWithoutRights)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if res {
		t.Errorf("Organization " + existingOrgWithoutRights + " must not be accessible")
	}
}

func TestCheckRepoScopeOfToken(t *testing.T) {

	token := os.Getenv("GITHUB_TOKEN_NIGHTLY")
	if token == "" {
		t.Errorf("Test failed because it requires an environment variable `GITHUB_TOKEN_NIGHTLY` for the GitHub personal access token")
		t.FailNow()
	}

	res, err := CheckRepoScopeOfToken(token)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if !res {
		t.Errorf("Used token has required rights")
	}
}
