package cmd

import (
	"fmt"
	"os"
	"testing"
)

// TestGenerateKeptnService tests the service template generation of keptn
func TestGenerateKeptnService(t *testing.T) {
	cmd := fmt.Sprintf("generate keptn-service --service=%s --image=%s --events=sh.keptn.events.configuration-changed", "tempService", "tempImage")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
	defer os.RemoveAll("tempService")
}
