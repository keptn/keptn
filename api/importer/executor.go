package importer

import (
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

type KeptnAPIExecutor struct {
}

func (kae *KeptnAPIExecutor) Execute(task *model.ManifestTask) (any, error) {
	panic("Implement me!")
}

type APITaskExecution struct {
	payload    io.Reader
	project    string
	endpointID string
}