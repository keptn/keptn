package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/internal/common"
	"github.com/keptn/keptn/shipyard-controller/internal/handler/fake"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/stretchr/testify/assert"
)

func TestEvaluationParamsValidator(t *testing.T) {
	type args struct {
		params *models.CreateEvaluationParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "End and time frame both set",
			args: args{params: &models.CreateEvaluationParams{
				End:       "t1",
				Timeframe: "d1",
			}},
			wantErr: assert.Error,
		},
		{
			name: "End and Timeframe both not set",
			args: args{params: &models.CreateEvaluationParams{
				Start: "t1",
			}},
			wantErr: assert.Error,
		},
		{
			name: "End set but start not set",
			args: args{params: &models.CreateEvaluationParams{
				End: "t1",
			}},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EvaluationParamsValidator{}
			tt.wantErr(t, e.Validate(tt.args.params), fmt.Sprintf("validateEvaluationParams(%v)", tt.args.params))
		})
	}
}

func TestEvaluationHandler_CreateEvaluation(t *testing.T) {
	type fields struct {
		EvaluationManager IEvaluationManager
	}
	tests := []struct {
		name                             string
		fields                           fields
		jsonPayload                      string
		expectCreateEvaluationToBeCalled bool
		expectCreateEvaluationParams     *models.CreateEvaluationParams
		expectHttpStatus                 int
		expectJSONResponse               *models.CreateEvaluationResponse
		expectJSONError                  *models.Error
	}{
		{
			name:                             "correctly pass time values (start and end) from json payload",
			expectCreateEvaluationToBeCalled: true,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload: `{"start":"2021-01-02T15:00:00", "end":"2021-01-02T15:10:00"}`,
			expectCreateEvaluationParams: &models.CreateEvaluationParams{
				Labels:    nil,
				Start:     "2021-01-02T15:00:00",
				End:       "2021-01-02T15:10:00",
				Timeframe: "",
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &models.CreateEvaluationResponse{KeptnContext: "test-context"},
			expectJSONError:    nil,
		},
		{
			name:                             "correctly pass time values (start and time) from json payload",
			expectCreateEvaluationToBeCalled: true,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload: `{"start":"2021-01-02T15:00:00", "timeframe":"5m"}`,
			expectCreateEvaluationParams: &models.CreateEvaluationParams{
				Labels:    nil,
				Start:     "2021-01-02T15:00:00",
				End:       "",
				Timeframe: "5m",
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &models.CreateEvaluationResponse{KeptnContext: "test-context"},
			expectJSONError:    nil,
		},
		{
			name:                             "correct params - just timeframe",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload: `{"timeframe":"5m"}`,
			expectCreateEvaluationParams: &models.CreateEvaluationParams{
				Labels:    nil,
				Start:     "",
				End:       "",
				Timeframe: "5m",
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &models.CreateEvaluationResponse{KeptnContext: "test-context"},
			expectJSONError:    nil,
		},
		{
			name:                             "no time specifications",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload: `{}`,
			expectCreateEvaluationParams: &models.CreateEvaluationParams{
				Labels:    nil,
				Start:     "",
				End:       "",
				Timeframe: "",
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &models.CreateEvaluationResponse{KeptnContext: "test-context"},
			expectJSONError:    nil,
		},
		{
			name:                             "invalid params - end and timeframe together",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload:      `{"start":"2021-01-02T15:00:00", "end":"2021-01-02T15:10:00", "timeframe":"5m"}`,
			expectHttpStatus: http.StatusBadRequest,
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: common.Stringp("Invalid request format: timeframe and end time specifications cannot be set together"),
			},
		},
		{
			name:                             "invalid params - end and timeframe together without start",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload:      `{"end":"2021-01-02T15:10:00", "timeframe":"5m"}`,
			expectHttpStatus: http.StatusBadRequest,
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: common.Stringp("Invalid request format: timeframe and end time specifications cannot be set together"),
			},
		},
		{
			name:                             "invalid params - end without start",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload:      `{"end":"2021-01-02T15:10:00"}`,
			expectHttpStatus: http.StatusBadRequest,
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: common.Stringp("Invalid request format: end time specifications cannot be set without start parameter"),
			},
		},
		{
			name:                             "invalid payload - do not call evaluationManager.CreateEvaluation",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return &models.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload:      `invalid`,
			expectHttpStatus: http.StatusBadRequest,
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: common.Stringp("Invalid request format: invalid character 'i' looking for beginning of value"),
			},
		},
		{
			name:                             "internal error when creating evaluation - return 500",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *models.CreateEvaluationParams) (*models.CreateEvaluationResponse, *models.Error) {
						return nil, &models.Error{
							Code:    evaluationErrSendEventFailed,
							Message: common.Stringp("failed to send event"),
						}
					},
				},
			},
			jsonPayload:      `{"start":"2021-01-02T15:00:00", "end":"2021-01-02T15:10:00"}`,
			expectHttpStatus: http.StatusInternalServerError,
			expectJSONError: &models.Error{
				Code:    evaluationErrSendEventFailed,
				Message: common.Stringp("failed to send event"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "test-project"},
				gin.Param{Key: "stage", Value: "test-stage"},
				gin.Param{Key: "service", Value: "test-service"},
			}

			eh := NewEvaluationHandler(tt.fields.EvaluationManager)

			eh.CreateEvaluation(c)

			mockEvaluationManager := tt.fields.EvaluationManager.(*fake.IEvaluationManagerMock)

			if tt.expectCreateEvaluationToBeCalled {
				if len(mockEvaluationManager.CreateEvaluationCalls()) != 1 {
					t.Errorf("evaluationManager.CreateEvaluation() was called %d times, expected %d", len(mockEvaluationManager.CreateEvaluationCalls()), 1)
				}

				assert.Equal(t, tt.expectCreateEvaluationParams, mockEvaluationManager.CreateEvaluationCalls()[0].Params)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)

			responseBytes, _ := ioutil.ReadAll(w.Body)
			if tt.expectJSONResponse != nil {
				response := &models.CreateEvaluationResponse{}
				_ = json.Unmarshal(responseBytes, response)

				assert.Equal(t, tt.expectJSONResponse, response)
			} else if tt.expectJSONError != nil {
				errorResponse := &models.Error{}

				_ = json.Unmarshal(responseBytes, errorResponse)

				assert.Equal(t, tt.expectJSONError, errorResponse)
			}
		})
	}
}

func Test_getHTTPStatusForError(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "service not available - return 400",
			args: args{
				code: evaluationErrServiceNotAvailable,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "invalid timeframe - return 400",
			args: args{
				code: evaluationErrInvalidTimeframe,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "default - return 500",
			args: args{
				code: evaluationErrSendEventFailed,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHTTPStatusForError(tt.args.code); got != tt.want {
				t.Errorf("getHTTPStatusForError() = %v, want %v", got, tt.want)
			}
		})
	}
}
