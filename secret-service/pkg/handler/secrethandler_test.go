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
		request            *http.Request
		expectedHTTPStatus int
	}{
		{
			name: "POST Create Secret - SUCCESS",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			request:            httptest.NewRequest("POST", "/secret", bytes.NewBuffer([]byte(`{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`))),
			expectedHTTPStatus: http.StatusCreated,
		},
		{
			name: "POST Create Secret - Secret already exists",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return backend.ErrSecretAlreadyExists },
				},
			},
			request:            httptest.NewRequest("POST", "/secret", bytes.NewBuffer([]byte(`{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`))),
			expectedHTTPStatus: http.StatusConflict,
		},
		{
			name: "POST Create Secret - Backend some error",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return fmt.Errorf("failed to store secret in backend") },
				},
			},
			request:            httptest.NewRequest("POST", "/secret", bytes.NewBuffer([]byte(`{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`))),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name: "POST Create Secret - Input INVALID",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			request:            httptest.NewRequest("POST", "/secret", bytes.NewBuffer([]byte(`SOME_WEIRD_INPUT`))),
			expectedHTTPStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secretsHandler := handler.NewSecretHandler(tt.fields.Backend)
			handler := func(w http.ResponseWriter, r *http.Request) {
				c, _ := gin.CreateTestContext(w)
				c.Request = r
				secretsHandler.CreateSecret(c)
			}

			w := httptest.NewRecorder()
			handler(w, tt.request)

			resp := w.Result()
			assert.Equal(t, tt.expectedHTTPStatus, resp.StatusCode)

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
		expectedHTTPStatus int
		request            *http.Request
	}{
		{
			name: "DELETE Secret - SUCCESS",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			request:            httptest.NewRequest("DELETE", "/secret?name=my-secret&scope=my-scope", nil),
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "DELETE Secret - Backend some error",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return fmt.Errorf("failed to delete secret in backend") },
				},
			},
			request:            httptest.NewRequest("DELETE", "/secrets?name=my-secret&scope=my-scope", nil),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
		{
			name: "DELETE Secret - Secret not found",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					DeleteSecretFunc: func(secret model.Secret) error { return backend.ErrSecretNotFound },
				},
			},
			request:            httptest.NewRequest("DELETE", "/secrets?name=my-secret&scope=my-scope", nil),
			expectedHTTPStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secretsHandler := handler.NewSecretHandler(tt.fields.Backend)
			handler := func(w http.ResponseWriter, r *http.Request) {
				c, _ := gin.CreateTestContext(w)
				c.Request = r
				secretsHandler.DeleteSecret(c)
			}

			w := httptest.NewRecorder()
			handler(w, tt.request)

			resp := w.Result()
			assert.Equal(t, tt.expectedHTTPStatus, resp.StatusCode)

		})
	}
}
