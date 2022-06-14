package handlers

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"

	logger "github.com/sirupsen/logrus"
)

type ImportHandler struct {
	tempStorageDir string
}

func GetImportHandlerFunc() func(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	ih := getImportHandlerInstance()
	return ih.HandleImport
}

func getImportHandlerInstance() *ImportHandler {
	return new(ImportHandler)
}

func (ih *ImportHandler) HandleImport(
	params import_operations.ImportParams, principal *models.Principal,
) middleware.Responder {
	// TODO verify the Content-type of the request ?
	file, err := ioutil.TempFile(ih.tempStorageDir, "importConfig*.zip")
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
		return import_operations.NewImportBadRequest()
	}

	zipReader, err := zip.NewReader(file, tempFileSize)
	if err != nil {
		logger.Errorf("Error opening import archive %s: %v", file.Name(), err)
		message := "Error opening import archive: " + err.Error()
		error := models.Error{
			Code:    http.StatusUnsupportedMediaType,
			Message: &message,
		}
		return import_operations.NewImportUnsupportedMediaType().WithPayload(&error)
	}

	var totalUncompressedSize uint64
	for _, zippedFile := range zipReader.File {
		logger.Debugf("Inspecting file %+v from import archive...", zippedFile)
		totalUncompressedSize += zippedFile.UncompressedSize64
	}

	logger.Infof("Total uncompressed size of import archive %d bytes", totalUncompressedSize)

	return import_operations.NewImportOK()
}
