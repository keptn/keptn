package importer

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/keptn/keptn/api/importer/model"
)

type YAMLManifestUnMarshaler struct {
}

func (ymum *YAMLManifestUnMarshaler) Parse(input io.Reader) (*model.ImportManifest, error) {
	manifestBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %w", err)
	}
	im := new(model.ImportManifest)
	err = yaml.Unmarshal(manifestBytes, im)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml manifest: %w", err)
	}
	return im, err
}
