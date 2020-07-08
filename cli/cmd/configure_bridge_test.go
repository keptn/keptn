package cmd

import (
	"net/http"
	"net/http/httptest"
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
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(``))
		}),
	)
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
			name: "Expose",
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
