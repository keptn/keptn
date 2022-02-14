package common_models

import "testing"

func TestGitCredentials_Validate(t *testing.T) {
	type fields struct {
		User      string
		Token     string
		RemoteURI string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid credentials",
			fields: fields{
				User:      "my-user",
				Token:     "my-token",
				RemoteURI: "https://my-repo",
			},
			wantErr: false,
		},
		{
			name: "empty token",
			fields: fields{
				User:      "my-user",
				Token:     "",
				RemoteURI: "https://my-repo",
			},
			wantErr: true,
		},
		{
			name: "invalid URI",
			fields: fields{
				User:      "my-user",
				Token:     "my-token",
				RemoteURI: "https://my:repo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GitCredentials{
				User:      tt.fields.User,
				Token:     tt.fields.Token,
				RemoteURI: tt.fields.RemoteURI,
			}
			if err := g.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
