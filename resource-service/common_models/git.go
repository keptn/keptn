package common_models

const masterBranch = "master"

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
