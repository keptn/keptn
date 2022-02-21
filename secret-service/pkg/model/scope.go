package model

type Scopes struct {
	Scopes map[string]Scope `yaml:"scopes"`
}

type Scope struct {
	Capabilities map[string]Capability `yaml:"capabilities"`
}

type Capability struct {
	Permissions []string `yaml:"permissions"`
}

type GetScopesResponse struct {
	Scopes []string `json:"scopes"`
}
