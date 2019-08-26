package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	models "github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/service_resource"
)

type notFoundError int

const (
	projectNotFound notFoundError = iota
	stageNotFound
	serviceNotFound
)

// DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource
/*
func DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
}
*/

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	serviceFoundErr := checkServiceAvailability(params.ProjectName, params.StageName, params.ServiceName)
	if serviceFoundErr >= 0 {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(400).WithPayload(getNotFoundErrorMessage(serviceFoundErr, params.ProjectName, params.StageName, params.ServiceName))
	}
	resourceHandler := keptnutils.NewResourceHandler("configuration-service:8080")

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
	serviceFoundErr := checkServiceAvailability(params.ProjectName, params.StageName, params.ServiceName)
	if serviceFoundErr >= 0 {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceDefault(400).WithPayload(getNotFoundErrorMessage(serviceFoundErr, params.ProjectName, params.StageName, params.ServiceName))
	}
	resourceHandler := keptnutils.NewResourceHandler("configuration-service:8080")

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

func getNotFoundErrorMessage(notFoundErrorcode notFoundError, project string, stage string, service string) *models.Error {
	errorMsg := &models.Error{
		Code: 400,
	}
	switch notFoundErrorcode {
	case projectNotFound:
		errorMsg.Message = swag.String("Project " + project + " does not exist.")
		break
	case stageNotFound:
		errorMsg.Message = swag.String("Stage " + stage + " does not exist within project " + project)
		break
	case serviceNotFound:
		errorMsg.Message = swag.String("Service " + service + " does not exist within stage " + stage + " of project " + project)
		break
	default:
		errorMsg.Message = swag.String("Unknown error")
		break
	}

	return errorMsg
}

func checkServiceAvailability(project string, stage string, service string) notFoundError {
	// check if project exists
	status, err := doGet("http://configuration-service:8080/v1/project/" + project)

	if err != nil || status == 404 {
		return projectNotFound
	}
	status, err = doGet("http://configuration-service:8080/v1/project/" + project + "/stage/" + stage)
	if err != nil || status == 404 {
		return stageNotFound
	}
	status, err = doGet("http://configuration-service:8080/v1/project/" + project + "/stage/" + stage + "/service/" + service)
	if err != nil || status == 404 {
		return serviceNotFound
	}
	return -1
}

func doGet(uri string) (int, error) {
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}
