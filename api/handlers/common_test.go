package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime/middleware"

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
			if got := utils.SanitizeURL(tt.in); got != tt.want {
				t.Errorf("sanitizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func verifyHTTPResponse(got middleware.Responder, expectedStatus int, t *testing.T) {
	producer := &mockProducer{}
	recorder := &httptest.ResponseRecorder{}
	got.WriteResponse(recorder, producer)

	if recorder.Result().StatusCode != expectedStatus {
		t.Errorf("Returned HTTP status = %v, want %v", recorder.Result().StatusCode, expectedStatus)
	}
}
