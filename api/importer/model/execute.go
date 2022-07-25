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

type ManifestExecution struct {
	Inputs       map[string]string
	Tasks        map[string]TaskExecution
	TaskSequence []string
}

func (mc ManifestExecution) GetProject() string {
	return mc.Inputs[projectInputContextKey]
}

func NewManifestExecution(project string) *ManifestExecution {
	return &ManifestExecution{
		Inputs:       map[string]string{projectInputContextKey: project},
		Tasks:        map[string]TaskExecution{},
		TaskSequence: nil,
	}
}
