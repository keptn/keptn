package repository

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_ReadFromFileBasedScopesRepository(t *testing.T) {
	testScopes := testScopes()
	fakeReader := func(filename string) ([]byte, error) {
		return asYAML(testScopes), nil
	}

	repository := FileBasedScopesRepository{
		FileReader: fakeReader,
		Decoder:    yaml.Unmarshal,
	}
	scopes, err := repository.Read()
	assert.Nil(t, err)
	assert.Equal(t, testScopes, scopes)
}

func Test_ReadFromFileBasedRepository_InvalidContent(t *testing.T) {
	fakeReader := func(filename string) ([]byte, error) {
		return []byte("invalid stuff"), nil
	}
	repository := FileBasedScopesRepository{
		FileReader: fakeReader,
		Decoder:    yaml.Unmarshal,
	}
	scopes, err := repository.Read()
	assert.NotNil(t, err)
	assert.Equal(t, model.Scopes{}, scopes)
}

func Test_ReadFromFileBasedRepository_ReadFails(t *testing.T) {
	fakeReader := func(filename string) ([]byte, error) {
		return []byte{}, fmt.Errorf("read failed")
	}

	repository := FileBasedScopesRepository{
		FileReader: fakeReader,
		Decoder:    yaml.Unmarshal,
	}
	scopes, err := repository.Read()
	assert.NotNil(t, err)
	assert.Equal(t, model.Scopes{}, scopes)
}

type failingReader struct{}

func (r failingReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func testScopes() model.Scopes {
	return model.Scopes{
		Scopes: map[string]model.Scope{
			"my-scope": {
				Capabilities: map[string]model.Capability{
					"my-scope-read-secrets": {
						Permissions: []string{"read"},
					},
					"my-scope-manage-secrets": {
						Permissions: []string{"create", "read", "update"},
					},
				},
			},
		},
	}
}

func asYAML(scopes model.Scopes) []byte {
	yamled, _ := yaml.Marshal(scopes)
	return yamled
}
