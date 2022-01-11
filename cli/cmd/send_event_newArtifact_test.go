package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/keptn/cli/pkg/logging"
)

const shipyardResourceMockResponse = `{
      "resourceContent": "YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjAiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXNvY2tzaG9wIgpzcGVjOgogIHN0YWdlczoKICAgIC0gbmFtZTogImRldiIKICAgICAgc2VxdWVuY2VzOgogICAgICAgIC0gbmFtZTogImFydGlmYWN0LWRlbGl2ZXJ5IgogICAgICAgICAgdGFza3M6CiAgICAgICAgICAgIC0gbmFtZTogImRlcGxveW1lbnQiCiAgICAgICAgICAgICAgcHJvcGVydGllczoKICAgICAgICAgICAgICAgIGRlcGxveW1lbnRzdHJhdGVneTogImRpcmVjdCIKICAgICAgICAgICAgLSBuYW1lOiAidGVzdCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgdGVzdHN0cmF0ZWd5OiAiZnVuY3Rpb25hbCIKICAgICAgICAgICAgLSBuYW1lOiAiZXZhbHVhdGlvbiIKICAgICAgICAgICAgLSBuYW1lOiAicmVsZWFzZSIKICAgICAgICAtIG5hbWU6ICJhcnRpZmFjdC1kZWxpdmVyeS1kYiIKICAgICAgICAgIHRhc2tzOgogICAgICAgICAgICAtIG5hbWU6ICJkZXBsb3ltZW50IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBkZXBsb3ltZW50c3RyYXRlZ3k6ICJkaXJlY3QiCiAgICAgICAgICAgIC0gbmFtZTogInJlbGVhc2Ui",
      "resourceURI": "shipyard.yaml"
}`

const metadataMockResponse = `{"bridgeversion":"v1","keptnlabel":"keptn","keptnversion":"0.8.0","namespace":"keptn"}`

const serviceMockResponse = `{
  "nextPageKey": "0",
  "services": [
    {
      "creationDate": "1638796449959720167",
      "deployedImage": "ghcr.io/podtato-head/podtatoserver:v0.1.2",
      "lastEventTypes": {},
      "openRemediations": null,
      "serviceName": "carts"
    }
  ],
  "totalCount": 1
}`

const projectMockResponse = `{
	"creationDate": "1638796448951137480",
	"projectName": "sockshop",
	"shipyard": "",
	"shipyardVersion": "spec.keptn.sh/0.2.0",
	"stages": [
	  {
		"services": [
		  {
			"creationDate": "1638796449959720167",
			"deployedImage": "podtatoserver:v0.1.2",
			"openRemediations": null,
			"serviceName": "carts"
		  }
		],
		"stageName": "dev"
	  }
	]
}`

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestNewArtifact tests the new-artifact command.
func TestNewArtifact(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	receivedEvent := make(chan bool)
	//mocking = true
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
			} else if strings.Contains(r.RequestURI, "/service") {
				defer r.Body.Close()
				w.Write([]byte(serviceMockResponse))
				return
			} else if strings.Contains(r.RequestURI, "/project") {
				defer r.Body.Close()
				w.Write([]byte(projectMockResponse))
				return
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
