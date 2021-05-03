package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStateHandler_GetState(t *testing.T) {
	type fields struct {
		StateRepo *db_mock.StateRepoMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantStatus int
	}{
		{
			name: "state repo returns states",
			fields: fields{
				StateRepo: &db_mock.StateRepoMock{
					FindStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return &models.SequenceStates{
							States: []models.SequenceState{
								{
									Name:           "delivery",
									Service:        "my-service",
									Project:        "my-project",
									Time:           time.Now().String(),
									Shkeptncontext: "my-context",
									State:          "triggered",
									Stages:         nil,
								},
							},
							NextPageKey: 0,
							PageSize:    1,
							TotalCount:  1,
						}, nil
					},
				},
			},
			request:    httptest.NewRequest("POST", "/state/my-project", nil),
			wantStatus: http.StatusOK,
		},
		{
			name: "state repo returns error",
			fields: fields{
				StateRepo: &db_mock.StateRepoMock{
					FindStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("POST", "/state/my-project", nil),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sh := NewStateHandler(tt.fields.StateRepo)

			handler := func(w http.ResponseWriter, r *http.Request) {
				c, _ := gin.CreateTestContext(w)
				c.Request = r
				sh.GetState(c)
			}

			w := httptest.NewRecorder()
			handler(w, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			require.Equal(t, 1, len(tt.fields.StateRepo.FindStatesCalls()))
			require.Equal(t, "my-project", tt.fields.StateRepo.FindStatesCalls()[0].Filter.Project)
		})
	}
}
