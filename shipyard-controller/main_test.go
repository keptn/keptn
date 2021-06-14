package main

import "testing"

func Test_getLogTTLDurationInSeconds(t *testing.T) {
	type args struct {
		logsTTL string
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "get default value",
			args: args{
				logsTTL: "",
			},
			want: 432000,
		},
		{
			name: "get configured value",
			args: args{
				logsTTL: "10s",
			},
			want: 10,
		},
		{
			name: "get configured value",
			args: args{
				logsTTL: "2m",
			},
			want: 120,
		},
		{
			name: "get configured value",
			args: args{
				logsTTL: "1h30m",
			},
			want: 5400,
		},
		{
			name: "get default value because of invalid config",
			args: args{
				logsTTL: "invalid",
			},
			want: 432000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLogTTLDurationInSeconds(tt.args.logsTTL); got != tt.want {
				t.Errorf("getLogTTLDurationInSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
