package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/utils"
)

func Test_sanitizeURL(t *testing.T) {
	tests := []struct {
		name string
		want string
		in   string
	}{
		{
			name: "get from env var without https:// or http:// prefix",
			want: "http://localhost",
			in:   "localhost",
		},
		{
			name: "get from env var with https:// prefix",
			want: "https://localhost",
			in:   "https://localhost",
		},
		{
			name: "get from env var with http:// prefix",
			want: "http://localhost",
			in:   "http://localhost",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, utils.SanitizeURL(tt.in))
		})
	}
}

func verifyHTTPResponse(got middleware.Responder, expectedStatus int, t *testing.T) {
	producer := &mockProducer{}
	recorder := &httptest.ResponseRecorder{}
	got.WriteResponse(recorder, producer)

	require.Equal(t, expectedStatus, recorder.Result().StatusCode)
}
