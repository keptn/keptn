package handlers

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service"
	"gopkg.in/yaml.v3"
)

// PostProjectProjectNameStageStageNameServiceHandlerFunc creates a new service
func PostProjectProjectNameStageStageNameServiceHandlerFunc(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	servicePath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.Service.ServiceName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}

	if common.ServiceExists(params.ProjectName, params.StageName, params.Service.ServiceName) {
		return service.NewPostProjectProjectNameStageStageNameServiceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service already exists")})
	}

	logger.Debugf("Creating new service in: %s", servicePath)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}

	err = os.MkdirAll(servicePath, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not create service directory")})
	}

	newServiceMetadata := &common.ServiceMetadata{
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
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	servicePath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exists")})
	}

	logger.Debug(fmt.Sprintf("Deleting service %s of project %s in stage %s", params.ServiceName, params.ProjectName, params.StageName))
	err := common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch of project %s", params.StageName, params.ProjectName)
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}
	err = os.RemoveAll(servicePath)
	if err != nil {
		logger.Error(err.Error())
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete service directory: " + err.Error())})
	}

	_, err = common.StageAndCommitAll(params.ProjectName, "Deleted service: "+params.ServiceName)
	if err != nil {
		logger.Error(err.Error())
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete service directory: " + err.Error())})
	}

	return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameNoContent()
}
