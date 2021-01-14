package api

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/statistics-service/controller"
	"github.com/keptn/keptn/statistics-service/operations"
	"net/http"
)

// HandleEvent godoc
// @Summary Handle event
// @Description Handle incoming cloud event
// @Tags Events
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   event     body    operations.Event     true        "Event type"
// @Success 200 "ok"
// @Failure 400 {object} operations.Error "Invalid payload"
// @Failure 500 {object} operations.Error "Internal error"
// @Router /event [post]
func HandleEvent(c *gin.Context) {
	event := &operations.Event{}
	if err := c.ShouldBindJSON(event); err != nil {
		c.JSON(http.StatusBadRequest, operations.Error{
			ErrorCode: 400,
			Message:   "Invalid request format",
		})
	}

	sb := controller.GetStatisticsBucketInstance()

	sb.AddEvent(*event)

	c.Status(http.StatusOK)
}
