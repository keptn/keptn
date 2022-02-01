package cmd

import (
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
	testInvalidInputHelper("configure bridge someUnknownCommand --user=user --password=pass", "unknown command \"someUnknownCommand\" for \"keptn configure bridge\"", t)
}

// TestConfigureBridgeUnknownParameter
func TestConfigureBridgeUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("configure bridge --userr=user --password=pass", "unknown flag: --userr", t)
}
