package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestUpdateProjectCmd tests the default use of the update project command
func TestUpdateProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("update project sockshop --git-token=token --git-remote-url=https://some.url --mock")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestUpdateProjectCmd tests the use of the update project command with a git user set up
func TestUpdateProjectCmdWithUser(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("update project sockshop --git-user=user --git-token=token --git-remote-url=https://some.url --mock")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestUpdateProjectIncorrectProjectNameCmd tests whether the update project command aborts
// due to a project name with upper case character
func TestUpdateProjectIncorrectProjectNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("update project Sockshop --git-token=token --git-user=user --git-remote-url=https://github.com/user/upstream.git --mock")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "contains upper case letter(s) or special character(s)") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestUpdateProjectUnknownCommand
func TestUpdateProjectUnknownCommand(t *testing.T) {
	testInvalidInputHelper("update project sockshop someUnknownCommand --git-user=user --git-token=token --git-remote-url=http://some.url", "too many arguments set", t)
}

// TestUpdateProjectUnknownParameter
func TestUpdateProjectUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("update project sockshop --git-userr=user --git-token=token --git-remote-url=http://some.url", "unknown flag: --git-userr", t)
}

// TestUpdateProjectCmdTokenAndKey
func TestUpdateProjectCmdTokenAndKey(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("update project sockshop --git-user=user --git-remote-url=https://someurl.com --mock --git-private-key=key --git-token=token")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "Access token and private key cannot be set together") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// IMPORTANT NOTE: tests below are disabled due to broken cli, which is unrepairable adn needs to be rewritten
// TestUpdateProjectCmdProxyAndSSH
// func TestUpdateProjectCmdProxyAndSSH(t *testing.T) {
// 	credentialmanager.MockAuthCreds = true

// 	cmd := fmt.Sprintf("update project sockshop --git-user=user --git-remote-url=ssh://someurl.com --mock --git-private-key=key --git-proxy-url=ip-address")
// 	_, err := executeActionCommandC(cmd)

// 	if !errorContains(err, "Proxy cannot be set with SSH") {
// 		t.Errorf("missing expected error, but got %v", err)
// 	}
// }

// // TestUpdateProjectCmdProxyNoScheme
// func TestUpdateProjectCmdProxyNoScheme(t *testing.T) {
// 	credentialmanager.MockAuthCreds = true

// 	cmd := fmt.Sprintf("update project sockshop --git-user=user --git-remote-url=https://someurl.com --mock --git-token=token --git-proxy-url=ip-address")
// 	_, err := executeActionCommandC(cmd)

// 	if !errorContains(err, "Proxy cannot be set without scheme") {
// 		t.Errorf("missing expected error, but got %v", err)
// 	}
// }
