package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/api-service/handler"
)

type MetadataController struct {
	MetadataHandler handler.IMetadataHandler
}

func NewMetadataController(metadataHandler handler.IMetadataHandler) *MetadataController {
	return &MetadataController{MetadataHandler: metadataHandler}
}

func (controller MetadataController) Inject(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/metadata", controller.MetadataHandler.GetMetadata)
}
