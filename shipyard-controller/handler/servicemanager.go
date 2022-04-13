package handler

import (
	"errors"
	"fmt"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
)

const (
	serviceNameMaxLen = 53
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/servicemanager.go . IServiceManager
type IServiceManager interface {
	CreateService(projectName string, params *models.CreateServiceParams) error
	DeleteService(projectName, serviceName string) error
	GetService(projectName, stageName, serviceName string) (*models.ExpandedService, error)
	GetAllServices(projectName, stageName string) ([]*models.ExpandedService, error)
}

type serviceManager struct {
	projectMVRepo      db.ProjectMVRepo
	configurationStore common.ConfigurationStore
	uniformRepo        db.UniformRepo
}

func NewServiceManager(servicesDBOperations db.ProjectMVRepo, configurationStore common.ConfigurationStore, uniformRepo db.UniformRepo) *serviceManager {
	return &serviceManager{
		projectMVRepo:      servicesDBOperations,
		configurationStore: configurationStore,
		uniformRepo:        uniformRepo,
	}
}

func (sm *serviceManager) GetAllStages(projectName string) ([]*models.ExpandedStage, error) {
	project, err := sm.projectMVRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	return project.Stages, nil

}

func (sm *serviceManager) GetService(projectName, stageName, serviceName string) (*models.ExpandedService, error) {
	project, err := sm.projectMVRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			for _, svc := range stg.Services {
				if svc.ServiceName == serviceName {
					return svc, nil
				}
			}
			return nil, ErrServiceNotFound
		}
	}
	return nil, ErrStageNotFound
}

func (sm *serviceManager) GetAllServices(projectName, stageName string) ([]*models.ExpandedService, error) {
	project, err := sm.projectMVRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	for _, stg := range project.Stages {
		if stg.StageName == stageName {
			return stg.Services, nil
		}
	}
	return nil, ErrStageNotFound
}

func (sm *serviceManager) CreateService(projectName string, params *models.CreateServiceParams) error {
	log.Infof("Received request to create service %s in project %s", *params.ServiceName, projectName)

	// check service name length
	log.Infof("Validating service %s", *params.ServiceName)
	if err := validateServiceName(*params.ServiceName); err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not create service %s for project %s: %s", *params.ServiceName, projectName, err.Error()))
	}

	stages, err := sm.GetAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not get stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		log.Infof("Checking if service %s already exists in project %s", *params.ServiceName, projectName)
		// check if the service exists, do not continue if yes
		service, _ := sm.GetService(projectName, stage.StageName, *params.ServiceName)
		if service != nil {
			log.Infof("Service %s already exists in project %s", *params.ServiceName, projectName)
			return ErrServiceAlreadyExists
		}

		log.Infof("Creating service %s in project %s", *params.ServiceName, projectName)

		if err := sm.configurationStore.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		if err := sm.projectMVRepo.CreateService(projectName, stage.StageName, *params.ServiceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not create service %s in stage %s of project %s: %s", *params.ServiceName, stage.StageName, projectName, err.Error()))
		}
		log.Infof("Created service %s in stage %s of project %s", *params.ServiceName, stage.StageName, projectName)
	}

	return nil
}

func (sm *serviceManager) DeleteService(projectName, serviceName string) error {
	log.Infof("Deleting service %s from project %s", serviceName, projectName)

	stages, err := sm.GetAllStages(projectName)
	if err != nil {
		return sm.logAndReturnError(fmt.Sprintf("could not retrieve stages of project %s: %s", projectName, err.Error()))
	}

	for _, stage := range stages {
		log.Infof("Deleting service %s from stage %s", serviceName, stage.StageName)
		if err := sm.configurationStore.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			// If we get a ErrServiceNotFound, we can proceed with deleting the service from the db.
			// For other types of errors (e.g. due to a temporary upstream repo connection issue), we return without deleting it from the db.
			// Otherwise, it could be that the service directory is still present in the configuration service, but gone from the db, which means we cannot
			// retry the deletion via the bridge (since the service won't show up anymore), and recreating the service will fail because we'll get a 409 from
			// the configuration service
			if errors.Is(err, ErrServiceNotFound) {
				log.Infof("Service %s has already been deleted from stage %s", serviceName, stage.StageName)
			} else {
				return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
			}
		}
		if err := sm.projectMVRepo.DeleteService(projectName, stage.StageName, serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
		if err := sm.uniformRepo.DeleteServiceFromSubscriptions(serviceName); err != nil {
			return sm.logAndReturnError(fmt.Sprintf("could not delete service %s from stage %s: %s", serviceName, stage.StageName, err.Error()))
		}
	}
	log.Infof("deleted service %s from project %s", serviceName, projectName)

	return nil
}

// validateServiceName validates that the service name is less than 43 characters (this is a requirement of helm-service)
func validateServiceName(serviceName string) error {
	// helm-service creates release names that have the service name and the string -generated in them
	// this means that we need to ensure in here that service names are not too long
	allowedLength := serviceNameMaxLen - len("generated") // = 43

	if len(serviceName) > allowedLength {
		return fmt.Errorf("service name needs to be less than %d characters", allowedLength)
	}

	return nil
}

func (sm *serviceManager) logAndReturnError(msg string) error {
	log.Error(msg)
	return errors.New(msg)
}
