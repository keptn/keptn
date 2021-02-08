package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
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
						return nil, errProjectNotFound
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
						return nil, errStageNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET stage error from deatabase",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetStageFunc: func(projectName string, stageName string) (*models.ExpandedStage, error) {
						return nil, errors.New("whoops...")
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
				gin.Param{Key: "projectName", Value: "my-project"},
				gin.Param{Key: "stageName", Value: "my-stage"},
			}
			handler := NewStageHandler(tt.fields.StageManager)
			handler.GetStage(c)

			assert.Equal(t, tt.expectHttpStatus, w.Code)
		})
	}
}

func TestGetStages(t *testing.T) {

	type fields struct {
		StageManager IStageManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET stages project not found",
			fields: fields{
				StageManager: &fake.IStageManagerMock{
					GetAllStagesFunc: func(projectName string) ([]*models.ExpandedStage, error) {
						return nil, errProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/?nextPageKey=2", bytes.NewBuffer([]byte{}))
			c.Params = gin.Params{
				gin.Param{Key: "projectName", Value: "my-project"},
			}

			handler := NewStageHandler(tt.fields.StageManager)
			handler.GetAllStages(c)

			assert.Equal(t, tt.expectHttpStatus, w.Code)
		})
	}
}
