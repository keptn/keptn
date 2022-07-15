package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/shipyard-controller/handler"
)

type DebugController struct {
	DebugHandler handler.IDebugHandler
}

func NewDebugController(debugHandler handler.IDebugHandler) Controller {
	return &DebugController{DebugHandler: debugHandler}
}

func (controller DebugController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.Static("/ui", "./debug-ui")

	apiGroup.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/ui")
	})

	seq := apiGroup.Group("/sequence")
	{
		seq.GET("/project", controller.DebugHandler.GetAllProjects)
		seq.GET("/project/:project", controller.DebugHandler.GetAllSequencesForProject)
		seq.GET("/project/:project/shkeptncontext/:shkeptncontext", controller.DebugHandler.GetSequenceByID)
		seq.GET("/project/:project/shkeptncontext/:shkeptncontext/event", controller.DebugHandler.GetAllEvents)
		seq.GET("/project/:project/shkeptncontext/:shkeptncontext/event/:event_id", controller.DebugHandler.GetEventByID)
	}
}
