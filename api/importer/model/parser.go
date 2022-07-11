package model

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type YAMLManifestUnMarshaler struct {
}

func (ymum *YAMLManifestUnMarshaler) Parse(input io.Reader) (*ImportManifest, error) {
	manifestBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %w", err)
	}
	im := new(ImportManifest)
	err = yaml.Unmarshal(manifestBytes, im)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling yaml manifest: %w", err)
	}
	return im, err
}
