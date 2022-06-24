package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/api/importer"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"

	logger "github.com/sirupsen/logrus"
)

const defaultImportArchiveExtension = ".zip"

//go:generate moq -pkg fake --skip-ensure -out ./fake/projectchecker_mock.go . projectChecker:ProjectCheckerMock

type projectChecker interface {
	ProjectExists(projectName string) (bool, error)
}

// ParseArchiveFunction is the function called to parse the uploaded file
type ParseArchiveFunction func(string, uint64) (*importer.ZippedPackage, error)

// ImportHandler is the rest handler for the /import endpoint
type ImportHandler struct {
	checker                    projectChecker
	tempStorageDir             string
	maxUncompressedPackageSize uint64
	parseArchive               ParseArchiveFunction
}

// GetImportHandlerFunc will instantiate a configured ImportHandler and return the method that can be used for
// handling http requests to the endpoint. See restapi.configureAPI for usage
func GetImportHandlerFunc(storagePath string, checker projectChecker, maxPackageSize uint64) func(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	ih := getImportHandlerInstance(storagePath, checker, maxPackageSize, importer.NewPackage)
	return ih.HandleImport
}

func getImportHandlerInstance(
	storagePath string, checker projectChecker, maxPackageSize uint64,
	parserFunction ParseArchiveFunction,
) *ImportHandler {
	return &ImportHandler{
		checker:                    checker,
		tempStorageDir:             storagePath,
		maxUncompressedPackageSize: maxPackageSize,
		parseArchive:               parserFunction,
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

	return import_operations.NewImportOK()
}
