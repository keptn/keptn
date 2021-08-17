package model

const DefaultSecretScope = "keptn-default"

// Secret secret
// swagger:model secret
type Secret struct {
	SecretMetadata
	Data Data `json:"data"`
}

type Data map[string]string

type SecretMetadata struct {
	Name string `json:"name" binding:"required"`
	// Scope determines the scope of the secret (default="keptn-default")
	Scope string `json:"scope,omitempty"`
}

type GetSecretsResponse struct {
	Secrets []SecretMetadata
}
