package cmd

import (
	"encoding/json"
	"io/ioutil"
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
			name: "action=expose should succeed",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{Action: stringp("expose")},
			},
			wantErr: false,
		},
		{
			name: "action=lockdown should succeed",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{Action: stringp("lockdown")},
			},
			wantErr: false,
		},
		{
			name: "action=invalid should not succeed",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{Action: stringp("invalid")},
			},
			wantErr: true,
		},
		{
			name: "action=nil should not succeed",
			args: args{
				configureBridgeParams: &configureBridgeCmdParams{Action: nil},
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
			param := &exposeBridgeAPIPayload{}
			body, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(body, param)

			if param.Expose {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(`bridge.keptn.test`))
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write([]byte(``))
			}
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
				params:   &configureBridgeCmdParams{Action: stringp("expose")},
			},
			wantErr: false,
		},
		{
			name: "Lockdown",
			args: args{
				endpoint: ts.URL,
				apiToken: "",
				params:   &configureBridgeCmdParams{Action: stringp("lockdown")},
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
