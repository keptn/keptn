package handlers

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	handlers_mock "github.com/keptn/keptn/api/handlers/fake"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"
)

type testReader struct {
	data         []byte
	allowedLoops int
	loopCount    int
	pos          int
	throwError   bool
}

func (tr *testReader) Read(p []byte) (n int, err error) {
	// check if we can reset
	if tr.pos >= len(tr.data) && tr.loopCount < tr.allowedLoops {
		tr.allowedLoops++
		tr.pos = 0
	}

	if tr.pos >= len(tr.data) {
		err := io.EOF

		if tr.throwError {
			err = errors.New("testReader error")
		}

		// end of buffer
		return 0, err
	}

	// try to copy as much as possible in the buffer
	copied := copy(p, tr.data[tr.pos:])
	tr.pos += copied
	return copied, nil
}

func TestErrorNonExistingProject(t *testing.T) {
	var actualCheckedProject string
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return false, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker)
	projectName := "this_project_doesn't_exist"
	actualResponder := sut.HandleImport(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: nil,
			Project:       projectName,
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportNotFound{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportNotFound).Payload
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestErrorUnableToCheckProject(t *testing.T) {
	var actualCheckedProject string
	prjCheckerErrDesc := "some obscure project checker error"
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, errors.New(prjCheckerErrDesc)
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker)
	projectName := "this_project_existence_cannot_be_checked"
	actualResponder := sut.HandleImport(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: nil,
			Project:       projectName,
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportNotFound{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportNotFound).Payload
	assert.NotEmpty(t, actualPayload.Message)
	assert.Contains(t, *actualPayload.Message, prjCheckerErrDesc)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestErrorImportNoValidZip(t *testing.T) {
	contentReader := io.NopCloser(bytes.NewReader([]byte("this is clearly not a zip file")))

	var actualCheckedProject string
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker)
	projectName := "foobar"
	actualResponder := sut.HandleImport(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: contentReader,
			Project:       projectName,
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportUnsupportedMediaType{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportUnsupportedMediaType).Payload
	require.NotNil(t, actualPayload)
	assert.Equal(t, int64(http.StatusUnsupportedMediaType), actualPayload.Code)
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestErrorImportBrokenReader(t *testing.T) {

	contentReader := io.NopCloser(
		&testReader{
			data:         []byte("some bytes before the error"),
			allowedLoops: 0,
			throwError:   true,
		},
	)

	var actualCheckedProject string
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker)

	projectName := "foobarbaz"
	actualResponder := sut.HandleImport(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: contentReader,
			Project:       projectName,
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportBadRequest{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportBadRequest).Payload
	require.NotNil(t, actualPayload)
	assert.Equal(t, int64(http.StatusBadRequest), actualPayload.Code)
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestExtractZipFileHappyPath(t *testing.T) {

	tempDir, err := ioutil.TempDir("", "test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempZipFile, err := ioutil.TempFile(tempDir, "test-archive")
	require.NoError(t, err)
	writeZip(tempZipFile, "../test/data/import/samplemanifest")

	err = tempZipFile.Close()
	require.NoError(t, err)

	tempZipFile, err = os.Open(tempZipFile.Name())
	require.NoError(t, err)

	zipFileStat, err := os.Stat(tempZipFile.Name())
	require.NoError(t, err)

	zipFileReader, err := zip.NewReader(tempZipFile, zipFileStat.Size())
	require.NoError(t, err)

	extractedPath := path.Join(tempDir, "extracted")
	err = os.Mkdir(extractedPath, 0700)
	require.NoError(t, err)

	err = extractZipFile(zipFileReader, extractedPath)
	assert.NoError(t, err)

	assertDirEqual(t, "../test/data/import/samplemanifest", extractedPath)
}

func assertDirEqual(t *testing.T, expected string, actual string) {
	filepath.WalkDir(
		expected, func(walkedPath string, d os.DirEntry, err error) error {

			relPath, err := filepath.Rel(expected, walkedPath)
			require.NoError(t, err)

			actualPath := path.Join(actual, relPath)
			if d.IsDir() {
				assert.DirExists(t, actualPath)
			} else {
				assert.FileExists(t, actualPath)
				actualFileBytes, err := ioutil.ReadFile(actualPath)
				assert.NoError(t, err)
				actualFileHash := sha256.Sum256(actualFileBytes)

				expectedFileBytes, err := ioutil.ReadFile(walkedPath)
				assert.NoError(t, err)
				expectedFileHash := sha256.Sum256(expectedFileBytes)

				assert.Equalf(
					t, hex.EncodeToString(expectedFileHash[:]), hex.EncodeToString(actualFileHash[:]),
					"files %s and %s should have the same content but their SHA256 is different!", walkedPath,
					actualPath,
				)
			}

			return nil
		},
	)
}

func writeZip(file *os.File, filePaths ...string) error {
	zipWriter := zip.NewWriter(file)
	for _, filePath := range filePaths {
		fileInfo, err := os.Lstat(filePath)
		if err != nil {
			return fmt.Errorf("error getting info on %s: %w", filePath, err)
		}
		if fileInfo.IsDir() {
			err = addDirectoryContentToZip(zipWriter, filePath)
		} else {
			err = addFileToZip(zipWriter, filePath, path.Base(filePath))
		}

		if err != nil {
			return fmt.Errorf("error adding %s to zip: %w", filePath, err)
		}
	}

	return zipWriter.Close()

}

func addDirectoryToZip(writer *zip.Writer, dirName string) error {
	if !strings.HasSuffix(dirName, "/") {
		dirName = dirName + "/"
	}
	_, err := writer.Create(dirName)
	if err != nil {
		return fmt.Errorf("error creating directory %s in zip file: %w", dirName, err)
	}
	return nil
}

func addFileToZip(writer *zip.Writer, srcPath string, dstPath string) error {
	fileWriter, err := writer.Create(dstPath)
	if err != nil {
		return fmt.Errorf("error creating file %s in zip file: %w", dstPath, err)
	}
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", srcPath, err)
	}

	defer src.Close()
	_, err = io.Copy(fileWriter, src)
	if err != nil {
		return fmt.Errorf(
			"error writing file content from %s to %s (in zip file): %w",
			srcPath, dstPath, err,
		)
	}
	return nil
}

func addDirectoryContentToZip(writer *zip.Writer, baseDir string) error {
	return filepath.WalkDir(
		baseDir, func(walkedPath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			relativePath, err := filepath.Rel(baseDir, walkedPath)

			if err != nil {
				return fmt.Errorf("error calculating path %s as relative path of %s", walkedPath, baseDir)
			}

			if relativePath == "." {
				return nil
			}

			if d.IsDir() {
				return addDirectoryToZip(writer, relativePath)
			} else {
				return addFileToZip(writer, walkedPath, relativePath)
			}
		},
	)
}
