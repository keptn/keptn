package _import

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/pkg/import/fake"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"
)

const testArchiveSize20MB uint64 = 20 * 1024 * 1204

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
	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return false, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB)
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

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB)
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

func TestErrorImportBrokenReader(t *testing.T) {

	contentReader := io.NopCloser(
		&testReader{
			data:         []byte("some bytes before the error"),
			allowedLoops: 0,
			throwError:   true,
		},
	)

	var actualCheckedProject string
	mockedprojectChecker := &fake.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB)

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
