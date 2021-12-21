package common

// GitCredentials contains git credentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

// IGit provides functions to interact with the git repository of a project
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/git_mock.go . IGit
type IGit interface {
	ProjectExists(gitContext GitContext) bool
	CloneRepo(gitContext GitContext) (bool, error)
	StageAndCommitAll(gitContext GitContext, message string) error
	Push(gitContext GitContext) error
	Pull(gitContext GitContext) error
	CreateBranch(gitContext GitContext, branch string, sourceBranch string) error
	CheckoutBranch(gitContext GitContext, branch string) error
	GetDefaultBranch(gitContext GitContext)
}
