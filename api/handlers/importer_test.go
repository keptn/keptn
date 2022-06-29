package handlers

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/importer"
	"github.com/keptn/keptn/api/importer/fake"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"
	"github.com/keptn/keptn/api/test/utils"
)

const testArchiveSize20MB uint64 = 20 * 1024 * 1204

func TestErrorNonExistingProject(t *testing.T) {
	var actualCheckedProject string
	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return false, nil
		},
	}

	sut := GetImportHandlerFunc("", mockedprojectChecker, testArchiveSize20MB)
	projectName := "this_project_doesn't_exist"
	actualResponder := sut(
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
	assert.Equal(t, int64(http.StatusNotFound), actualPayload.Code)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestErrorUnableToCheckProject(t *testing.T) {
	var actualCheckedProject string
	prjCheckerErrDesc := "some obscure project checker error"
	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, errors.New(prjCheckerErrDesc)
		},
	}

	sut := GetImportHandlerFunc("", mockedprojectChecker, testArchiveSize20MB)
	projectName := "this_project_existence_cannot_be_checked"
	actualResponder := sut(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: nil,
			Project:       projectName,
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportFailedDependency{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportFailedDependency).Payload
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, int64(http.StatusFailedDependency), actualPayload.Code)
	assert.Contains(t, *actualPayload.Message, prjCheckerErrDesc)
	assert.Equal(t, projectName, actualCheckedProject)
}

func TestErrorImportBrokenReader(t *testing.T) {

	contentReader := io.NopCloser(
		utils.NewTestReader([]byte("some bytes before the error"), 0, true),
	)

	var actualCheckedProject string
	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, nil
		},
	}

	sut := GetImportHandlerFunc("", mockedprojectChecker, testArchiveSize20MB)

	projectName := "foobarbaz"
	actualResponder := sut(
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

func TestErrorSaveArchiveFromUpload(t *testing.T) {

	// The following code should create a folder that we cannot wtite but in Ci we run tests in a docker container as
	// root so we cannot use this
	// tempDir, err := ioutil.TempDir("", "test-handler-save-")
	// t.Logf("using %s as temp folder", tempDir)
	// require.NoError(t, err)
	// defer os.RemoveAll(tempDir)
	// // set the directory permissions as r-x
	// err := os.Chmod(tempDir, 0500)
	// require.NoError(t, err)
	//
	// // restore permission before cleanup
	// defer os.Chmod(tempDir, 0700)

	tempDir := "this-directory-must-not-exist"

	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	importHandlerFunc := GetImportHandlerFunc(tempDir, mockedprojectChecker, testArchiveSize20MB)
	actualResponder := importHandlerFunc(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: io.NopCloser(bytes.NewReader([]byte("some payload bytes here"))),
			Project:       "projectName",
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportBadRequest{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportBadRequest).Payload
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, int64(http.StatusBadRequest), actualPayload.Code)
}

func TestErrorCreateNewZippedArchiveFromUpload(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-handler-save-")
	t.Logf("using %s as temp folder", tempDir)
	require.NoError(t, err)

	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	packageContent := []byte("some payload bytes here")

	errorParsingPackage := func(file string, maxSize uint64) (*importer.ZippedPackage, error) {
		require.FileExists(t, file)

		return nil, errors.New("error parsing package")
	}

	importHandlerFunc := getImportHandlerInstance(
		tempDir, mockedprojectChecker, testArchiveSize20MB, errorParsingPackage,
	).HandleImport

	actualResponder := importHandlerFunc(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: io.NopCloser(bytes.NewReader(packageContent)),
			Project:       "projectName",
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportUnsupportedMediaType{}, actualResponder)
	actualPayload := actualResponder.(*import_operations.ImportUnsupportedMediaType).Payload
	assert.NotEmpty(t, actualPayload.Message)
	assert.Equal(t, int64(http.StatusUnsupportedMediaType), actualPayload.Code)
}

func TestImportHandlerHappyPath(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-handler-save-")
	t.Logf("using %s as temp folder", tempDir)
	require.NoError(t, err)

	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	packageContent := []byte("some payload bytes here")

	parsingPackage := func(file string, maxSize uint64) (*importer.ZippedPackage, error) {
		require.FileExists(t, file)
		packageBytes, err := ioutil.ReadFile(file)
		require.NoError(t, err)
		assert.Equal(t, packageContent, packageBytes)
		return &importer.ZippedPackage{}, nil
	}

	importHandlerFunc := getImportHandlerInstance(
		tempDir, mockedprojectChecker, testArchiveSize20MB, parsingPackage,
	).HandleImport

	actualResponder := importHandlerFunc(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: io.NopCloser(bytes.NewReader(packageContent)),
			Project:       "projectName",
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportOK{}, actualResponder)
}
