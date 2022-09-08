package common

import "testing"

func TestEncodeKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode dots",
			args: args{
				key: "sh.keptn.event",
			},
			want: "sh~pkeptn~pevent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeKey(tt.args.key); got != tt.want {
				t.Errorf("EncodeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode dots",
			args: args{
				key: "sh~pkeptn~pevent",
			},
			want: "sh.keptn.event",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeKey(tt.args.key); got != tt.want {
				t.Errorf("DecodeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
