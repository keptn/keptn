package common_models

// GitCredentials contains git credentials info
type GitCredentials struct {
	User             string `json:"user,omitempty"`
	Token            string `json:"token,omitempty"`
	RemoteURI        string `json:"remoteURI,omitempty"`
	PrivateKey       string `json:"privateKey,omitempty"`
	GitProxyUrl      string `json:"gitProxyUrl,omitempty"`
	GitProxyScheme   string `json:"gitProxyScheme,omitempty"`
	GitProxyUser     string `json:"gitProxyUser,omitempty"`
	GitProxyPassword string `json:"gitProxyPassword,omitempty"`
	GitPublicCert    string `json:"gitPublicCert,omitempty"`
}

type GitContext struct {
	Project     string
	Credentials *GitCredentials
}

func (g GitCredentials) Validate() error {
	// _, err := url.Parse(g.RemoteURI)
	// if err != nil {
	// 	return kerrors.ErrCredentialsInvalidRemoteURI
	// }
	// if g.Token == "" {
	// 	return kerrors.ErrCredentialsTokenMustNotBeEmpty
	// }
	return nil
}
