package importer

import (
	"fmt"
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock ManifestParser:ManifestParserMock

type ImportPackage interface {
	io.Closer
	GetResource(resourceName string) (io.ReadCloser, error)
}

type ManifestParser interface {
	Parse(input io.Reader) (*model.ImportManifest, error)
}

type ImportPackageProcessor struct {
	parser ManifestParser
}

func NewImportPackageProcessor(mp ManifestParser) *ImportPackageProcessor {
	return &ImportPackageProcessor{
		parser: mp,
	}
}

const manifestFileName = "manifest.yaml"

func (ipp *ImportPackageProcessor) Process(ip ImportPackage) error {

	defer ip.Close()

	_, err := ip.GetResource(manifestFileName)

	if err != nil {
		return fmt.Errorf("error accessing manifest: %w", err)
	}

	return nil
}
