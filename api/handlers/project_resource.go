package handlers

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	configmodels "github.com/keptn/go-utils/pkg/configuration-service/models"
	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"

	models "github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/project_resource"
)

// PostProjectProjectNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameResourceHandlerFunc(params project_resource.PostProjectProjectNameResourceParams, principal *models.Principal) middleware.Responder {
	resourceHandler := configutils.NewResourceHandler(getConfigurationServiceURL())

	resourcesToUpload := []*configmodels.Resource{}
	for _, resource := range params.Resources.Resources {
		decodedStrBytes, err := b64.StdEncoding.DecodeString(*resource.ResourceContent)
		if err != nil {
			return project_resource.NewPostProjectProjectNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not base64-decode content of resource " + *resource.ResourceURI)})
		}
		descodedString := string(decodedStrBytes)
		resourcesToUpload = append(resourcesToUpload, &configmodels.Resource{
			ResourceContent: descodedString,
			ResourceURI:     resource.ResourceURI,
		})
	}

	_, err := resourceHandler.CreateProjectResources(params.ProjectName, resourcesToUpload)
	if err != nil {
		errorObj := &models.Error{}
		json.Unmarshal([]byte(err.Error()), errorObj)
		return project_resource.NewPostProjectProjectNameResourceDefault(500).WithPayload(errorObj)
	}

	return project_resource.NewPostProjectProjectNameResourceCreated()
}
