package handlers

import (
	"encoding/base64"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/keptn/keptn/configuration-service/restapi/operations/stage_resource"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_resource"
	archive "github.com/mholt/archiver/v3"
	"github.com/otiai10/copy"
)

// GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc get list of resources for the service
func GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceNotFound().
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}

	err := common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}

	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)
	result := common.GetPaginatedResources(serviceConfigPath, params.PageSize, params.NextPageKey)
	return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceOK().WithPayload(result)
}

// GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc gets the specified resource
func GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(
	params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)
	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not unescape resource name")})
	}
	resourcePath := serviceConfigPath + "/" + unescapedResourceName
	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}

	err = common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}

	// archive the Helm chart
	if strings.Contains(resourcePath, "helm") && strings.Contains(params.ResourceURI, ".tgz") {
		logger.Debug("Archive the Helm chart: " + params.ResourceURI)

		chartDir := strings.Replace(resourcePath, ".tgz", "", -1)
		if !common.FileExists(chartDir) {
			return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
				WithPayload(&models.Error{Code: 404, Message: swag.String("Service resource not found")})
		}
		if err := archive.Archive([]string{chartDir}, resourcePath); err != nil {
			logger.Error(err.Error())
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String("Could not archive the Helm chart directory")})
		}
	}

	if !common.FileExists(resourcePath) {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURINotFound().
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service resource not found")})
	}

	resourcePath = filepath.Clean(resourcePath)
	dat, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		logger.Error(err.Error())
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not read file")})
	}

	// remove Helm chart .tgz file
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

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewDeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(404).
			WithPayload(&models.Error{Code: 404, Message: swag.String("Service not found")})
	}

	err := common.PullUpstream(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not check out %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}

	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)
	unescapedResourceName, err := url.QueryUnescape(params.ResourceURI)
	if err != nil {
		return service_resource.NewGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).
			WithPayload(&models.Error{Code: 500, Message: swag.String("Could not unescape resource name")})
	}
	serviceResourcePath := serviceConfigPath + "/" + unescapedResourceName

	err = common.DeleteFile(serviceResourcePath)
	if err != nil {
		logger.Error(err.Error())
		return service_resource.NewDeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String("Could not delete file")})
	}

	logger.Debug("Staging Changes")
	err = common.StageAndCommitAll(params.ProjectName, "Updated resource: "+unescapedResourceName, true)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch for project %s", params.StageName, params.ProjectName)
		return service_resource.NewDeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debugf("Successfully updated resource: %s", unescapedResourceName)

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	return stage_resource.NewPutProjectProjectNameStageStageNameResourceResourceURICreated().WithPayload(metadata)
}

func AddUpdateHelmResources(resource string, actionType string, filePath string, projectName string) error {
	logger.Debugf("%sing resource: %s", actionType, filePath)
	err := common.WriteBase64EncodedFile(filePath, resource)
	if err != nil {
		logger.WithError(err).Errorf("Could not %s resource %s to project %s", actionType, filePath, projectName)
		return err
	}

	return nil
}

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc creates a new resource
func PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	if !common.ProjectExists(params.ProjectName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Project " + params.ProjectName + " does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Stage " + params.StageName + " does not exist within project " + params.ProjectName)})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Service " + params.ServiceName + " does not exist within stage " + params.StageName + " of project " + params.ProjectName)})
	}
	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)

	logger.Debug("Creating new resource(s) in: " + serviceConfigPath + " in stage " + params.StageName)

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI

		if err := AddUpdateHelmResources(res.ResourceContent, "add", filePath, params.ProjectName); err != nil {
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String(common.CannotAddResourceErrorMsg)})
		}

		if strings.Contains(filePath, "helm") && strings.HasSuffix(*res.ResourceURI, ".tgz") {
			if err := extractHelmArchiveResource(params.ProjectName, params.StageName, filePath, res); err != nil {
				logger.Errorf("Could not extract helm archive: %v", err)
				return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
					WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not extract helm archive")})
			}
		}
	}

	logger.Debug("Staging Changes")
	err := common.StageAndCommitAll(params.ProjectName, "Added resources", true)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully added resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	metadata.Branch = params.StageName
	return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceCreated().
		WithPayload(metadata)
}

func untarHelm(res *models.Resource, filePath string) error {
	// unarchive the Helm chart
	logger.Debug("Unarchive the Helm chart: " + *res.ResourceURI)
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			logger.WithError(err).Errorf("Could not remove directory %s", tmpDir)
		}
	}()

	tarGz := archive.NewTarGz()
	tarGz.OverwriteExisting = true
	if err := tarGz.Unarchive(filePath, tmpDir); err != nil {
		return fmt.Errorf("could not unarchive Helm chart: %w", err)
	}

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("could not read unpacked files: %w", err)
	}

	if len(files) != 1 {
		return errors.New("unexpected amount of unpacked files")
	}

	uri := *res.ResourceURI
	folderName := filepath.Join(tmpDir, uri[strings.LastIndex(uri, "/")+1:len(uri)-4])
	oldPath := filepath.Join(tmpDir, files[0].Name())
	if oldPath != folderName {
		if err := os.Rename(oldPath, folderName); err != nil {
			return fmt.Errorf("could not rename unpacked folder: %w", err)
		}
	}

	dir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return fmt.Errorf("patch of helm chart is invalid: %w", err)
	}

	if err := copy.Copy(tmpDir, dir); err != nil {
		return fmt.Errorf("could not copy folder: %w", err)
	}

	// remove Helm chart .tgz file
	logger.Debug("Remove the Helm chart: " + *res.ResourceURI)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("could not delete helm chart package: %w", err)
	}
	return nil
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc updates a list of resources
func PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(
	params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {

	if !common.ProjectExists(params.ProjectName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Project " + params.ProjectName + " does not exist")})
	}

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.StageExists(params.ProjectName, params.StageName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Stage " + params.StageName + " does not exist within project " + params.ProjectName)})
	}

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Service " + params.ServiceName + " does not exist within stage " + params.StageName + " of project " + params.ProjectName)})
	}
	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)

	for _, res := range params.Resources.Resources {
		filePath := serviceConfigPath + "/" + *res.ResourceURI

		if err := AddUpdateHelmResources(res.ResourceContent, "update", filePath, params.ProjectName); err != nil {
			return service_resource.NewPostProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
				WithPayload(&models.Error{Code: 400, Message: swag.String(common.CannotUpdateResourceErrorMsg)})
		}

		if strings.Contains(filePath, "helm") && strings.HasSuffix(*res.ResourceURI, ".tgz") {
			if err := extractHelmArchiveResource(params.ProjectName, params.StageName, filePath, res); err != nil {
				logger.Errorf("Could not extract helm archive: %v", err)
				return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
					WithPayload(&models.Error{Code: http.StatusBadRequest, Message: swag.String("Could not extract helm archive")})
			}
		}
	}

	logger.Debug("Staging Changes")
	err := common.StageAndCommitAll(params.ProjectName, "Updated resources", true)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resources")

	metadata := common.GetResourceMetadata(params.ProjectName)
	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not determine default branch of project %s", params.ProjectName)
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	metadata.Branch = defaultBranch
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceCreated().WithPayload(metadata)
}

func extractHelmArchiveResource(project, stage, filePath string, res *models.Resource) error {
	rollbackFunc := func() {
		// restore previously deleted helm/resourceURI folder using git reset
		if err := common.Reset(project); err != nil {
			logger.Errorf("Could not reset current branch '%s': %v", stage, err)
		}
		if err := common.DeleteFile(filePath); err != nil {
			logger.Errorf("Could not delete file %s: %v", filePath, err)
		}
	}
	// remove previous helm/resourceURI folder
	targetFolderPath := strings.TrimSuffix(filePath, ".tgz")
	if err := os.RemoveAll(targetFolderPath); err != nil {
		rollbackFunc()
		return fmt.Errorf("could not delete existing folder %s, %v", targetFolderPath, err)
	}
	if err := untarHelm(res, filePath); err != nil {
		rollbackFunc()
		return err
	}
	return nil
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc updates a specified resource
func PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(
	params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {

	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	if !common.ServiceExists(params.ProjectName, params.StageName, params.ServiceName) {
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String(common.ServiceDoesNotExistErrorMsg)})
	}

	serviceConfigPath := common.GetServiceConfigPath(params.ProjectName, params.StageName, params.ServiceName)

	filePath := serviceConfigPath + "/" + params.ResourceURI
	common.WriteBase64EncodedFile(filePath, params.Resource.ResourceContent)

	logger.Debug("Staging Changes")
	err := common.StageAndCommitAll(params.ProjectName, "Updated resource: "+params.ResourceURI, true)
	if err != nil {
		logger.WithError(err).Errorf("Could not commit to %s branch of project %s", params.StageName, params.ProjectName)
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest().
			WithPayload(&models.Error{Code: 400, Message: swag.String("Could not commit changes")})
	}
	logger.Debug("Successfully updated resource: " + params.ResourceURI)

	metadata := common.GetResourceMetadata(params.ProjectName)
	defaultBranch, err := common.GetDefaultBranch(params.ProjectName)
	if err != nil {
		logger.WithError(err).Errorf("Could not determine default branch of project %s", params.ProjectName)
		return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(common.CannotCheckOutBranchErrorMsg)})
	}
	if defaultBranch == "" {
		defaultBranch = "master"
	}
	metadata.Branch = defaultBranch
	return service_resource.NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated().
		WithPayload(metadata)
}
