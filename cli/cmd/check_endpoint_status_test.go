package cmd

import (
	"context"
	"testing"
	"time"
)

func Test_checkEndPointStatus(t *testing.T) {
	ts := getTestAPI()

	type args struct {
		endPoint string
		timeout  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simulate timeout", args{endPoint: ts.URL, timeout: 1 * time.Microsecond}, true},
		{"does not simulate timeout", args{endPoint: ts.URL, timeout: 5 * time.Second}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//override maxHTTPTimeout to simulate http timeout
			maxHTTPTimeout = tt.args.timeout

			if err := checkEndPointStatus(tt.args.endPoint); (err != nil) != tt.wantErr && err != context.DeadlineExceeded {
				t.Errorf("checkEndPointStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
