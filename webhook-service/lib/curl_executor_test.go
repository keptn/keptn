package lib_test

import (
	"github.com/keptn/keptn/webhook-service/lib"
	"github.com/keptn/keptn/webhook-service/lib/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCmdCurlExecutor_InvalidCommand(t *testing.T) {
	executor := lib.NewCmdCurlExecutor(&fake.ICommandExecutorMock{ExecuteCommandFunc: func(cmd string, args ...string) (string, error) {
		return "", nil
	}})

	output, err := executor.Curl("invalid command")

	require.NotNil(t, err)
	require.Empty(t, output)
}

func TestNewCmdCurlExecutor_UnAllowedURL(t *testing.T) {
	executor := lib.NewCmdCurlExecutor(&fake.ICommandExecutorMock{ExecuteCommandFunc: func(cmd string, args ...string) (string, error) {
		return "", nil
	}}, lib.WithUnAllowedURLs([]string{"kube-api"}))

	output, err := executor.Curl("curl http://kube-api")

	require.NotNil(t, err)
	require.Empty(t, output)
}

func TestNewCmdCurlExecutor_EmptyCommand(t *testing.T) {
	executor := lib.NewCmdCurlExecutor(&fake.ICommandExecutorMock{ExecuteCommandFunc: func(cmd string, args ...string) (string, error) {
		return "", nil
	}})

	output, err := executor.Curl("")

	require.NotNil(t, err)
	require.Empty(t, output)
}

func TestCmdCurlExecutor_Curl(t *testing.T) {
	type fields struct {
		commandExecutor *fake.ICommandExecutorMock
	}
	type args struct {
		curlCmd string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           string
		shouldExecute  bool
		wantPassedArgs []string
		wantErr        bool
	}{
		{
			name: "valid request",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!\"}' https://my.hook.com/foo`,
			},
			want:          "success",
			shouldExecute: true,
			wantPassedArgs: []string{
				"-X", "POST", "-H", "Content-type: application/json", "--data", `{\"text\":\"Hello, World!\"}`, "https://my.hook.com/foo",
			},
			wantErr: false,
		},
		{
			name: "try to inject command - should return error",
			args: args{
				curlCmd: `curl -X POST -H 'token: $(kubectl exec)' --data '{\"text\":\"Hello, World!\"}' https://my.hook.com/foo`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (2)",
			args: args{
				curlCmd: `curl -X POST -H 'token: abcd' --data '{\"text\":\"Hello, World!\"}' https://my.hook.com/foo | pwd`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to download to file - should return error",
			args: args{
				curlCmd: `curl -X POST -H 'token: abcd' --data '{\"text\":\"Hello, World!\"}' https://my.hook.com/foo -o somefile`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "unclosed quote",
			args: args{
				curlCmd: `curl -X POST -H 'token: abcd' --data '{\"text\":\"Hello, World!\"} https://my.hook.com/foo -o somefile`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeCommandExecutor := &fake.ICommandExecutorMock{
				ExecuteCommandFunc: func(cmd string, args ...string) (string, error) {
					return "success", nil
				},
			}

			if tt.fields.commandExecutor != nil {
				fakeCommandExecutor = tt.fields.commandExecutor
			}

			ce := lib.NewCmdCurlExecutor(fakeCommandExecutor)

			got, err := ce.Curl(tt.args.curlCmd)

			if tt.wantErr {
				assert.NotNil(t, err)
			}

			if tt.shouldExecute {
				require.NotEmpty(t, fakeCommandExecutor.ExecuteCommandCalls())
				require.Equal(t, tt.wantPassedArgs, fakeCommandExecutor.ExecuteCommandCalls()[0].Args)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
