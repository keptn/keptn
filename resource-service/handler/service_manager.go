package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/keptn/keptn/resource-service/common"
	"github.com/keptn/keptn/resource-service/common_models"
	kerrors "github.com/keptn/keptn/resource-service/errors"
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
	fileSystem       common.IFileSystem
	stageContext     IConfigurationContext
}

func NewServiceManager(git common.IGit, credentialReader common.CredentialReader, fileWriter common.IFileSystem, stageContext IConfigurationContext) *ServiceManager {
	serviceManager := &ServiceManager{
		git:              git,
		credentialReader: credentialReader,
		fileSystem:       fileWriter,
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

	_, resultErr := s.createService(gitContext, params.ServiceName, servicePath)

	// if there are conflicting changes first pull then try again
	if errors.Is(resultErr, kerrors.ErrNonFastForwardUpdate) || errors.Is(resultErr, kerrors.ErrForceNeeded) {
		_ = retry.Retry(func() error {
			err := s.git.Pull(*gitContext)
			if err != nil {
				resultErr = err
				// return nil at this point because retry does not make sense in that case
				return nil
			}

			_, err = s.createService(gitContext, params.ServiceName, servicePath)
			if err != nil {
				if errors.Is(err, kerrors.ErrNonFastForwardUpdate) || errors.Is(err, kerrors.ErrForceNeeded) {
					return err
				}
				resultErr = err
				// return nil at this point because retry does not make sense in that case
				return nil
			}
			resultErr = err
			return nil
		}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	}
	//return fmt.Errorf("could not initialize service %s: %w", params.ServiceName, err)
	return resultErr
}

func (s ServiceManager) DeleteService(params models.DeleteServiceParams) error {
	common.LockProject(params.ProjectName)
	defer common.UnlockProject(params.ProjectName)

	gitContext, servicePath, err := s.establishServiceContext(params.Project, params.Stage, params.Service)
	if err != nil {
		return err
	}

	_, resultErr := s.deleteService(gitContext, params.ServiceName, servicePath)
	// if there are conflicting changes first pull then try again
	if errors.Is(resultErr, kerrors.ErrNonFastForwardUpdate) || errors.Is(resultErr, kerrors.ErrForceNeeded) {
		_ = retry.Retry(func() error {
			err := s.git.Pull(*gitContext)
			if err != nil {
				resultErr = err
				// return nil at this point because retry does not make sense in that case
				return nil
			}

			_, err = s.deleteService(gitContext, params.ServiceName, servicePath)
			if err != nil {
				if errors.Is(err, kerrors.ErrNonFastForwardUpdate) || errors.Is(err, kerrors.ErrForceNeeded) {
					return err
				}
				resultErr = err
				// return nil at this point because retry does not make sense in that case
				return nil
			}
			resultErr = err
			return nil
		}, retry.NumberOfRetries(5), retry.DelayBetweenRetries(1*time.Second))
	}
	return resultErr
}

func (s ServiceManager) deleteService(gitContext *common_models.GitContext, serviceName, servicePath string) (string, error) {

	if !s.fileSystem.FileExists(servicePath) {
		return "", kerrors.ErrServiceNotFound
	}
	if err := s.fileSystem.DeleteFile(servicePath); err != nil {
		return "", err
	}

	return s.git.StageAndCommitAll(*gitContext, "Removed service: "+serviceName)
}

func (s ServiceManager) establishServiceContext(project models.Project, stage models.Stage, service models.Service) (*common_models.GitContext, string, error) {
	credentials, err := s.credentialReader.GetCredentials(project.ProjectName)
	if err != nil {
		return nil, "", fmt.Errorf(kerrors.ErrMsgCouldNotRetrieveCredentials, project.ProjectName, err)
	}

	gitContext := common_models.GitContext{
		Project:     project.ProjectName,
		Credentials: credentials,
	}

	if !s.git.ProjectExists(gitContext) {
		return nil, "", kerrors.ErrProjectNotFound
	}

	configPath, err := s.stageContext.Establish(common_models.ConfigurationContextParams{
		Project:                 project,
		Stage:                   &stage,
		Service:                 &service,
		GitContext:              gitContext,
		CheckConfigDirAvailable: false,
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not check out branch %s of project %s: %w", stage.StageName, project.ProjectName, err)
	}

	return &gitContext, configPath, nil
}

func (s ServiceManager) createService(gitContext *common_models.GitContext, serviceName, servicePath string) (string, error) {
	if s.fileSystem.FileExists(servicePath) {
		return "", kerrors.ErrServiceAlreadyExists
	}
	if err := s.fileSystem.MakeDir(servicePath); err != nil {
		return "", fmt.Errorf("could not create directory for service %s: %w", serviceName, err)
	}

	newServiceMetadata := &common.ServiceMetadata{
		ServiceName:       serviceName,
		CreationTimestamp: time.Now().UTC().String(),
	}

	metadataString, err := yaml.Marshal(newServiceMetadata)
	if err = s.fileSystem.WriteFile(servicePath+"/metadata.yaml", metadataString); err != nil {
		return "", fmt.Errorf("could not create metadata file for service %s: %w", serviceName, err)
	}
	return s.git.StageAndCommitAll(*gitContext, "Added service: "+serviceName)
}
