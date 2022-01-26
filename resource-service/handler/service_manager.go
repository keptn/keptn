package handler

import (
	"fmt"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	"github.com/keptn/keptn/resource-service/errors"
	"github.com/keptn/keptn/resource-service/models"
	"gopkg.in/yaml.v3"
	"time"
)

//IServiceManager provides an interface for stage CRUD operations
//go:generate moq -pkg handler_mock -skip-ensure -out ./fake/service_manager_mock.go . IServiceManager
type IServiceManager interface {
	CreateService(params models.CreateServiceParams) error
	DeleteService(params models.DeleteServiceParams) error
}

type ServiceManager struct {
	git              common.IGit
	credentialReader common.CredentialReader
	fileWriter       common.IFileSystem
	stageContext     IStageContext
}

func NewServiceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem, stageContext IStageContext) *ServiceManager {
	serviceManager := &ServiceManager{
		git:              git,
		credentialReader: credentialReader,
		fileWriter:       fileWriter,
		stageContext:     stageContext,
	}
	return serviceManager
}

func (s ServiceManager) CreateService(params models.CreateServiceParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, servicePath, err := s.establishServiceContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return err
	}

	if s.fileWriter.FileExists(servicePath) {
		return errors.ErrServiceAlreadyExists
	}
	if err := s.fileWriter.MakeDir(servicePath); err != nil {
		return fmt.Errorf("could not create directory for service %s: %w", params.ServiceName, err)
	}

	newServiceMetadata := &common.ServiceMetadata{
		ServiceName:       params.Service.ServiceName,
		CreationTimestamp: time.Now().UTC().String(),
	}

	metadataString, err := yaml.Marshal(newServiceMetadata)
	if err = s.fileWriter.WriteFile(servicePath+"/metadata.yaml", metadataString); err != nil {
		return fmt.Errorf("could not create metadata file for service %s: %w", params.ServiceName, err)
	}

	if _, err := s.git.StageAndCommitAll(*gitContext, "Added service: "+params.Service.ServiceName); err != nil {
		return fmt.Errorf("could not initialize service %s: %w", params.ServiceName, err)
	}

	return nil
}

func (s ServiceManager) DeleteService(params models.DeleteServiceParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, servicePath, err := s.establishServiceContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return err
	}

	if !s.fileWriter.FileExists(servicePath) {
		return errors.ErrServiceNotFound
	}
	if err := s.fileWriter.DeleteFile(servicePath); err != nil {
		return err
	}

	if _, err := s.git.StageAndCommitAll(*gitContext, "Removed service: "+params.Service.ServiceName); err != nil {
		return fmt.Errorf("could not remove service %s: %w", params.ServiceName, err)
	}

	return nil
}

func (s ServiceManager) establishServiceContext(project models.Project, stage models.Stage, service models.Service) (*common_models.GitContext, string, error) {
	credentials, err := s.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, "", fmt.Errorf(errors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !s.git.ProjectExists(gitContext) {
		return nil, "", errors.ErrProjectNotFound
	}

	configPath, err := s.stageContext.Establish(project, &stage, &service, gitContext)
	if err != nil {
		return nil, "", fmt.Errorf("could not check out branch %s of project %s: %w", stage.StageName, project.ProjectName, err)
	}

	return &gitContext, configPath, nil
}
