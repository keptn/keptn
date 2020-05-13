package common

import (
	"errors"
	"fmt"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"time"
)

var ProjectNotFoundErr = errors.New("project not found")
var StageNotFoundErr = errors.New("stage not found")
var ServiceNotFoundErr = errors.New("service not found")
var OpenApprovalNotFoundErr = errors.New("open approval not found")

var instance *projectsMaterializedView

type projectsMaterializedView struct {
	ProjectRepo ProjectRepo
	Logger      keptn.LoggerInterface
}

func GetProjectsMaterializedView() *projectsMaterializedView {
	if instance == nil {
		instance = &projectsMaterializedView{
			ProjectRepo: &MongoDBProjectRepo{},
			Logger:      keptn.NewLogger("", "", "configuration-service"),
		}
	}
	return instance
}

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

func (mv *projectsMaterializedView) UpdateShipyard(projectName string, shipyardContent string) error {
	existingProject, err := mv.GetProject(projectName)
	if err != nil {
		return err
	}

	existingProject.Shipyard = shipyardContent

	return mv.updateProject(existingProject)
}

func (mv *projectsMaterializedView) GetProjects() ([]*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProjects()
}

func (mv *projectsMaterializedView) GetProject(projectName string) (*models.ExpandedProject, error) {
	return mv.ProjectRepo.GetProject(projectName)
}

func (mv *projectsMaterializedView) DeleteProject(projectName string) error {
	return mv.ProjectRepo.DeleteProject(projectName)
}

func (mv *projectsMaterializedView) CreateStage(project string, stage string) error {
	fmt.Println("Adding stage " + stage + " to project " + project)
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
					break
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
			}
			break
		}
	}
	mv.Logger.Info("Service " + service + " already exists in stage " + stage + " in project " + project)
	return nil
}

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

func (mv *projectsMaterializedView) UpdateEventOfService(event interface{}, eventType string, keptnContext string, eventID string) error {

	keptnBase := &keptn.KeptnBase{}
	err := mapstructure.Decode(event, keptnBase)
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
		if eventType == keptn.DeploymentFinishedEventType {
			if keptnBase.Image != nil && keptnBase.Tag != nil {
				service.DeployedImage = *keptnBase.Image + ":" + *keptnBase.Tag
			}
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

func (mv *projectsMaterializedView) CreateOpenApproval(project, stage, service string, approval *models.Approval) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could create approval for service " + service + " in stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return ProjectNotFoundErr
	}
	err = updateServiceInStage(existingProject, stage, service, func(service *models.ExpandedService) error {
		if service.OpenApprovals == nil {
			service.OpenApprovals = []*models.Approval{}
		}
		service.OpenApprovals = append(service.OpenApprovals, approval)
		return nil
	})
	return mv.updateProject(existingProject)
}

func (mv *projectsMaterializedView) CloseOpenApproval(project, stage, service, approvalEventID string) error {
	existingProject, err := mv.GetProject(project)
	if err != nil {
		mv.Logger.Error("Could not close approval for service " + service + " in stage " + stage + " in project " + project + ". Could not load project: " + err.Error())
		return ProjectNotFoundErr
	}
	if approvalEventID == "" {
		mv.Logger.Debug("No approvalEventID has been set.")
		return errors.New("no approvalEventID has been set")
	}

	err = updateServiceInStage(existingProject, stage, service, func(service *models.ExpandedService) error {
		foundApproval := false
		updatedApprovals := []*models.Approval{}
		for _, approval := range service.OpenApprovals {
			if approval.EventID == approvalEventID {
				foundApproval = true
				continue
			}
			updatedApprovals = append(updatedApprovals, approval)
		}

		if !foundApproval {
			return OpenApprovalNotFoundErr
		}
		service.OpenApprovals = updatedApprovals
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
