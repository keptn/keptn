package models

type GetUniformIntegrationsParams struct {
	Name      string `form:"name" json:"name"`
	ID        string `form:"id" json:"id"`
	Namespace string `form:"namespace" json:"namespace"`
	HostName  string `form:"hostname" json:"hostname"`
	Project   string `form:"project" json:"project"`
	Stage     string `form:"stage" json:"stage"`
	Service   string `form:"service" json:"service"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

type CreateSubscriptionResponse struct {
	ID string `json:"id"`
}
type UnregisterResponse struct{}

type DeleteSubscriptionResponse struct{}
