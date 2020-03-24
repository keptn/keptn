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
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"send",
		"event",
		"new-artifact",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--image=%s", "docker.io/keptnexamples/carts"),
		fmt.Sprintf("--tag=%s", "0.9.1"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}
}

type DockerImg struct {
	Image string
	Tag   string
}

func TestCheckImageAvailability(t *testing.T) {

	validImgs := []DockerImg{DockerImg{"docker.io/keptnexamples/carts", "0.7.0"},
		DockerImg{"docker.io/keptnexamples/carts:0.7.0", ""},
		DockerImg{"keptnexamples/carts", ""},
		DockerImg{"keptnexamples/carts", "0.7.0"},
		DockerImg{"keptnexamples/carts:0.7.0", ""},
		DockerImg{"10.10.10.10:10/keptnexamples/carts", "0.7.5"},
		DockerImg{"10.10.10.10:10/keptnexamples/carts:0.7.5", ""},
		DockerImg{"httpd", ""}}

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
			t.Errorf("unexpected error, got '%v'", err)
		}
	}
}

func TestCheckImageNonAvailability(t *testing.T) {

	invalidImgs := []DockerImg{DockerImg{"docker.io/keptnexamples/carts:0.7.5", ""}}

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
