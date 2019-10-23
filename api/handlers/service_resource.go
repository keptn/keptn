package handlers

import (
	b64 "encoding/base64"
	"encoding/json"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	configmodels "github.com/keptn/go-utils/pkg/configuration-service/models"
	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"

	models "github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service_resource"
)

type notFoundError int

const (
	projectNotFound notFoundError = iota
	stageNotFound
	serviceNotFound
)

func getConfigurationServiceURL() string {
	return "http://" + os.Getenv("CONFIGURATION_URI")
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource
/*
func DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
}
*/

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	resourceHandler := configutils.NewResourceHandler(getConfigurationServiceURL())

	resourcesToUpload := []*configmodels.Resource{}
	for _, resource := range params.Resources.Resources {
		decodedStrBytes, err := b64.StdEncoding.DecodeString(*resource.ResourceContent)
		if err != nil {
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not base64-decode content of resource " + *resource.ResourceURI)})
		}
		descodedString := string(decodedStrBytes)
		resourcesToUpload = append(resourcesToUpload, &configmodels.Resource{
			ResourceContent: descodedString,
			ResourceURI:     resource.ResourceURI,
		})
	}

	_, err := resourceHandler.CreateServiceResources(params.ProjectName, params.StageName, params.ServiceName, resourcesToUpload)
	if err != nil {
		errorObj := &models.Error{}
		json.Unmarshal([]byte(err.Error()), errorObj)
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(errorObj)
	}

	return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceCreated()
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc updates a list of resources
func PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	resourceHandler := configutils.NewResourceHandler(getConfigurationServiceURL())

	resourcesToUpload := []*configmodels.Resource{}
	for _, resource := range params.Resources.Resources {
		decodedStrBytes, err := b64.StdEncoding.DecodeString(*resource.ResourceContent)
		if err != nil {
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not base64-decode content of resource " + *resource.ResourceURI)})
		}
		descodedString := string(decodedStrBytes)
		resourcesToUpload = append(resourcesToUpload, &configmodels.Resource{
			ResourceContent: descodedString,
			ResourceURI:     resource.ResourceURI,
		})
	}

	_, err := resourceHandler.CreateServiceResources(params.ProjectName, params.StageName, params.ServiceName, resourcesToUpload)
	if err != nil {
		errorObj := &models.Error{}
		json.Unmarshal([]byte(err.Error()), errorObj)
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(errorObj)
	}

	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceCreated()
}
