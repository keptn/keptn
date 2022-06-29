package importer

import (
	"fmt"
	"io"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock

type ImportPackage interface {
	io.Closer
	GetResource(resourceName string) (io.ReadCloser, error)
}

type ImportPackageProcessor struct {
}

func NewImportPackageProcessor() *ImportPackageProcessor {
	return new(ImportPackageProcessor)
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
