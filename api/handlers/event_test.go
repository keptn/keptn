package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"

	datastoremodels "github.com/keptn/go-utils/pkg/api/models"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/utils"
)

func Test_createOrApplyKeptnContext(t *testing.T) {
	type args struct {
		eventKeptnContext string
	}
	tests := []struct {
		name         string
		args         args
		keepProvided bool
		idempotent   bool
	}{
		{
			name: "Keep provided keptnContext UUID",
			args: args{
				eventKeptnContext: "3f709bb7-0246-403e-83a0-1f436f7c6c09",
			},
			keepProvided: true,
			idempotent:   true,
		},
		{
			name: "Generate new random UUID",
			args: args{
				eventKeptnContext: "",
			},
			keepProvided: false,
			idempotent:   false,
		},
		{
			name: "Derive UUID from provided value (<16 chars)",
			args: args{
				eventKeptnContext: "a",
			},
			keepProvided: false,
			idempotent:   true,
		},
		{
			name: "Derive UUID from provided value (>=16 chars)",
			args: args{
				eventKeptnContext: "1234567890123456",
			},
			keepProvided: false,
			idempotent:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createOrApplyKeptnContext(tt.args.eventKeptnContext)

			_, err := uuid.Parse(got)

			if err != nil {
				t.Errorf("createOrApplyKeptnContext(): generated value %v is not a valid UUID", got)
			}

			if tt.keepProvided {
				if got != tt.args.eventKeptnContext {
					t.Errorf("createOrApplyKeptnContext() = %v, want %v", got, tt.args.eventKeptnContext)
				}
			} else if !tt.keepProvided && got == tt.args.eventKeptnContext {
				t.Errorf("createOrApplyKeptnContext() = %v, want != %v", got, tt.args.eventKeptnContext)
			}

			got2 := createOrApplyKeptnContext(tt.args.eventKeptnContext)

			_, err = uuid.Parse(got2)

			if err != nil {
				t.Errorf("createOrApplyKeptnContext(): generated value %v is not a valid UUID", got2)
			}

			if tt.idempotent {
				if got != got2 {
					t.Errorf("createOrApplyKeptnContext() = %v, want = %v", got2, got)
				}
			} else if !tt.idempotent && got == got2 {
				t.Errorf("createOrApplyKeptnContext() = %v, want != %v", got2, got)
			}

		})
	}
}

func TestPostEventHandlerFunc(t *testing.T) {
	_ = os.Setenv("SECRET_TOKEN", "testtesttesttesttest")

	returnedStatus := 200

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(returnedStatus)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	_ = os.Setenv("EVENTBROKER_URI", ts.URL)

	type args struct {
		params    event.PostEventParams
		principal *models.Principal
	}
	tests := []struct {
		name                  string
		args                  args
		wantStatus            int
		statusFromEventBroker int
	}{
		{
			name: "Send event",
			args: args{
				params: event.PostEventParams{
					HTTPRequest: nil,
					Body: &models.KeptnContextExtendedCE{
						Contenttype:    "application/json",
						Data:           map[string]interface{}{},
						Extensions:     nil,
						ID:             "",
						Shkeptncontext: "",
						Source:         stringp("test-source"),
						Specversion:    "1.0",
						Time:           strfmt.DateTime{},
						Type:           stringp(keptnevents.ConfigureMonitoringEventType),
					},
				},
				principal: nil,
			},
			wantStatus:            200,
			statusFromEventBroker: 200,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			returnedStatus = tt.statusFromEventBroker
			got := PostEventHandlerFunc(tt.args.params, tt.args.principal)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

type mockProducer struct {
}

func (m *mockProducer) Produce(io.Writer, interface{}) error {
	return nil
}

func stringp(s string) *string {
	return &s
}

func TestGetEventHandlerFunc(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			e := &datastoremodels.Events{
				Events: []*datastoremodels.KeptnContextExtendedCE{
					{
						Contenttype:    "",
						Data:           nil,
						Extensions:     nil,
						ID:             "",
						Source:         stringp(""),
						Specversion:    "",
						Time:           strfmt.DateTime{},
						Type:           stringp(""),
						Shkeptncontext: "",
					},
				},
				NextPageKey: "",
				PageSize:    0,
				TotalCount:  0,
			}

			marshal, _ := json.Marshal(e)
			w.Write(marshal)
		}),
	)
	defer ts.Close()

	_ = os.Setenv("DATASTORE_URI", ts.URL)
	type args struct {
		params    event.GetEventParams
		principal *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Get events",
			args: args{
				params: event.GetEventParams{
					HTTPRequest:  nil,
					KeptnContext: "",
					Type:         keptnevents.ConfigureMonitoringEventType,
				},
				principal: nil,
			},
			wantStatus: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEventHandlerFunc(tt.args.params, tt.args.principal)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}

func Test_getDatastoreURL(t *testing.T) {
	tests := []struct {
		name            string
		want            string
		datastoreURLEnv string
	}{
		{
			name:            "get from env var without https:// or http:// prefix",
			want:            "http://localhost",
			datastoreURLEnv: "localhost",
		},
		{
			name:            "get from env var with https:// prefix",
			want:            "https://localhost",
			datastoreURLEnv: "https://localhost",
		},
		{
			name:            "get from env var with http:// prefix",
			want:            "http://localhost",
			datastoreURLEnv: "http://localhost",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DATASTORE_URI", tt.datastoreURLEnv)

			if got := utils.GetDatastoreURL(); got != tt.want {
				t.Errorf("getDatastoreURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
