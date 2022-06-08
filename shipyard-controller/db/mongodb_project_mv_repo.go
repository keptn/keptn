package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	goutilsmodels "github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var instance *MongoDBProjectMVRepo

// EventsRetriever defines the interface for fetching events from the data store
type EventsRetriever interface {
	GetEvents(filter *goutils.EventFilter) ([]*goutilsmodels.KeptnContextExtendedCE, *goutilsmodels.Error)
}

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/projectmvrepo_mock.go . ProjectMVRepo
type ProjectMVRepo interface {
	CreateProject(prj *apimodels.ExpandedProject) error
	UpdateShipyard(projectName string, shipyardContent string) error
	UpdateProject(prj *apimodels.ExpandedProject) error
	UpdateUpstreamInfo(projectName string, uri, user string) error
	UpdatedShipyard(projectName string, shipyard string) error
	DeleteUpstreamInfo(projectName string) error
	GetProjects() ([]*apimodels.ExpandedProject, error)
	GetProject(projectName string) (*apimodels.ExpandedProject, error)
	DeleteProject(projectName string) error
	CreateStage(project string, stage string) error
	DeleteStage(project string, stage string) error
	CreateService(project string, stage string, service string) error
	GetService(projectName, stageName, serviceName string) (*apimodels.ExpandedService, error)
	DeleteService(project string, stage string, service string) error
	UpdateEventOfService(e apimodels.KeptnContextExtendedCE) error
	CreateRemediation(project, stage, service string, remediation *apimodels.Remediation) error
	CloseOpenRemediations(project, stage, service, keptnContext string) error
	OnSequenceTaskStarted(event apimodels.KeptnContextExtendedCE)
	OnSequenceTaskFinished(event apimodels.KeptnContextExtendedCE)
}

type MongoDBProjectMVRepo struct {
	projectRepo ProjectRepo
	eventRepo   EventRepo
}

func NewProjectMVRepo(projectRepo ProjectRepo, eventRepo EventRepo) *MongoDBProjectMVRepo {
	if instance == nil {
		instance = &MongoDBProjectMVRepo{
			projectRepo: projectRepo,
			eventRepo:   eventRepo,
		}
	}
	return instance
}

// CreateProject creates a project
func (mv *MongoDBProjectMVRepo) CreateProject(prj *apimodels.ExpandedProject) error {
	existingProject, err := mv.GetProject(prj.ProjectName)
	if existingProject != nil {
		return nil
	}

	updatedProject, err := generateStageInfo(*prj)
	if err != nil {
		return err
	}
	return mv.createProject(&updatedProject)
}

// UpdatedShipyard updates the shipyard of a project
func (mv *MongoDBProjectMVRepo) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent
	if err := setShipyardVersion(existingProject); err != nil {
		log.Errorf("could not update shipyard version of project %s: %s"+projectName, err.Error())
	}

	updatedProject, err := generateStageInfo(*existingProject)
	if err != nil {
		log.Errorf("could not update stage information of project %s: %s"+projectName, err.Error())
	}

	return mv.projectRepo.UpdateProject(&updatedProject)
}

func generateStageInfo(project apimodels.ExpandedProject) (apimodels.ExpandedProject, error) {
	shipyard, err := common.UnmarshalShipyard(project.Shipyard)
	if err != nil {
		return project, err
	}

	for index := range project.Stages {
		project.Stages[index].ParentStages = getParentStages(project.Stages[index].StageName, shipyard)
	}
	return project, nil
}

func getParentStages(stageName string, shipyard *keptnv2.Shipyard) []string {
	parentStages := []string{}
	for _, stage := range shipyard.Spec.Stages {
		if stage.Name != stageName {
			continue
		}

		for _, sequence := range stage.Sequences {
			for _, trigger := range sequence.TriggeredOn {
				// trigger events have the format <stage>.<sequenceName>.finished
				split := strings.Split(trigger.Event, ".")
				if len(split) == 3 {
					newParentStageName := split[0]
					if newParentStageName == stageName {
						// do not add the stage itself as a parent
						continue
					}
					parentStageAvailable := false
					for _, parentStage := range parentStages {
						if parentStage == newParentStageName {
							parentStageAvailable = true
							break
						}
					}
					if !parentStageAvailable {
						parentStages = append(parentStages, newParentStageName)
					}
				}
			}
		}
		break
	}
	return parentStages
}

// UpdateProject updates a project
func (mv *MongoDBProjectMVRepo) UpdateProject(prj *apimodels.ExpandedProject) error {
	return mv.projectRepo.UpdateProject(prj)
}

func setShipyardVersion(existingProject *apimodels.ExpandedProject) error {
	const previousShipyardVersion = "spec.keptn.sh/0.1.7"
	if existingProject.Shipyard == "" {
		// if the field is not set, it can only be 0.1.7, since in Keptn 0.8 we ensure that the shipyard content is always included in the materialized view
		existingProject.ShipyardVersion = previousShipyardVersion
		return nil
	}
	shipyard := &keptnv2.Shipyard{}
	if err := yaml.Unmarshal([]byte(existingProject.Shipyard), shipyard); err != nil {
		return errors.New("could not parse shipyard file content to shipyard struct: " + err.Error())
	}
	if shipyard.ApiVersion != "" {
		existingProject.ShipyardVersion = shipyard.ApiVersion
	} else {
		existingProject.ShipyardVersion = previousShipyardVersion
	}
	return nil
}

// UpdateUpstreamInfo updates the Upstream Repository URL and git user of a project
func (mv *MongoDBProjectMVRepo) UpdateUpstreamInfo(projectName string, uri, user string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	if existingProject.GitCredentials.RemoteURL != uri || existingProject.GitCredentials.User != user {
		existingProject.GitCredentials.RemoteURL = uri
		existingProject.GitCredentials.User = user
		if err := mv.projectRepo.UpdateProject(existingProject); err != nil {
			log.Errorf("could not update upstream credentials of project %s: %s", projectName, err.Error())
			return err
		}
	}
	return nil
}

func (mv *MongoDBProjectMVRepo) UpdatedShipyard(projectName string, shipyard string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}

	if existingProject.Shipyard != shipyard {
		existingProject.Shipyard = shipyard
		err = mv.projectRepo.UpdateProject(existingProject)
		if err != nil {
			return err
		}
		if err != nil {
			log.Errorf("could not update shipyard of project %s: %s", projectName, err.Error())
			return nil
		}
	}
	return nil

}

// DeleteUpstreamInfo deletes the Upstream Repository URL and git user of a project
func (mv *MongoDBProjectMVRepo) DeleteUpstreamInfo(projectName string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	existingProject.GitCredentials.User = ""
	existingProject.GitCredentials.RemoteURL = ""
	if err := mv.projectRepo.UpdateProject(existingProject); err != nil {
		log.Errorf("could not delete upstream credentials of project %s: %s", projectName, err.Error())
		return err
	}
	return nil
}

// GetProjects returns all projects
func (mv *MongoDBProjectMVRepo) GetProjects() ([]*apimodels.ExpandedProject, error) {
	projects, err := mv.projectRepo.GetProjects()
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		if err := setShipyardVersion(project); err != nil {
			// log the error but continue
			log.Errorf("could not set shipyard version of project %s: %s", project.ProjectName, err.Error())
		}
	}
	return projects, nil
}

// GetProject returns a project by its name
func (mv *MongoDBProjectMVRepo) GetProject(projectName string) (*apimodels.ExpandedProject, error) {
	project, err := mv.projectRepo.GetProject(projectName)
	if err != nil {
		return nil, err
	}
	if project != nil {
		if err := setShipyardVersion(project); err != nil {
			// log the error but continue
			log.Errorf("could not set shipyard version of project %s: %s", project.ProjectName, err.Error())
		}
	}
	return project, nil
}

// DeleteProject deletes a project
func (mv *MongoDBProjectMVRepo) DeleteProject(projectName string) error {
	return mv.projectRepo.DeleteProject(projectName)
}

// CreateStage creates a stage
func (mv *MongoDBProjectMVRepo) CreateStage(project string, stage string) error {
	log.Infof("Adding stage %s to project %s ", stage, project)
	prj, err := mv.GetProject(project)

	if err != nil {
		log.Errorf("Could not add stage %s to project %s : %s\n", stage, project, err.Error())
		return err
	}

	stageAlreadyExists := false
	for _, stg := range prj.Stages {
		if stg.StageName == stage {
			stageAlreadyExists = true
			break
		}
	}

	if stageAlreadyExists {
		log.Infof("Stage %s already exists in project %s", stage, project)
		return nil
	}

	prj.Stages = append(prj.Stages, &apimodels.ExpandedStage{
		Services:  []*apimodels.ExpandedService{},
		StageName: stage,
	})

	err = mv.projectRepo.UpdateProject(prj)
	if err != nil {
		return err
	}

	log.Infof("Added stage %s to project %s", stage, project)
	return nil
}

func (mv *MongoDBProjectMVRepo) createProject(project *apimodels.ExpandedProject) error {

	err := mv.projectRepo.CreateProject(project)
	if err != nil {
		log.Errorf("Could not create project %s: %s", project.ProjectName, err.Error())
		return err
	}
	return nil
}

// DeleteStage deletes a stage
func (mv *MongoDBProjectMVRepo) DeleteStage(project string, stage string) error {
	log.Infof("Deleting stage %s from project %s", stage, project)
	prj, err := mv.GetProject(project)

	if err != nil {
		return fmt.Errorf("could not delete stage %s from project %s: %w", stage, project, err)
	}

	stageIndex := -1

	for idx, stg := range prj.Stages {
		if stg.StageName == stage {
			stageIndex = idx
			break
		}
	}
	if stageIndex < 0 {
		return nil
	}

	copy(prj.Stages[stageIndex:], prj.Stages[stageIndex+1:])
	prj.Stages[len(prj.Stages)-1] = nil
	prj.Stages = prj.Stages[:len(prj.Stages)-1]

	err = mv.projectRepo.UpdateProject(prj)
	return nil
}

// CreateService creates a service
func (mv *MongoDBProjectMVRepo) CreateService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not add service %s to stage %s in project %s. Could not load project: %s", service, stage, project, err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			for _, svc := range stg.Services {
				if svc.ServiceName == service {
					log.Infof("Service %s already exists in stage %s in project %s", service, stage, project)
					return nil
				}
			}
			stg.Services = append(stg.Services, &apimodels.ExpandedService{
				CreationDate:  strconv.FormatInt(time.Now().UnixNano(), 10),
				DeployedImage: "",
				ServiceName:   service,
			})
			log.Infof("Adding %s to stage %s in project %s in database", service, stage, project)
			err := mv.projectRepo.UpdateProject(existingProject)
			if err != nil {
				log.Errorf("Could not add service %s to stage %s in project %s. Could not update project: %s", service, stage, project, err.Error())
				return err
			}
			log.Infof("Service %s has been added to stage %s in project %s", service, stage, project)
			break
		}
	}
	return nil
}

func (mv *MongoDBProjectMVRepo) GetService(projectName, stageName, serviceName string) (*apimodels.ExpandedService, error) {
	project, err := mv.GetProject(projectName)
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

// DeleteService deletes a service
func (mv *MongoDBProjectMVRepo) DeleteService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not delete service %s from stage %s in project %s. Could not load project: %s", service, stage, project, err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			serviceIndex := -1
			for idx, svc := range stg.Services {
				if svc.ServiceName == service {
					serviceIndex = idx
				}
			}
			if serviceIndex < 0 {
				log.Infof("Could not delete service %s from stage %s in project %s. Service not found in database", service, stage, project)
				return nil
			}
			copy(stg.Services[serviceIndex:], stg.Services[serviceIndex+1:])
			stg.Services[len(stg.Services)-1] = nil
			stg.Services = stg.Services[:len(stg.Services)-1]
			break
		}
	}
	err = mv.projectRepo.UpdateProject(existingProject)
	if err != nil {
		log.Errorf("Could not delete service %s from stage %s in project %s: %s", service, stage, project, err.Error())
		return err
	}
	log.Infof("Deleted service %s from stage %s in project %s", service, stage, project)
	return nil
}

// UpdateEventOfService updates a service event
func (mv *MongoDBProjectMVRepo) UpdateEventOfService(e apimodels.KeptnContextExtendedCE) error {
	if e.Type == nil {
		return errors.New("event type must be set")
	}
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(e.Data, eventData)
	if err != nil {
		log.Errorf("Could not parse event data: %s", err.Error())
		return err
	}

	existingProject, err := mv.GetProject(eventData.Project)
	if err != nil {
		log.Errorf("Could not update service %s in stage %s in project %s. Could not load project: %s", eventData.Service, eventData.Stage, eventData.Project, err.Error())
		return err
	} else if existingProject == nil {
		log.Errorf("Could not update service %s in stage %s in project %s: Project not found.", eventData.Service, eventData.Stage, eventData.Project)
		return ErrProjectNotFound
	}

	contextInfo := &apimodels.EventContextInfo{
		EventID:      e.ID,
		KeptnContext: e.Shkeptncontext,
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	err = updateServiceInStage(existingProject, eventData.Stage, eventData.Service, func(service *apimodels.ExpandedService) error {
		if service.LastEventTypes == nil {
			service.LastEventTypes = map[string]apimodels.EventContextInfo{}
		}
		service.LastEventTypes[*e.Type] = *contextInfo

		// for events of type "deployment.finished", find the correlating
		// "deployment.triggered" event to update the deployed image name
		if *e.Type == keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName) {

			events, errObj := mv.getAllDeploymentTriggeredEvents(eventData, e.Triggeredid, e.Shkeptncontext)
			if errObj != nil {
				return err
			}
			if events == nil || len(events) == 0 {
				return errors.New("No deployment.triggered events could be found for keptn context " + e.Shkeptncontext)
			}

			triggeredData := &keptnv2.DeploymentTriggeredEventData{}
			err := keptnv2.Decode(events[0].Data, triggeredData)
			if err != nil {
				return errors.New("unable to decode deployment.triggered event data: " + err.Error())
			}

			deployedImage := common.ExtractImageOfDeploymentEvent(*triggeredData)
			service.DeployedImage = deployedImage
		}
		return nil
	})

	if err != nil {
		log.Errorf("Could not update image of service %s: %s", eventData.Service, err.Error())
		return err
	}
	err = mv.projectRepo.UpdateProject(existingProject)
	if err != nil {
		log.Errorf("Could not update project %s: %s", eventData.Project, err.Error())
		return err
	}
	return nil
}

// CreateRemediation creates a remediation action
func (mv *MongoDBProjectMVRepo) CreateRemediation(project, stage, service string, remediation *apimodels.Remediation) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not create remediation for service %s in stage %s in project%s. Could not load project: %s", service, stage, project, err.Error())
		return ErrProjectNotFound
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *apimodels.ExpandedService) error {
		if service.OpenRemediations == nil {
			service.OpenRemediations = []*apimodels.Remediation{}
		}
		service.OpenRemediations = append(service.OpenRemediations, remediation)
		return nil
	})
	return mv.projectRepo.UpdateProject(existingProject)
}

// CloseOpenRemediations closes a open remediation actions for a given keptnContext
func (mv *MongoDBProjectMVRepo) CloseOpenRemediations(project, stage, service, keptnContext string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not close remediation for service %s in stage %s in project %s. Could not load project: %s", service, stage, project, err.Error())
		return ErrProjectNotFound
	}
	if keptnContext == "" {
		log.Warn("No keptnContext has been set.")
		return errors.New("no keptnContext has been set")
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *apimodels.ExpandedService) error {
		foundRemediation := false
		updatedRemediations := []*apimodels.Remediation{}
		for _, approval := range service.OpenRemediations {
			if approval.KeptnContext == keptnContext {
				foundRemediation = true
				continue
			}
			updatedRemediations = append(updatedRemediations, approval)
		}

		if !foundRemediation {
			return ErrOpenRemediationNotFound
		}
		service.OpenRemediations = updatedRemediations
		return nil
	})

	if err != nil {
		return err
	}

	return mv.projectRepo.UpdateProject(existingProject)
}

func (mv *MongoDBProjectMVRepo) OnSequenceTaskTriggered(event apimodels.KeptnContextExtendedCE) {
	err := mv.UpdateEventOfService(event)
	if err != nil {
		log.WithError(err).Errorf("Could not update lastEvent property for task.started event")
	}
}

func (mv *MongoDBProjectMVRepo) OnSequenceTaskStarted(event apimodels.KeptnContextExtendedCE) {
	err := mv.UpdateEventOfService(event)
	if err != nil {
		log.WithError(err).Errorf("could not update lastEvent property for task.started event")
	}
}

func (mv *MongoDBProjectMVRepo) OnSequenceTaskFinished(event apimodels.KeptnContextExtendedCE) {
	err := mv.UpdateEventOfService(event)
	if err != nil {
		log.WithError(err).Errorf("could not update lastEvent property for task.finished event")
	}
}

func (mv *MongoDBProjectMVRepo) getAllDeploymentTriggeredEvents(eventData *keptnv2.EventData, triggeredID string, keptnContext string) ([]apimodels.KeptnContextExtendedCE, error) {
	stage := eventData.GetStage()
	service := eventData.GetService()
	events, errObj := mv.eventRepo.GetEvents(eventData.GetProject(), common.EventFilter{
		Type:         keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:        &stage,
		Service:      &service,
		ID:           &triggeredID,
		KeptnContext: &keptnContext,
	})
	return events, errObj
}

type serviceUpdateFunc func(service *apimodels.ExpandedService) error

func updateServiceInStage(project *apimodels.ExpandedProject, stage string, service string, fn serviceUpdateFunc) error {
	if project == nil {
		return errors.New("cannot update service in nil project")
	}
	for _, stg := range project.Stages {
		if stg.StageName == stage {
			serviceIndex := -1
			for idx, svc := range stg.Services {
				if svc.ServiceName == service {
					serviceIndex = idx
				}
			}
			if serviceIndex < 0 {
				return errors.New("service not found")
			}
			return fn(stg.Services[serviceIndex])
		}
	}
	return errors.New("stage not found")
}
