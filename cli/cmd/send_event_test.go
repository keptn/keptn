package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestSendEvent tests the functionality to send an event defined in JSON file.
func TestSendEvent(t *testing.T) {

	time := types.Timestamp{Time: time.Now()}

	newArtifactEvent := `{"id":"` + uuid.New().String() + `",
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

	const tmpCE = "ce.json"
	err := ioutil.WriteFile(tmpCE, []byte(newArtifactEvent), 0644)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

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
	err = rootCmd.Execute()

	os.Remove(tmpCE)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
}

// TestSendEventAndOpenWebSocket tests the functionality to send an event defined in JSON file and to open a WebSocket communication
func TestSendEventAndOpenWebSocket(t *testing.T) {

	time := types.Timestamp{Time: time.Now()}

	newArtifactEvent := `{"id":"` + uuid.New().String() + `",
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

	const tmpCE = "ce.json"
	err := ioutil.WriteFile(tmpCE, []byte(newArtifactEvent), 0644)
	if err != nil {
		log.Fatalf("An error occured %v", err)
	}

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"send",
		"event",
		fmt.Sprintf("--file=%s", tmpCE),
		"--stream-websocket",
		"--mock",
	}
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	os.Remove(tmpCE)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
}
