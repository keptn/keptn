package models

import (
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
)

type GetUniformIntegrationsParams struct {
	Name    string `form:"name" json:"name"`
	ID      string `form:"id" json:"id"`
	Project string `form:"project" json:"project"`
	Stage   string `form:"stage" json:"stage"`
	Service string `form:"service" json:"service"`
}

type GetUniformIntegrationByIDParams struct {
	ID string `form:"id" json:"id"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}
type UnregisterResponse struct{}

type Integration keptnmodels.Integration
