package models

type ProvisioningData struct {
	GitRemoteURL string `json:"gitRemoteURL"`
	GitToken     string `json:"gitToken"`
	GitUser      string `json:"gitUser"`
}
