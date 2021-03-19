package repository

import (
	"github.com/ghodss/yaml"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"io/ioutil"
)

const ScopesConfigurationFile = "/scopes.yaml"

//go:generate moq -pkg fake -out ./fake/scopesrepository_mock.go . ScopesRepository
type ScopesRepository interface {
	Read() (model.Scopes, error)
}

type FileBasedScopesRepository struct {
	FileReader func(filename string) ([]byte, error)
	Decoder    func(y []byte, o interface{}) error
}

func NewFileBasedScopesRepository() *FileBasedScopesRepository {
	return &FileBasedScopesRepository{
		FileReader: ioutil.ReadFile,
		Decoder:    yaml.Unmarshal,
	}
}

func (s FileBasedScopesRepository) Read() (model.Scopes, error) {
	content, err := s.FileReader(ScopesConfigurationFile)
	if err != nil {
		return model.Scopes{}, err
	}

	scopes := model.Scopes{}
	if err := s.Decoder(content, &scopes); err != nil {
		return model.Scopes{}, err
	}

	return scopes, nil
}
