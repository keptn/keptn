package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

type DockerImg struct {
	Image string
	Tag   string
}

// TestNewArtifact tests the new-artifact command.
func TestDelivery(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	receivedEvent := make(chan bool)
	mocking = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if strings.Contains(r.RequestURI, "shipyard.yaml") {
				w.Write([]byte(shipyardResourceMockResponse))
				return
			} else if strings.Contains(r.RequestURI, "v1/event") {
				defer r.Body.Close()
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("could not read received event payload: %s", err.Error())
				}
				event := &apimodels.KeptnContextExtendedCE{}
				if err := json.Unmarshal(bytes, event); err != nil {
					t.Errorf("could not decode received event: %s", err.Error())
				}
				if *event.Type != keptnv2.GetTriggeredEventType("dev.artifact-delivery") {
					t.Errorf("did not receive correct event: %s", err.Error())
				}
				go func() {
					receivedEvent <- true
				}()
			} else if strings.Contains(r.RequestURI, "/v1/metadata") {
				defer r.Body.Close()
				w.Write([]byte(metadataMockResponse))
				return
			}
			return
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	cmd := fmt.Sprintf("trigger delivery --project=%s --service=%s --stage=%s --sequence=%s "+
		"--image=%s --tag=%s --values=a.b.c=d --mock --values=c.d=e", "sockshop", "carts", "dev", "artifact-delivery", "docker.io/keptnexamples/carts", "0.9.1")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	select {
	case <-receivedEvent:
		t.Log("event has been sent successfully")
		break
	case <-time.After(5 * time.Second):
		t.Error("event was not sent")
	}
}

func TestCheckImageAvailabilityD(t *testing.T) {

	validImgs := []DockerImg{{"docker.io/keptnexamples/carts", "0.7.0"},
		{"docker.io/keptnexamples/carts:0.7.0", ""},
		{"keptnexamples/carts", ""},
		{"keptnexamples/carts", "0.7.0"},
		{"keptnexamples/carts:0.7.0", ""},
		{"127.0.0.1:10/keptnexamples/carts", "0.7.5"},
		{"127.0.0.1:10/keptnexamples/carts:0.7.5", ""},
		{"httpd", ""}}

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	for _, validImg := range validImgs {
		*delivery.Project = "sockshop"
		*delivery.Service = "carts"
		*delivery.Image = validImg.Image
		*delivery.Tag = validImg.Tag

		err := triggerDeliveryCmd.PreRunE(triggerDeliveryCmd, []string{})

		if err != nil {
			t.Errorf(unexpectedErrMsg, err)
		}
	}
}

func TestCheckImageNonAvailabilityD(t *testing.T) {

	invalidImgs := []DockerImg{{"docker.io/keptnexamples/carts:0.7.5", ""}}

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	for _, validImg := range invalidImgs {
		*delivery.Project = "sockshop"
		*delivery.Service = "carts"
		*delivery.Image = validImg.Image
		*delivery.Tag = validImg.Tag

		err := triggerDeliveryCmd.PreRunE(triggerDeliveryCmd, []string{})

		Expected := "Provided image not found: Tag not found"
		if err == nil || err.Error() != Expected {
			t.Errorf("Error actual = %v, and Expected = %v.", err, Expected)
		}
	}
}
