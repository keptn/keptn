package cmd

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestAuthCmd tests the auth command.
func TestAuthCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	endPoint, apiToken, err := credentialmanager.NewCredentialManager(false).GetCreds(namespace)
	if err != nil {
		t.Error(err)
		return
	}

	cmd := fmt.Sprintf("auth --endpoint=%s --api-token=%s --mock", endPoint.String(), apiToken)
	_, err = executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestLookupHostname(t *testing.T) {
	var tests = []struct {
		in  string
		out bool
	}{
		{"xip.io", true},
		{"127.0.0.1.xip.io", true},
		{"127.0.0.2.xip.io", true},
		{"192.168.0.0.xip.io", true},
		{"api.keptn.192.168.0.0.xip.io", true},
		{"a.b.c.d", false},
		{"test.com", true},
		{"keptn.github.io", true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			s := lookupHostname(tt.in, mockedHostnameResolveFn, mockedSleepFn)
			if s != tt.out {
				t.Errorf("lookupHostname(%s): got %v, want %v", tt.in, s, tt.out)
			}
		})
	}
}

func mockedHostnameResolveFn(hostname string) ([]string, error) {
	if hostname == "a.b.c.d" {
		return []string{}, errors.New("Unable to resolve " + hostname)
	}

	return []string{"0.0.0.0"}, nil
}

func mockedSleepFn(d time.Duration) {
	//no-op
}
