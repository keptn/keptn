package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllProjects(t *testing.T) {

	s1 := &models.ExpandedStage{
		StageName: "s1",
	}
	expandedStages := []*models.ExpandedStage{s1}

	p1 := &models.ExpandedProject{
		Stages: expandedStages,
	}

	type fields struct {
		ProjectManager IProjectManager
		EventSender    keptn.EventSender
	}

	tests := []struct {
		name               string
		fields             fields
		jsonPayload        string
		expectHttpStatus   int
		expectJSONResponse *models.ExpandedProjects
	}{
		{
			name: "Get all projects DB access fails",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return nil, errors.New("Whoops...")
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "Get all projects",
			fields: fields{
				ProjectManager: &fake.IProjectManagerMock{
					GetFunc: func() ([]*models.ExpandedProject, error) {
						return []*models.ExpandedProject{p1}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectHttpStatus: http.StatusOK,
			expectJSONResponse: &models.ExpandedProjects{
				NextPageKey: "0",
				Projects:    []*models.ExpandedProject{p1},
				TotalCount:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := NewProjectHandler(tt.fields.ProjectManager, tt.fields.EventSender)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))

			handler.GetAllProjects(c)

			if tt.expectJSONResponse != nil {
				response := &models.ExpandedProjects{}
				responseBytes, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(responseBytes, response)
				assert.Equal(t, tt.expectJSONResponse, response)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}
