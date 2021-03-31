package model

type Secret struct {
	Name  string `json:"name" binding:"required"`
	Scope string `json:"scope" binding:"required"`
	Data  Data   `json:"data"`
}

type Data map[string]string
