package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/api/importer"
	"github.com/keptn/keptn/api/importer/model"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"

	logger "github.com/sirupsen/logrus"
)

const defaultImportArchiveExtension = ".zip"

//go:generate moq -pkg handlers_mock --skip-ensure -out ./fake/projectchecker_mock.go . projectChecker:ProjectCheckerMock importPackageProcessor:MockImportPackageProcessor

type projectChecker interface {
	ProjectExists(projectName string) (bool, error)
}

type importPackageProcessor interface {
	Process(project string, ip importer.ImportPackage) error
}

// ParseArchiveFunction is the function called to parse the uploaded file
type ParseArchiveFunction func(string, uint64) (importer.ImportPackage, error)

// ImportHandler is the rest handler for the /import endpoint
type ImportHandler struct {
	checker                    projectChecker
	tempStorageDir             string
	maxUncompressedPackageSize uint64
	parseArchive               ParseArchiveFunction
	processor                  importPackageProcessor
}

// GetImportHandlerFunc will instantiate a configured ImportHandler and return the method that can be used for
// handling http requests to the endpoint. See restapi.configureAPI for usage
func GetImportHandlerFunc(storagePath string, checker projectChecker, maxPackageSize uint64) func(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	ih := getImportHandlerInstance(
		storagePath, checker, maxPackageSize, importer.NewPackage,
		importer.NewImportPackageProcessor(
			new(model.YAMLManifestUnMarshaler), nil,
		), // TODO replace with correct task executor
	)
	return ih.HandleImport
}

func getImportHandlerInstance(
	storagePath string, checker projectChecker, maxPackageSize uint64,
	parserFunction ParseArchiveFunction, processor importPackageProcessor,
) *ImportHandler {
	return &ImportHandler{
		checker:                    checker,
		tempStorageDir:             storagePath,
		maxUncompressedPackageSize: maxPackageSize,
		parseArchive:               parserFunction,
		processor:                  processor,
	}
}

// HandleImport is the method invoked when a POST request is received on the import endpoint.
// This method will check that the project passed as parameter already exists in Keptn (
// return a 404 immediately if that is not the case),
// save the import package on the scratch storage and parse its contents.
func (ih *ImportHandler) HandleImport(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {

	// Check if the project exists
	projectExists, err := ih.checker.ProjectExists(params.Project)

	if err != nil {
		message := fmt.Sprintf("error checking for project %s existence : %v", params.Project, err)
		mError := models.Error{
			Code:    http.StatusFailedDependency,
			Message: &message,
		}
		return import_operations.NewImportFailedDependency().WithPayload(&mError)
	}

	if !projectExists {
		message := fmt.Sprintf("project %s does not exist", params.Project)

		mError := models.Error{
			Code:    http.StatusNotFound,
			Message: &message,
		}
		return import_operations.NewImportNotFound().WithPayload(&mError)
	}

	file, err := ioutil.TempFile(ih.tempStorageDir, "importConfig*"+defaultImportArchiveExtension)
	if err != nil {
		logger.Errorf("Error saving import archive: %v", err)
		message := fmt.Sprintf("Error saving import archive: %s", err)
		mError := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&mError)
	}

	defer func() {
		if errDefer := file.Close(); errDefer != nil {
			logger.Warnf("Error closing temporary import archive %s: %v", file.Name(), errDefer)
		}

		if errDefer := os.Remove(file.Name()); errDefer != nil {
			logger.Errorf("Error deleting temporary import archive %s: %v", file.Name(), errDefer)
		}
	}()

	_, err = io.Copy(file, params.ConfigPackage)

	if err != nil {
		logger.Errorf("Error saving import archive: %v", err)
		message := fmt.Sprintf("Error saving import archive: %s", err)
		mError := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&mError)
	}

	m, err := ih.parseArchive(file.Name(), ih.maxUncompressedPackageSize)

	if err != nil {
		logger.Errorf("Error opening import archive: %v", err)
		message := fmt.Sprintf("Error opening import archive: %s", err)
		mError := models.Error{
			Code:    http.StatusUnsupportedMediaType,
			Message: &message,
		}
		return import_operations.NewImportUnsupportedMediaType().WithPayload(&mError)
	}

	defer func() {
		manifestCloseErr := m.Close()
		if manifestCloseErr != nil {
			logger.Warnf("Error closing manifest %+v: %s", m, manifestCloseErr)
		}
	}()

	err = ih.processor.Process(params.Project, m)
	if err != nil {
		logger.Errorf("Error processing import archive: %v", err)
		message := fmt.Sprintf("Error processing import archive: %s", err)
		mError := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&mError)
	}

	return import_operations.NewImportOK()
}
