package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStage(t *testing.T) {

	type fields struct {
		StageManager IStageManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET stage project not found",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetStageFunc: func(projectName string, stageName string) (*models.ExpandedStage, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET stage stage not found",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetStageFunc: func(projectName string, stageName string) (*models.ExpandedStage, error) {
						return nil, ErrStageNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET stage error from database",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetStageFunc: func(projectName string, stageName string) (*models.ExpandedStage, error) {
						return nil, errors.New("whoops")
					},
				},
			},
			expectHttpStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "my-project"},
				gin.Param{Key: "stageName", Value: "my-stage"},
			}
			handler := NewStageHandler(tt.fields.StageManager)
			handler.GetStage(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}

func TestGetStages(t *testing.T) {

	s1 := &models.ExpandedStage{
		StageName: "s1",
	}
	s2 := &models.ExpandedStage{
		StageName: "s2",
	}
	s3 := &models.ExpandedStage{
		StageName: "s3",
	}
	expandedStages := []*models.ExpandedStage{s1, s2, s3}

	type fields struct {
		StageManager IStageManager
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *models.Stages
		pageSizeQueryParam string
	}{
		{
			name: "GET stages project not found",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET stages GetStages fails",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return nil, errors.New("whoops")
					},
				},
			},
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "GET stages",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return expandedStages, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.Stages{
				NextPageKey: "0",
				PageSize:    0,
				Stages:      expandedStages,
				TotalCount:  3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
			c.Params = gin.Params{
				gin.Param{Key: "projectName", Value: "my-project"},
			}

			handler := NewStageHandler(tt.fields.StageManager)
			handler.GetAllStages(c)

			if tt.expectJSONResponse != nil {
				response := &models.Stages{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)
		})
	}
}

func TestGetStagesWithPagination(t *testing.T) {

	s1 := &models.ExpandedStage{
		StageName: "s1",
	}
	s2 := &models.ExpandedStage{
		StageName: "s2",
	}
	s3 := &models.ExpandedStage{
		StageName: "s3",
	}
	expandedStages := []*models.ExpandedStage{s1, s2, s3}

	type fields struct {
		StageManager IStageManager
	}

	tests := []struct {
		name               string
		fields             fields
		expectHttpStatus   int
		expectJSONResponse *models.Stages
		url                string
	}{
		{
			name: "GET stages With Pagination",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return expandedStages, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.Stages{
				NextPageKey: "1",
				PageSize:    0,
				Stages:      expandedStages[0:1],
				TotalCount:  3,
			},
			url: "/?pageSize=1",
		},
		{
			name: "GET stages With Pagination2",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return expandedStages, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.Stages{
				NextPageKey: "2",
				PageSize:    0,
				Stages:      expandedStages[0:2],
				TotalCount:  3,
			},
			url: "/?pageSize=2",
		},
		{
			name: "GET stages With Pagination2",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return expandedStages, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.Stages{
				NextPageKey: "0",
				PageSize:    0,
				Stages:      expandedStages[2:3],
				TotalCount:  3,
			},
			url: "/?pageSize=1&nextPageKey=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, tt.url, bytes.NewBuffer([]byte{}))
			c.Params = gin.Params{
				gin.Param{Key: "projectName", Value: "my-project"},
			}

			handler := NewStageHandler(tt.fields.StageManager)
			handler.GetAllStages(c)

			if tt.expectJSONResponse != nil {
				response := &models.Stages{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}
			assert.Equal(t, tt.expectHttpStatus, w.Code)
		})
	}
}

func createExpandedStages() []*models.ExpandedStage {
	s1 := &models.ExpandedStage{
		StageName: "s1",
	}
	s2 := &models.ExpandedStage{
		StageName: "s2",
	}
	s3 := &models.ExpandedStage{
		StageName: "s3",
	}
	return []*models.ExpandedStage{s1, s2, s3}

}
