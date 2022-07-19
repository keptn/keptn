package importer

import (
	"fmt"
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock ManifestParser:ManifestParserMock TaskExecutor:TaskExecutorMock
//go:generate moq -pkg fake --skip-ensure -out ./fake/stage_retriever_mock.go . ProjectStageRetriever:MockStageRetriever

type RenderFunction func(raw io.ReadCloser, context any) (io.ReadCloser, error)

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
	render         RenderFunction
}

func NewImportPackageProcessor(
	mp ManifestParser, ex TaskExecutor, retriever ProjectStageRetriever,
) *ImportPackageProcessor {
	return newImportPackageProcessor(mp, ex, retriever, RenderContent)
}

func newImportPackageProcessor(
	mp ManifestParser, ex TaskExecutor, retriever ProjectStageRetriever, renderF RenderFunction,
) *ImportPackageProcessor {
	return &ImportPackageProcessor{
		parser:         mp,
		executor:       ex,
		stageRetriever: retriever,
		render:         renderF,
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

	mCtx := model.NewManifestExecution(project)

	for _, task := range manifest.Tasks {

		switch task.Type {
		case apiTaskType:
			if err = ipp.processAPITask(mCtx, ip, task); err != nil {
				return err
			}
		case resourceTaskType:
			if err = ipp.processResourceTask(mCtx, ip, task); err != nil {
				return err
			}
		default:
			return fmt.Errorf("task of type %s not implemented", task.Type)
		}
	}

	return nil
}

func (ipp *ImportPackageProcessor) processResourceTask(
	mCtx *model.ManifestExecution, ip ImportPackage, task *model.ManifestTask,
) error {
	if task.ResourceTask == nil {
		return fmt.Errorf("malformed task of type resource: %+v", task)
	}
	var stages []string
	var err error
	if task.ResourceTask.Stage == "" {
		stages, err = ipp.stageRetriever.GetStages(mCtx.GetProject())
		if err != nil {
			return fmt.Errorf("error retrieving stages for project %s: %w", mCtx.GetProject(), err)
		}
	} else {
		stages = []string{task.Stage}
	}

	for _, stage := range stages {
		resourcePush, err := ipp.mapResourcePush(mCtx, stage, ip, task)
		if err != nil {
			return fmt.Errorf("error setting up resource push for task ID %s: %w", task.ID, err)
		}
		response, err := ipp.executor.PushResource(resourcePush)
		if err != nil {
			return fmt.Errorf("resource task id %s failed: %w", task.ID, err)
		}
		// TODO what should we store for multiple stage resources upload ?
		mCtx.Tasks[task.ID] = model.TaskExecution{
			TaskContext: resourcePush.Context,
			Response:    response,
		}
	}
	return nil
}

func (ipp *ImportPackageProcessor) processAPITask(
	mCtx *model.ManifestExecution, ip ImportPackage, task *model.ManifestTask,
) error {
	if task.APITask == nil {
		return fmt.Errorf("malformed task of type api: %+v", task)
	}
	apiTaskExecution, err := ipp.mapAPITask(mCtx, ip, task)
	if err != nil {
		return fmt.Errorf("error setting up API task ID %s: %w", task.ID, err)
	}
	response, err := ipp.executor.ExecuteAPI(apiTaskExecution)
	if err != nil {
		return fmt.Errorf("execution of task %s failed: %w", task.ID, err)
	}

	// store task context and response into the manifest context
	mCtx.Tasks[task.ID] = model.TaskExecution{
		TaskContext: apiTaskExecution.Context,
		Response:    response,
	}

	return nil
}

func (ipp *ImportPackageProcessor) mapResourcePush(
	mCtx *model.ManifestExecution, stage string, ip ImportPackage,
	task *model.ManifestTask,
) (model.ResourcePush,
	error) {
	resource, err := ip.GetResource(task.ResourceTask.File)
	if err != nil {
		return model.ResourcePush{}, fmt.Errorf("error accessing resource content: %w", err)
	}

	taskContext := newTaskContext(mCtx, task)

	renderedResource, err := ipp.render(resource, taskContext)
	if err != nil {
		return model.ResourcePush{}, fmt.Errorf("error rendering resource content: %w", err)
	}

	ret := model.ResourcePush{
		Content:     renderedResource,
		ResourceURI: task.ResourceTask.RemoteURI,
		Stage:       stage,
		Service:     task.ResourceTask.Service,
		Context:     taskContext,
	}

	return ret, nil
}

func (ipp *ImportPackageProcessor) mapAPITask(
	mCtx *model.ManifestExecution, ip ImportPackage,
	task *model.ManifestTask,
) (model.APITaskExecution, error) {

	payload, err := ip.GetResource(task.APITask.PayloadFile)

	if err != nil {
		return model.APITaskExecution{}, fmt.Errorf(
			"error accessing payload %s: %w", task.APITask.PayloadFile,
			err,
		)
	}

	taskContext := newTaskContext(mCtx, task)

	renderedPayload, err := ipp.render(payload, taskContext)
	if err != nil {
		return model.APITaskExecution{}, fmt.Errorf(
			"error rendering payload %s: %w", task.APITask.PayloadFile,
			err,
		)
	}

	ret := model.APITaskExecution{
		Payload:    renderedPayload,
		EndpointID: task.APITask.Action,
		Context:    taskContext,
	}

	return ret, nil
}

func newTaskContext(mCtx *model.ManifestExecution, task *model.ManifestTask) model.TaskContext {
	renderedContext, _ := renderContext(mCtx, task.Context)
	// TODO handle error

	return model.TaskContext{
		Project: mCtx.GetProject(),
		Task:    task,
		Context: renderedContext,
	}
}

func renderContext(mCtx *model.ManifestExecution, context map[string]string) (map[string]string, error) {
	// TODO implement rendering of values
	return context, nil
}
