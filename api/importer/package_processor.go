package importer

import (
	"fmt"
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock ManifestParser:ManifestParserMock TaskExecutor:TaskExecutorMock

type ImportPackage interface {
	io.Closer
	GetResource(resourceName string) (io.ReadCloser, error)
}

type ManifestParser interface {
	Parse(input io.Reader) (*model.ImportManifest, error)
}

type TaskExecutor interface {
	ExecuteAPI(ate model.APITaskExecution) (any, error)
}

type ImportPackageProcessor struct {
	parser   ManifestParser
	executor TaskExecutor
}

func NewImportPackageProcessor(mp ManifestParser, ex TaskExecutor) *ImportPackageProcessor {
	return &ImportPackageProcessor{
		parser:   mp,
		executor: ex,
	}
}

const manifestFileName = "manifest.yaml"
const apiTaskType = "api"
const resourceTaskType = "resource"

func (ipp *ImportPackageProcessor) Process(project string, ip ImportPackage) error {

	defer ip.Close()

	manifestReader, err := ip.GetResource(manifestFileName)

	if err != nil {
		return fmt.Errorf("error accessing manifest: %w", err)
	}

	defer manifestReader.Close()

	manifest, err := ipp.parser.Parse(manifestReader)

	if err != nil {
		return fmt.Errorf("error parsing manifest: %w", err)
	}
	// TODO add manifestValidation

	for _, task := range manifest.Tasks {

		switch task.Type {
		case apiTaskType:
			apiTaskExecution, err := mapAPITask(project, ip, task)
			if err != nil {
				return fmt.Errorf("error setting up API task ID %s: %w", task.ID, err)
			}
			_, err = ipp.executor.ExecuteAPI(apiTaskExecution)
			if err != nil {
				return fmt.Errorf("execution of task %s failed: %w", task.ID, err)
			}
		default:
			return fmt.Errorf("task of type %s not implemented", task.Type)
		}
	}

	return nil
}

func mapAPITask(project string, ip ImportPackage, task *model.ManifestTask) (model.APITaskExecution, error) {

	resource, err := ip.GetResource(task.APITask.PayloadFile)

	// TODO who should close the resource ?

	if err != nil {
		return model.APITaskExecution{}, fmt.Errorf("error accessing resource %s: %w", task.APITask.PayloadFile, err)
	}

	ret := model.APITaskExecution{
		Payload:    resource,
		EndpointID: task.APITask.Action,
		Context: model.TaskContext{
			Project: project,
			Task:    task,
		},
	}

	return ret, nil
}
