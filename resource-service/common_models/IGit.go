package common_models

const masterBranch = "master"

// GitCredentials contains git credentials info
type GitCredentials struct {
	User      string `json:"user,omitempty"`
	Token     string `json:"token,omitempty"`
	RemoteURI string `json:"remoteURI,omitempty"`
}

type GitContext struct {
	Project     string //TODO would it make sense to store branch ? no
	Credentials *GitCredentials
}

// IGit provides functions to interact with the git repository of a project
//go:generate moq -pkg common_mock -skip-ensure -out ./fake/git_mock.go . IGit
type IGit interface {
	ProjectExists(gitContext GitContext) bool
	CloneRepo(gitContext GitContext) (bool, error)
	StageAndCommitAll(gitContext GitContext, message string) (string, error)
	Push(gitContext GitContext) error
	Pull(gitContext GitContext) error
	CreateBranch(gitContext GitContext, branch string, sourceBranch string) error
	CheckoutBranch(gitContext GitContext, branch string) error
	GetFileRevision(gitContext GitContext, revision string, file string) ([]byte, error)
	GetDefaultBranch(gitContext GitContext) (string, error)
	GetCurrentRevision(gitContext GitContext) (string, error)
}
