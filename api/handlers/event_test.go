package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

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
			require.NoError(t, err)

			if tt.keepProvided {
				require.Equal(t, tt.args.eventKeptnContext, got)
			} else {
				require.NotEqual(t, tt.args.eventKeptnContext, got)
			}

			got2 := createOrApplyKeptnContext(tt.args.eventKeptnContext)

			_, err = uuid.Parse(got2)
			require.NoError(t, err)

			if tt.idempotent {
				require.Equal(t, got, got2)
			} else {
				require.NotEqual(t, got, got2)
			}
		})
	}
}

func TestPostEventHandlerFunc(t *testing.T) {

	topicName := "my-topic"
	natsServer, shutdown := runNATSServer()

	defer shutdown()

	err := os.Setenv("NATS_URL", natsServer.ClientURL())
	require.NoError(t, err)

	natsClient, err := nats.Connect(natsServer.ClientURL())
	require.Nil(t, err)

	receivedMessage := false

	_, err = natsClient.Subscribe(topicName, func(msg *nats.Msg) {
		receivedMessage = true
	})

	require.Nil(t, err)

	params := event.PostEventParams{
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
			Type:           &topicName,
		},
	}

	got := PostEventHandlerFunc(params, nil)

	verifyHTTPResponse(got, http.StatusOK, t)

	require.Eventually(t, func() bool {
		return receivedMessage
	}, 1*time.Second, 10*time.Millisecond)
}

type mockProducer struct {
}

func (m *mockProducer) Produce(io.Writer, interface{}) error {
	return nil
}

func stringp(s string) *string {
	return &s
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
			err := os.Setenv("DATASTORE_URI", tt.datastoreURLEnv)
			require.NoError(t, err)

			require.Equal(t, tt.want, utils.GetDatastoreURL())
		})
	}
}

func runNATSServer() (*server.Server, func()) {
	svr := natstest.RunRandClientPortServer()
	return svr, func() { svr.Shutdown() }
}
