package model

import (
	"io"
)

type TaskContext struct {
	Project string
	Task    *ManifestTask
	Context map[string]string
}

type APITaskExecution struct {
	Payload    io.ReadCloser
	EndpointID string
	Context    TaskContext
}

type ResourcePush struct {
	Content     io.ReadCloser
	ResourceURI string
	Stage       string
	Service     string
	Context     TaskContext
}

type TaskExecution struct {
	TaskContext
	Response any
}

const projectInputContextKey = "project"

const (
	CreateServiceAction = "keptn-api-v1-create-service"
	CreateSecretAction  = "keptn-api-v1-uniform-create-secret"
	CreateWebhookAction = "keptn-api-v1-uniform-create-webhook-subscription"
)

var AllActions = []string{CreateServiceAction, CreateSecretAction, CreateWebhookAction}

type ManifestExecution struct {
	Inputs map[string]string
	Tasks  map[string]TaskExecution
}

func (mc ManifestExecution) GetProject() string {
	return mc.Inputs[projectInputContextKey]
}

func NewManifestExecution(project string) *ManifestExecution {
	return &ManifestExecution{
		Inputs: map[string]string{projectInputContextKey: project},
		Tasks:  map[string]TaskExecution{},
	}
}
