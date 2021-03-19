package model

type Secret struct {
	Name  string `json:"name"`
	Scope string `json:"scope,omitempty"`
	Data  Data   `json:"data"`
}

type Data map[string]string
