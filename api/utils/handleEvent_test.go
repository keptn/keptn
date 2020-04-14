package utils

import (
	"os"
	"testing"
)

func Test_getEventBrokerURL(t *testing.T) {
	tests := []struct {
		name              string
		want              string
		eventbrokerURLEnv string
	}{
		{
			name:              "get from env var without https:// or http:// prefix",
			want:              "http://localhost",
			eventbrokerURLEnv: "localhost",
		},
		{
			name:              "get from env var with https:// prefix",
			want:              "https://localhost",
			eventbrokerURLEnv: "https://localhost",
		},
		{
			name:              "get from env var with http:// prefix",
			want:              "http://localhost",
			eventbrokerURLEnv: "http://localhost",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("EVENTBROKER_URI", tt.eventbrokerURLEnv)
			if got := getEventBrokerURL(); got != tt.want {
				t.Errorf("getEventBrokerURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
