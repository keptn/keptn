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

func TestConfigureDomainNoParamCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("configure domain")
	_, err := executeActionCommandC(cmd)
	if err == nil || err.Error() != "Requires a domain as argument" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestConfigureDomainCmdWithVersion(t *testing.T) {

	cmd := fmt.Sprintf("configure domain my.keptn.domain.com --mock")

	r := newRedirector()
	r.redirectStdOut()

	Version = "0.6.1"
	*configureDomainParams.ConfigVersion = ""
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	out := r.revertStdOut()

	if !strings.Contains(out, "Used version for manifests: 0.6.1") {
		t.Errorf("unexpected used version: %s", out)
	}
}

func TestConfigureDomainCmdWithoutVersion(t *testing.T) {

	cmd := fmt.Sprintf("configure domain my.keptn.domain.com --mock")

	r := newRedirector()
	r.redirectStdOut()

	Version = ""
	*configureDomainParams.ConfigVersion = ""
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	out := r.revertStdOut()

	if !strings.Contains(out, "Used version for manifests: master") {
		t.Errorf("unexpected used version: %s", out)
	}
}
