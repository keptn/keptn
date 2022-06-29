package importer

import (
	"fmt"
	"io"

	"github.com/keptn/keptn/api/importer/model"
)

//go:generate moq -pkg fake --skip-ensure -out ./fake/package_processor_mock.go . ImportPackage:ImportPackageMock ManifestParser:ManifestParserMock TaskExecutor:TaskExecutorMock

type ImportPackage interface {
	io.Closer
	GetResource(resourceName string) (io.ReadCloser, error)
}

type ManifestParser interface {
	Parse(input io.Reader) (*model.ImportManifest, error)
}

type TaskExecutor interface {
	Execute(task *model.ManifestTask) (any, error)
}

type ImportPackageProcessor struct {
	parser   ManifestParser
	executor TaskExecutor
}

func NewImportPackageProcessor(mp ManifestParser, ex TaskExecutor) *ImportPackageProcessor {
	return &ImportPackageProcessor{
		parser:   mp,
		executor: ex,
	}
}

const manifestFileName = "manifest.yaml"

func (ipp *ImportPackageProcessor) Process(ip ImportPackage) error {

	defer ip.Close()

	manifestReader, err := ip.GetResource(manifestFileName)

	if err != nil {
		return fmt.Errorf("error accessing manifest: %w", err)
	}

	defer manifestReader.Close()

	manifest, err := ipp.parser.Parse(manifestReader)

	if err != nil {
		return fmt.Errorf("error parsing manifest: %w", err)
	}
	// TODO add manifestValidation

	for _, task := range manifest.Tasks {
		_, err := ipp.executor.Execute(task)
		if err != nil {
			return fmt.Errorf("execution of task %s failed: %w", task.ID, err)
		}
	}

	return nil
}
