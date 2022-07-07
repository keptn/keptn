package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const approvalTriggeredMockResponse = `{
    "events": [
        {
		  "contenttype": "application/json",
		  "data": {
			"deploymentURILocal": "http://carts.sockshop-dev",
			"deploymentURIPublic": "http://carts.sockshop-dev",
			"deploymentstrategy": "direct",
			"image": "docker.io/keptnexamples/carts",
			"labels": null,
			"project": "sockshop",
			"service": "carts",
			"stage": "dev",
			"tag": "0.10.1",
			"teststrategy": "functional"
		  },
		  "id": "test-event-id-1",
		  "source": "helm-service",
		  "specversion": "1.0",
		  "time": "2020-04-14T08:11:27.484Z",
		  "type": "sh.keptn.event.approval.triggered",
		  "shkeptncontext": "test-event-context-1"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

const evaluationFinishedMockResponse = `{
    "events": [
        {
		  "contenttype": "application/json",
		  "data": {
			"deploymentstrategy": "blue_green_service",
			"evaluation": {
			  "result": "pass",
			  "score": 100,
			},
			"labels": null,
			"project": "sockshop",
			"result": "pass",
			"service": "carts",
			"stage": "staging",
			"teststrategy": "performance",
		  },
		  "id": "123",
		  "source": "lighthouse-service",
		  "specversion": "1.0",
		  "time": "2020-06-02T12:28:54.642Z",
		  "type": "sh.keptn.event.evaluation.finished",
		  "shkeptncontext": "test-event-context-1"
		}
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

func Test_sendApprovalFinishedEvent(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			if strings.Contains(r.RequestURI, keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName)) {
				w.Write([]byte(approvalTriggeredMockResponse))
				return
			} else if strings.Contains(r.RequestURI, keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName)) {
				w.Write([]byte(evaluationFinishedMockResponse))
				return
			} else if strings.Contains(r.RequestURI, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName)) {
				w.Write([]byte(evaluationFinishedMockResponse))
				return
			}
			return
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	type args struct {
		sendApprovalFinishedOptions sendApprovalFinishedStruct
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "send approval finished event for ID",
			args: args{
				sendApprovalFinishedOptions: sendApprovalFinishedStruct{
					Project: stringp("sockshop"),
					Stage:   stringp("dev"),
					Service: stringp(""),
					ID:      stringp("test-event-id-1"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendApprovalFinishedEvent(tt.args.sendApprovalFinishedOptions); (err != nil) != tt.wantErr {
				t.Errorf("sendApprovalFinishedEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_selectApprovalOption(t *testing.T) {
	type args struct {
		nrOfOptions int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Select 1",
			args: args{
				nrOfOptions: 2,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, oldStdin := createMockStdIn("1")

			defer os.Remove(tmpfile.Name())        // clean up
			defer func() { os.Stdin = oldStdin }() // Restore original Stdin

			os.Stdin = tmpfile
			got, err := selectApprovalOption(tt.args.nrOfOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectApprovalOption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("selectApprovalOption() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_approveOrDecline(t *testing.T) {
	tests := []struct {
		name      string
		want      bool
		userInput string
	}{
		{
			name:      "select approval",
			want:      true,
			userInput: "a",
		},
		{
			name:      "select decline",
			want:      false,
			userInput: "d",
		},
	}
	for _, tt := range tests {
		tmpfile, oldStdin := createMockStdIn(tt.userInput)

		defer os.Remove(tmpfile.Name())        // clean up
		defer func() { os.Stdin = oldStdin }() // Restore original Stdin

		os.Stdin = tmpfile
		t.Run(tt.name, func(t *testing.T) {
			if got := approveOrDecline(); got != tt.want {
				t.Errorf("approveOrDecline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createMockStdIn(userInput string) (*os.File, *os.File) {
	content := []byte(userInput + "\n")
	tmpfile, err := ioutil.TempFile("", "test_select_option_tmp")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	return tmpfile, oldStdin
}

// TestSendEventApprovalFinishedUnknownCommand
func TestSendEventApprovalFinishedUnknownCommand(t *testing.T) {
	testInvalidInputHelper("send event approval.finished someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn send event approval.finished\"", t)
}

// TestSendEventApprovalFinishedUnknownParameter
func TestSendEventApprovalFinishedUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("send event approval.finished --projectt=sockshop", "unknown flag: --projectt", t)
}
