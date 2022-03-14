package handler

import (
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
)

//go:generate moq  -pkg fake -out ./fake/resource_handler_mock.go . IResourceHandler
type IResourceHandler interface {
	GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error)
}
