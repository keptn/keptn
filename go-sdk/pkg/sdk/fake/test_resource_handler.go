package fake

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
)

type TestResourceHandler struct {
	Resource models.Resource
}

func (t TestResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s/%s", project, stage, service, resourceURI)), nil
}

func (t TestResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s", project, stage, resourceURI)), nil
}

func (t TestResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s", project, resourceURI)), nil
}

func newResourceFromFile(filename string) *models.Resource {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}

type StringResourceHandler struct {
	ResourceContent string
}

func (s StringResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

func (s StringResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

func (s StringResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

type FailingResourceHandler struct {
}

func (f FailingResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}

func (f FailingResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}

func (f FailingResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}
