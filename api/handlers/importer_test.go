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

	handlers_mock "github.com/keptn/keptn/api/handlers/fake"
	"github.com/keptn/keptn/api/importer"
	"github.com/keptn/keptn/api/importer/model"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/import_operations"
	"github.com/keptn/keptn/api/test/utils"
)

const testArchiveSize20MB uint64 = 20 * 1024 * 1204

func TestErrorNonExistingProject(t *testing.T) {
	var actualCheckedProject string
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return false, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB, nil, nil)
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
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, errors.New(prjCheckerErrDesc)
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB, nil, nil)

	projectName := "this_project_existence_cannot_be_checked"
	actualResponder := sut.HandleImport(
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
	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			actualCheckedProject = projectName
			return true, nil
		},
	}

	sut := getImportHandlerInstance("", mockedprojectChecker, testArchiveSize20MB, nil, nil)

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

	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	importHandler := getImportHandlerInstance(tempDir, mockedprojectChecker, testArchiveSize20MB, nil, nil)
	actualResponder := importHandler.HandleImport(
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

	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	packageContent := []byte("some payload bytes here")

	errorParsingPackage := func(file string, maxSize uint64) (*ZippedPackage, error) {
		require.FileExists(t, file)

		return nil, errors.New("error parsing package")
	}

	mockedProcessor := &handlers_mock.MockImportPackageProcessor{
		ProcessFunc: func(
			project string, ip importer.ImportPackage,
		) (*model.ManifestExecution, error) {
			return nil, nil
		},
	}

	importHandlerFunc := getImportHandlerInstance(
		tempDir, mockedprojectChecker, testArchiveSize20MB, errorParsingPackage, mockedProcessor,
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
	assert.Len(t, mockedProcessor.ProcessCalls(), 0)
}

func TestErrorProcessingImportPackage(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-handler-save-")
	t.Logf("using %s as temp folder", tempDir)
	require.NoError(t, err)

	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	packageContent := []byte("some payload bytes here")

	parsingPackage := func(file string, maxSize uint64) (*ZippedPackage, error) {
		require.FileExists(t, file)
		packageBytes, err := ioutil.ReadFile(file)
		require.NoError(t, err)
		assert.Equal(t, packageContent, packageBytes)
		return &ZippedPackage{}, nil
	}

	mockedProcessor := &handlers_mock.MockImportPackageProcessor{
		ProcessFunc: func(
			project string, ip importer.ImportPackage,
		) (*model.ManifestExecution, error) {
			return nil, errors.New("error processing archive")
		},
	}

	importHandlerFunc := getImportHandlerInstance(
		tempDir, mockedprojectChecker, testArchiveSize20MB, parsingPackage, mockedProcessor,
	).HandleImport

	actualResponder := importHandlerFunc(
		import_operations.ImportParams{
			HTTPRequest:   nil,
			ConfigPackage: io.NopCloser(bytes.NewReader(packageContent)),
			Project:       "projectName",
		},
		new(models.Principal),
	)

	require.IsType(t, &import_operations.ImportBadRequest{}, actualResponder)
	assert.Len(t, mockedProcessor.ProcessCalls(), 1)
}

func TestImportHandlerHappyPath(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-handler-save-")
	t.Logf("using %s as temp folder", tempDir)
	require.NoError(t, err)

	mockedprojectChecker := &handlers_mock.ProjectCheckerMock{
		ProjectExistsFunc: func(projectName string) (bool, error) {
			return true, nil
		},
	}

	packageContent := []byte("some payload bytes here")

	parsingPackage := func(file string, maxSize uint64) (*ZippedPackage, error) {
		require.FileExists(t, file)
		packageBytes, err := ioutil.ReadFile(file)
		require.NoError(t, err)
		assert.Equal(t, packageContent, packageBytes)
		return &ZippedPackage{}, nil
	}

	mockedProcessor := &handlers_mock.MockImportPackageProcessor{
		ProcessFunc: func(
			project string, ip importer.ImportPackage,
		) (*model.ManifestExecution, error) {
			return nil, nil
		},
	}

	importHandlerFunc := getImportHandlerInstance(
		tempDir, mockedprojectChecker, testArchiveSize20MB, parsingPackage, mockedProcessor,
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
	assert.Len(t, mockedProcessor.ProcessCalls(), 1)
}

func Test_mapManifestExecution(t *testing.T) {

	taskaExecution := model.TaskExecution{
		TaskContext: model.TaskContext{
			Project: "test-prj",
			Task: &model.ManifestTask{
				APITask: &model.APITask{
					Action:      "",
					PayloadFile: "",
				},
				ID:      "taska",
				Type:    "api",
				Name:    "Task A",
				Context: map[string]string{"ctxKey1": "ctxValue1"},
			},
			Context: map[string]string{"ctxKey1": "ctxValue1"},
		},
		Response: map[string]any{"respkey1": "value1", "respkey2": "value2"},
	}

	taskbExecution := model.TaskExecution{
		TaskContext: model.TaskContext{
			Project: "test-prj",
			Task: &model.ManifestTask{
				APITask: &model.APITask{
					Action:      "",
					PayloadFile: "",
				},
				ID:      "taskb",
				Type:    "api",
				Name:    "Task B",
				Context: map[string]string{"ctxBKey1": "[[.Tasks.taska.Context.ctxKey1]]"},
			},
			Context: map[string]string{"ctxBKey1": "ctxValue1"},
		},
		Response: map[string]any{"foo": "bar"},
	}

	type args struct {
		exec *model.ManifestExecution
	}
	tests := []struct {
		name string
		args args
		want *models.ImportSummary
	}{
		{
			name: "Nil manifest execution - empty summary",
			args: args{exec: nil},
			want: &models.ImportSummary{
				Message: "",
				Outcome: models.ImportSummaryOutcomeSuccess,
				Tasks:   nil,
			},
		},
		{
			name: "Empty manifest execution - empty summary",
			args: args{
				exec: &model.ManifestExecution{
					Inputs:       map[string]string{},
					Tasks:        map[string]model.TaskExecution{},
					TaskSequence: nil,
				},
			},
			want: &models.ImportSummary{
				Message: "",
				Outcome: models.ImportSummaryOutcomeSuccess,
				Tasks:   []*models.Task{},
			},
		},
		{
			name: "Pair of tasks in manifest execution - simple summary",
			args: args{
				exec: &model.ManifestExecution{
					Inputs: map[string]string{},
					Tasks: map[string]model.TaskExecution{
						"taska": taskaExecution,
						"taskb": taskbExecution,
					},
					TaskSequence: []string{"taska", "taskb"},
				},
			},
			want: &models.ImportSummary{
				Message: "",
				Outcome: models.ImportSummaryOutcomeSuccess,
				Tasks: []*models.Task{
					{
						Response: taskaExecution.Response,
						Task:     taskaExecution.TaskContext,
					},
					{
						Response: taskbExecution.Response,
						Task:     taskbExecution.TaskContext,
					},
				},
			},
		},
		{
			name: "Unknown task in manifest execution - nil task in summary",
			args: args{
				exec: &model.ManifestExecution{
					Inputs: map[string]string{},
					Tasks: map[string]model.TaskExecution{
						"taska": taskaExecution,
						"taskb": taskbExecution,
					},
					TaskSequence: []string{"taska", "taskc", "taskb"},
				},
			},
			want: &models.ImportSummary{
				Message: "",
				Outcome: models.ImportSummaryOutcomeSuccess,
				Tasks: []*models.Task{
					{
						Response: taskaExecution.Response,
						Task:     taskaExecution.TaskContext,
					},
					nil, // this is where taskc should have been
					{
						Response: taskbExecution.Response,
						Task:     taskbExecution.TaskContext,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, mapManifestExecution(tt.args.exec), "mapManifestExecution(%v)", tt.args.exec)
			},
		)
	}
}
