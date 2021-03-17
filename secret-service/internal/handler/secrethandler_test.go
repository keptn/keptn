package handler_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/internal/backend"
	"github.com/keptn/keptn/secret-service/internal/backend/fake"
	"github.com/keptn/keptn/secret-service/internal/handler"
	"github.com/keptn/keptn/secret-service/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateNewHandler(t *testing.T) {
	secretStore := fake.SecretStoreMock{}
	handler := handler.NewSecretHandler(&secretStore)
	assert.NotNil(t, handler)
}

func Test_CreateSecret(t *testing.T) {

	jsonPayload := `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/secrets", bytes.NewBuffer([]byte(jsonPayload)))

	secretStore := fake.SecretStoreMock{}

	secretStore.CreateSecretFunc = func(secret model.Secret) error {
		return nil
	}

	secretHandler := handler.SecretHandler{
		SecretBackend: &secretStore,
	}
	secretHandler.CreateSecret(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_CreateSecret(t *testing.T) {

	type fields struct {
		Backend backend.SecretStore
	}

	tests := []struct {
		name             string
		fields           fields
		payload          string
		expectHttpStatus int
	}{
		{
			name: "POST Create Secret - SUCCESS",
			fields: fields{
				Backend: &fake.SecretStoreMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			payload:          `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`,
			expectHttpStatus: http.StatusOK,
		},
		{
			name: "POST Create Secret - Backend FAILED",
			fields: fields{
				Backend: &fake.SecretStoreMock{
					CreateSecretFunc: func(secret model.Secret) error { return fmt.Errorf("Failed to store secret in backend") },
				},
			},
			payload:          `{"name":"my-secret","scope":"my-scope","data":{"username":"keptn"}}`,
			expectHttpStatus: http.StatusInternalServerError,
		},
		{
			name: "POST Create Secret - Input INVALID",
			fields: fields{
				Backend: &fake.SecretStoreMock{
					CreateSecretFunc: func(secret model.Secret) error { return nil },
				},
			},
			payload:          `SOME_WEIRD_INPUT`,
			expectHttpStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/secrets", bytes.NewBuffer([]byte(tt.payload)))
			handler := handler.NewSecretHandler(tt.fields.Backend)
			handler.CreateSecret(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)
		})
	}
}
