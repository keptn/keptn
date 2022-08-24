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

type GetSecretResponseItem struct {
	SecretMetadata
	Keys []string `json:"keys"`
}

type GetSecretsResponse struct {
	Secrets []GetSecretResponseItem `json:"Secrets"`
}

type GetSecretQueryParams struct {
	Name  string `form:"name,omitempty"`
	Scope string `form:"scope,omitempty"`
}

type DeleteSecretQueryParams struct {
	Name  string `form:"name" binding:"required"`
	Scope string `form:"scope" binding:"required"`
}
