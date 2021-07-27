package middleware

import (
	"github.com/keptn/keptn/api/models"
	"os"
	"reflect"
	"testing"
)

func TestValidateToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name            string
		args            args
		want            models.Principal
		configuredToken string
		wantErr         bool
	}{
		{
			name: "token valid",
			args: args{
				token: "my-token",
			},
			configuredToken: "my-token",
			want:            models.Principal("my-token"),
			wantErr:         false,
		},
		{
			name: "token invalid",
			args: args{
				token: "my-invalid-token",
			},
			configuredToken: "my-token",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("SECRET_TOKEN", tt.configuredToken)
			tv := &BasicTokenValidator{}
			got, err := tv.ValidateToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" {
				if !reflect.DeepEqual(*got, tt.want) {
					t.Errorf("ValidateToken() got = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
