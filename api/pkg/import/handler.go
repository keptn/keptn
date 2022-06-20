package _import

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

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
	checker        projectChecker
	tempStorageDir string
}

func GetImportHandlerFunc(storagePath string, checker projectChecker) func(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	ih := getImportHandlerInstance(storagePath, checker)
	return ih.HandleImport
}

func getImportHandlerInstance(storagePath string, checker projectChecker) *ImportHandler {
	return &ImportHandler{
		checker:        checker,
		tempStorageDir: storagePath,
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

	file, err := ioutil.TempFile(ih.tempStorageDir, "importConfig*"+defaultImportArchiveExtension)
	if err != nil {
		return import_operations.NewImportBadRequest()
	}

	defer func() {
		errDefer := file.Close()
		if errDefer != nil {
			logger.Warnf("Error closing temporary import archive %s: %v", file.Name(), errDefer)
		}

		errDefer = os.Remove(file.Name())
		if errDefer != nil {
			logger.Errorf("Error deleting temporary import archive %s: %v", file.Name(), errDefer)
		}
	}()

	tempFileSize, err := io.Copy(file, params.ConfigPackage)
	if err != nil {
		logger.Errorf("Error saving import archive: %v", err)
		message := fmt.Sprintf("Error reading import archive: %s", err)
		error := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&error)
	}

	zipReader, err := zip.NewReader(file, tempFileSize)
	if err != nil {
		logger.Errorf("Error opening import archive %s: %v", file.Name(), err)
		message := fmt.Sprintf("Error opening import archive: %s", err)
		error := models.Error{
			Code:    http.StatusUnsupportedMediaType,
			Message: &message,
		}
		return import_operations.NewImportUnsupportedMediaType().WithPayload(&error)
	}

	extractionDir := strings.TrimSuffix(file.Name(), defaultImportArchiveExtension)
	err = os.Mkdir(extractionDir, os.ModeDir|os.ModePerm)
	if err != nil {
		logger.Errorf("Error creating folder %s for zip extraction: %v", extractionDir, err)
		message := fmt.Sprintf("Error reading import archive: %s", err)
		error := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		return import_operations.NewImportBadRequest().WithPayload(&error)
	}

	defer func() {
		errCleanup := os.RemoveAll(extractionDir)
		if errCleanup != nil {
			logger.Warnf("Error cleaning up extraction directory %s: %v", extractionDir, errCleanup)
		}
	}()

	err = extractZipFile(zipReader, extractionDir)

	if err != nil {
		logger.Errorf("Error extracting zip %s: %v", file.Name(), err)
		message := fmt.Sprintf("Error extracting archive: %s", err)
		error := models.Error{
			Code:    http.StatusBadRequest,
			Message: &message,
		}
		// TODO check if we want to return something != 400
		return import_operations.NewImportBadRequest().WithPayload(&error)
	}

	// TODO validate and start mainfest import

	return import_operations.NewImportOK()
}

func extractZipFile(reader *zip.Reader, outputDir string) error {

	// TODO add constraint on maximum extracted archive size

	for _, zippedFile := range reader.File {
		logger.Debugf("Extracting file %+v to %s...", reader, outputDir)

		if zippedFile.FileInfo().IsDir() {
			fullOuputDirectoryName := path.Join(outputDir, zippedFile.Name)
			err := os.Mkdir(fullOuputDirectoryName, 0755)
			if err != nil {
				return fmt.Errorf("error creating directory %s: %w", fullOuputDirectoryName, err)
			}
		} else {
			err := func() error {
				src, err := zippedFile.Open()
				if err != nil {
					return fmt.Errorf("error reading %s from archive: %w", zippedFile.Name, err)
				}
				defer src.Close()

				dst, err := os.Create(path.Join(outputDir, zippedFile.Name))
				if err != nil {
					return fmt.Errorf(
						"error creating file %s in output directory %s: %w", zippedFile.Name, outputDir, err,
					)
				}
				defer dst.Close()

				// TODO add safety checks on the actual written file size
				_, err = io.Copy(dst, src)

				if err != nil {
					return fmt.Errorf("error extracting %s from archive into %s: %w", zippedFile.Name, dst.Name(), err)
				}

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}

	return nil
}
