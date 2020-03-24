package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestConfigureDomainNoParamCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("configure domain")
	_, err := executeActionCommandC(cmd)
	if err == nil || err.Error() != "Requires a domain as argument" {
		t.Errorf("unexpected error, got '%v'", err)
	}
}

func TestConfigureDomainCmdWithVersion(t *testing.T) {

	cmd := fmt.Sprintf("configure domain my.keptn.domain.com --mock")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Version = "0.6.1"
	*configureDomainParams.ConfigVersion = ""
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !strings.Contains(out, "Used version for manifests: 0.6.1") {
		t.Errorf("unexpected used version: %s", out)
	}
}

func TestConfigureDomainCmdWithoutVersion(t *testing.T) {

	cmd := fmt.Sprintf("configure domain my.keptn.domain.com --mock")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Version = ""
	*configureDomainParams.ConfigVersion = ""
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if !strings.Contains(out, "Used version for manifests: master") {
		t.Errorf("unexpected used version: %s", out)
	}
}
