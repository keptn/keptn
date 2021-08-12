package fake

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
	"log"
)

type TestResourceHandler struct {
}

func (t TestResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s/%s", project, stage, service, resourceURI)), nil
}

func (t TestResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s", project, stage, resourceURI)), nil
}

func newResourceFromFile(filename string) *models.Resource {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to locate resources requested by the service: %s", err.Error())
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}
