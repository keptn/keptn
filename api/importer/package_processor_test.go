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

func TestImportPackageManifestRetrievedAndPackageClosed(t *testing.T) {
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}
	sut := NewImportPackageProcessor(parserMock)
	closeCallCount := 0
	requestedResources := map[string]int{}

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			closeCallCount++
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {
			requestedResources[resourceName] = requestedResources[resourceName] + 1
			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
	}
	err := sut.Process(importPackageMock)
	require.NoError(t, err)
	assert.Equal(t, 1, closeCallCount)
	assert.Equal(t, 1, requestedResources["manifest.yaml"])
	assert.Len(t, parserMock.ParseCalls(), 0)
}

func TestErrorImportPackageWhenManifestCannotBeRetrieved(t *testing.T) {
	parserMock := &fake.ManifestParserMock{
		ParseFunc: func(input io.Reader) (*model.ImportManifest, error) {
			return new(model.ImportManifest), nil
		},
	}

	sut := NewImportPackageProcessor(parserMock)

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
}
