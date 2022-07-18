package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/keptn/shipyard-controller/handler/fake"
	"github.com/stretchr/testify/assert"
)

func TestGetAllEvents(t *testing.T) {

	type fields struct {
		DebugManager IDebugManager
	}

	tests := []struct {
		name             string
		fields           fields
		expectHttpStatus int
	}{
		{
			name: "GET events project not found",
			fields: fields{
				DebugManager: &fake.IDebugManagerMock{
					GetAllEventsFunc: func(projectName string, shkeptncontext string) ([]*models.KeptnContextExtendedCE, error) {
						var test []*models.KeptnContextExtendedCE
						return test, ErrProjectNotFound
					},
				},
			},
			expectHttpStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "", bytes.NewBuffer([]byte{}))
			c.Params = gin.Params{
				gin.Param{Key: "project", Value: "invalid"},
				gin.Param{Key: "shkeptncontext", Value: "asdf"},
			}

			handler := NewDebugHandler(tt.fields.DebugManager)
			//handler.GetStage(c)
			handler.GetAllEvents(c)
			assert.Equal(t, tt.expectHttpStatus, w.Code)

		})
	}
}
