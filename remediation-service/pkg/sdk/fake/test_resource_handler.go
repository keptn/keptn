package fake

import "github.com/keptn/go-utils/pkg/api/models"

type TestResourceHandler struct {
}

func (t TestResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	panic("implement me")
}

func (t TestResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	panic("implement me")
}
