package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func newArtifactEvent(t *testing.T, eventFileName string, eventContent string) func() {
	time := types.Timestamp{Time: time.Now()}

	if eventContent == "" {
		eventContent = `{"id":"` + uuid.New().String() + `",
  "type":"sh.keptn.event.configuration.change",
  "specversion":"0.2",
  "source": "https://github.com/keptn/keptn/cli",
  "time":"` + time.String() + `",
  "contenttype":"application/json",
  "data":{
	   "project":"sockshop",
	   "service":"carts",
	   "valuesCanary": {
		   "image": "docker.io/keptnexamples/carts:0.9.1"
	  }
  }
}`
	}

	ioutil.WriteFile(eventFileName, []byte(eventContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	return func() {
		os.Remove(eventFileName)
	}
}

// TestSendEvent tests the functionality to send an event defined in JSON file.
func TestSendEvent(t *testing.T) {
	const tmpCE = "newArtifactCE.json"
	defer newArtifactEvent(t, tmpCE, "")()

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"send",
		"event",
		fmt.Sprintf("--file=%s", tmpCE),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}
}
