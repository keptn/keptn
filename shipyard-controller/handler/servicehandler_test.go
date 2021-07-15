package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/v2/event"
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

func TestServiceHandler_CreateService(t *testing.T) {
	testServiceName := "my-service"
	type fields struct {
		serviceManager IServiceManager
		EventSender    common.EventSender
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
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
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
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
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
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
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
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
					},
				},
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
				EventSender: &fake.IEventSenderMock{
					SendEventFunc: func(eventMoqParam event.Event) error {
						return nil
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
				EventSender:    tt.fields.EventSender,
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

func TestServiceHandler_DeleteService(t *testing.T) {
	type fields struct {
		serviceManager IServiceManager
		EventSender    common.EventSender
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name                          string
		fields                        fields
		expectDeleteServiceToBeCalled bool
		expectHttpStatus              int
		expectJSONResponse            *operations.DeleteServiceResponse
		expectJSONError               *models.Error
	}{
		{
			name: "delete service",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					DeleteServiceFunc: func(projectName string, serviceName string) error {
						return nil
					},
				},
				EventSender: &fake.IEventSenderMock{SendEventFunc: func(eventMoqParam event.Event) error {
					return nil
				}},
			},
			expectDeleteServiceToBeCalled: true,
			expectHttpStatus:              http.StatusOK,
			expectJSONResponse:            &operations.DeleteServiceResponse{},
			expectJSONError:               nil,
		},
		{
			name: "delete service failed - expect 500",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					DeleteServiceFunc: func(projectName string, serviceName string) error {
						return errors.New("internal error")
					},
				},
				EventSender: &fake.IEventSenderMock{SendEventFunc: func(eventMoqParam event.Event) error {
					return nil
				}},
			},
			expectDeleteServiceToBeCalled: true,
			expectHttpStatus:              http.StatusInternalServerError,
			expectJSONResponse:            nil,
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
				gin.Param{Key: "service", Value: "test-service"},
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte{}))

			sh := &ServiceHandler{
				serviceManager: tt.fields.serviceManager,
				EventSender:    tt.fields.EventSender,
			}

			sh.DeleteService(c)

			mockServiceManager := tt.fields.serviceManager.(*fake.IServiceManagerMock)

			if tt.expectDeleteServiceToBeCalled {
				if len(mockServiceManager.DeleteServiceCalls()) != 1 {
					t.Errorf("serviceManager.DeleteService() was called %d times, expected %d", len(mockServiceManager.DeleteServiceCalls()), 1)
				}

				assert.Equal(t, "test-project", mockServiceManager.DeleteServiceCalls()[0].ProjectName)
				assert.Equal(t, "test-service", mockServiceManager.DeleteServiceCalls()[0].ServiceName)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)
			responseBytes, _ := ioutil.ReadAll(w.Body)
			if tt.expectJSONResponse != nil {
				response := &operations.DeleteServiceResponse{}
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

func TestServiceHandler_GetService(t *testing.T) {
	type fields struct {
		serviceManager IServiceManager
		EventSender    common.EventSender
	}
	tests := []struct {
		name                       string
		fields                     fields
		expectGetServiceToBeCalled bool
		expectHttpStatus           int
		expectJSONResponse         *models.ExpandedService
		expectJSONError            *models.Error
	}{
		{
			name: "get available service - expect 200",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
						return &models.ExpandedService{
							ServiceName: "test-service",
						}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusOK,
			expectJSONResponse: &models.ExpandedService{
				ServiceName: "test-service",
			},
			expectJSONError: nil,
		},
		{
			name: "get unavailable service - expect 404",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
						return nil, errServiceNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusNotFound,
			expectJSONResponse:         nil,
			expectJSONError: &models.Error{
				Code:    http.StatusNotFound,
				Message: stringp(errServiceNotFound.Error()),
			},
		},
		{
			name: "get unavailable stage - expect 404",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
						return nil, errStageNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusNotFound,
			expectJSONResponse:         nil,
			expectJSONError: &models.Error{
				Code:    http.StatusNotFound,
				Message: stringp(errStageNotFound.Error()),
			},
		},
		{
			name: "get unavailable project - expect 404",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
						return nil, errProjectNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusNotFound,
			expectJSONResponse:         nil,
			expectJSONError: &models.Error{
				Code:    http.StatusNotFound,
				Message: stringp(errProjectNotFound.Error()),
			},
		},
		{
			name: "internal error - expect 500",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetServiceFunc: func(projectName string, stageName string, serviceName string) (*models.ExpandedService, error) {
						return nil, errors.New("internal error")
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusInternalServerError,
			expectJSONResponse:         nil,
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
				gin.Param{Key: "stage", Value: "test-stage"},
				gin.Param{Key: "service", Value: "test-service"},
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte{}))

			sh := NewServiceHandler(tt.fields.serviceManager, tt.fields.EventSender)

			sh.GetService(c)

			mockServiceManager := tt.fields.serviceManager.(*fake.IServiceManagerMock)

			if tt.expectGetServiceToBeCalled {
				if len(mockServiceManager.GetServiceCalls()) != 1 {
					t.Errorf("serviceManager.GetService() was called %d times, expected %d", len(mockServiceManager.GetServiceCalls()), 1)
				}

				assert.Equal(t, "test-project", mockServiceManager.GetServiceCalls()[0].ProjectName)
				assert.Equal(t, "test-stage", mockServiceManager.GetServiceCalls()[0].StageName)
				assert.Equal(t, "test-service", mockServiceManager.GetServiceCalls()[0].ServiceName)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)
			responseBytes, _ := ioutil.ReadAll(w.Body)
			if tt.expectJSONResponse != nil {
				response := &models.ExpandedService{}
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

func TestServiceHandler_GetServices(t *testing.T) {
	type fields struct {
		serviceManager IServiceManager
		EventSender    common.EventSender
	}
	tests := []struct {
		name                       string
		fields                     fields
		expectGetServiceToBeCalled bool
		expectHttpStatus           int
		expectJSONResponse         *models.ExpandedServices
		expectJSONError            *models.Error
	}{
		{
			name: "get available service - expect 200",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetAllServicesFunc: func(projectName string, stageName string) ([]*models.ExpandedService, error) {
						return []*models.ExpandedService{
							{
								ServiceName: "test-service",
							},
						}, nil
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusOK,
			expectJSONResponse: &models.ExpandedServices{
				NextPageKey: "0",
				PageSize:    0,
				Services: []*models.ExpandedService{
					{
						ServiceName: "test-service",
					},
				},
				TotalCount: 1,
			},
			expectJSONError: nil,
		},
		{
			name: "get unavailable stage - expect 404",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetAllServicesFunc: func(projectName string, stageName string) ([]*models.ExpandedService, error) {
						return nil, errStageNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusNotFound,
			expectJSONResponse:         nil,
			expectJSONError: &models.Error{
				Code:    http.StatusNotFound,
				Message: stringp(errStageNotFound.Error()),
			},
		},
		{
			name: "get unavailable project - expect 404",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetAllServicesFunc: func(projectName string, stageName string) ([]*models.ExpandedService, error) {
						return nil, errProjectNotFound
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusNotFound,
			expectJSONResponse:         nil,
			expectJSONError: &models.Error{
				Code:    http.StatusNotFound,
				Message: stringp(errProjectNotFound.Error()),
			},
		},
		{
			name: "internal error - expect 500",
			fields: fields{
				serviceManager: &fake.IServiceManagerMock{
					GetAllServicesFunc: func(projectName string, stageName string) ([]*models.ExpandedService, error) {
						return nil, errors.New("internal error")
					},
				},
				EventSender: &fake.IEventSenderMock{},
			},
			expectGetServiceToBeCalled: true,
			expectHttpStatus:           http.StatusInternalServerError,
			expectJSONResponse:         nil,
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
				gin.Param{Key: "stage", Value: "test-stage"},
				gin.Param{Key: "service", Value: "test-service"},
			}

			c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer([]byte{}))

			sh := NewServiceHandler(tt.fields.serviceManager, tt.fields.EventSender)

			sh.GetServices(c)

			mockServiceManager := tt.fields.serviceManager.(*fake.IServiceManagerMock)

			if tt.expectGetServiceToBeCalled {
				if len(mockServiceManager.GetAllServicesCalls()) != 1 {
					t.Errorf("serviceManager.GetAllServices() was called %d times, expected %d", len(mockServiceManager.GetAllServicesCalls()), 1)
				}

				assert.Equal(t, "test-project", mockServiceManager.GetAllServicesCalls()[0].ProjectName)
				assert.Equal(t, "test-stage", mockServiceManager.GetAllServicesCalls()[0].StageName)
			}

			assert.Equal(t, tt.expectHttpStatus, w.Code)
			responseBytes, _ := ioutil.ReadAll(w.Body)
			if tt.expectJSONResponse != nil {
				response := &models.ExpandedServices{}
				_ = json.Unmarshal(responseBytes, response)

				assert.EqualValues(t, tt.expectJSONResponse, response)
			} else if tt.expectJSONError != nil {
				errorResponse := &models.Error{}

				_ = json.Unmarshal(responseBytes, errorResponse)

				assert.Equal(t, tt.expectJSONError, errorResponse)
			}
		})
	}
}
