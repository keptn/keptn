package importer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/keptn/keptn/api/importer/fake"
	"github.com/keptn/keptn/api/importer/model"
)

func TestImportPackageEmptyManifestRetrievedAndPackageClosed(t *testing.T) {
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	require.NoError(t, err)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 0)
}

func TestErrorImportPackageWhenManifestCannotBeRetrieved(t *testing.T) {
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	errorManifestAccess := errors.New("error retrieving manifest.yaml")

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {

			if resourceName == "manifest.yaml" {
				return nil, errorManifestAccess
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	assert.ErrorIs(t, err, errorManifestAccess)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 0)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 0)
}

func TestErrorImportPackageWhenManifestCannotBeParsed(t *testing.T) {
	parsingError := errors.New("error parsing manifest")
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return nil, parsingError
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	assert.ErrorIs(t, err, parsingError)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 0)
}

func TestErrorImportPackageWhenManifestResourceNotFound(t *testing.T) {

	const missingFileName = "non-existing-file.json"
	taskWithMissingResource := &model.ManifestTask{
		APITask: &model.APITask{
			Action:      "keptn-api-v1-uniform-create-webhook-subscription",
			PayloadFile: missingFileName,
		},
		ResourceTask: nil,
		ID:           "missingresourcetask",
		Type:         "api",
		Name:         "Missing Resource Task",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					taskWithMissingResource,
				},
			}, nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	resourceError := errors.New("error retrieving resource manifest")

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			if resourceName == missingFileName {
				return nil, resourceError
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	assert.ErrorIs(t, err, resourceError)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}, {ResourceName: missingFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 0)
}

func TestErrorImportPackageWhenUnknownManifestTaskType(t *testing.T) {

	taskWithUnknownType := &model.ManifestTask{
		APITask:      nil,
		ResourceTask: nil,
		ID:           "misterytask",
		Type:         "weirdunknowntasktype",
		Name:         "Mistery Task",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					taskWithUnknownType,
				},
			}, nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	assert.ErrorContains(t, err, "task of type weirdunknowntasktype not implemented")
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 0)
}

func TestErrorImportPackageWhenTaskFails(t *testing.T) {

	firstTask := &model.ManifestTask{
		// random filler task to check that we execute in order until the failure
		APITask: &model.APITask{
			Action:      "keptn-api-v1-create-service",
			PayloadFile: "okfile/someok.json",
		},
		ResourceTask: nil,
		ID:           "firsttask",
		Type:         "api",
		Name:         "FirstTask",
	}

	failingTask := &model.ManifestTask{
		APITask: &model.APITask{
			Action:      "keptn-api-v1-uniform-create-webhook-subscription",
			PayloadFile: "somefile/somewhere.json",
		},
		ResourceTask: nil,
		ID:           "sometask",
		Type:         "api",
		Name:         "SomeTask",
	}

	neverExecutedTask := &model.ManifestTask{
		// random filler task to check that we stop executing at the first failure
		APITask: &model.APITask{
			Action:      "keptn-api-v1-uniform-create-webhook-subscription",
			PayloadFile: "somefile/somewhere.json",
		},
		ResourceTask: nil,
		ID:           "neverexecuted",
		Type:         "api",
		Name:         "NeverExecuted",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					firstTask,
					failingTask,
					neverExecutedTask,
				},
			}, nil
		},
	}

	taskError := errors.New("api task failed")

	var apiTasksExecuted []string

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			apiTasksExecuted = append(apiTasksExecuted, ate.Context.Task.ID)
			if ate.Context.Task.Type == "api" && ate.Context.Task.APITask.Action == "keptn-api-v1-uniform-create-webhook-subscription" {
				return nil, taskError
			}

			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}
	err := sut.Process("project", importPackageMock)
	assert.ErrorIs(t, err, taskError)
	assert.ErrorContains(t, err, "execution of task sometask failed")
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{
			{ResourceName: manifestFileName},
			{ResourceName: "okfile/someok.json"},
			{ResourceName: "somefile/somewhere.json"},
		},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteAPICalls(), 2)
	assert.Equal(t, []string{"firsttask", "sometask"}, apiTasksExecuted)
}

func TestImportPackageProcessor_Process_ResourceTask(t *testing.T) {
	const resourceFileName = "somelocation/somefile.pcap"

	resourceTask := &model.ManifestTask{
		ResourceTask: &model.ResourceTask{
			File:      resourceFileName,
			RemoteURI: "/wireshark/capture.pcap",
			Stage:     "dev",
			Service:   "service",
		},
		ID:   "res_task",
		Type: "resource",
		Name: "ResTask",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					resourceTask,
				},
			}, nil
		},
	}

	const project = "somekeptnproject"

	resourceContentString := "some fancy binary content here"

	taskExecutor := &fake.TaskExecutorMock{
		PushResourceFunc: func(rp model.ResourcePush) (any, error) {
			assert.Equal(t, resourceTask.ResourceTask.Service, rp.Service)
			assert.Equal(t, resourceTask.ResourceTask.Stage, rp.Stage)
			assert.Equal(t, resourceTask.ResourceTask.RemoteURI, rp.ResourceURI)
			assert.Equal(t, project, rp.Context.Project)
			assert.Equal(t, resourceTask, rp.Context.Task)
			actualResourceContent, err := io.ReadAll(rp.Content)
			require.NoError(t, err)
			assert.Equal(t, []byte(resourceContentString), actualResourceContent)

			return &struct{}{}, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	resourceContentReader := io.NopCloser(strings.NewReader(resourceContentString))
	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			if resourceName == resourceFileName {
				return resourceContentReader, nil
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	err := sut.Process(project, importPackageMock)

	assert.NoError(t, err)
	assert.Len(t, taskExecutor.PushResourceCalls(), 1)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{
			{ResourceName: manifestFileName},
			{ResourceName: resourceFileName},
		},
	)
}

func TestImportPackageProcessor_Process_ResourceTask_AllStages(t *testing.T) {
	const resourceFileName = "somelocation/somefile.pcap"

	resourceTask := &model.ManifestTask{
		ResourceTask: &model.ResourceTask{
			File:      resourceFileName,
			RemoteURI: "/wireshark/capture.pcap",
			Stage:     "",
			Service:   "service",
		},
		ID:   "res_task",
		Type: "resource",
		Name: "ResTask",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					resourceTask,
				},
			}, nil
		},
	}

	const project = "somekeptnproject"

	stages := []string{"dev", "test", "prod"}
	var actualStages []string

	taskExecutor := &fake.TaskExecutorMock{
		PushResourceFunc: func(rp model.ResourcePush) (any, error) {
			assert.NotEmpty(t, rp.Stage)
			actualStages = append(actualStages, rp.Stage)
			return &struct{}{}, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{
		GetStagesFunc: func(project string) ([]string, error) {
			return stages, nil
		},
	}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	err := sut.Process(project, importPackageMock)

	assert.NoError(t, err)
	assert.Len(t, taskExecutor.PushResourceCalls(), len(stages))
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	expectedGetResourceArgs := []struct {
		ResourceName string
	}{
		{ResourceName: manifestFileName},
	}
	for i := 0; i < len(stages); i++ {
		expectedGetResourceArgs = append(
			expectedGetResourceArgs,
			struct {
				ResourceName string
			}{ResourceName: resourceFileName},
		)
	}

	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), expectedGetResourceArgs,
	)
}

func TestImportPackageProcessor_Process_WebhookConfigWithTemplating(t *testing.T) {

	const resourceFileName = "webhook.yaml"

	const rawWebhookConfigResourceFile = "../test/data/import/sample-package/resources/webhook.yaml"
	const webhookConfigResourceRenderContext = "../test/data/import/rendered-sample-package/simple-task-rendering/resources/webhook.context.yaml"
	const renderedWebhookConfigResourceFile = "../test/data/import/rendered-sample-package/simple-task-rendering/resources/webhook.yaml"

	context := map[string]string{}

	contextbytes, err := ioutil.ReadFile(webhookConfigResourceRenderContext)
	require.NoError(t, err)

	err = yaml.Unmarshal(contextbytes, &context)
	require.NoError(t, err)

	resourceTask := &model.ManifestTask{
		ResourceTask: &model.ResourceTask{
			File:      rawWebhookConfigResourceFile,
			RemoteURI: resourceFileName,
			Stage:     "dev",
			Service:   "service",
		},
		ID:      "res_task",
		Type:    "resource",
		Name:    "ResTask",
		Context: context,
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					resourceTask,
				},
			}, nil
		},
	}

	const project = "somekeptnproject"

	taskExecutor := &fake.TaskExecutorMock{
		PushResourceFunc: func(rp model.ResourcePush) (any, error) {
			assert.NotNil(t, rp.Content)
			defer rp.Content.Close()

			renderedBytes, err := io.ReadAll(rp.Content)
			require.NoError(t, err)

			expectedRenderedBytes, err := ioutil.ReadFile(renderedWebhookConfigResourceFile)
			require.NoError(t, err)

			assert.Equal(t, string(expectedRenderedBytes), string(renderedBytes))
			// For debugging purposes it may be easier to look at the YAML/JSON comparison
			// assert.YAMLEq(t, string(expectedRenderedBytes), string(renderedBytes))
			// assert.JSONEq(t, string(expectedRenderedBytes), string(renderedBytes))
			return nil, nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			if resourceName == manifestFileName {
				return io.NopCloser(bytes.NewReader([]byte{})), nil
			}
			return os.Open(resourceName)
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	err = sut.Process(project, importPackageMock)

	assert.NoError(t, err)
	assert.Len(t, taskExecutor.PushResourceCalls(), 1)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	expectedGetResourceArgs := []struct {
		ResourceName string
	}{
		{ResourceName: manifestFileName},
		{ResourceName: rawWebhookConfigResourceFile},
	}

	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), expectedGetResourceArgs,
	)
}

func TestImportPackageProcessor_Process_APITaskWithTemplating(t *testing.T) {

	const rawPayloadDir = "../test/data/import/sample-package/api/"
	const renderedPayloadDir = "../test/data/import/rendered-sample-package/simple-task-rendering/api/"

	tests := []struct {
		name                     string
		rawPayloadFile           string
		payloadRenderContextFile string
		action                   string
	}{
		{
			name:                     "Template create-service request",
			rawPayloadFile:           "create-service.json",
			payloadRenderContextFile: "create-service.context.yaml",
			action:                   "keptn-api-v1-create-service",
		},
		{
			name:                     "Template create-subscription request",
			rawPayloadFile:           "create-subscription.json",
			payloadRenderContextFile: "create-subscription.context.yaml",
			action:                   "keptn-api-v1-uniform-create-webhook-subscription",
		},
	}

	for i, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				context := map[string]string{}

				contextbytes, err := ioutil.ReadFile(path.Join(renderedPayloadDir, tt.payloadRenderContextFile))
				require.NoError(t, err)

				err = yaml.Unmarshal(contextbytes, &context)
				require.NoError(t, err)

				rawPayloadFullPath := path.Join(rawPayloadDir, tt.rawPayloadFile)
				apiTask := &model.ManifestTask{
					ID:      fmt.Sprintf("api_task_%d", i),
					Type:    "api",
					Name:    fmt.Sprintf("API Task No. %d", i),
					Context: context,
					APITask: &model.APITask{
						Action:      tt.action,
						PayloadFile: rawPayloadFullPath,
					},
				}

				parserMock := &fake.ManifestParserMock{
					ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
						return &model.ImportManifest{
							ApiVersion: "v1beta1",
							Tasks: []*model.ManifestTask{
								apiTask,
							},
						}, nil
					},
				}

				const project = "somekeptnproject"

				taskExecutor := &fake.TaskExecutorMock{
					ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
						assert.NotNil(t, ate.Payload)
						defer ate.Payload.Close()

						renderedBytes, err := io.ReadAll(ate.Payload)
						require.NoError(t, err)

						expectedRenderedBytes, err := ioutil.ReadFile(path.Join(renderedPayloadDir, tt.rawPayloadFile))
						require.NoError(t, err)

						assert.Equal(t, string(expectedRenderedBytes), string(renderedBytes))
						// For debugging purposes it may be easier to look at the YAML/JSON comparison
						// assert.YAMLEq(t, string(expectedRenderedBytes), string(renderedBytes))
						// assert.JSONEq(t, string(expectedRenderedBytes), string(renderedBytes))
						return nil, nil
					},
				}

				stageRetriever := &fake.MockStageRetriever{}
				sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

				importPackageMock := &fake.ImportPackageMock{
					CloseFunc: func() error {
						return nil
					},
					GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
						if resourceName == manifestFileName {
							return io.NopCloser(bytes.NewReader([]byte{})), nil
						}
						return os.Open(resourceName)
					},
					CheckIfResourceExistsFunc: func(resourceName string) error {
						return nil
					},
				}

				err = sut.Process(project, importPackageMock)

				assert.NoError(t, err)
				assert.Len(t, taskExecutor.ExecuteAPICalls(), 1)
				assert.Len(t, importPackageMock.CloseCalls(), 1)
				expectedGetResourceArgs := []struct {
					ResourceName string
				}{
					{ResourceName: manifestFileName},
					{ResourceName: rawPayloadFullPath},
				}

				assert.ElementsMatch(
					t, importPackageMock.GetResourceCalls(), expectedGetResourceArgs,
				)
			},
		)
	}
}

func TestImportPackageProcessor_ProcessResourceTask_ErrorGettingResource(t *testing.T) {
	const resourceFileName = "somelocation/somefile.pcap"

	resourceTask := &model.ManifestTask{
		ResourceTask: &model.ResourceTask{
			File:      resourceFileName,
			RemoteURI: "/wireshark/capture.pcap",
			Stage:     "dev",
			Service:   "service",
		},
		ID:   "res_task",
		Type: "resource",
		Name: "ResTask",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					resourceTask,
				},
			}, nil
		},
	}

	const project = "somekeptnproject"

	taskExecutor := &fake.TaskExecutorMock{}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			if resourceName == resourceFileName {
				return nil, errors.New("error getting resource")
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	err := sut.Process(project, importPackageMock)

	assert.Error(t, err)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{
			{ResourceName: manifestFileName},
			{ResourceName: resourceFileName},
		},
	)
}

func TestImportPackageProcessor_Process_ResourceTask_ErrorExecutingTask(t *testing.T) {
	const resourceFileName = "somelocation/somefile.pcap"

	resourceTask := &model.ManifestTask{
		ResourceTask: &model.ResourceTask{
			File:      resourceFileName,
			RemoteURI: "/wireshark/capture.pcap",
			Stage:     "dev",
			Service:   "service",
		},
		ID:   "res_task",
		Type: "resource",
		Name: "ResTask",
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return &model.ImportManifest{
				ApiVersion: "v1beta1",
				Tasks: []*model.ManifestTask{
					resourceTask,
				},
			}, nil
		},
	}

	const project = "somekeptnproject"

	resourceContentString := "some fancy binary content here"

	taskExecutor := &fake.TaskExecutorMock{
		PushResourceFunc: func(rp model.ResourcePush) (any, error) {
			return nil, errors.New("error executing resource push")
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	resourceContentReader := io.NopCloser(strings.NewReader(resourceContentString))
	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			if resourceName == resourceFileName {
				return resourceContentReader, nil
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	err := sut.Process(project, importPackageMock)

	assert.Error(t, err)
	assert.Len(t, taskExecutor.PushResourceCalls(), 1)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{
			{ResourceName: manifestFileName},
			{ResourceName: resourceFileName},
		},
	)
}

func TestImportPackageProcessor_Process_ErrorMalformedTasks(t *testing.T) {
	tests := []struct {
		name string
		task *model.ManifestTask
	}{
		{
			name: "malformed resource task",
			task: &model.ManifestTask{
				ResourceTask: nil,
				ID:           "res_task",
				Type:         "resource",
				Name:         "ResTask",
			},
		},
		{
			name: "malformed api task",
			task: &model.ManifestTask{
				APITask: nil,
				ID:      "api_task",
				Type:    "api",
				Name:    "APITask",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				parserMock := &fake.ManifestParserMock{
					ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
						return &model.ImportManifest{
							ApiVersion: "v1beta1",
							Tasks: []*model.ManifestTask{
								tt.task,
							},
						}, nil
					},
				}

				taskExecutor := &fake.TaskExecutorMock{}

				stageRetriever := &fake.MockStageRetriever{}
				sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)
				importPackageMock := &fake.ImportPackageMock{
					CloseFunc: func() error {
						return nil
					},
					GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte{})), nil
					},
					CheckIfResourceExistsFunc: func(resourceName string) error {
						return nil
					},
				}

				err := sut.Process("test-project", importPackageMock)
				assert.Error(t, err)
				assert.ErrorContains(t, err, fmt.Sprintf("empty %s definition not supported", tt.task.Type))
			},
		)
	}
}

func TestImportPackageProcessor_Process_ErrorRenderingContext(t *testing.T) {

	tests := []struct {
		name string
		task *model.ManifestTask
	}{
		{
			name: "Error rendering context for api task",
			task: &model.ManifestTask{
				ID:      "api_task",
				Type:    "api",
				Name:    "API Task",
				Context: map[string]string{"foo": "bar"},
				APITask: &model.APITask{
					Action:      "keptn-api-v1-uniform-create-webhook-subscription",
					PayloadFile: "payload.json",
				},
			},
		},
		{
			// this will fail during content rendering because there is no key in context ;)
			name: "Error rendering payload for api task",
			task: &model.ManifestTask{
				ID:      "api_task",
				Type:    "api",
				Name:    "API Task",
				Context: map[string]string{},
				APITask: &model.APITask{
					Action:      "keptn-api-v1-uniform-create-webhook-subscription",
					PayloadFile: "payload.json",
				},
			},
		},
		{
			name: "Error rendering context for resource task",
			task: &model.ManifestTask{
				ResourceTask: &model.ResourceTask{
					File:      "somefile.json",
					RemoteURI: "somefile.json",
					Stage:     "dev",
				},
				ID:      "res_task",
				Type:    "resource",
				Name:    "Resource Task",
				Context: map[string]string{"foo": "bar"},
			},
		},
		{
			// this will fail during content rendering because there is no key in context ;)
			name: "Error rendering content for resource task",
			task: &model.ManifestTask{
				ResourceTask: &model.ResourceTask{
					File:      "somefile.json",
					RemoteURI: "somefile.json",
					Stage:     "dev",
				},
				ID:      "res_task",
				Type:    "resource",
				Name:    "Resource Task",
				Context: map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				const renderErrorMessage = "error rendering some weird template"
				renderer := &fake.MockRenderer{
					RenderStringFunc: func(raw string, context any) (string, error) {
						return "", errors.New(renderErrorMessage)
					},
					RenderContentFunc: func(raw io.ReadCloser, context any) (io.ReadCloser, error) {
						return nil, errors.New(renderErrorMessage)
					},
				}

				task := tt.task

				parserMock := &fake.ManifestParserMock{
					ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
						return &model.ImportManifest{
							ApiVersion: "v1beta1",
							Tasks: []*model.ManifestTask{
								task,
							},
						}, nil
					},
				}

				taskExecutor := &fake.TaskExecutorMock{}
				stageRetriever := &fake.MockStageRetriever{}
				sut := newImportPackageProcessor(parserMock, taskExecutor, stageRetriever, renderer)

				importPackageMock := &fake.ImportPackageMock{
					CloseFunc: func() error {
						return nil
					},
					GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
						return io.NopCloser(bytes.NewReader([]byte{})), nil
					},
					CheckIfResourceExistsFunc: func(resourceName string) error {
						return nil
					},
				}

				err := sut.Process("test-project", importPackageMock)
				assert.Error(t, err)
				assert.ErrorContains(t, err, renderErrorMessage)
			},
		)
	}
}

// This test is written as a "component" (not really but close enough) test for the import core business logic
// (everything is the real thing except interaction with external systems which relies on canned responses through
// mocks). This is by necessity unwieldy and more complicated than a regular unit test,
// however it still delivers valuable information on import package behavior and it can be used as a testbed for
// manifest processing/rendering verification. Match of api and resource tasks is based on action and remoteURI,
// content/payload is checked against expectation and a not empty response is returned if defined, nil otherwise.
// The inputs, expectations and canned responses are defined in a .yaml file like
// ../../test/data/import/rendered-sample-package/input.yaml
func TestImportPackageProcessor_Process_FullManifestRendering(t *testing.T) {

	type apiCallExpectation struct {
		Action   string `yaml:"action"`
		Expected string `yaml:"expected"`
		Response string `yaml:"response"`
	}

	type resourceUploadExpectation struct {
		ResourceURI string `yaml:"uri"`
		Expected    string `yaml:"expected"`
	}

	type testData struct {
		API      map[string][]apiCallExpectation        `yaml:"api"`
		Resource map[string][]resourceUploadExpectation `yaml:"resource"`
	}

	type testInput struct {
		Project  string   `yaml:"project"`
		TestData testData `yaml:"test-data"`
	}

	const rawPackageDir = "../test/data/import/sample-package/"
	const renderedPackageDir = "../test/data/import/rendered-sample-package/full-manifest-rendering/"

	tests := []struct {
		name               string
		rawPackageDir      string
		renderedPackageDir string
		testInput          string
	}{
		{
			name:               "Template sample package",
			rawPackageDir:      rawPackageDir,
			renderedPackageDir: renderedPackageDir,
			testInput:          "input.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {

				inputFile, err := os.Open(path.Join(tt.renderedPackageDir, tt.testInput))
				require.NoError(t, err)
				defer inputFile.Close()

				inputs := testInput{}
				decoder := yaml.NewDecoder(inputFile)
				decoder.KnownFields(true)
				err = decoder.Decode(&inputs)
				require.NoError(t, err)

				taskExecutor := &fake.TaskExecutorMock{
					ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
						t.Logf("Handling API execution %+v", ate)
						expectations, ok := inputs.TestData.API[ate.EndpointID]
						if !ok || len(expectations) == 0 {
							t.Fatalf("API call %+v has no expectations available in inputs %+v", ate, inputs)
						}
						inputs.TestData.API[ate.EndpointID] = expectations[1:]
						expectation := expectations[0]
						t.Logf("Consuming expectation %+v", expectation)
						assert.NotNil(t, ate.Payload)
						defer ate.Payload.Close()

						renderedBytes, err := io.ReadAll(ate.Payload)
						require.NoError(t, err)

						expectedRenderedBytes, err := ioutil.ReadFile(
							path.Join(
								tt.renderedPackageDir, expectation.Expected,
							),
						)
						require.NoError(t, err)

						assert.Equal(t, string(expectedRenderedBytes), string(renderedBytes))
						// For debugging purposes it may be easier to look at the YAML/JSON comparison
						// assert.YAMLEq(t, string(expectedRenderedBytes), string(renderedBytes))
						// assert.JSONEq(t, string(expectedRenderedBytes), string(renderedBytes))
						var retval any
						if expectation.Response != "" {
							err = json.Unmarshal([]byte(expectation.Response), &retval)
							require.NoErrorf(
								t, err, "canned response for %s is not parseable JSON: %s",
								expectation.Action, expectation.Response,
							)
						}
						return retval, nil
					},
					PushResourceFunc: func(rp model.ResourcePush) (any, error) {
						t.Logf("Handling resource push %+v", rp)
						expectations, ok := inputs.TestData.Resource[rp.ResourceURI]
						if !ok || len(expectations) == 0 {
							t.Fatalf("Resource upload %+v has no expectations available in inputs %+v", rp, inputs)
						}
						inputs.TestData.Resource[rp.ResourceURI] = expectations[1:]
						expectation := expectations[0]

						t.Logf("Consuming expectation %+v", expectation)

						assert.NotNil(t, rp.Content)
						defer rp.Content.Close()

						renderedBytes, err := io.ReadAll(rp.Content)
						require.NoError(t, err)

						expectedRenderedBytes, err := ioutil.ReadFile(
							path.Join(
								tt.renderedPackageDir, expectation.Expected,
							),
						)
						require.NoError(t, err)

						assert.Equal(t, string(expectedRenderedBytes), string(renderedBytes))
						// For debugging purposes it may be easier to look at the YAML/JSON comparison
						// assert.YAMLEq(t, string(expectedRenderedBytes), string(renderedBytes))
						// assert.JSONEq(t, string(expectedRenderedBytes), string(renderedBytes))

						return nil, nil
					},
				}

				stageRetriever := &fake.MockStageRetriever{}
				parser := &model.YAMLManifestUnMarshaler{}
				sut := NewImportPackageProcessor(parser, taskExecutor, stageRetriever)

				importPackageMock := &fake.ImportPackageMock{
					CloseFunc: func() error {
						return nil
					},
					GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
						return os.Open(path.Join(tt.rawPackageDir, resourceName))
					},
					CheckIfResourceExistsFunc: func(resourceName string) error {
						return nil
					},
				}

				err = sut.Process(inputs.Project, importPackageMock)

				assert.NoError(t, err)
				assert.Len(t, importPackageMock.CloseCalls(), 1)
			},
		)
	}
}

func TestPackageValidationEmptyId(t *testing.T) {
	manifest := &model.ImportManifest{
		ApiVersion: "v1beta1",
		Tasks: []*model.ManifestTask{
			{
				ID:      "",
				Type:    "api",
				Context: nil,
				APITask: &model.APITask{
					Action:      "keptn-api-v1-create-service",
					PayloadFile: "",
				},
			},
		},
	}

	expectedErrorMsg := "task id cannot be empty"

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	err := sut.validateManifest(manifest, importPackageMock)

	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func TestPackageValidationIdAlphaNumeric(t *testing.T) {
	manifests := []model.ImportManifest{
		{
			ApiVersion: "v1beta1",
			Tasks: []*model.ManifestTask{
				{
					ID:      "valid_id",
					Type:    "api",
					Context: nil,
					APITask: &model.APITask{
						Action:      "keptn-api-v1-create-service",
						PayloadFile: "",
					},
				},
			},
		},
		{
			ApiVersion: "v1beta1",
			Tasks: []*model.ManifestTask{
				{
					ID:      "invalid-id",
					Type:    "api",
					Context: nil,
					APITask: &model.APITask{
						Action:      "keptn-api-v1-create-service",
						PayloadFile: "",
					},
				},
			},
		},
		{
			ApiVersion: "v1beta1",
			Tasks: []*model.ManifestTask{
				{
					ID:      "invalid~char#",
					Type:    "api",
					Context: nil,
					APITask: &model.APITask{
						Action:      "keptn-api-v1-create-service",
						PayloadFile: "",
					},
				},
			},
		},
	}

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	tests := []struct {
		name               string
		manifest           model.ImportManifest
		expectedErrMessage string
		expectError        bool
	}{
		{
			name:        "valid id",
			manifest:    manifests[0],
			expectError: false,
		},
		{
			name:               "containing invalid dash",
			manifest:           manifests[1],
			expectedErrMessage: fmt.Sprintf("task id %s can only consist of alphnumeric characters and underscores", manifests[1].Tasks[0].ID),
			expectError:        true,
		},
		{
			name:               "containing invalid character",
			manifest:           manifests[2],
			expectedErrMessage: fmt.Sprintf("task id %s can only consist of alphnumeric characters and underscores", manifests[2].Tasks[0].ID),
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := sut.validateManifest(&tt.manifest, importPackageMock)
				if tt.expectError {
					assert.EqualErrorf(t, err, tt.expectedErrMessage, "Error should be: %v, got: %v", tt.expectedErrMessage, err)
				} else {
					require.NoError(t, err)
				}
			})
	}

}

func TestPackageValidationInvalidType(t *testing.T) {
	manifest := &model.ImportManifest{
		ApiVersion: "v1beta1",
		Tasks: []*model.ManifestTask{
			{
				ID:      "testing_invalid_type",
				Type:    "invalid",
				Context: nil,
				APITask: &model.APITask{
					Action:      "keptn-api-v1-create-service",
					PayloadFile: "",
				},
			},
		},
	}

	expectedErrorMsg := fmt.Sprintf("task of type %s not implemented", manifest.Tasks[0].Type)

	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteAPIFunc: func(ate model.APITaskExecution) (any, error) {
			return nil, nil
		},
	}

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
		CheckIfResourceExistsFunc: func(resourceName string) error {
			return nil
		},
	}

	stageRetriever := &fake.MockStageRetriever{}
	sut := NewImportPackageProcessor(parserMock, taskExecutor, stageRetriever)

	err := sut.validateManifest(manifest, importPackageMock)

	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
