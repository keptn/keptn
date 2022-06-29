package importer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/importer/fake"
)

func TestImportPackageManifestRetrievedAndPackageClosed(t *testing.T) {
	sut := NewImportPackageProcessor()
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
}

func TestErrorImportPackageWhenManifestCannotBeRetrieved(t *testing.T) {
	sut := NewImportPackageProcessor()
	closeCallCount := 0
	requestedResources := map[string]int{}

	errorManifestAccess := errors.New("error retrieving manifest.yaml")

	importPackageMock := &fake.ImportPackageMock{
		CloseFunc: func() error {
			closeCallCount++
			return nil
		},
		GetResourceFunc: func(resourceName string) (io.ReadCloser, error) {

			requestedResources[resourceName] = requestedResources[resourceName] + 1

			if resourceName == "manifest.yaml" {
				return nil, errorManifestAccess
			}

			return io.NopCloser(bytes.NewReader([]byte{})), nil
		},
	}
	err := sut.Process(importPackageMock)
	assert.ErrorIs(t, err, errorManifestAccess)
	assert.Equal(t, 1, closeCallCount)
	assert.Equal(t, 1, requestedResources["manifest.yaml"])
}
