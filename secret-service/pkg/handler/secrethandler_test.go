package handler_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/internal/backend"
	"github.com/keptn/keptn/secret-service/internal/backend/fake"
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
			name: "POST Create Secret - Backend FAILED",
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
