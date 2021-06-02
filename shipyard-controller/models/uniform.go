package models

import (
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
)

type GetUniformIntegrationParams struct {
	Name    string `form:"name" json:"name"`
	ID      string `form:"id" json:"id"`
	Project string `form:"project" json:"project"`
	Stage   string `form:"stage" json:"stage"`
	Service string `form:"service" json:"service"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}
type UnregisterResponse struct{}

type Integration keptnmodels.Integration
