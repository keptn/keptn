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

const getProjectMockResponseSequence = `{
	"creationDate": "1646503404448416664",
	"projectName": "%s",
	"shipyard": "",
	"shipyardVersion": "spec.keptn.sh/0.2.2",
	"stages": [
	  {
		"services": [
		  {
			"creationDate": "1646503568818540826",
			"openRemediations": null,
			"openApprovals": null,
			"serviceName": "%s"
		  }
		],
		"stageName": "%s"
	  }
	]
}`

const getProjectMockResponseSequenceNotFound = `{
	"code": 404,
	"message": "Project not found: %s"
}`

const getSvcMockResponseSequence = `{
	"nextPageKey": "0",
	"services": [
	  {
		"creationDate": "1638796449959720167",
		"openRemediations": null,
		"openApprovals": null,
		"serviceName": "%s"
	  }
	],
	"totalCount": 1
}`

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestTriggerSequence tests the trigger sequence command.
func TestTriggerSequence(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	receivedEvent := make(chan bool)
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if strings.Contains(r.RequestURI, "v1/event") {
				defer r.Body.Close()
				bytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("could not read received event payload: %s", err.Error())
				}
				event := &apimodels.KeptnContextExtendedCE{}
				if err := json.Unmarshal(bytes, event); err != nil {
					t.Errorf("could not decode received event: %s", err.Error())
				}
				if *event.Type != keptnv2.GetTriggeredEventType("dev.hello") {
					t.Errorf("did not receive correct event: %s", err.Error())
				}
				go func() {
					receivedEvent <- true
				}()
			} else if strings.Contains(r.RequestURI, "/v1/metadata") {
				defer r.Body.Close()
				w.Write([]byte(metadataMockResponse))
				return
			} else if strings.Contains(r.RequestURI, "service") {
				res := fmt.Sprintf(getSvcMockResponseSequence, "demo")
				w.Write([]byte(res))
				return
			} else if strings.Contains(r.RequestURI, "/controlPlane/v1/project/") {
				res := fmt.Sprintf(getProjectMockResponseSequence, "hello-world", "demo", "dev")
				w.Write([]byte(res))
				return
			}
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	cmd := fmt.Sprintf("trigger sequence --sequence=%s --project=%s --service=%s --stage=%s --mock", "hello", "hello-world", "demo", "dev")
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

// TestTriggerSequenceNonExistingProject tests the trigger sequence
// with non-existing project.
func TestTriggerSequenceNonExistingProject(t *testing.T) {

	const nonExistingProject = "myproj"

	credentialmanager.MockAuthCreds = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			res := fmt.Sprintf(getProjectMockResponseSequenceNotFound, nonExistingProject)
			w.Write([]byte(res))
		}),
	)
	defer ts.Close()
	os.Setenv("MOCK_SERVER", ts.URL)

	tests := []struct {
		project string
		wantErr bool
	}{
		{
			project: nonExistingProject,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			cmd := fmt.Sprintf("trigger sequence --sequence=%s --project=%s --service=mysvc --stage=%s --mock",
				tt.project,
				"hello",
				"dev")
			_, err := executeActionCommandC(cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error: %t, got: %v", tt.wantErr, err)
			}
			msg := fmt.Sprintf("%v", err)
			if !strings.Contains(msg, "Project not found") {
				t.Errorf("wanted project not found")
			}
		})
	}
}

// TestTriggerSequenceNonExistingService tests the trigger sequence
// with non-existing service.
func TestTriggerSequenceNonExistingService(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	projectName := "hello-world"
	sequenceName := "hello"
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if strings.Contains(r.RequestURI, "/v1/metadata") {
				defer r.Body.Close()
				w.Write([]byte(metadataMockResponse))
				return
			} else if strings.Contains(r.RequestURI, "service") {
				res := fmt.Sprintf(getSvcMockResponseSequence, "helloservice")
				w.Write([]byte(res))
				return
			} else if strings.Contains(r.RequestURI, "/controlPlane/v1/project/") {
				res := fmt.Sprintf(getProjectMockResponseSequence, projectName, "demo", "dev")
				w.Write([]byte(res))
				return
			}
		}),
	)
	defer ts.Close()
	os.Setenv("MOCK_SERVER", ts.URL)

	tests := []struct {
		service string
		wantErr bool
	}{
		{
			service: "demo-service",
			wantErr: true,
		},
		{
			service: "helloservice",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			cmd := fmt.Sprintf("trigger sequence --sequence=%s --project=%s --service=%s --mock",
				sequenceName,
				projectName,
				tt.service,
			)
			_, err := executeActionCommandC(cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error: %t, got: %v", tt.wantErr, err)
			}
		})
	}
}
