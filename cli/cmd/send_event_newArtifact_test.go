package cmd

import (
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

const shipyardMockResponseContent = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "artifact-delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "artifact-delivery-db"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"`

const shipyardResourceMockResponse = `{
      "resourceContent": "YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjAiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXNvY2tzaG9wIgpzcGVjOgogIHN0YWdlczoKICAgIC0gbmFtZTogImRldiIKICAgICAgc2VxdWVuY2VzOgogICAgICAgIC0gbmFtZTogImFydGlmYWN0LWRlbGl2ZXJ5IgogICAgICAgICAgdGFza3M6CiAgICAgICAgICAgIC0gbmFtZTogImRlcGxveW1lbnQiCiAgICAgICAgICAgICAgcHJvcGVydGllczoKICAgICAgICAgICAgICAgIGRlcGxveW1lbnRzdHJhdGVneTogImRpcmVjdCIKICAgICAgICAgICAgLSBuYW1lOiAidGVzdCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgdGVzdHN0cmF0ZWd5OiAiZnVuY3Rpb25hbCIKICAgICAgICAgICAgLSBuYW1lOiAiZXZhbHVhdGlvbiIKICAgICAgICAgICAgLSBuYW1lOiAicmVsZWFzZSIKICAgICAgICAtIG5hbWU6ICJhcnRpZmFjdC1kZWxpdmVyeS1kYiIKICAgICAgICAgIHRhc2tzOgogICAgICAgICAgICAtIG5hbWU6ICJkZXBsb3ltZW50IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBkZXBsb3ltZW50c3RyYXRlZ3k6ICJkaXJlY3QiCiAgICAgICAgICAgIC0gbmFtZTogInJlbGVhc2Ui",
      "resourceURI": "shipyard.yaml"
}`

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestNewArtifact tests the new-artifact command.
func TestNewArtifact(t *testing.T) {

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
			}
			return
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	cmd := fmt.Sprintf("send event new-artifact --project=%s --service=%s --stage=%s --sequence=%s "+
		"--image=%s --tag=%s  --mock", "sockshop", "carts", "dev", "artifact-delivery", "docker.io/keptnexamples/carts", "0.9.1")
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
