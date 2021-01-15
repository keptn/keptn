package handlers

import (
	"encoding/base64"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/config"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_resource"
	"github.com/mholt/archiver"
	"github.com/otiai10/copy"
)

// GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc get list of resources for the service
func GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	return middleware.NotImplemented("operation service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResource has not yet been implemented")
}

// GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(
	params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName
	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not unescape resource name")})
	}
	resourcePath := serviceConfigPath + "/" + unescapedResourceName
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, *params.DisableUpstreamSync) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}

	logger.Debug("Checking out " + params.StageName + " branch")
	err = common.CheckoutBranch(params.ProjectName, params.StageName, *params.DisableUpstreamSync)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}

	// archive the Helm chart
	if strings.Contains(resourcePath, "helm") && strings.Contains(params.ResourceURI, ".tgz") {
		logger.Debug("Archive the Helm chart: " + params.ResourceURI)

		chartDir := strings.Replace(resourcePath, ".tgz", "", -1)
		if !common.FileExists(chartDir) {
			return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
				WithPayload(&models.Error{Code: 404, Message: swag.String("Service resource not found")})
		}
		if err := archiver.Archive([]string{chartDir}, resourcePath); err != nil {
			logger.Error(err.Error())
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String("Could not archive the Helm chart directory")})
		}
	}

	if !common.FileExists(resourcePath) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service resource not found")})
	}

	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read file")})
	}

	// remove Helch chart .tgz file
	if strings.Contains(resourcePath, "helm") && strings.HasSuffix(params.ResourceURI, ".tgz") {
		logger.Debug("Remove the Helm chart: " + params.ResourceURI)

		if err := os.Remove(resourcePath); err != nil {
			logger.Error(err.Error())
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete Helm chart package")})
		}
	}

	resourceContent := base64.StdEncoding.EncodeToString(dat)

	resource := &models.Resource{
		ResourceURI:     &params.ResourceURI,
		ResourceContent: resourceContent,
	}

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	resource.Metadata = metadata

	return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIOK().WithPayload(resource)
}

// DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc deletes the specified resource
func DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(
	params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return middleware.NotImplemented(
		"operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
}

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Project " + params.ProjectName + " does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName, false) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Stage " + params.StageName + " does not exist within project " + params.ProjectName)})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, false) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Service " + params.ServiceName + " does not exist within stage " + params.StageName + " of project " + params.ProjectName)})
	}
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI
		logger.Debug("Adding resource: " + filePath)

		common.WriteBase64EncodedFile(filePath, res.ResourceContent)

		if strings.Contains(filePath, "helm") && strings.HasSuffix(*res.ResourceURI, ".tgz") {
			if resp := untarHelm(res, logger, filePath); resp != nil {
				return resp
			}
		}
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Added resources", true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully added resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceCreated().
		WithPayload(metadata)
}

func untarHelm(res *models.Resource, logger *keptncommon.Logger, filePath string) middleware.Responder {
	// unarchive the Helm chart
	logger.Debug("Unarchive the Helm chart: " + *res.ResourceURI)
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tarGz := archiver.NewTarGz()
	tarGz.OverwriteExisting = true
	if err := tarGz.Unarchive(filePath, tmpDir); err != nil {
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not unarchive Helm chart")})
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not read unpacked files")})
	}

	if len(files) != 1 {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Unexpected amount of unpacked files")})
	}

	uri := *res.ResourceURI
	folderName := filepath.Join(tmpDir, uri[strings.LastIndex(uri, "/")+1:len(uri)-4])
	oldPath := filepath.Join(tmpDir, files[0].Name())
	if oldPath != folderName {
		if err := os.Rename(oldPath, folderName); err != nil {
			logger.Error(err.Error())
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String("Could not rename unpacked folder")})
		}
	}

	dir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Path of Helm chart is invalid")})
	}

	if err := copy.Copy(tmpDir, dir); err != nil {
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not copy folder")})
	}

	// remove Helm chart .tgz file
	logger.Debug("Remove the Helm chart: " + *res.ResourceURI)
	if err := os.Remove(filePath); err != nil {
		logger.Error(err.Error())
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not delete Helm chart package")})
	}
	return nil
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc updates a list of resources
func PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")
	if !common.ProjectExists(params.ProjectName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Project " + params.ProjectName + " does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName, false) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Stage " + params.StageName + " does not exist within project " + params.ProjectName)})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, false) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Service " + params.ServiceName + " does not exist within stage " + params.StageName + " of project " + params.ProjectName)})
	}
	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	logger.Debug("Updating resource(s) in: " + serviceConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI
		logger.Debug("Updating resource: " + filePath)
		common.WriteBase64EncodedFile(filePath, res.ResourceContent)
		if strings.Contains(filePath, "helm") && strings.HasSuffix(*res.ResourceURI, ".tgz") {
			if resp := untarHelm(res, logger, filePath); resp != nil {
				return resp
			}
		}
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resources", true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	metadata.Branch = defaultBranch
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(metadata)
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc updates a specified resource
func PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(
	params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	logger := keptncommon.NewLogger("", "", "configuration-service")

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName, false) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Service does not exist")})
	}

	serviceConfigPath := config.ConfigDir + "/" + params.ProjectName + "/" + params.ServiceName

	logger.Debug("updating resource(s) in: " + serviceConfigPath + " in stage " + params.StageName)
	logger.Debug("Checking out branch: " + params.StageName)
	err := common.CheckoutBranch(params.ProjectName, params.StageName, false)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not check out %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not check out branch")})
	}

	filePath := serviceConfigPath + "/" + params.ResourceURI
	common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI, true)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName))
		logger.Error(err.Error())
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resource: " + params.ResourceURI)

	metadata := common.GetResourceMetadata(params.ProjectName)
	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not determine default branch of project %s: %s", params.ProjectName, err.Error()))
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not check out branch")})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	metadata.Branch = defaultBranch
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated().
		WithPayload(metadata)
}
