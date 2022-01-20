package cmd

import (
	"bytes"
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

const getProjectMockResponse = `{
	"creationDate": "1638796448951137480",
	"projectName": "%s",
	"shipyard": "",
	"shipyardVersion": "spec.keptn.sh/0.2.0",
	"stages": [
	  {
		"services": [
		  {
			"creationDate": "1638796449959720167",
			"deployedImage": "podtatoserver:v0.1.2",
			"openRemediations": null,
			"serviceName": "%s"
		  }
		],
		"stageName": "%s"
	  }
	]
}`

const getProjectMockResponseNotFound = `{
	"code": 404,
	"message": "Project not found: %s"
}`

const getSvcMockResponse = `{
	"nextPageKey": "0",
	"services": [
	  {
		"creationDate": "1638796449959720167",
		"deployedImage": "ghcr.io/podtato-head/podtatoserver:v0.1.2",
		"openRemediations": null,
		"serviceName": "%s"
	  }
	],
	"totalCount": 1
}`

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

type DockerImg struct {
	Image string
	Tag   string
}

// TestTriggerDelivery tests the trigger delivery command.
func TestTriggerDelivery(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

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
			} else if strings.Contains(r.RequestURI, "service") {
				res := fmt.Sprintf(getSvcMockResponse, "carts")
				w.Write([]byte(res))
				return
			} else if strings.Contains(r.RequestURI, "/controlPlane/v1/project/") {
				res := fmt.Sprintf(getProjectMockResponse, "sockshop", "carts", "dev")
				w.Write([]byte(res))
				return
			}
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

// TestTriggerDelivery tests the trigger delivery command.
func TestTriggerDeliveryNoStageProvided(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

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
			} else if strings.Contains(r.RequestURI, "service") {
				res := fmt.Sprintf(getSvcMockResponse, "carts")
				w.Write([]byte(res))
				return
			} else if strings.Contains(r.RequestURI, "/controlPlane/v1/project/") {
				res := fmt.Sprintf(getProjectMockResponse, "sockshop", "carts", "dev")
				w.Write([]byte(res))
				return
			}
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	cmd := fmt.Sprintf("trigger delivery --project=%s --service=%s --sequence=%s "+
		"--image=%s --tag=%s --values=a.b.c=d --mock --values=c.d=e", "sockshop", "carts", "artifact-delivery", "docker.io/keptnexamples/carts", "0.9.1")
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

// TestTriggerDeliveryNonExistingProject tests the trigger delivery
// with non-existing project.
func TestTriggerDeliveryNonExistingProject(t *testing.T) {

	const nonExistingProject = "myproj"

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(400)
			res := fmt.Sprintf(getProjectMockResponseNotFound, nonExistingProject)
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
			cmd := fmt.Sprintf("trigger delivery --project=%s --service=mysvc --image=%s --tag=%s --mock",
				tt.project,
				"docker.io/keptnexamples/carts",
				"0.9.1")
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

// TestTriggerDeliveryNonExistingService tests the trigger delivery
// with non-existing service.
func TestTriggerDeliveryNonExistingService(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true
	projectName := "sockshop"
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if strings.Contains(r.RequestURI, "/v1/metadata") {
				defer r.Body.Close()
				w.Write([]byte(metadataMockResponse))
				return
			} else if strings.Contains(r.RequestURI, "service") {
				res := fmt.Sprintf(getSvcMockResponse, "helloservice")
				w.Write([]byte(res))
				return
			} else if strings.Contains(r.RequestURI, "/controlPlane/v1/project/") {
				res := fmt.Sprintf(getProjectMockResponse, projectName, "carts", "dev")
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
			service: "some-service",
			wantErr: true,
		},
		{
			service: "helloservice",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			cmd := fmt.Sprintf("trigger delivery --project=%s --service=%s --image=%s --tag=%s --values=a.b.c=d --mock",
				projectName,
				tt.service,
				"docker.io/keptnexamples/carts",
				"0.9.1")
			_, err := executeActionCommandC(cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error: %t, got: %v", tt.wantErr, err)
			}
		})
	}
}
