package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_retrieveDefaultBranchFromEnv(t *testing.T) {
	tests := []struct {
		name string
		env  EnvConfig
		want string
	}{
		{
			name: "not set",
			env: EnvConfig{
				DefaultRemoteGitRepositoryBranch: "",
			},
			want: "master",
		},
		{
			name: "master set",
			env: EnvConfig{
				DefaultRemoteGitRepositoryBranch: "master",
			},
			want: "master",
		},
		{
			name: "main set",
			env: EnvConfig{
				DefaultRemoteGitRepositoryBranch: "main",
			},
			want: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.env.RetrieveDefaultBranchFromEnv()
			require.Equal(t, tt.want, got)
		})
	}
}
