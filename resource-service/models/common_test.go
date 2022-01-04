package models

import "testing"

func Test_validateEntityName(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid name",
			args: args{
				s: "my-name",
			},
			wantErr: false,
		},
		{
			name: "invalid name",
			args: args{
				s: "my/name",
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			args: args{
				s: "my name",
			},
			wantErr: true,
		},
		{
			name: "invalid name",
			args: args{
				s: " ",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEntityName(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("validateEntityName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
