package handler_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/backend/fake"
	"github.com/keptn/keptn/secret-service/pkg/handler"
	"github.com/keptn/keptn/secret-service/pkg/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateNewHandler(t *testing.T) {
	secretsBackend := fake.SecretBackendMock{}
	secretsHandler := handler.NewSecretHandler(&secretsBackend)
	assert.NotNil(t, secretsHandler)
}

func TestHandler_CreateSecret(t *testing.T) {

	type fields struct {
		Backend backend.SecretBackend
	}

	tests := []struct {
		name               string
		fields             fields
		payload            string
		expectedHTTPStatus int
	}{
		{
			name: "POST Create Secret - SUCCESS",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			payload:            `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`,
			expectedHTTPStatus: http.StatusCreated,
		},
		{
			name: "POST Create Secret - Secret already exists",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return backend.ErrSecretAlreadyExists },
				},
			},
			payload:            `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`,
			expectedHTTPStatus: http.StatusConflict,
		},
		{
			name: "POST Create Secret - Backend some error",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return fmt.Errorf("failed to store secret in backend") },
				},
			},
			payload:            `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`,
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name: "POST Create Secret - Input INVALID",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			payload:            `SOME_WEIRD_INPUT`,
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/secrets", bytes.NewBuffer([]byte(tt.payload)))
			secretsHandler := handler.NewSecretHandler(tt.fields.Backend)
			secretsHandler.CreateSecret(c)
			assert.Equal(t, tt.expectedHTTPStatus, w.Code)
		})
	}
}

func TestHandler_DeleteSecret(t *testing.T) {

	type fields struct {
		Backend backend.SecretBackend
	}

	tests := []struct {
		name               string
		fields             fields
		scopeParam         string
		secretNameParam    string
		expectedHTTPStatus int
	}{
		{
			name: "DELETE Secret - SUCCESS",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "DELETE Secret - Backend some error",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return fmt.Errorf("failed to delete secret in backend") },
				},
			},
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name: "DELETE Secret - Secret not found",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return backend.ErrSecretNotFound },
				},
			},
			expectedHTTPStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodDelete, "/secrets&scope="+tt.scopeParam+"&name="+tt.secretNameParam, bytes.NewBuffer([]byte{}))
			secretsHandler := handler.NewSecretHandler(tt.fields.Backend)
			secretsHandler.DeleteSecret(c)
			assert.Equal(t, tt.expectedHTTPStatus, w.Code)
		})
	}
}
