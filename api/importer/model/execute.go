package model

import "io"

type TaskContext struct {
	Project string
	Task    *ManifestTask
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
