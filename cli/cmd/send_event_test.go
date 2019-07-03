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
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestNewArtifact tests the new-artifact command.
func TestSend(t *testing.T) {

	time := types.Timestamp{Time: time.Now()}

	newArtifactEvent := `{"id":"` + uuid.New().String() + `",
"specversion":"0.2",
"time":"` + time.String() + `",
"contenttype":"application/json",
"type":"sh.keptn.events.new-artifact",
"data":{
	"project":"sockshop",
	"service":"carts",
	"image":"docker.io/keptnexamples/carts",
	"tag":"0.6.0.latest"
}}`

	const tmpCE = "ce.json"
	ioutil.WriteFile(tmpCE, []byte(newArtifactEvent), 0644)

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

	os.Remove(tmpCE)

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
