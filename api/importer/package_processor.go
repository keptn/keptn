package importer

import (
	"fmt"
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock ManifestParser:ManifestParserMock TaskExecutor:TaskExecutorMock
//go:generate moq -pkg fake --skip-ensure -out ./fake/stage_retriever_mock.go . ProjectStageRetriever:MockStageRetriever
type ImportPackage interface {
	io.Closer
	GetResource(resourceName string) (io.ReadCloser, error)
}

type ManifestParser interface {
	Parse(input io.Reader) (*model.ImportManifest, error)
}

type TaskExecutor interface {
	ExecuteAPI(ate model.APITaskExecution) (any, error)
	PushResource(rp model.ResourcePush) (any, error)
}

type ProjectStageRetriever interface {
	GetStages(project string) ([]string, error)
}

type ImportPackageProcessor struct {
	parser         ManifestParser
	executor       TaskExecutor
	stageRetriever ProjectStageRetriever
}

func NewImportPackageProcessor(
	mp ManifestParser, ex TaskExecutor, retriever ProjectStageRetriever,
) *ImportPackageProcessor {
	return &ImportPackageProcessor{
		parser:         mp,
		executor:       ex,
		stageRetriever: retriever,
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

	for _, task := range manifest.Tasks {

		switch task.Type {
		case apiTaskType:
			if err = ipp.processAPITask(project, ip, task); err != nil {
				return err
			}
		case resourceTaskType:
			if err = ipp.processResourceTask(project, ip, task); err != nil {
				return err
			}
		default:
			return fmt.Errorf("task of type %s not implemented", task.Type)
		}
	}

	return nil
}

func (ipp *ImportPackageProcessor) processResourceTask(
	project string, ip ImportPackage, task *model.ManifestTask,
) error {
	if task.ResourceTask == nil {
		return fmt.Errorf("malformed task of type resource: %+v", task)
	}
	var stages []string
	var err error
	if task.ResourceTask.Stage == "" {
		stages, err = ipp.stageRetriever.GetStages(project)
		if err != nil {
			return fmt.Errorf("error retrieving stages for project %s: %w", project, err)
		}
	} else {
		stages = []string{task.Stage}
	}

	for _, stage := range stages {
		resourcePush, err := mapResourcePush(project, stage, ip, task)
		if err != nil {
			return fmt.Errorf("error setting up resource push for task ID %s: %w", task.ID, err)
		}
		_, err = ipp.executor.PushResource(resourcePush)
		if err != nil {
			return fmt.Errorf("resource task id %s failed: %w", task.ID, err)
		}
	}
	return nil
}

func (ipp *ImportPackageProcessor) processAPITask(project string, ip ImportPackage, task *model.ManifestTask) error {
	if task.APITask == nil {
		return fmt.Errorf("malformed task of type api: %+v", task)
	}
	apiTaskExecution, err := mapAPITask(project, ip, task)
	if err != nil {
		return fmt.Errorf("error setting up API task ID %s: %w", task.ID, err)
	}
	_, err = ipp.executor.ExecuteAPI(apiTaskExecution)
	if err != nil {
		return fmt.Errorf("execution of task %s failed: %w", task.ID, err)
	}
	return nil
}

func mapResourcePush(project string, stage string, ip ImportPackage, task *model.ManifestTask) (model.ResourcePush,
	error) {
	resource, err := ip.GetResource(task.ResourceTask.File)
	if err != nil {
		return model.ResourcePush{}, fmt.Errorf("error accessing resource content: %w", err)
	}

	ret := model.ResourcePush{
		Content:     resource,
		ResourceURI: task.ResourceTask.RemoteURI,
		Stage:       stage,
		Service:     task.ResourceTask.Service,
		Context: model.TaskContext{
			Project: project,
			Task:    task,
		},
	}

	return ret, nil
}

func mapAPITask(project string, ip ImportPackage, task *model.ManifestTask) (model.APITaskExecution, error) {

	resource, err := ip.GetResource(task.APITask.PayloadFile)

	if err != nil {
		return model.APITaskExecution{}, fmt.Errorf(
			"error accessing resource %s: %w", task.APITask.PayloadFile,
			err,
		)
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
