package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getReplaceSecretCommand(t *testing.T) {
	type args struct {
		cmdParams configureBridgeCmdParams
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "print output with placeholders",
			args: args{
				cmdParams: configureBridgeCmdParams{},
			},
			want: "kubectl create secret -n keptn generic bridge-credentials --from-literal=\"BASIC_AUTH_USERNAME=${BRIDGE_USER}\" --from-literal=\"BASIC_AUTH_PASSWORD=${BRIDGE_PASSWORD}\" -oyaml --dry-run=client | kubectl replace -f -\n",
		},
		{
			name: "print output with provided values",
			args: args{
				cmdParams: configureBridgeCmdParams{
					User:     stringp("my-user"),
					Password: stringp("my-password"),
				},
			},
			want: "kubectl create secret -n keptn generic bridge-credentials --from-literal=\"BASIC_AUTH_USERNAME=my-user\" --from-literal=\"BASIC_AUTH_PASSWORD=my-password\" -oyaml --dry-run=client | kubectl replace -f -\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getReplaceSecretCommand(tt.args.cmdParams)

			require.Contains(t, got, tt.want)
		})
	}
}

// TestConfigureBridgeUnknownCommand
func TestConfigureBridgeUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("configure bridge someUnknownCommand --user=user --password=pass")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn configure bridge\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestConfigureBridgeUnknownParameter
func TestConfigureBridgeUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("configure bridge --userr=user --password=pass")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown flag: --userr"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
