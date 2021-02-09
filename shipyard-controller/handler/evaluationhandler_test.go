package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEvaluationHandler_CreateEvaluation(t *testing.T) {
	type fields struct {
		EvaluationManager IEvaluationManager
	}
	tests := []struct {
		name                             string
		fields                           fields
		jsonPayload                      string
		expectCreateEvaluationToBeCalled bool
		expectCreateEvaluationParams     *operations.CreateEvaluationParams
		expectHttpStatus                 int
		expectJSONResponse               *operations.CreateEvaluationResponse
		expectJSONError                  *models.Error
	}{
		{
			name:                             "correctly pass time values (start and end) from json payload",
			expectCreateEvaluationToBeCalled: true,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
						return &operations.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload: `{"start":"2021-01-02T15:00:00", "end":"2021-01-02T15:10:00"}`,
			expectCreateEvaluationParams: &operations.CreateEvaluationParams{
				Labels:    nil,
				Start:     "2021-01-02T15:00:00",
				End:       "2021-01-02T15:10:00",
				Timeframe: "",
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &operations.CreateEvaluationResponse{KeptnContext: "test-context"},
			expectJSONError:    nil,
		},
		{
			name:                             "invalid payload - do not call evaluationManager.CreateEvaluation",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
						return &operations.CreateEvaluationResponse{
							KeptnContext: "test-context",
						}, nil
					},
				},
			},
			jsonPayload:      `invalid`,
			expectHttpStatus: http.StatusBadRequest,
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: common.Stringp("Invalid request format"),
			},
		},
		{
			name:                             "internal error when creating evaluation - return 500",
			expectCreateEvaluationToBeCalled: false,
			fields: fields{
				EvaluationManager: &fake.IEvaluationManagerMock{
					CreateEvaluationFunc: func(project string, stage string, service string, params *operations.CreateEvaluationParams) (*operations.CreateEvaluationResponse, *models.Error) {
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
				response := &operations.CreateEvaluationResponse{}
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
