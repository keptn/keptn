package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	utils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service"
	"gopkg.in/yaml.v2"
)

type serviceMetadata struct {
	ServiceName       string
	CreationTimestamp string
}

// GetProjectProjectNameStageStageNameServiceHandlerFunc get list of services
func GetProjectProjectNameStageStageNameServiceHandlerFunc(params service.GetProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service.NewGetProjectProjectNameStageStageNameServiceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}
	
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName, *params.DisableUpstreamSync) {
		return service.NewGetProjectProjectNameStageStageNameServiceNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
	}

	var payload = &models.Services{
		PageSize:    0,
		NextPageKey: "0",
		TotalCount:  0,
		Services:    []*models.Service{},
	}

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	err := common.CheckoutBranch(params.ProjectName, params.StageName, *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service.NewGetProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	files, err := ioutil.ReadDir(projectConfigPath)
	if err != nil {
		return service.NewGetProjectProjectNameStageStageNameServiceOK().WithPayload(payload)
	}

	filteredFiles := []os.FileInfo{}
	for _, file := range files {
		if file.IsDir() && common.FileExists(projectConfigPath+"/"+file.Name()+"/metadata.yaml") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	paginationInfo := common.Paginate(len(filteredFiles), params.PageSize, params.NextPageKey)

	totalCount := len(filteredFiles)
	if paginationInfo.NextPageKey < int64(totalCount) {
		for _, f := range filteredFiles[paginationInfo.NextPageKey:paginationInfo.EndIndex] {
			var service = &models.Service{ServiceName: f.Name()}
			payload.Services = append(payload.Services, service)
		}
	}

	payload.TotalCount = float64(totalCount)
	payload.NextPageKey = paginationInfo.NewNextPageKey
	return service.NewGetProjectProjectNameStageStageNameServiceOK().WithPayload(payload)
}

// GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc get the specified service
func GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.GetProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Project not found")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName, *params.DisableUpstreamSync) {
		return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Stage not found")})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, *params.DisableUpstreamSync) {
		return service.NewGetProjectProjectNameStageStageNameServiceServiceNameNotFound().WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}
	var serviceResponse = &models.Service{
		ServiceName: params.ServiceName,
	}
	return service.NewGetProjectProjectNameStageStageNameServiceServiceNameOK().WithPayload(serviceResponse)
}

// PostProjectProjectNameStageStageNameServiceHandlerFunc creates a new service
func PostProjectProjectNameStageStageNameServiceHandlerFunc(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	servicePath := projectConfigPath + "/" + params.Service.ServiceName

	if !common.StageExists(params.ProjectName, params.StageName, false) {
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}

	if common.ServiceExists(params.ProjectName, params.StageName, params.Service.ServiceName, false) {
		return service.NewPostProjectProjectNameStageStageNameServiceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service already exists")})
	}

	logger.Debug("Creating new resource(s) in: " + projectConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	
	err = os.MkdirAll(servicePath, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create service directory")})
	}

	newServiceMetadata := &serviceMetadata{
		ServiceName:       params.Service.ServiceName,
		CreationTimestamp: time.Now().String(),
	}

	metadataString, err := yaml.Marshal(newServiceMetadata)
	err = common.WriteFile(servicePath+"/metadata.yaml", metadataString)

	err = common.StageAndCommitAll(params.ProjectName, "Added service: "+params.Service.ServiceName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
	}
	return service.NewPostProjectProjectNameStageStageNameServiceNoContent()
}

// PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc updates a service
func PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.PutProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.PutProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc deletes a service
func DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.DeleteProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.DeleteProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}
