package handlers

import (
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/go-utils/pkg/utils"
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
	return middleware.NotImplemented("operation service.GetProjectProjectNameStageStageNameService has not yet been implemented")
}

// GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc get the specified service
func GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.GetProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.GetProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}

// PostProjectProjectNameStageStageNameServiceHandlerFunc creates a new service
func PostProjectProjectNameStageStageNameServiceHandlerFunc(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	logger := utils.NewLogger("", "", "configuration-service")
	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	servicePath := projectConfigPath + "/" + params.Service.ServiceName

	if !common.StageExists(params.ProjectName, params.StageName) {
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}
	logger.Debug("Creating new resource(s) in: " + projectConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName)
	if err != nil {
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

	common.StageAndCommitAll(params.ProjectName, "Added service: "+params.Service.ServiceName)
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
