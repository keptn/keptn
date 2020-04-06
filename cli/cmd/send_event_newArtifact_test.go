package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestNewArtifact tests the new-artifact command.
func TestNewArtifact(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("send event new-artifact --project=%s --service=%s "+
		"--image=%s --tag=%s  --mock", "sockshop", "carts", "docker.io/keptnexamples/carts", "0.9.1")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

type DockerImg struct {
	Image string
	Tag   string
}

func TestCheckImageAvailability(t *testing.T) {

	validImgs := []DockerImg{{"docker.io/keptnexamples/carts", "0.7.0"},
		{"docker.io/keptnexamples/carts:0.7.0", ""},
		{"keptnexamples/carts", ""},
		{"keptnexamples/carts", "0.7.0"},
		{"keptnexamples/carts:0.7.0", ""},
		{"127.0.0.1:10/keptnexamples/carts", "0.7.5"},
		{"127.0.0.1:10/keptnexamples/carts:0.7.5", ""},
		{"httpd", ""}}

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	for _, validImg := range validImgs {
		*newArtifact.Project = "sockshop"
		*newArtifact.Service = "carts"
		*newArtifact.Image = validImg.Image
		*newArtifact.Tag = validImg.Tag

		err := newArtifactCmd.PreRunE(newArtifactCmd, []string{})

		if err != nil {
			t.Errorf(unexpectedErrMsg, err)
		}
	}
}

func TestCheckImageNonAvailability(t *testing.T) {

	invalidImgs := []DockerImg{{"docker.io/keptnexamples/carts:0.7.5", ""}}

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	for _, validImg := range invalidImgs {
		*newArtifact.Project = "sockshop"
		*newArtifact.Service = "carts"
		*newArtifact.Image = validImg.Image
		*newArtifact.Tag = validImg.Tag

		err := newArtifactCmd.PreRunE(newArtifactCmd, []string{})

		Expected := "Provided image not found: Tag not found"
		if err == nil || err.Error() != Expected {
			t.Errorf("Error actual = %v, and Expected = %v.", err, Expected)
		}
	}
}
