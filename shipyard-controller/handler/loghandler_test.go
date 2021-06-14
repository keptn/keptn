package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogHandler_CreateLogEntries(t *testing.T) {
	myLogEntries := &models.CreateLogsRequest{
		Logs: []models.LogEntry{
			{
				IntegrationID: "my-id",
				Message:       "my message",
			},
		},
	}

	payload, _ := json.Marshal(myLogEntries)

	type fields struct {
		logManager *fake.ILogManagerMock
	}
	tests := []struct {
		name       string
		fields     fields
		request    *http.Request
		wantStatus int
	}{
		{
			name: "create log entry",
			fields: fields{
				logManager: &fake.ILogManagerMock{
					CreateLogEntriesFunc: func(entry models.CreateLogsRequest) error {
						return nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodPost, "/log", bytes.NewReader(payload)),
			wantStatus: http.StatusOK,
		},
		{
			name: "create log entry fails",
			fields: fields{
				logManager: &fake.ILogManagerMock{
					CreateLogEntriesFunc: func(entry models.CreateLogsRequest) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodPost, "/log", bytes.NewReader(payload)),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid payload",
			fields: fields{
				logManager: &fake.ILogManagerMock{
					CreateLogEntriesFunc: func(entry models.CreateLogsRequest) error {
						return errors.New("oops")
					},
				},
			},
			request:    httptest.NewRequest(http.MethodPost, "/log", bytes.NewReader([]byte("foo"))),
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := handler.NewLogHandler(tt.fields.logManager)

			router := gin.Default()
			router.POST("/log", func(c *gin.Context) {
				lh.CreateLogEntries(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

		})
	}
}

func TestLogHandler_GetLogEntries(t *testing.T) {
	type fields struct {
		logManager *fake.ILogManagerMock
	}
	tests := []struct {
		name              string
		fields            fields
		request           *http.Request
		wantStatus        int
		wantLogs          *models.GetLogsResponse
		wantGetLogsParams *models.GetLogParams
	}{
		{
			name: "get logs - no filter",
			fields: fields{
				&fake.ILogManagerMock{
					GetLogEntriesFunc: func(filter models.GetLogParams) (*models.GetLogsResponse, error) {
						return &models.GetLogsResponse{
							NextPageKey: 0,
							PageSize:    1,
							TotalCount:  1,
							Logs: []models.LogEntry{
								{
									IntegrationID: "my-id",
									Message:       "my message",
								},
							},
						}, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/log", nil),
			wantStatus: http.StatusOK,
			wantLogs: &models.GetLogsResponse{
				NextPageKey: 0,
				PageSize:    1,
				TotalCount:  1,
				Logs: []models.LogEntry{
					{
						IntegrationID: "my-id",
						Message:       "my message",
					},
				},
			},
			wantGetLogsParams: &models.GetLogParams{},
		},
		{
			name: "get logs - with filter",
			fields: fields{
				&fake.ILogManagerMock{
					GetLogEntriesFunc: func(filter models.GetLogParams) (*models.GetLogsResponse, error) {
						return &models.GetLogsResponse{
							NextPageKey: 0,
							PageSize:    1,
							TotalCount:  1,
							Logs: []models.LogEntry{
								{
									IntegrationID: "my-id",
									Message:       "my message",
								},
							},
						}, nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodGet, "/log?nextPageKey=1&pageSize=2&integrationId=my-id&fromTime=from&beforeTime=to", nil),
			wantStatus: http.StatusOK,
			wantLogs: &models.GetLogsResponse{
				NextPageKey: 0,
				PageSize:    1,
				TotalCount:  1,
				Logs: []models.LogEntry{
					{
						IntegrationID: "my-id",
						Message:       "my message",
					},
				},
			},
			wantGetLogsParams: &models.GetLogParams{
				NextPageKey: 1,
				PageSize:    2,
				LogFilter: models.LogFilter{
					IntegrationID: "my-id",
					FromTime:      "from",
					BeforeTime:    "to",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := handler.NewLogHandler(tt.fields.logManager)

			router := gin.Default()
			router.GET("/log", func(c *gin.Context) {
				lh.GetLogEntries(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantGetLogsParams != nil {
				require.Len(t, tt.fields.logManager.GetLogEntriesCalls(), 1)
				require.Equal(t, *tt.wantGetLogsParams, tt.fields.logManager.GetLogEntriesCalls()[0].Filter)
			}

			if tt.wantLogs != nil {
				logs := &models.GetLogsResponse{}
				err := json.Unmarshal(w.Body.Bytes(), logs)
				require.Nil(t, err)
				require.Equal(t, tt.wantLogs, logs)
			}
		})
	}
}

func TestLogHandler_DeleteLogEntries(t *testing.T) {
	type fields struct {
		logManager *fake.ILogManagerMock
	}
	tests := []struct {
		name                 string
		fields               fields
		request              *http.Request
		wantStatus           int
		wantDeleteLogsParams *models.DeleteLogParams
	}{
		{
			name: "delete logs - no filter",
			fields: fields{
				&fake.ILogManagerMock{
					DeleteLogEntriesFunc: func(params models.DeleteLogParams) error {
						return nil
					},
				},
			},
			request:              httptest.NewRequest(http.MethodDelete, "/log", nil),
			wantStatus:           http.StatusOK,
			wantDeleteLogsParams: &models.DeleteLogParams{},
		},
		{
			name: "delete logs - with filter",
			fields: fields{
				&fake.ILogManagerMock{
					DeleteLogEntriesFunc: func(filter models.DeleteLogParams) error {
						return nil
					},
				},
			},
			request:    httptest.NewRequest(http.MethodDelete, "/log?nextPageKey=1&pageSize=2&integrationId=my-id&fromTime=from&beforeTime=to", nil),
			wantStatus: http.StatusOK,
			wantDeleteLogsParams: &models.DeleteLogParams{
				LogFilter: models.LogFilter{
					IntegrationID: "my-id",
					FromTime:      "from",
					BeforeTime:    "to",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := handler.NewLogHandler(tt.fields.logManager)

			router := gin.Default()
			router.DELETE("/log", func(c *gin.Context) {
				lh.DeleteLogEntries(c)
			})
			w := performRequest(router, tt.request)

			require.Equal(t, tt.wantStatus, w.Code)

			if tt.wantDeleteLogsParams != nil {
				require.Len(t, tt.fields.logManager.DeleteLogEntriesCalls(), 1)
				require.Equal(t, *tt.wantDeleteLogsParams, tt.fields.logManager.DeleteLogEntriesCalls()[0].Params)
			}
		})
	}
}
