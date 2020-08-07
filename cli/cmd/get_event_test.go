package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestGetEventCmdEptyInput(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event --project=%s",
		"sockshop")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "please provide an event type as an argument") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

func TestGetTriggeredEvent(t *testing.T) {
	mocking = true
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			if strings.Contains(r.RequestURI, "/service") {
				w.WriteHeader(200)
				w.Write([]byte(allServicesInStageResponse))
				return
			} else if strings.Contains(r.RequestURI, "/event") {
				if strings.Contains(r.RequestURI, "sockshop") {
					w.WriteHeader(200)
					w.Write([]byte(eventsForID1Response))
					return
				}
			}

			w.WriteHeader(500)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	os.Setenv("MOCK_SERVER", ts.URL)

	var numOfPages *int
	numOfPages = new(int)
	*numOfPages = 1

	tests := []struct {
		name    string
		args    []string
		eventParam	GetEventStruct
		wantErr bool
	}{
		{
			name: "get evaluation-done events",
			args: []string{
				"sh.keptn.events.evaluation-done",
			},
			eventParam: GetEventStruct{
				Project: stringp("sockshop"),
				Stage:   stringp("staging"),
				Service: stringp("carts"),
				PageSize: stringp(""),
				Output: stringp("yaml"),
				KeptnContext: stringp(""),
				NumOfPages: numOfPages,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		if err := getEvent(tt.eventParam, tt.args); (err != nil) != tt.wantErr {
			t.Errorf("getApprovalTriggeredEvents() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}