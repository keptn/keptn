package common

import (
	"errors"
	"fmt"
	goutilsmodels "github.com/keptn/go-utils/pkg/api/models"
	goutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/configuration-service/models"
	"os"
	"strconv"
	"time"
)

// ErrProjectNotFound indicates that a project has not been found
var ErrProjectNotFound = errors.New("project not found")

// ErrStageNotFound indicates that a stage has not been found
var ErrStageNotFound = errors.New("stage not found")

// ErrServiceNotFound indicates that a service has not been found
var ErrServiceNotFound = errors.New("service not found")

// ErrOpenRemediationNotFound indeicates that no open remediation has been found
var ErrOpenRemediationNotFound = errors.New("open remediation not found")

var instance *projectsMaterializedView

// EventsRetriever defines the interface for fetching events from the data store
type EventsRetriever interface {
	GetEvents(filter *goutils.EventFilter) ([]*goutilsmodels.KeptnContextExtendedCE, *goutilsmodels.Error)
}

type projectsMaterializedView struct {
	ProjectRepo     ProjectRepo
	EventsRetriever EventsRetriever
	Logger          keptncommon.LoggerInterface
}

// GetProjectsMaterializedView returns the materialized view
func GetProjectsMaterializedView() *projectsMaterializedView {
	fmt.Println(instance)
	if instance == nil {
		instance = &projectsMaterializedView{
			ProjectRepo:     &MongoDBProjectRepo{},
			EventsRetriever: keptnapi.NewEventHandler(os.Getenv("DATASTORE")),
			Logger:          keptncommon.NewLogger("", "", "configuration-service"),
		}
	}
	return instance
}

// CreateProject creates a project
func (mv *projectsMaterializedView) CreateProject(prj *models.Project) error {
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
func (mv *projectsMaterializedView) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent

	return mv.updateProject(existingProject)
}

// UpdateUpstreamInfo updates the Upstream Repository URL and git user of a project
func (mv *projectsMaterializedView) UpdateUpstreamInfo(projectName string, uri, user string) error {
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
		if err := mv.updateProject(existingProject); err != nil {
			mv.Logger.Error(fmt.Sprintf("could not update upstream credentials of project %s: %s", projectName, err.Error()))
			return err
		}
	}
	return nil
}

// GetProjects returns all projects
func (mv *projectsMaterializedView) GetProjects() ([]*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProjects()
}

// GetProject returns a project by its name
func (mv *projectsMaterializedView) GetProject(projectName string) (*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProject(projectName)
}

// DeleteProject deletes a project
func (mv *projectsMaterializedView) DeleteProject(projectName string) error {
	return mv.ProjectRepo.DeleteProject(projectName)
}

// CreateStage creates a stage
func (mv *projectsMaterializedView) CreateStage(project string, stage string) error {
	mv.Logger.Info("Adding stage " + stage + " to project " + project)
	prj, err := mv.GetProject(project)

	if err != nil {
		mv.Logger.Error(fmt.Sprintf("Could not add stage %s to project %s : %s\n", stage, project, err.Error()))
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
		mv.Logger.Info("Stage " + stage + " already exists in project " + project)
		return nil
	}

	prj.Stages = append(prj.Stages, &models.ExpandedStage{
		Services:  []*models.ExpandedService{},
		StageName: stage,
	})

	err = mv.updateProject(prj)
	if err != nil {
		return err
	}

	mv.Logger.Info("Added stage " + stage + " to project " + project)
	return nil
}

func (mv *projectsMaterializedView) createProject(prj *models.Project) error {
	expandedProject := &models.ExpandedProject{
		CreationDate: strconv.FormatInt(time.Now().UnixNano(), 10),
		GitRemoteURI: prj.GitRemoteURI,
		GitUser:      prj.GitUser,
		ProjectName:  prj.ProjectName,
		Shipyard:     "",
		Stages:       nil,
	}

	err := mv.ProjectRepo.CreateProject(expandedProject)
	if err != nil {
		mv.Logger.Error("Could not create project " + prj.ProjectName + ": " + err.Error())
		return err
	}
	return nil
}

func (mv *projectsMaterializedView) updateProject(prj *models.ExpandedProject) error {
	return mv.ProjectRepo.UpdateProject(prj)
}

// DeleteStage deletes a stage
func (mv *projectsMaterializedView) DeleteStage(project string, stage string) error {
	mv.Logger.Info("Deleting stage " + stage + " from project " + project)
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

	err = mv.updateProject(prj)
	return nil
}

// CreateService creates a service
func (mv *projectsMaterializedView) CreateService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return err
	}

	for _, stg := range existingProject.Stages {
		if stg.StageName == stage {
			for _, svc := range stg.Services {
				if svc.ServiceName == service {
					mv.Logger.Info("Service " + service + " already exists in stage " + stage + " in project " + project)
					return nil
				}
			}
			stg.Services = append(stg.Services, &models.ExpandedService{
				CreationDate:  strconv.FormatInt(time.Now().UnixNano(), 10),
				DeployedImage: "",
				ServiceName:   service,
			})
			mv.Logger.Info("Adding " + service + " to stage " + stage + " in project " + project + " in database")
			err := mv.updateProject(existingProject)
			if err != nil {
				mv.Logger.Error("Could not add service " + service + " to stage " + stage + " in project " + project + ". Could not update project: " + err.Error())
				return err
			}
			mv.Logger.Info("Service " + service + " has been added to stage " + stage + " in project " + project)
			break
		}
	}

	return nil
}

// DeleteService deletes a service
func (mv *projectsMaterializedView) DeleteService(project string, stage string, service string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could not delete service " + service + " from stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
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
				mv.Logger.Info("Could not delete service " + service + " from stage " + stage + " in project " + project + ". Service not found in database")
				return nil
			}
			copy(stg.Services[serviceIndex:], stg.Services[serviceIndex+1:])
			stg.Services[len(stg.Services)-1] = nil
			stg.Services = stg.Services[:len(stg.Services)-1]
			break
		}
	}
	err = mv.updateProject(existingProject)
	if err != nil {
		mv.Logger.Error("Could not delete service " + service + " from stage " + stage + " in project " + project + ": " + err.Error())
		return err
	}
	mv.Logger.Info("Deleted service " + service + " from stage " + stage + " in project " + project)
	return nil
}

// UpdateEventOfService updates a service event
func (mv *projectsMaterializedView) UpdateEventOfService(event interface{}, eventType string, keptnContext string, eventID string, triggeredID string) error {

	keptnBase := &keptnv2.EventData{}
	err := keptnv2.DecodeKeptnEventData(event, keptnBase)
	//err := mapstructure.Decode(event, keptnBase)
	if err != nil {
		mv.Logger.Error("Could not parse event data: " + err.Error())
		return err
	}

	existingProject, err := mv.GetProject(keptnBase.Project)
	if err != nil {
		mv.Logger.Error("Could not update service " + keptnBase.Service + " in stage " + keptnBase.Stage + " in project " + keptnBase.Project + ". Could not load project: " + err.Error())
		return err
	}

	contextInfo := &models.EventContext{
		EventID:      eventID,
		KeptnContext: keptnContext,
		Time:         strconv.FormatInt(time.Now().UnixNano(), 10),
	}
	err = updateServiceInStage(existingProject, keptnBase.Stage, keptnBase.Service, func(service *models.ExpandedService) error {
		if service.LastEventTypes == nil {
			service.LastEventTypes = map[string]models.EventContext{}
		}

		if eventType == keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName) {

			events, errObj := mv.EventsRetriever.GetEvents(&keptnapi.EventFilter{
				Project:      keptnBase.GetProject(),
				Stage:        keptnBase.GetStage(),
				Service:      keptnBase.GetService(),
				EventType:    keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName),
				KeptnContext: keptnContext,
			})

			if errObj != nil || events == nil || len(events) == 0 {
				return errors.New(*errObj.Message)
			}

			var matchingTriggeredEvent *goutilsmodels.KeptnContextExtendedCE = nil
			for _, e := range events {
				if e.Triggeredid == triggeredID {
					matchingTriggeredEvent = e
					break
				}
			}

			triggeredData := keptnv2.DeploymentTriggeredEventData{}
			err := keptnv2.DecodeKeptnEventData(matchingTriggeredEvent.Data, &triggeredData)
			if err != nil {
				return err
			}

			deployedImage := triggeredData.ConfigurationChange.Values["image"]
			service.DeployedImage = fmt.Sprintf("%v", deployedImage)
		}
		service.LastEventTypes[eventType] = *contextInfo
		return nil
	})

	if err != nil {
		mv.Logger.Error("Could not update image of service " + keptnBase.Service + ": " + err.Error())
		return err
	}
	err = mv.updateProject(existingProject)
	if err != nil {
		mv.Logger.Error("Could not update " + keptnBase.Project + ": " + err.Error())
		return err
	}
	return nil
}

// CreateRemediation creates a remediation action
func (mv *projectsMaterializedView) CreateRemediation(project, stage, service string, remediation *models.Remediation) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could not create remediation for service " + service + " in stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return ErrProjectNotFound
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *models.ExpandedService) error {
		if service.OpenRemediations == nil {
			service.OpenRemediations = []*models.Remediation{}
		}
		service.OpenRemediations = append(service.OpenRemediations, remediation)
		return nil
	})
	return mv.updateProject(existingProject)
}

// CloseOpenRemediations closes a open remediation actions for a given keptnContext
func (mv *projectsMaterializedView) CloseOpenRemediations(project, stage, service, keptnContext string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could not close remediation for service " + service + " in stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return ErrProjectNotFound
	}
	if keptnContext == "" {
		mv.Logger.Debug("No keptnContext has been set.")
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

	return mv.updateProject(existingProject)
}

type serviceUpdateFunc func(service *models.ExpandedService) error

func updateServiceInStage(project *models.ExpandedProject, stage string, service string, fn serviceUpdateFunc) error {
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
