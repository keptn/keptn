package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
)

func TestDebughandlerGetAllSequencesForProject(t *testing.T) {
	type fields struct {
		DebugManager *fake.IDebugManagerMock
	}

	tests := []struct {
		name         string
		fields       fields
		request      *http.Request
		wantStatus   int
		projectName  string
		wantResponse []models.SequenceExecution
	}{
		{
			name: "get all sequences ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) ([]models.SequenceExecution, error) {
						return []models.SequenceExecution{}, nil
					},
				},
			},
			request:      httptest.NewRequest("GET", "/sequences/project/projectname", nil),
			wantStatus:   http.StatusOK,
			wantResponse: []models.SequenceExecution{},
			projectName:  "projectname",
		},
		{
			name: "get all sequences project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) ([]models.SequenceExecution, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request:      httptest.NewRequest("GET", "/sequences/project/projectname", nil),
			wantStatus:   http.StatusNotFound,
			wantResponse: nil,
			projectName:  "projectname",
		},
		{
			name: "get all sequences internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllSequencesForProjectFunc: func(projectName string) ([]models.SequenceExecution, error) {
						return nil, common.ErrInternalError
					},
				},
			},
			request:      httptest.NewRequest("GET", "/sequences/project/projectname", nil),
			wantStatus:   http.StatusInternalServerError,
			wantResponse: nil,
			projectName:  "projectname",
		},
	}

	for _, tt := range tests {
		dh := handler.NewDebugHandler(tt.fields.DebugManager)

		router := gin.Default()
		router.GET("/sequences/project/:project", func(c *gin.Context) {
			dh.GetAllSequencesForProject(c)
		})

		w := performRequest(router, tt.request)

		require.Equal(t, tt.wantStatus, w.Code)

		if tt.wantStatus == http.StatusOK {
			var object []models.SequenceExecution
			err := json.Unmarshal(w.Body.Bytes(), &object)
			require.Nil(t, err)
			require.Equal(t, object, tt.wantResponse)
		}

		require.Equal(t, tt.projectName, tt.fields.DebugManager.GetAllSequencesForProjectCalls()[0].ProjectName)
	}
}

func TestDebughandlerGetSequenceByID(t *testing.T) {
	type fields struct {
		DebugManager *fake.IDebugManagerMock
	}

	tests := []struct {
		name           string
		fields         fields
		request        *http.Request
		wantStatus     int
		projectName    string
		shkeptncontext string
		wantResponse   *apimodels.SequenceState
	}{
		{
			name: "get sequence by id ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName, shkeptncontext string) (*apimodels.SequenceState, error) {
						return &apimodels.SequenceState{}, nil
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context", nil),
			wantStatus:     http.StatusOK,
			wantResponse:   &apimodels.SequenceState{},
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get sequence by id project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get sequence by id sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, common.ErrSequenceNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get sequence by id internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetSequenceByIDFunc: func(projectName, shkeptncontext string) (*apimodels.SequenceState, error) {
						return nil, common.ErrInternalError
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context", nil),
			wantStatus:     http.StatusInternalServerError,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
	}

	for _, tt := range tests {
		dh := handler.NewDebugHandler(tt.fields.DebugManager)

		router := gin.Default()
		router.GET("/sequences/project/:project/shkeptncontext/:shkeptncontext", func(c *gin.Context) {
			dh.GetSequenceByID(c)
		})

		w := performRequest(router, tt.request)

		require.Equal(t, tt.wantStatus, w.Code)

		if tt.wantStatus == http.StatusOK {
			var object apimodels.SequenceState
			err := json.Unmarshal(w.Body.Bytes(), &object)
			require.Nil(t, err)

			require.Equal(t, object, *tt.wantResponse)
		}

		require.Equal(t, tt.projectName, tt.fields.DebugManager.GetSequenceByIDCalls()[0].ProjectName)
		require.Equal(t, tt.shkeptncontext, tt.fields.DebugManager.GetSequenceByIDCalls()[0].Shkeptncontext)
	}
}

func TestDebughandlerGetEventByID(t *testing.T) {

	type fields struct {
		DebugManager *fake.IDebugManagerMock
	}

	tests := []struct {
		name           string
		fields         fields
		request        *http.Request
		wantStatus     int
		projectName    string
		shkeptncontext string
		eventId        string
		wantResponse   *apimodels.KeptnContextExtendedCE
	}{
		{
			name: "get event by id ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName, shkeptncontext, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return &apimodels.KeptnContextExtendedCE{}, nil
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event/eventid", nil),
			wantStatus:     http.StatusOK,
			wantResponse:   &apimodels.KeptnContextExtendedCE{},
			projectName:    "projectname",
			shkeptncontext: "context",
			eventId:        "eventid",
		},
		{
			name: "get event by id project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName, shkeptncontext, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event/eventid", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
			eventId:        "eventid",
		},
		{
			name: "get event by id sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName, shkeptncontext, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrSequenceNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event/eventid", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
			eventId:        "eventid",
		},
		{
			name: "get event by id internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetEventByIDFunc: func(projectName, shkeptncontext, eventId string) (*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrInternalError
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event/eventid", nil),
			wantStatus:     http.StatusInternalServerError,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
			eventId:        "eventid",
		},
	}

	for _, tt := range tests {
		dh := handler.NewDebugHandler(tt.fields.DebugManager)

		router := gin.Default()
		router.GET("/sequences/project/:project/shkeptncontext/:shkeptncontext/event/:eventId", func(c *gin.Context) {
			dh.GetEventByID(c)
		})

		w := performRequest(router, tt.request)

		require.Equal(t, tt.wantStatus, w.Code)

		if tt.wantStatus == http.StatusOK {
			var object apimodels.KeptnContextExtendedCE
			err := json.Unmarshal(w.Body.Bytes(), &object)
			require.Nil(t, err)
			require.Equal(t, &object, tt.wantResponse)
		}

		require.Equal(t, tt.projectName, tt.fields.DebugManager.GetEventByIDCalls()[0].ProjectName)
		require.Equal(t, tt.shkeptncontext, tt.fields.DebugManager.GetEventByIDCalls()[0].Shkeptncontext)
		require.Equal(t, tt.eventId, tt.fields.DebugManager.GetEventByIDCalls()[0].EventId)
	}
}

func TestDebughandlerGetAllEvents(t *testing.T) {

	type fields struct {
		DebugManager *fake.IDebugManagerMock
	}

	var expected = &apimodels.Events{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Events:      []*apimodels.KeptnContextExtendedCE{},
	}

	tests := []struct {
		name           string
		fields         fields
		request        *http.Request
		wantStatus     int
		projectName    string
		shkeptncontext string
		wantResponse   *apimodels.Events
	}{
		{
			name: "get all events ok",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return []*apimodels.KeptnContextExtendedCE{}, nil
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event", nil),
			wantStatus:     http.StatusOK,
			wantResponse:   expected,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get all events project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrProjectNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get all events sequence not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrSequenceNotFound
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event", nil),
			wantStatus:     http.StatusNotFound,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
		{
			name: "get all events internal error",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName, shkeptncontext string) ([]*apimodels.KeptnContextExtendedCE, error) {
						return nil, common.ErrInternalError
					},
				},
			},
			request:        httptest.NewRequest("GET", "/sequences/project/projectname/shkeptncontext/context/event", nil),
			wantStatus:     http.StatusInternalServerError,
			wantResponse:   nil,
			projectName:    "projectname",
			shkeptncontext: "context",
		},
	}

	for _, tt := range tests {
		dh := handler.NewDebugHandler(tt.fields.DebugManager)

		router := gin.Default()
		router.GET("/sequences/project/:project/shkeptncontext/:shkeptncontext/event", func(c *gin.Context) {
			dh.GetAllEvents(c)
		})

		w := performRequest(router, tt.request)

		require.Equal(t, tt.wantStatus, w.Code)

		if tt.wantStatus == http.StatusOK {
			var object apimodels.Events
			err := json.Unmarshal(w.Body.Bytes(), &object)
			require.Nil(t, err)
			require.Equal(t, object, *tt.wantResponse)
		}

		require.Equal(t, tt.projectName, tt.fields.DebugManager.GetAllEventsCalls()[0].ProjectName)
		require.Equal(t, tt.shkeptncontext, tt.fields.DebugManager.GetAllEventsCalls()[0].Shkeptncontext)
	}
}
