package common_models

import "testing"

func TestGitCredentials_Validate(t *testing.T) {
	type fields struct {
		User           string
		Token          string
		PrivateKey     string
		RemoteURI      string
		GitProxyURL    string
		GitProxyScheme string
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
			wantErr: false,
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
		{
			name: "empty PrivateKey",
			fields: fields{
				User:       "my-user",
				PrivateKey: "",
				RemoteURI:  "ssh://my:repo",
			},
			wantErr: true,
		},
		{
			name: "valid PrivateKey",
			fields: fields{
				User:       "my-user",
				PrivateKey: "privatekey",
				RemoteURI:  "ssh://my:repo",
			},
			wantErr: false,
		},
		{
			name: "PrivateKey with https",
			fields: fields{
				User:       "my-user",
				PrivateKey: "",
				RemoteURI:  "https://my:repo",
			},
			wantErr: true,
		},
		{
			name: "token with ssh",
			fields: fields{
				User:      "my-user",
				Token:     "token",
				RemoteURI: "ssh://my-repo",
			},
			wantErr: true,
		},
		{
			name: "http proxy",
			fields: fields{
				User:           "my-user",
				Token:          "my-token",
				RemoteURI:      "https://my-repo",
				GitProxyURL:    "1.1.1.1:12",
				GitProxyScheme: "http",
			},
			wantErr: false,
		},
		{
			name: "https proxy",
			fields: fields{
				User:           "my-user",
				Token:          "token",
				RemoteURI:      "https://my-repo",
				GitProxyURL:    "1.1.1.1:12",
				GitProxyScheme: "https",
			},
			wantErr: false,
		},
		{
			name: "proxy invalid scheme",
			fields: fields{
				User:           "my-user",
				Token:          "token",
				RemoteURI:      "https://my-repo",
				GitProxyURL:    "1.1.1.1:12",
				GitProxyScheme: "fddd",
			},
			wantErr: true,
		},
		{
			name: "proxy URL without port",
			fields: fields{
				User:           "my-user",
				Token:          "token",
				RemoteURI:      "https://my-repo",
				GitProxyURL:    "1.1.1.1",
				GitProxyScheme: "https",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GitCredentials{
				User:           tt.fields.User,
				Token:          tt.fields.Token,
				GitPrivateKey:  tt.fields.PrivateKey,
				RemoteURI:      tt.fields.RemoteURI,
				GitProxyURL:    tt.fields.GitProxyURL,
				GitProxyScheme: tt.fields.GitProxyScheme,
			}
			if err := g.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
