package config

import (
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/sirupsen/logrus"
)

var Global EnvConfig

const GitImplementationLibGit2 = "libgit2"
const GitImplementationGoGit = "go-git"

type EnvConfig struct {
	LogLevel                         string `envconfig:"LOG_LEVEL" default:"info"`
	DirectoryStageStructure          bool   `envconfig:"DIRECTORY_STAGE_STRUCTURE" default:"false"`
	DefaultRemoteGitRepositoryBranch string `envconfig:"DEFAULT_REMOTE_GIT_BRANCH" default:"master"`
	GitImplementation                string `envconfig:"GIT_IMPLEMENTATION" default:"go-git"`
}

func (e EnvConfig) RetrieveDefaultBranchFromEnv() string {
	if e.DefaultRemoteGitRepositoryBranch == "" {
		logrus.Debugf("Could not determine default remote git repository branch from env variable")
		return common_models.GitInitDefaultBranchName
	}
	return e.DefaultRemoteGitRepositoryBranch
}
