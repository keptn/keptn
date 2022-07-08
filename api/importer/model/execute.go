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
