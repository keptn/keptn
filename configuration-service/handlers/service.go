package handlers

import (
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"os"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service"
	"gopkg.in/yaml.v3"
)

type serviceMetadata struct {
	ServiceName       string
	CreationTimestamp string
}

// PostProjectProjectNameStageStageNameServiceHandlerFunc creates a new service
func PostProjectProjectNameStageStageNameServiceHandlerFunc(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", common.ConfigurationServiceName)

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
		return service.NewPostProjectProjectNameStageStageNameServiceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
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

	common.StageAndCommitAll(params.ProjectName, "Added service: "+params.Service.ServiceName, true)

	return service.NewPostProjectProjectNameStageStageNameServiceNoContent()
}

// PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc updates a service
func PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.PutProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	return middleware.NotImplemented("operation service.PutProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc deletes a service
func DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(params service.DeleteProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", common.ConfigurationServiceName)

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	projectConfigPath := config.ConfigDir + "/" + params.ProjectName
	servicePath := projectConfigPath + "/" + params.ServiceName

	if !common.StageExists(params.ProjectName, params.StageName, false) {
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(400).WithPayload(&models.Error{Code: 400, Message: swag.String("Stage  " + params.StageName + " does not exist.")})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, false) {
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exists")})
	}

	logger.Debug(fmt.Sprintf("Deleting service %s of project %s in stage %s", params.ServiceName, params.ProjectName, params.StageName))
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s: %s", params.StageName, params.ProjectName, err.Error()))
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}

	err = os.RemoveAll(servicePath)
	if err != nil {
		logger.Error(err.Error())
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete service directory: " + err.Error())})
	}

	err = common.StageAndCommitAll(params.ProjectName, "Deleted service: "+params.ServiceName, true)
	if err != nil {
		logger.Error(err.Error())
		return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete service directory: " + err.Error())})
	}

	return service.NewDeleteProjectProjectNameStageStageNameServiceServiceNameNoContent()
}
