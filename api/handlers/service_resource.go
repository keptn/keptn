package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	models "github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service_resource"
)

// DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource
/*
func DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
}
*/

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	resourceHandler := keptnutils.NewResourceHandler("configuration-service")

	resourcesToUpload := []*keptnutils.Resource{}
	for _, resource := range params.Resources.Resources {
		resourcesToUpload = append(resourcesToUpload, &keptnutils.Resource{
			ResourceContent: *resource.ResourceContent,
			ResourceURI:     *resource.ResourceURI,
		})
	}
	newVersion, err := resourceHandler.CreateServiceResources(params.ProjectName, params.StageName, params.ServiceName, resourcesToUpload)
	if err != nil {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc updates a list of resources
func PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	resourceHandler := keptnutils.NewResourceHandler("configuration-service")

	resourcesToUpload := []*keptnutils.Resource{}
	for _, resource := range params.Resources.Resources {
		resourcesToUpload = append(resourcesToUpload, &keptnutils.Resource{
			ResourceContent: *resource.ResourceContent,
			ResourceURI:     *resource.ResourceURI,
		})
	}
	newVersion, err := resourceHandler.CreateServiceResources(params.ProjectName, params.StageName, params.ServiceName, resourcesToUpload)
	if err != nil {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
	}
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(&models.Version{
		Version: newVersion,
	})
}
