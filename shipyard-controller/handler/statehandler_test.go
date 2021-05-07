package handler_test

import (
	"errors"
	"github.com/gin-gonic/gin"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStateHandler_GetState(t *testing.T) {
	type fields struct {
		StateRepo *db_mock.SequenceStateRepoMock
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
				StateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
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
			request:    httptest.NewRequest("GET", "/state/my-project", nil),
			wantStatus: http.StatusOK,
		},
		{
			name: "state repo returns error",
			fields: fields{
				StateRepo: &db_mock.SequenceStateRepoMock{
					FindSequenceStatesFunc: func(filter models.StateFilter) (*models.SequenceStates, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest("GET", "/state/my-project", nil),
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sh := handler.NewStateHandler(tt.fields.StateRepo)

			router := gin.Default()
			router.GET("/state/:project", func(c *gin.Context) {
				sh.GetSequenceState(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			require.Equal(t, 1, len(tt.fields.StateRepo.FindSequenceStatesCalls()))
			require.Equal(t, "my-project", tt.fields.StateRepo.FindSequenceStatesCalls()[0].Filter.Project)
		})
	}
}

func performRequest(r http.Handler, request *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	return w
}
