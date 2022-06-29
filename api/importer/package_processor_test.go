package importer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		ExecuteFunc: func(task *model.ManifestTask) (any, error) {
			return nil, nil
		},
	}
	sut := NewImportPackageProcessor(parserMock, taskExecutor)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
	}
	err := sut.Process(importPackageMock)
	require.NoError(t, err)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteCalls(), 0)
}

func TestErrorImportPackageWhenManifestCannotBeRetrieved(t *testing.T) {
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteFunc: func(task *model.ManifestTask) (any, error) {
			return nil, nil
		},
	}

	sut := NewImportPackageProcessor(parserMock, taskExecutor)

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
	}
	err := sut.Process(importPackageMock)
	assert.ErrorIs(t, err, errorManifestAccess)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 0)
	assert.Len(t, taskExecutor.ExecuteCalls(), 0)
}

func TestErrorImportPackageWhenManifestCannotBeParsed(t *testing.T) {
	parsingError := errors.New("error parsing manifest")
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return nil, parsingError
		},
	}

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteFunc: func(task *model.ManifestTask) (any, error) {
			return nil, nil
		},
	}

	sut := NewImportPackageProcessor(parserMock, taskExecutor)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
	}
	err := sut.Process(importPackageMock)
	assert.ErrorIs(t, err, parsingError)
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteCalls(), 0)
}

func TestErrorImportPackageWhenTaskFails(t *testing.T) {

	firstTask := &model.ManifestTask{
		// random filler task to check that we execute in order until the failure
		APITask:      &model.APITask{},
		ResourceTask: nil,
		ID:           "firsttask",
		Type:         "api",
		Name:         "FirstTask",
	}

	failingTask := &model.ManifestTask{
		APITask: &model.APITask{
			Action:      "fail",
			PayloadFile: "somefile/somewhere.json",
		},
		ResourceTask: nil,
		ID:           "sometask",
		Type:         "api",
		Name:         "SomeTask",
	}

	neverExecutedTask := &model.ManifestTask{
		// random filler task to check that we stop executing at the first failure
		APITask:      &model.APITask{},
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

	taskExecutor := &fake.TaskExecutorMock{
		ExecuteFunc: func(task *model.ManifestTask) (any, error) {
			if task.Type == "api" && task.APITask.Action == "fail" {
				return nil, taskError
			}

			return nil, nil
		},
	}
	sut := NewImportPackageProcessor(parserMock, taskExecutor)

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
	}
	err := sut.Process(importPackageMock)
	assert.ErrorIs(t, err, taskError)
	assert.ErrorContains(t, err, "execution of task sometask failed")
	assert.Len(t, importPackageMock.CloseCalls(), 1)
	assert.ElementsMatch(
		t, importPackageMock.GetResourceCalls(), []struct {
			ResourceName string
		}{{ResourceName: manifestFileName}},
	)
	assert.Len(t, parserMock.ParseCalls(), 1)
	assert.Len(t, taskExecutor.ExecuteCalls(), 2)
	assert.Equal(t, []struct{ Task *model.ManifestTask }{{firstTask}, {failingTask}}, taskExecutor.ExecuteCalls())
}
