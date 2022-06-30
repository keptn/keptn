package model

import "io"

type TaskContext struct {
	Project string
	Task    *ManifestTask
}

type APITaskExecution struct {
	Payload    io.Reader
	EndpointID string
	Context    TaskContext
}
