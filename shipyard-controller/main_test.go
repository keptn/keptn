package main

import (
	"os"
	"testing"
	"time"
)

func Test_getDurationFromEnvVar(t *testing.T) {
	type args struct {
		envVarValue string
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "get default value",
			args: args{
				envVarValue: "",
			},
			want: 432000 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "10s",
			},
			want: 10 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "2m",
			},
			want: 120 * time.Second,
		},
		{
			name: "get configured value",
			args: args{
				envVarValue: "1h30m",
			},
			want: 5400 * time.Second,
		},
		{
			name: "get default value because of invalid config",
			args: args{
				envVarValue: "invalid",
			},
			want: 432000 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LOG_TTL", tt.args.envVarValue)
			if got := getDurationFromEnvVar("LOG_TTL", envVarLogsTTLDefault); got != tt.want {
				t.Errorf("getLogTTLDurationInSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
