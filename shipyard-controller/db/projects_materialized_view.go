package db

import (
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	goutilsmodels "github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/keptn/keptn/shipyard-controller/models"
	"strconv"
	"time"
)

// ErrProjectNotFound indicates that a project has not been found
var ErrProjectNotFound = errors.New("project not found")

// ErrStageNotFound indicates that a stage has not been found
var ErrStageNotFound = errors.New("stage not found")

// ErrServiceNotFound indicates that a service has not been found
var ErrServiceNotFound = errors.New("service not found")

// ErrOpenRemediationNotFound indicates that no open remediation has been found
var ErrOpenRemediationNotFound = errors.New("open remediation not found")

var instance *ProjectsMaterializedView

// EventsRetriever defines the interface for fetching events from the data store
type EventsRetriever interface {
	GetEvents(filter *goutils.EventFilter) ([]*goutilsmodels.KeptnContextExtendedCE, *goutilsmodels.Error)
}

type ProjectsMaterializedView struct {
	ProjectRepo     ProjectRepo
	EventsRetriever EventRepo
}

// GetProjectsMaterializedView returns the materialized view
//TODO:
func GetProjectsMaterializedView() *ProjectsMaterializedView {
	if instance == nil {
		instance = &ProjectsMaterializedView{
			ProjectRepo:     &MongoDBProjectsRepo{},
			EventsRetriever: nil, //TODO
		}
	}
	return instance
}

// CreateProject creates a project
func (mv *ProjectsMaterializedView) CreateProject(prj *models.ExpandedProject) error {
	existingProject, err := mv.GetProject(prj.ProjectName)
	if existingProject != nil {
		return nil
	}
	err = mv.createProject(prj)
	if err != nil {
		return err
	}
	return nil
}

// UpdatedShipyard updates the shipyard of a project
func (mv *ProjectsMaterializedView) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent
	if err := setShipyardVersion(existingProject); err != nil {
		log.Errorf("could not update shipyard version fo project %s: %s"+projectName, err.Error())
	}

	shipyard, err := common.UnmarshalShipyard(shipyardContent)

	for index := range existingProject.Stages {
		existingProject.Stages[index].ParentStages = getParentStages(existingProject.Stages[index].StageName, shipyard)
	}

	return mv.ProjectRepo.UpdateProject(existingProject)
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
func (mv *ProjectsMaterializedView) UpdateProject(prj *models.ExpandedProject) error {
	return mv.ProjectRepo.UpdateProject(prj)
}

func setShipyardVersion(existingProject *models.ExpandedProject) error {
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
func (mv *ProjectsMaterializedView) UpdateUpstreamInfo(projectName string, uri, user string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	if existingProject.GitRemoteURI != uri || existingProject.GitUser != user {
		existingProject.GitRemoteURI = uri
		existingProject.GitUser = user
		if err := mv.ProjectRepo.UpdateProject(existingProject); err != nil {
			log.Errorf("could not update upstream credentials of project %s: %s", projectName, err.Error())
			return err
		}
	}
	return nil
}

func (mv *ProjectsMaterializedView) UpdatedShipyard(projectName string, shipyard string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}

	if existingProject.Shipyard != shipyard {
		existingProject.Shipyard = shipyard
		mv.ProjectRepo.UpdateProject(existingProject)
		if err != nil {
			log.Errorf("could not update shipyard of project %s: %s", projectName, err.Error())
			return nil
		}
	}
	return nil

}

// DeleteUpstreamInfo deletes the Upstream Repository URL and git user of a project
func (mv *ProjectsMaterializedView) DeleteUpstreamInfo(projectName string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}
	if existingProject == nil {
		return nil
	}
	existingProject.GitUser = ""
	existingProject.GitRemoteURI = ""
	if err := mv.ProjectRepo.UpdateProject(existingProject); err != nil {
		log.Errorf("could not delete upstream credentials of project %s: %s", projectName, err.Error())
		return err
	}
	return nil
}

// GetProjects returns all projects
func (mv *ProjectsMaterializedView) GetProjects() ([]*models.ExpandedProject, error) {
	projects, err := mv.ProjectRepo.GetProjects()
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
func (mv *ProjectsMaterializedView) GetProject(projectName string) (*models.ExpandedProject, error) {
	project, err := mv.ProjectRepo.GetProject(projectName)
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
func (mv *ProjectsMaterializedView) DeleteProject(projectName string) error {
	return mv.ProjectRepo.DeleteProject(projectName)
}

// CreateStage creates a stage
func (mv *ProjectsMaterializedView) CreateStage(project string, stage string) error {
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

	prj.Stages = append(prj.Stages, &models.ExpandedStage{
		Services:  []*models.ExpandedService{},
		StageName: stage,
	})

	err = mv.ProjectRepo.UpdateProject(prj)
	if err != nil {
		return err
	}

	log.Infof("Added stage %s to project %s", stage, project)
	return nil
}

func (mv *ProjectsMaterializedView) createProject(project *models.ExpandedProject) error {

	err := mv.ProjectRepo.CreateProject(project)
	if err != nil {
		log.Errorf("Could not create project %s: %s", project.ProjectName, err.Error())
		return err
	}
	return nil
}

// DeleteStage deletes a stage
func (mv *ProjectsMaterializedView) DeleteStage(project string, stage string) error {
	log.Infof("Deleting stage %s from project %s", stage, project)
	prj, err := mv.GetProject(project)

	if err != nil {
		fmt.Sprintf("Could not delete stage %s from project %s : %s\n", stage, project, err.Error())
		return err
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

	err = mv.ProjectRepo.UpdateProject(prj)
	return nil
}

// CreateService creates a service
func (mv *ProjectsMaterializedView) CreateService(project string, stage string, service string) error {
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
			stg.Services = append(stg.Services, &models.ExpandedService{
				CreationDate:  strconv.FormatInt(time.Now().UnixNano(), 10),
				DeployedImage: "",
				ServiceName:   service,
			})
			log.Infof("Adding %s to stage %s in project %s in database", service, stage, project)
			err := mv.ProjectRepo.UpdateProject(existingProject)
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

func (mv *ProjectsMaterializedView) GetService(projectName, stageName, serviceName string) (*models.ExpandedService, error) {
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
func (mv *ProjectsMaterializedView) DeleteService(project string, stage string, service string) error {
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
	err = mv.ProjectRepo.UpdateProject(existingProject)
	if err != nil {
		log.Errorf("Could not delete service %s from stage %s in project %s: %s", service, stage, project, err.Error())
		return err
	}
	log.Infof("Deleted service %s from stage %s in project %s", service, stage, project)
	return nil
}

// UpdateEventOfService updates a service event
func (mv *ProjectsMaterializedView) UpdateEventOfService(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error {
	eventData := &keptnv2.EventData{}
	err := keptnv2.Decode(event, eventData)
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

	contextInfo := &models.EventContext{
		EventID:      eventID,
		KeptnContext: keptnContext,
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	err = updateServiceInStage(existingProject, eventData.Stage, eventData.Service, func(service *models.ExpandedService) error {
		if service.LastEventTypes == nil {
			service.LastEventTypes = map[string]models.EventContext{}
		}
		service.LastEventTypes[eventType] = *contextInfo

		// for events of type "deployment.finished", find the correlating
		// "deployment.triggered" event to update the deployed image name
		if eventType == keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName) {

			events, errObj := mv.getAllDeploymentTriggeredEvents(eventData, triggeredID, keptnContext)
			if errObj != nil {
				return err
			}
			if events == nil || len(events) == 0 {
				return errors.New("No deployment.triggered events could be found for keptn context " + keptnContext)
			}

			triggeredData := &keptnv2.DeploymentTriggeredEventData{}
			err := keptnv2.Decode(events[0].Data, triggeredData)
			if err != nil {
				return errors.New("unable to decode deployment.triggered event data: " + err.Error())
			}

			if deployedImage := triggeredData.ConfigurationChange.Values["image"]; deployedImage != nil {
				service.DeployedImage = fmt.Sprintf("%v", deployedImage)
			}
		}
		return nil
	})

	if err != nil {
		log.Errorf("Could not update image of service %s: %s", eventData.Service, err.Error())
		return err
	}
	err = mv.ProjectRepo.UpdateProject(existingProject)
	if err != nil {
		log.Errorf("Could not update project %s: %s", eventData.Project, err.Error())
		return err
	}
	return nil
}

// CreateRemediation creates a remediation action
func (mv *ProjectsMaterializedView) CreateRemediation(project, stage, service string, remediation *models.Remediation) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not create remediation for service %s in stage %s in project%s. Could not load project: %s", service, stage, project, err.Error())
		return ErrProjectNotFound
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *models.ExpandedService) error {
		if service.OpenRemediations == nil {
			service.OpenRemediations = []*models.Remediation{}
		}
		service.OpenRemediations = append(service.OpenRemediations, remediation)
		return nil
	})
	return mv.ProjectRepo.UpdateProject(existingProject)
}

// CloseOpenRemediations closes a open remediation actions for a given keptnContext
func (mv *ProjectsMaterializedView) CloseOpenRemediations(project, stage, service, keptnContext string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		log.Errorf("Could not close remediation for service %s in stage %s in project %s. Could not load project: %s", service, stage, project, err.Error())
		return ErrProjectNotFound
	}
	if keptnContext == "" {
		log.Warn("No keptnContext has been set.")
		return errors.New("no keptnContext has been set")
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *models.ExpandedService) error {
		foundRemediation := false
		updatedRemediations := []*models.Remediation{}
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

	return mv.ProjectRepo.UpdateProject(existingProject)
}

func (mv *ProjectsMaterializedView) getAllDeploymentTriggeredEvents(eventData *keptnv2.EventData, triggeredID string, keptnContext string) ([]models.Event, error) {
	stage := eventData.GetStage()
	service := eventData.GetService()
	events, errObj := mv.EventsRetriever.GetEvents(eventData.GetProject(), common.EventFilter{
		Type:         keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
		Stage:        &stage,
		Service:      &service,
		ID:           &triggeredID,
		KeptnContext: &keptnContext,
	})
	return events, errObj
}

type serviceUpdateFunc func(service *models.ExpandedService) error

func updateServiceInStage(project *models.ExpandedProject, stage string, service string, fn serviceUpdateFunc) error {
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
