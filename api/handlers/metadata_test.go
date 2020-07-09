package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/metadata"
)

func TestGetMetadataHandlerFunc(t *testing.T) {
	type args struct {
		params metadata.MetadataParams
		p      *models.Principal
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Get metadata",
			args: args{
				params: metadata.MetadataParams{
					HTTPRequest: nil,
				},
				p: nil,
			},
			wantStatus: 200,
		},
	}

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMetadataHandlerFunc(tt.args.params, tt.args.p)

			verifyHTTPResponse(got, tt.wantStatus, t)
		})
	}
}
