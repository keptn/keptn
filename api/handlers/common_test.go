package handlers

import (
	"testing"
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
			if got := sanitizeURL(tt.in); got != tt.want {
				t.Errorf("sanitizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
