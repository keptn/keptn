package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceHandler_CreateService(t *testing.T) {
	testServiceName := "my-service"
	type fields struct {
		serviceManager IServiceManager
	}
	tests := []struct {
		name                          string
		fields                        fields
		jsonPayload                   string
		expectCreateServiceToBeCalled bool
		expectCreateServiceParams     *operations.CreateServiceParams
		expectHttpStatus              int
		expectJSONResponse            *operations.CreateServiceResponse
		expectJSONError               *models.Error
	}{
		{
			name: "create service - return 200",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					CreateServiceFunc: func(projectName string, params *operations.CreateServiceParams) error {
						return nil
					},
				},
			},
			jsonPayload:                   `{"serviceName":"my-service"}`,
			expectCreateServiceToBeCalled: true,
			expectCreateServiceParams: &operations.CreateServiceParams{
				ServiceName: &testServiceName,
			},
			expectHttpStatus:   http.StatusOK,
			expectJSONResponse: &operations.CreateServiceResponse{},
			expectJSONError:    nil,
		},
		{
			name: "service already exists - return 409",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					CreateServiceFunc: func(projectName string, params *operations.CreateServiceParams) error {
						return errServiceAlreadyExists
					},
				},
			},
			jsonPayload:                   `{"serviceName":"my-service"}`,
			expectCreateServiceToBeCalled: true,
			expectCreateServiceParams: &operations.CreateServiceParams{
				ServiceName: &testServiceName,
			},
			expectHttpStatus:   http.StatusConflict,
			expectJSONResponse: &operations.CreateServiceResponse{},
			expectJSONError: &models.Error{
				Code:    http.StatusConflict,
				Message: stringp(errServiceAlreadyExists.Error()),
			},
		},
		{
			name: "invalid payload - return 400",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{},
			},
			jsonPayload:                   `invalid`,
			expectCreateServiceToBeCalled: false,
			expectCreateServiceParams: &operations.CreateServiceParams{
				ServiceName: &testServiceName,
			},
			expectHttpStatus:   http.StatusBadRequest,
			expectJSONResponse: &operations.CreateServiceResponse{},
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: stringp("Invalid request format"),
			},
		},
		{
			name: "invalid service name - return 400",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{},
			},
			jsonPayload:                   `{"serviceName":"my/service"}`,
			expectCreateServiceToBeCalled: false,
			expectCreateServiceParams: &operations.CreateServiceParams{
				ServiceName: &testServiceName,
			},
			expectHttpStatus:   http.StatusBadRequest,
			expectJSONResponse: &operations.CreateServiceResponse{},
			expectJSONError: &models.Error{
				Code:    http.StatusBadRequest,
				Message: stringp("Service name contains special character(s). \" +\n\t\t\t\"The service name has to be a valid Unix directory name. For details see \" +\n\t\t\t\"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/"),
			},
		},
		{
			name: "internal error - return 500",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					CreateServiceFunc: func(projectName string, params *operations.CreateServiceParams) error {
						return errors.New("internal error")
					},
				},
			},
			jsonPayload:                   `{"serviceName":"my-service"}`,
			expectCreateServiceToBeCalled: false,
			expectCreateServiceParams: &operations.CreateServiceParams{
				ServiceName: &testServiceName,
			},
			expectHttpStatus:   http.StatusInternalServerError,
			expectJSONResponse: &operations.CreateServiceResponse{},
			expectJSONError: &models.Error{
				Code:    http.StatusInternalServerError,
				Message: stringp("internal error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "test-project"},
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte(tt.jsonPayload)))

			sh := &ServiceHandler{
				serviceManager: tt.fields.serviceManager,
				logger:         keptncommon.NewLogger("", "", ""),
			}

			sh.CreateService(c)

			mockServiceManager := tt.fields.serviceManager.(*fake.IServiceManagerMock)

			if tt.expectCreateServiceToBeCalled {
				if len(mockServiceManager.CreateServiceCalls()) != 1 {
					t.Errorf("serviceManager.CreateService() was called %d times, expected %d", len(mockServiceManager.CreateServiceCalls()), 1)
				}

				assert.Equal(t, tt.expectCreateServiceParams, mockServiceManager.CreateServiceCalls()[0].Params)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)
			responseBytes, _ := ioutil.ReadAll(w.Body)
			if tt.expectJSONResponse != nil {
				response := &operations.CreateServiceResponse{}
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

func stringp(s string) *string {
	return &s
}
