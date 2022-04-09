package cmd

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
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
			s := LookupHostname(tt.in, mockedHostnameResolveFn, mockedSleepFn)
			if s != tt.out {
				t.Errorf("lookupHostname(%s): got %v, want %v", tt.in, s, tt.out)
			}
		})
	}
}

func TestSmartFetchKeptnAuthParameters(t *testing.T) {
	var endPoint = "keptn.github.io"
	var apiToken = "someApiToken"
	var falseValue = false
	var authParams = authCmdParams{
		endPoint:             &endPoint,
		apiToken:             &apiToken,
		exportConfig:         &falseValue,
		acceptContext:        true,
		secure:               &falseValue,
		skipNamespaceListing: &falseValue,
	}

	var smartKeptnAuth = smartKeptnAuthParams{
		ingressName:    "api-keptn-ingress",
		serviceName:    "api-gateway-nginx",
		secretName:     "keptn-api-token",
		insecurePrefix: "http://",
	}

	err := smartFetchKeptnAuthParameters(&authParams, smartKeptnAuth)
	if err != nil {
		t.Errorf("TestSmartFetchKeptnAuthParameters: %v", err)
	}

	if !strings.HasPrefix(*authParams.endPoint, "http://") {
		t.Errorf("TestSmartFetchKeptnAuthParameters: endpoint %s does not have the required http prefix", *authParams.endPoint)
	}
}

func TestAddCorrectHttpPrefix(t *testing.T) {
	var falseValue = false
	var trueValue = true
	var endpoint = []string{"http://some.url", "https://some.url", "some.url"}
	var tests = []struct {
		in  authCmdParams
		out string
	}{
		{authCmdParams{secure: &trueValue, endPoint: &endpoint[0]}, "https://some.url"},
		{authCmdParams{secure: &falseValue, endPoint: &endpoint[0]}, "http://some.url"},
		{authCmdParams{secure: &trueValue, endPoint: &endpoint[1]}, "https://some.url"},
		{authCmdParams{secure: &falseValue, endPoint: &endpoint[1]}, "https://some.url"},
		{authCmdParams{secure: &trueValue, endPoint: &endpoint[2]}, "https://some.url"},
		{authCmdParams{secure: &falseValue, endPoint: &endpoint[2]}, "http://some.url"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := addCorrectHTTPPrefix(&tt.in)
			if s != tt.out {
				t.Errorf("addCorrectHTTPPrefix(): got %s, want %s", s, tt.out)
			}
		})
	}
}

// TestAuthUnknownCommand
func TestAuthUnknownCommand(t *testing.T) {
	testInvalidInputHelper("auth someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn auth\"", t)
}

// TestAuthUnknownParameter
func TestAuthUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("auth --projectt=sockshop", "unknown flag: --projectt", t)
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
