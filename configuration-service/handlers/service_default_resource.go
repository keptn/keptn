package handlers

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_default_resource"
)

// GetProjectProjectNameServiceServiceNameResourceHandlerFunc get list of default resources for the service
func GetProjectProjectNameServiceServiceNameResourceHandlerFunc(params service_default_resource.GetProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation service_default_resource.GetProjectProjectNameServiceServiceNameResource has not yet been implemented")
}

// GetProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc gets a specified default resource
func GetProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(params service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented("operation service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURI has not yet been implemented")
}

// PostProjectProjectNameServiceServiceNameResourceHandlerFunc creates a list of new default resources
func PostProjectProjectNameServiceServiceNameResourceHandlerFunc(params service_default_resource.PostProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	branches, err := common.GetBranches(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not get stages for project")})
	}

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	for _, branch := range branches {
		if branch == defaultBranch {
			continue
		}
		if !common.ServiceExists(params.ProjectName, branch, params.ServiceName, false) {
			return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
		}
		serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

		logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + branch)
		logger.Debug("Checking out branch: " + branch)
		err := common.CheckoutBranch(params.ProjectName, branch, false)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
		}

		for _, res := range params.Resources.Resources {
			filePath := serviceConfigPath + "/" + *res.ResourceURI
			logger.Debug("Adding resource: " + filePath)
			common.WriteBase64EncodedFile(filePath, res.ResourceContent)
		}

		logger.Debug("Staging Changes")
		err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceBadRequest().WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
		}
		logger.Debug("Successfully added resources")
	}

	return service_default_resource.NewPostProjectProjectNameServiceServiceNameResourceCreated()
}

// PutProjectProjectNameServiceServiceNameResourceHandlerFunc updates a list of default resources
func PutProjectProjectNameServiceServiceNameResourceHandlerFunc(params service_default_resource.PutProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	branches, err := common.GetBranches(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not get stages for project")})
	}

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	for _, branch := range branches {
		if branch == defaultBranch {
			continue
		}
		if !common.ServiceExists(params.ProjectName, branch, params.ServiceName, false) {
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
		}
		serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

		logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + branch)
		logger.Debug("Checking out branch: " + branch)
		err := common.CheckoutBranch(params.ProjectName, branch, false)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
		}

		for _, res := range params.Resources.Resources {
			filePath := serviceConfigPath + "/" + *res.ResourceURI
			logger.Debug("Adding resource: " + filePath)
			common.WriteBase64EncodedFile(filePath, res.ResourceContent)
		}

		logger.Debug("Staging Changes")
		err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
		}
		logger.Debug("Successfully added resources")
	}

	return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceCreated()
}

// PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc updates the specified resource for the service
func PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(params service_default_resource.PutProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	branches, err := common.GetBranches(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not get stages for project")})
	}

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	for _, branch := range branches {
		if branch == defaultBranch {
			continue
		}
		if !common.ServiceExists(params.ProjectName, branch, params.ServiceName, false) {
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
		}
		serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

		logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + branch)
		logger.Debug("Checking out branch: " + branch)
		err := common.CheckoutBranch(params.ProjectName, branch, false)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
		}

		filePath := serviceConfigPath + "/" + params.ResourceURI
		common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

		logger.Debug("Staging Changes")
		err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
		}
		logger.Debug("Successfully added resources")
	}

	return service_default_resource.NewPutProjectProjectNameServiceServiceNameResourceCreated()
}

// DeleteProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource from the service
func DeleteProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(params service_default_resource.DeleteProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Project does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	branches, err := common.GetBranches(params.ProjectName)
	if err != nil {
		logger.Error(err.Error())
		return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not get stages for project")})
	}

	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	for _, branch := range branches {
		if branch == defaultBranch {
			continue
		}
		if !common.ServiceExists(params.ProjectName, branch, params.ServiceName, false) {
			return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(404).WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
		}
		serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

		logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + branch)
		logger.Debug("Checking out branch: " + branch)
		err := common.CheckoutBranch(params.ProjectName, branch, false)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
		}

		filePath := serviceConfigPath + "/" + params.ResourceURI
		err = common.DeleteFile(filePath)
		if err != nil {
			logger.Error(err.Error())
		}

		logger.Debug("Staging Changes")
		err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", branch, params.ProjectName))
			logger.Error(err.Error())
			return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
		}
		logger.Debug("Successfully deletes resource")
	}

	return service_default_resource.NewDeleteProjectProjectNameServiceServiceNameResourceResourceURINoContent()
}
