package _import

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"

	logger "github.com/sirupsen/logrus"
)

const defaultImportArchiveExtension = ".zip"

//go:generate moq -pkg fake --skip-ensure -out ./fake/projectchecker_mock.go . projectChecker:ProjectCheckerMock

type projectChecker interface {
	ProjectExists(projectName string) (bool, error)
}

type ImportHandler struct {
	checker                    projectChecker
	tempStorageDir             string
	maxUncompressedPackageSize uint64
}

func GetImportHandlerFunc(storagePath string, checker projectChecker, maxPackageSize uint64) func(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	ih := getImportHandlerInstance(storagePath, checker, maxPackageSize)
	return ih.HandleImport
}

func getImportHandlerInstance(storagePath string, checker projectChecker, maxPackageSize uint64) *ImportHandler {
	return &ImportHandler{
		checker:                    checker,
		tempStorageDir:             storagePath,
		maxUncompressedPackageSize: maxPackageSize,
	}
}

func (ih *ImportHandler) HandleImport(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {

	// Check if the project exists
	projectExists, err := ih.checker.ProjectExists(params.Project)

	if err != nil || !projectExists {
		message := fmt.Sprintf("project %s does not exist", params.Project)
		if err != nil {
			message = fmt.Sprintf("error checking for project %s existence : %v", params.Project, err)
		}

		error := models.Error{
			Code:    http.StatusNotFound,
			Message: &message,
		}
		return import_operations.NewImportNotFound().WithPayload(&error)
	}

	tempFileName, err := func() (string, error) {
		file, err := ioutil.TempFile(ih.tempStorageDir, "importConfig*"+defaultImportArchiveExtension)
		if err != nil {
			return "", err
		}

		defer func() {
			errDefer := file.Close()
			if errDefer != nil {
				logger.Warnf("Error closing temporary import archive %s: %v", file.Name(), errDefer)
			}
		}()

		_, err = io.Copy(file, params.ConfigPackage)
		return file.Name(), err
	}()

	if tempFileName != "" {
		defer func() {
			errDefer := os.Remove(tempFileName)
			if errDefer != nil {
				logger.Errorf("Error deleting temporary import archive %s: %v", tempFileName, errDefer)
			}
		}()
	}

	if err != nil {
		logger.Errorf("Error saving import archive: %v", err)
		message := fmt.Sprintf("Error saving import archive: %s", err)
		error := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&error)
	}

	m, err := NewPackage(tempFileName, ih.maxUncompressedPackageSize)

	if err != nil {
		logger.Errorf("Error opening import archive: %v", err)
		message := fmt.Sprintf("Error opening import archive: %s", err)
		error := models.Error{
			Code:    http.StatusUnsupportedMediaType,
			Message: &message,
		}
		return import_operations.NewImportUnsupportedMediaType().WithPayload(&error)
	}

	defer func() {
		manifestCloseErr := m.Close()
		if manifestCloseErr != nil {
			logger.Warnf("Error closing manifest %+v: %s", m, manifestCloseErr)
		}
	}()

	return import_operations.NewImportOK()
}
