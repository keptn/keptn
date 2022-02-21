package handler_test

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/backend/fake"
	"github.com/keptn/keptn/secret-service/pkg/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetScopes(t *testing.T) {
	type fields struct {
		Backend backend.ScopeBackend
	}

	tests := []struct {
		name               string
		fields             fields
		expectedHTTPStatus int
		request            *http.Request
	}{
		{
			name: "GET Scope - SUCCESS",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					GetScopesFunc: func() ([]string, error) {
						return []string{}, nil
					},
				},
			},
			request:            httptest.NewRequest("GET", "/scope", nil),
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "GET Scope - Backend some error",
			fields: fields{
				Backend: &fake.SecretBackendMock{
					GetScopesFunc: func() ([]string, error) {
						return nil, errors.New("oops")
					},
				},
			},
			request:            httptest.NewRequest("GET", "/scope", nil),
			expectedHTTPStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			secretsHandler := handler.NewScopeHandler(tt.fields.Backend)
			handler := func(w http.ResponseWriter, r *http.Request) {
				c, _ := gin.CreateTestContext(w)
				c.Request = r
				secretsHandler.GetScopes(c)
			}

			w := httptest.NewRecorder()
			handler(w, tt.request)

			resp := w.Result()
			assert.Equal(t, tt.expectedHTTPStatus, resp.StatusCode)

		})
	}
}
