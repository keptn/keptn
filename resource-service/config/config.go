package config

var Global EnvConfig

type EnvConfig struct {
	LogLevel                         string `envconfig:"LOG_LEVEL" default:"info"`
	DirectoryStageStructure          bool   `envconfig:"DIRECTORY_STAGE_STRUCTURE" default:"false"`
	DefaultRemoteGitRepositoryBranch string `envconfig:"DEFAULT_REMOTE_GIT_BRANCH" default:"master"`
}
