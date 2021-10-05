package lib_test

import (
	"errors"
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
			name: "try to inject command - should return error (3)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!\"}' https://attack.domain || pwd`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (4)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!\"}' https://attack.domain & $(pwd)`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (5)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!\"}' https://orf.at && $(pwd)`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (6)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!\"}' https://attack.domain ; $(pwd)`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (7)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json';$(pwd) #' --data '{\"text\":\"Hello, World!\"}' localhost:8000`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (8)",
			args: args{
				curlCmd: `curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Hello, World!}'; $(pwd) #\"}' https://attack.domain`,
			},
			want:          "",
			shouldExecute: false,
			wantErr:       true,
		},
		{
			name: "try to inject command - should return error (8)",
			args: args{
				curlCmd: "curl -X POST -H 'Content-type:' `whoami` #'--data '{\"text\":\"Hello, World!\"}' localhost:8000",
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
			name: "try to upload file - should return error",
			args: args{
				curlCmd: `curl -X POST -H 'token: abcd' --data '{\"text\":\"Hello, World!\"}' https://my.hook.com/foo -F 'data=@path/to/local/file'`,
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

func TestIsNoCommandError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no command error",
			args: args{
				err: lib.NewCurlError(errors.New("oops"), lib.NoCommandError),
			},
			want: true,
		},
		{
			name: "any error",
			args: args{
				err: errors.New("oops"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.IsNoCommandError(tt.args.err); got != tt.want {
				t.Errorf("IsNoCommandError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalidCommandError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "invalid command error",
			args: args{
				err: lib.NewCurlError(errors.New("oops"), lib.InvalidCommandError),
			},
			want: true,
		},
		{
			name: "any error",
			args: args{
				err: errors.New("oops"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.IsInvalidCommandError(tt.args.err); got != tt.want {
				t.Errorf("IsInvalidCommandError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnallowedURLError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "unallowed URL error",
			args: args{
				err: lib.NewCurlError(errors.New("oops"), lib.UnallowedURLError),
			},
			want: true,
		},
		{
			name: "any error",
			args: args{
				err: errors.New("oops"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.IsUnallowedURLError(tt.args.err); got != tt.want {
				t.Errorf("IsUnallowedURLError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRequestError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "request execution error",
			args: args{
				err: lib.NewCurlError(errors.New("oops"), lib.RequestError),
			},
			want: true,
		},
		{
			name: "any error",
			args: args{
				err: errors.New("oops"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lib.IsRequestError(tt.args.err); got != tt.want {
				t.Errorf("IsRequestError() = %v, want %v", got, tt.want)
			}
		})
	}
}
