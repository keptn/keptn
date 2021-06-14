package model

type Secret struct {
	SecretMetadata
	Data Data `json:"data"`
}

type Data map[string]string

type SecretMetadata struct {
	Name  string `json:"name" binding:"required"`
	Scope string `json:"scope,omitempty" binding:"required"`
}

type GetSecretsResponse struct {
	Secrets []SecretMetadata
}
