package cmd

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_verifyConfigureBridgeParams(t *testing.T) {
	type args struct {
		configureBridgeParams *configureBridgeCmdParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should succeed",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{
					User:     stringp("user"),
					Password: stringp("password"),
				},
			},
			wantErr: false,
		}, {
			name: "should not succeed if no credentials are provided",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyConfigureBridgeParams(tt.args.configureBridgeParams); (err != nil) != tt.wantErr {
				t.Errorf("verifyConfigureBridgeParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}

func Test_configureBridge(t *testing.T) {
	ts := getTestAPI()
	defer ts.Close()

	type args struct {
		endpoint string
		apiToken string
		params   *configureBridgeCmdParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set credentials",
			args: args{
				endpoint: ts.URL,
				apiToken: "",
				params: &configureBridgeCmdParams{
					User:     stringp("user"),
					Password: stringp("password"),
				},
			},
			wantErr: false,
		},
		{
			name: "Get credentials",
			args: args{
				endpoint: ts.URL,
				apiToken: "",
				params: &configureBridgeCmdParams{
					Read: boolp(true),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configureBridge(tt.args.endpoint, tt.args.apiToken, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("configureBridge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func getTestAPI() *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			if r.Method == http.MethodPost {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(``))
				return
			} else if r.Method == http.MethodGet {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`{"user":"user", "password":"password"}`))
				return
			}
		}),
	)
	return ts
}

func boolp(b bool) *bool {
	return &b
}

func Test_retrieveBridgeCredentials(t *testing.T) {
	ts := getTestAPI()
	defer ts.Close()

	type args struct {
		endpoint string
		apiToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *configureBridgeAPIPayload
		wantErr bool
	}{
		{
			name: "Get credentials",
			args: args{
				endpoint: ts.URL,
				apiToken: "",
			},
			wantErr: false,
			want:    &configureBridgeAPIPayload{Password: "password", User: "user"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveBridgeCredentials(tt.args.endpoint, tt.args.apiToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrieveBridgeCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("retrieveBridgeCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}
