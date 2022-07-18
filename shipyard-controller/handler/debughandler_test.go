package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/stretchr/testify/assert"
)

func TestDebughandlerGetAllProjects(t *testing.T) {
	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET projects ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
						projects := []*apimodels.ExpandedProject{
							{
								CreationDate:     "string",
								LastEventContext: &apimodels.EventContextInfo{},
								ProjectName:      "project1",
								Shipyard:         "shipyard",
								ShipyardVersion:  "shipyard version",
								Stages:           []*apimodels.ExpandedStage{},
								GitCredentials:   &apimodels.GitAuthCredentialsSecure{},
							},
							{
								CreationDate:     "string",
								LastEventContext: &apimodels.EventContextInfo{},
								ProjectName:      "project2",
								Shipyard:         "shipyard",
								ShipyardVersion:  "shipyard version",
								Stages:           []*apimodels.ExpandedStage{},
								GitCredentials:   &apimodels.GitAuthCredentialsSecure{},
							},
						}

						return projects, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET projects empty",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
						var projects []*apimodels.ExpandedProject
						return projects, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		}, {
			name: "GET projects internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllProjectsFunc: func() ([]*apimodels.ExpandedProject, error) {
						return nil, fmt.Errorf("error")
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
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: ""},
				gin.Param{Key: "shkeptncontext", Value: ""},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			handler.GetAllProjects(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}

func TestDebughandlerGetAllSequencesForProject(t *testing.T) {
	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET sequences for project ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) (*apimodels.SequenceStates, error) {
						sequences := &apimodels.SequenceStates{
							States:      []apimodels.SequenceState{},
							NextPageKey: 0,
							PageSize:    0,
							TotalCount:  0,
						}
						return sequences, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET sequences for project empty",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) (*apimodels.SequenceStates, error) {
						var sequences *apimodels.SequenceStates
						return sequences, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET sequences for project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) (*apimodels.SequenceStates, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET sequences for project internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) (*apimodels.SequenceStates, error) {
						return nil, fmt.Errorf("error")
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
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: ""},
				gin.Param{Key: "shkeptncontext", Value: ""},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			handler.GetAllSequencesForProject(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}

}

func TestDebughandlerGetSequenceByID(t *testing.T) {
	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET sequenceByID ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
						event := &apimodels.SequenceState{
							Name:           "sequence1",
							Service:        "service1",
							Project:        "project1",
							Time:           "string",
							Shkeptncontext: "context",
							State:          "state",
							Stages:         []apimodels.SequenceStateStage{},
							ProblemTitle:   "problemtitle",
						}
						return event, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET sequenceByID project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		}, {
			name: "GET sequencetByID sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, ErrSequenceNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET sequenceByID internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName string, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, fmt.Errorf("error")
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
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: ""},
				gin.Param{Key: "shkeptncontext", Value: ""},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			handler.GetSequenceByID(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}

}

func TestDebughandlerGetEventByID(t *testing.T) {
	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET eventByID ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						eventSource := ""
						eventType := ""
						event := &apimodels.KeptnContextExtendedCE{
							Contenttype:        "contenttype",
							Data:               "data",
							Extensions:         "extensions",
							ID:                 "id",
							Shkeptncontext:     "shkeptncontext",
							Shkeptnspecversion: "Shkeptnspecversion",
							Source:             &eventSource,
							Specversion:        "specversion",
							Time:               time.Time{},
							Triggeredid:        "triggeredid",
							GitCommitID:        "gitcommitid",
							Type:               &eventType,
						}
						return event, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET eventByID project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		}, {
			name: "GET eventByID sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, ErrSequenceNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET eventByID event not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, ErrNoMatchingEvent
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET eventByID internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName string, shkeptncontext string, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, fmt.Errorf("error")
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
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: ""},
				gin.Param{Key: "shkeptncontext", Value: ""},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			handler.GetEventByID(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}

}

func TestDebughandlerGetAllEvents(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET events ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						eventSource := ""
						eventType := ""
						events := []*apimodels.KeptnContextExtendedCE{
							{
								Contenttype:        "contenttype",
								Data:               "data",
								Extensions:         "extensions",
								ID:                 "id",
								Shkeptncontext:     "shkeptncontext",
								Shkeptnspecversion: "Shkeptnspecversion",
								Source:             &eventSource,
								Specversion:        "specversion",
								Time:               time.Time{},
								Triggeredid:        "triggeredid",
								GitCommitID:        "gitcommitid",
								Type:               &eventType,
							},
						}
						return events, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET events empty",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, nil
					},
				},
			},
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "GET events project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		}, {
			name: "GET events sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, ErrSequenceNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
		{
			name: "GET events internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, fmt.Errorf("error")
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
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: ""},
				gin.Param{Key: "shkeptncontext", Value: ""},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			handler.GetAllEvents(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}

}
